package sqlite

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/common/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
)

const (
	dbName = "data001.sqlite3"
)

// Store is a sqlite based storage
type Store struct {
	db                *sqlx.DB
	replicateInterval time.Duration
	republishInterval time.Duration

	// for stats
	dbFilePath string
	rwMtx      *sync.RWMutex
}

// Record is a data record
type Record struct {
	Key           string
	Data          []byte
	IsOriginal    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReplicatedAt  time.Time
	RepublishedAt time.Time
}

// NewStore returns a new store
func NewStore(ctx context.Context, dataDir string, replicate time.Duration, republish time.Duration) (*Store, error) {
	s := &Store{
		rwMtx: &sync.RWMutex{},
	}

	log.P2P().WithContext(ctx).Debugf("p2p data dir: %v", dataDir)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0750); err != nil {
			return nil, fmt.Errorf("mkdir %q: %w", dataDir, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cannot create data folder: %w", err)
	}

	dbFile := path.Join(dataDir, dbName)
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite database: %w", err)
	}

	s.db = db
	s.replicateInterval = replicate
	s.republishInterval = republish
	s.dbFilePath = dbFile

	if !s.checkStore() {
		if err = s.migrate(); err != nil {
			return nil, fmt.Errorf("cannot create table(s) in sqlite database: %w", err)
		}
	}

	if !s.checkReplicateStore() {
		if err = s.migrateReplication(); err != nil {
			return nil, fmt.Errorf("cannot create table(s) in sqlite database: %w", err)
		}
	}

	return s, nil
}

func (s *Store) checkStore() bool {
	s.rwMtx.RLock()
	defer s.rwMtx.RUnlock()

	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='data'`
	var name string
	err := s.db.Get(&name, query)
	return err == nil
}

func (s *Store) checkReplicateStore() bool {
	s.rwMtx.RLock()
	defer s.rwMtx.RUnlock()

	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='replication_info'`
	var name string
	err := s.db.Get(&name, query)
	return err == nil
}

