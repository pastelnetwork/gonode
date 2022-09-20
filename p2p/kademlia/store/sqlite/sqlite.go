package sqlite

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"time"

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
	s := &Store{}

	log.P2P().WithContext(ctx).Debugf("p2p data dir: %v", dataDir)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0750); err != nil {
			return nil, fmt.Errorf("mkdir %q: %w", dataDir, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cannot create data folder: %w", err)
	}

	dbFile := path.Join(dataDir, dbName)
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("%s?%s", dbFile, "_journal_mode=WAL"))
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

	return s, nil
}

func (s *Store) checkStore() bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='data'`
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

// Store will store a key/value pair for the local node
func (s *Store) Store(_ context.Context, key []byte, value []byte) error {
	hkey := hex.EncodeToString(key)

	now := time.Now().UTC()
	r := Record{Key: hkey, Data: value, UpdatedAt: now}
	res, err := s.db.NamedExec(`INSERT INTO data(key, data) values(:key, :data) ON CONFLICT(key) DO UPDATE SET updatedAt=:updatedat`, r)
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
	hkey := hex.EncodeToString(key)

	res, err := s.db.Exec("DELETE FROM data WHERE key = ?", hkey)
	if err != nil {
		log.P2P().WithContext(ctx).Errorf("cannot delete record by key %s: %v", hkey, err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		log.P2P().WithContext(ctx).Errorf("failed to delete record by key %s: %v", hkey, err)
	} else if rowsAffected == 0 {
		log.P2P().WithContext(ctx).Errorf("failed to delete record by key %s", hkey)
	}
}

// GetKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) GetKeysForReplication(ctx context.Context) [][]byte {
	now := time.Now().UTC()
	after := now.Add(-s.replicateInterval).UTC()

	var keys [][]byte
	if err := s.db.Select(&keys, `SELECT key FROM data WHERE updatedAt < ? and (replicatedAt is NULL OR replicatedAt < ?)`, after, after); err != nil {
		log.P2P().WithContext(ctx).Errorf("failed to get records for replication older then %s: %v", after, err)
		return nil
	}
	log.P2P().WithContext(ctx).WithField("keyslength", len(keys)).Debugf("Replication keys found: %+v", keys)
	if len(keys) > 0 {
		_, err := s.db.Exec(`UPDATE data SET replicatedAt = ? WHERE updatedAt < ? and (replicatedAt is NULL OR replicatedAt < ?)`, now, after, after)
		if err != nil {
			log.P2P().WithContext(ctx).Errorf("failed to update replicated records: %v", err)
		}
	}
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
		stats["p2p_db_size"] = fi.Size()
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
	var count int
	err := s.db.Get(&count, `SELECT COUNT(*) FROM data`)
	if err != nil {
		return -1, fmt.Errorf("failed to get count of records: %w", err)
	}

	return count, nil
}

// DeleteAll the records in store
func (s *Store) DeleteAll(_ context.Context /*, type RecordType*/) error {
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
