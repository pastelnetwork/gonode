package sqlite

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

type nodeReplicationInfo struct {
	LastReplicated *time.Time `db:"lastReplicatedAt"`
	UpdatedAt      time.Time  `db:"updatedAt"`
	CreatedAt      time.Time  `db:"createdAt"`
	Active         bool       `db:"is_active"`
	Adjusted       bool       `db:"is_adjusted"`
	IP             string     `db:"ip"`
	Port           int        `db:"port"`
	ID             string     `db:"id"`
	LastSeen       time.Time  `db:"last_seen"`
}

type repKeys struct {
	Key       string    `db:"key"`
	UpdatedAt time.Time `db:"updatedAt"`
	IP        string    `db:"ip"`
	Port      int       `db:"port"`
	ID        string    `db:"id"`
}

func (r repKeys) toDomain() (domain.ToRepKey, error) {
	key, err := hex.DecodeString(r.Key)
	if err != nil {
		return domain.ToRepKey{}, fmt.Errorf("error decoding key: %w", err)
	}

	return domain.ToRepKey{
		Key:       key,
		UpdatedAt: r.UpdatedAt,
		IP:        r.IP,
		Port:      r.Port,
		ID:        r.ID,
	}, nil
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

func (s *Store) migrateRepKeys() error {
	replicateQuery := `
    CREATE TABLE IF NOT EXISTS replication_keys(
        key TEXT PRIMARY KEY,
		id TEXT NOT NULL,
        ip TEXT NOT NULL,
		port INTEGER NOT NULL,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `

	if _, err := s.db.Exec(replicateQuery); err != nil {
		return fmt.Errorf("failed to create table 'replication_keys': %w", err)
	}

	return nil
}

func (s *Store) checkReplicateStore() bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='replication_info'`
	var name string
	err := s.db.Get(&name, query)
	return err == nil
}

func (s *Store) ensureLastSeenColumn() error {
	rows, err := s.db.Query("PRAGMA table_info(replication_info)")
	if err != nil {
		return fmt.Errorf("failed to fetch table 'replication_info' info: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid, notnull, pk int
		var name, dtype string
		var dfltValue *string
		err = rows.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		if name == "last_seen" {
			return nil
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during iteration: %w", err)
	}

	_, err = s.db.Exec(`ALTER TABLE data ADD COLUMN last_seen DATETIME DEFAULT CURRENT_TIMESTAMP`)
	if err != nil {
		return fmt.Errorf("failed to add column 'last_seen' to table 'data': %w", err)
	}

	return nil
}

func (s *Store) checkReplicateKeysStore() bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='replication_keys'`
	var name string
	err := s.db.Get(&name, query)
	return err == nil
}

// GetAllToDoRepKeys returns all keys that need to be replicated
func (s *Store) GetAllToDoRepKeys() (retKeys domain.ToRepKeys, err error) {
	var keys []repKeys
	if err := s.db.Select(&keys, "SELECT * FROM replication_keys LIMIT 5000"); err != nil {
		return nil, fmt.Errorf("error reading all keys from database: %w", err)
	}

	retKeys = make(domain.ToRepKeys, 0, len(keys))

	for _, key := range keys {
		domainKey, err := key.toDomain()
		if err != nil {
			return nil, fmt.Errorf("error converting key to domain: %w", err)
		}
		retKeys = append(retKeys, domainKey)
	}

	return retKeys, nil
}

// DeleteRepKey will delete a key from the replication_keys table
func (s *Store) DeleteRepKey(key []byte) error {
	hkey := hex.EncodeToString(key)

	res, err := s.db.Exec("DELETE FROM replication_keys WHERE key = ?", hkey)
	if err != nil {
		return fmt.Errorf("cannot delete record from replication_keys table by key %s: %v", hkey, err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("cannot delete record from replication_keys table by key %s: %v", hkey, err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("cannot delete record from replication_keys table by key %s ", hkey)
	}

	return nil
}

// StoreBatchRepKeys will store a batch of values with their SHA256 hash as the key
func (s *Store) StoreBatchRepKeys(values [][]byte, id string, ip string, port int) error {
	operation := func() error {
		tx, err := s.db.Beginx()
		if err != nil {
			return fmt.Errorf("cannot begin transaction: %w", err)
		}

		// Prepare insert statement
		stmt, err := tx.PrepareNamed(`INSERT INTO replication_keys(key, id, ip, port, updatedAt) values(:key, :id, :ip, :port, :updatedAt) ON CONFLICT(key) DO UPDATE SET id=:id,ip=:ip,port=:port,updatedAt=:updatedAt`)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("statement preparation failed, rollback failed: %v, original error: %w", rollbackErr, err)
			}
			return fmt.Errorf("cannot prepare statement: %w", err)

		}
		defer stmt.Close()

		// For each value, calculate its hash and insert into DB
		now := time.Now().UTC()
		for i := 0; i < len(values); i++ {
			hkey := hex.EncodeToString(values[i])
			r := repKeys{Key: hkey, UpdatedAt: now, ID: id, IP: ip, Port: port}

			// Execute the insert statement
			_, err = stmt.Exec(r)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("cannot insert or update replication key record with key %s: %w", hkey, err)
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return fmt.Errorf("cannot commit transaction replication key update: %w", err)
		}

		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second

	err := backoff.Retry(operation, b)
	if err != nil {
		return fmt.Errorf("error storing rep keys record data: %w", err)
	}

	return nil
}

// GetKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) GetKeysForReplication(ctx context.Context, from time.Time, to time.Time) (retkeys [][]byte) {
	var keys []string
	if err := s.db.Select(&keys, `SELECT key FROM data WHERE createdAt > ? and createdAt < ?`, from, to); err != nil {
		log.P2P().WithError(err).WithContext(ctx).Errorf("failed to get records for replication older than %s", from)
		return nil
	}

	for i := 0; i < len(keys); i++ {
		str, err := hex.DecodeString(keys[i])
		if err != nil {
			log.WithContext(ctx).WithField("key", keys[i]).Error("replicate failed to hex decode key")
			continue
		}

		retkeys = append(retkeys, str)
	}

	return retkeys
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
	_, err := s.db.Exec(`INSERT INTO replication_info(id, ip, is_active, is_adjusted, lastReplicatedAt, updatedAt, port, last_seen) values(?, ?, ?, ?, ?, ?, ?, ?)`,
		string(rep.ID), rep.IP, rep.Active, rep.IsAdjusted, rep.LastReplicatedAt, rep.UpdatedAt, rep.Port, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert replicate record: %v", err)
	}

	return err
}

// UpdateLastSeen updates last seen
func (s *Store) UpdateLastSeen(_ context.Context, id string) error {
	_, err := s.db.Exec(`UPDATE replication_info SET last_seen = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last_seen of node: %s: - err: %v", id, err)
	}

	return err
}