func (s *Store) migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS data(
        key TEXT PRIMARY KEY,
        data BLOB NOT NULL,
        is_original BOOL DEFAULT FALSE,
        createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        replicatedAt DATETIME,
        republishedAt DATETIME
    );
    `

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create table 'data': %w", err)
	}

	return nil
}

func (s *Store) migrateReplication() error {
	replicateQuery := `
    CREATE TABLE IF NOT EXISTS replication_info(
        id TEXT PRIMARY KEY,
        ip TEXT NOT NULL,
		port INTEGER NOT NULL,
        is_active BOOL DEFAULT FALSE,
		is_adjusted BOOL DEFAULT FALSE,
        lastReplicatedAt DATETIME,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
        updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `

	if _, err := s.db.Exec(replicateQuery); err != nil {
		return fmt.Errorf("failed to create table 'replication_info': %w", err)
	}

	return nil
}

type nodeReplicationInfo struct {
	LastReplicated *time.Time `db:"lastReplicatedAt"`
	UpdatedAt      time.Time  `db:"updatedAt"`
	CreatedAt      time.Time  `db:"createdAt"`
	Active         bool       `db:"is_active"`
	Adjusted       bool       `db:"is_adjusted"`
	IP             string     `db:"ip"`
	Port           int        `db:"port"`
	ID             string     `db:"id"`
}

func (n *nodeReplicationInfo) toDomain() domain.NodeReplicationInfo {
	return domain.NodeReplicationInfo{
		LastReplicatedAt: n.LastReplicated,
		UpdatedAt:        n.UpdatedAt,
		CreatedAt:        n.CreatedAt,
		Active:           n.Active,
		IsAdjusted:       n.Adjusted,
		IP:               n.IP,
		Port:             n.Port,
		ID:               []byte(n.ID),
	}
}

// Store will store a key/value pair for the local node
func (s *Store) Store(_ context.Context, key []byte, value []byte) error {
	s.rwMtx.Lock()
	defer s.rwMtx.Unlock()

	if len(key) == 0 {
		return fmt.Errorf("key cannot be empty")
	}

	if len(value) == 0 {
		return fmt.Errorf("value cannot be empty")
	}

	hkey := hex.EncodeToString(key)

	now := time.Now().UTC()
	r := Record{Key: hkey, Data: value, UpdatedAt: now}
	res, err := s.db.NamedExec(`INSERT INTO data(key, data) values(:key, :data) ON CONFLICT(key) DO UPDATE SET data=:data,updatedAt=:updatedat`, r)
	if err != nil {
		return fmt.Errorf("cannot insert or update record with key %s: %w", hkey, err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("failed to insert/update record with key %s: %w", hkey, err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("failed to insert/update record with key %s", hkey)
	}

	return nil
}

// Retrieve will return the local key/value if it exists
func (s *Store) Retrieve(_ context.Context, key []byte) ([]byte, error) {
	s.rwMtx.RLock()
	defer s.rwMtx.RUnlock()

	hkey := hex.EncodeToString(key)

	r := Record{}
	err := s.db.Get(&r, `SELECT data FROM data WHERE key = ?`, hkey)
	if err != nil {
		return nil, fmt.Errorf("failed to get record by key %s: %w", hkey, err)
	}
	return r.Data, nil
}

// Delete a key/value pair from the Store
func (s *Store) Delete(ctx context.Context, key []byte) {
	s.rwMtx.Lock()
	defer s.rwMtx.Unlock()

	hkey := hex.EncodeToString(key)

	res, err := s.db.Exec("DELETE FROM data WHERE key = ?", hkey)
	if err != nil {
		log.P2P().WithContext(ctx).Debugf("cannot delete record by key %s: %v", hkey, err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		log.P2P().WithContext(ctx).Debugf("failed to delete record by key %s: %v", hkey, err)
	} else if rowsAffected == 0 {
		log.P2P().WithContext(ctx).Debugf("failed to delete record by key %s", hkey)
	}
}

// UpdateKeyReplication updates the replication time for a key
func (s *Store) UpdateKeyReplication(_ context.Context, key []byte) error {
	s.rwMtx.Lock()
	defer s.rwMtx.Unlock()

	keyStr := hex.EncodeToString(key)
	_, err := s.db.Exec(`UPDATE data SET replicatedAt = ? WHERE key = ?`, time.Now(), keyStr)
	if err != nil {
		return fmt.Errorf("failed to update replicated records: %v", err)
	}

	return err
}

// GetKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) GetKeysForReplication(ctx context.Context, from time.Time) [][]byte {
	s.rwMtx.RLock()
	defer s.rwMtx.RUnlock()

	var keys [][]byte
	if err := s.db.Select(&keys, `SELECT key FROM data WHERE createdAt > ?`, from); err != nil {
		log.P2P().WithError(err).WithContext(ctx).Errorf("failed to get records for replication older than %s", from)
		return nil
	}
	log.P2P().WithContext(ctx).WithField("keyslength", len(keys)).Debugf("Replication keys found: %+v", keys)

	var unhexedKeys [][]byte
	for _, key := range keys {
		dst := make([]byte, hex.DecodedLen(len(key)))
		_, err := hex.Decode(dst, key)
		if err != nil {
			log.P2P().WithContext(ctx).WithField("hkey", key).Errorf("failed to properly unhex hkey: %v", err)
			continue
		}

		unhexedKeys = append(unhexedKeys, dst)
	}

	return unhexedKeys
}

// Close the store
func (s *Store) Close(ctx context.Context) {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.P2P().WithContext(ctx).Errorf("Failed to close database: %s", err)
		}
	}
}

// Stats returns stats of store
func (s *Store) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	fi, err := os.Stat(s.dbFilePath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get p2p db size")
	} else {
		stats["p2p_db_size"] = utils.BytesToMB(uint64(fi.Size()))
	}

	if count, err := s.Count(ctx); err == nil {
		stats["p2p_db_records_count"] = count
	} else {
		log.WithContext(ctx).WithError(err).Error("failed to get p2p records count")
	}

	return stats, nil
}

// Count the records in store
func (s *Store) Count(_ context.Context /*, type RecordType*/) (int, error) {
	s.rwMtx.RLock()
	defer s.rwMtx.RUnlock()

	var count int
	err := s.db.Get(&count, `SELECT COUNT(*) FROM data`)
	if err != nil {
		return -1, fmt.Errorf("failed to get count of records: %w", err)
	}

	return count, nil
}

// DeleteAll the records in store
func (s *Store) DeleteAll(_ context.Context /*, type RecordType*/) error {
	s.rwMtx.Lock()
	defer s.rwMtx.Unlock()

	res, err := s.db.Exec("DELETE FROM data")
	if err != nil {
		return fmt.Errorf("cannot delete ALL records: %w", err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("failed to delete ALL records: %w", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("failed to delete ALL records")
	}
	return nil
}

// GetAllReplicationInfo returns all records in replication table
func (s *Store) GetAllReplicationInfo(_ context.Context) ([]domain.NodeReplicationInfo, error) {
	r := []nodeReplicationInfo{}
	err := s.db.Select(&r, "SELECT * FROM replication_info")
	if err != nil {
		return nil, fmt.Errorf("failed to get replication info %w", err)
	}

	list := []domain.NodeReplicationInfo{}
	for _, v := range r {
		list = append(list, v.toDomain())
	}

	return list, nil
}

// UpdateReplicationInfo updates replication info
func (s *Store) UpdateReplicationInfo(_ context.Context, rep domain.NodeReplicationInfo) error {
	_, err := s.db.Exec(`UPDATE replication_info SET ip = ?, is_active = ?, is_adjusted = ?, lastReplicatedAt = ?, updatedAt =?, port = ? WHERE id = ?`,
		rep.IP, rep.Active, rep.IsAdjusted, rep.LastReplicatedAt, rep.UpdatedAt, rep.Port, string(rep.ID))
	if err != nil {
		return fmt.Errorf("failed to update replicated records: %v", err)
	}

	return err
}

// AddReplicationInfo adds replication info
func (s *Store) AddReplicationInfo(_ context.Context, rep domain.NodeReplicationInfo) error {
	_, err := s.db.Exec(`INSERT INTO replication_info(id, ip, is_active, is_adjusted, lastReplicatedAt, updatedAt, port) values(?, ?, ?, ?, ?, ?, ?)`,
		string(rep.ID), rep.IP, rep.Active, rep.IsAdjusted, rep.LastReplicatedAt, rep.UpdatedAt, rep.Port)
	if err != nil {
		return fmt.Errorf("failed to update replicate record: %v", err)
	}

	return err
}
