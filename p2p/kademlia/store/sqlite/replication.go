package sqlite

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
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
	LastSeen       *time.Time `db:"last_seen"`
}

type repKeys struct {
	Key       string    `db:"key"`
	UpdatedAt time.Time `db:"updatedAt"`
	IP        string    `db:"ip"`
	Port      int       `db:"port"`
	ID        string    `db:"id"`
	Attempts  int       `db:"attempts"`
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
		Attempts:  r.Attempts,
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
		attempts INTEGER DEFAULT 0,
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

func (s *Store) ensureAttempsColumn() error {
	rows, err := s.db.Query("PRAGMA table_info(replication_keys)")
	if err != nil {
		return fmt.Errorf("failed to fetch table 'replication_keys' info: %w", err)
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

		if name == "attempts" {
			return nil
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during iteration: %w", err)
	}

	_, err = s.db.Exec(`ALTER TABLE replication_keys ADD COLUMN attempts INTEGER DEFAULT 0`)
	if err != nil {
		return fmt.Errorf("failed to add column 'attempts' to table 'data': %w", err)
	}

	return nil
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

	_, err = s.db.Exec(`ALTER TABLE replication_info ADD COLUMN last_seen DATETIME DEFAULT NULL`)
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
func (s *Store) GetAllToDoRepKeys(minAttempts, maxAttempts int) (retKeys domain.ToRepKeys, err error) {
	var keys []repKeys
	if err := s.db.Select(&keys, "SELECT * FROM replication_keys where attempts <= ? AND attempts >= ? LIMIT 5000", maxAttempts, minAttempts); err != nil {
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
	b.MaxElapsedTime = 1 * time.Minute

	err := backoff.Retry(operation, b)
	if err != nil {
		return fmt.Errorf("error storing rep keys record data: %w", err)
	}

	return nil
}

// GetKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) GetKeysForReplication(ctx context.Context, from time.Time, to time.Time) [][]byte {
	var keys []string
	if err := s.db.Select(&keys, `SELECT key FROM data WHERE createdAt > ? and createdAt < ?`, from, to); err != nil {
		log.P2P().WithError(err).WithContext(ctx).Errorf("failed to get records for replication older than %s", from)
		return nil
	}

	retkeys := make([][]byte, len(keys))
	var wg sync.WaitGroup
	for i, key := range keys {
		wg.Add(1)
		go func(i int, key string) {
			defer wg.Done()
			str, err := hex.DecodeString(key)
			if err != nil {
				log.WithContext(ctx).WithField("key", key).Error("replicate failed to hex decode key")
			}
			retkeys[i] = str
		}(i, key)
	}

	wg.Wait()
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

// RetrieveBatchNotExist returns a list of keys (hex-decoded) that do not exist in the table
func (s *Store) RetrieveBatchNotExist(ctx context.Context, keys [][]byte, batchSize int) ([][]byte, error) {
	var nonExistingKeys [][]byte
	keyMap := make(map[string]bool)

	// Convert the keys to hex and map them for easier lookups
	for i := 0; i < len(keys); i++ {
		hexKey := hex.EncodeToString(keys[i])
		keyMap[hexKey] = true
	}

	batchCount := (len(keys) + batchSize - 1) / batchSize // Round up division

	for i := 0; i < batchCount; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(keys) {
			end = len(keys)
		}
		batchKeys := keys[start:end]

		placeholders := make([]string, len(batchKeys))
		args := make([]interface{}, len(batchKeys))

		for j, key := range batchKeys {
			placeholders[j] = "?"
			args[j] = hex.EncodeToString(key)
		}

		query := fmt.Sprintf(`SELECT key FROM data WHERE key IN (%s)`, strings.Join(placeholders, ","))
		rows, err := s.db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve records: %w", err)
		}

		for rows.Next() {
			var key string
			if err := rows.Scan(&key); err != nil {
				rows.Close()
				return nil, fmt.Errorf("failed to scan key: %w", err)
			}

			// Remove existing keys from the map
			delete(keyMap, key)
		}

		if err := rows.Err(); err != nil {
			rows.Close() // ensure rows are closed in case of error
			return nil, fmt.Errorf("rows processing error: %w", err)
		}

		rows.Close() // explicitly close the rows

	}

	// Convert remaining keys in the map (which do not exist in the DB) to hex-decoded byte slices
	for key := range keyMap {
		hexDecodedKey, err := hex.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex key %s: %w", key, err)
		}
		nonExistingKeys = append(nonExistingKeys, hexDecodedKey)
	}

	return nonExistingKeys, nil
}

// RetrieveBatchValues returns a list of values (hex-decoded) for the given keys (hex-encoded)
func (s *Store) RetrieveBatchValues(ctx context.Context, keys []string) ([][]byte, int, error) {
	placeholders := make([]string, len(keys))
	args := make([]interface{}, len(keys))
	keyToIndex := make(map[string]int)

	for i := 0; i < len(keys); i++ {
		placeholders[i] = "?"
		args[i] = keys[i]
		keyToIndex[keys[i]] = i
	}

	log.WithContext(ctx).WithField("len(keys)", len(args)).WithField("len(placeholders)", len(placeholders)).
		WithField("args[len(args)]-1", fmt.Sprint(args[len(args)-1])).WithField("args[0", fmt.Sprint(args[0])).
		Info("RetrieveBatchValues db operation")

	query := fmt.Sprintf(`SELECT key, data FROM data WHERE key IN (%s)`, strings.Join(placeholders, ","))
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve records: %w", err)
	}
	defer rows.Close()

	values := make([][]byte, len(keys))
	keysFound := 0
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, keysFound, fmt.Errorf("failed to scan key and value: %w", err)
		}

		if idx, found := keyToIndex[key]; found {
			values[idx] = value
			keysFound++
		}
	}

	if err := rows.Err(); err != nil {
		return nil, keysFound, fmt.Errorf("rows processing error: %w", err)
	}

	return values, keysFound, nil
}

// BatchDeleteRepKeys will delete a list of keys from the replication_keys table
func (s *Store) BatchDeleteRepKeys(keys []string) error {
	var placeholders []string
	var arguments []interface{}
	for _, key := range keys {
		placeholders = append(placeholders, "?")
		arguments = append(arguments, key)
	}

	query := fmt.Sprintf("DELETE FROM replication_keys WHERE key IN (%s)", strings.Join(placeholders, ","))

	res, err := s.db.Exec(query, arguments...)
	if err != nil {
		return fmt.Errorf("cannot batch delete records from replication_keys table: %w", err)
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("cannot get affected rows after batch deleting records from replication_keys table: %w", err)
	} else if rowsAffected == 0 {
		return fmt.Errorf("no record deleted (batch) from replication_keys table")
	}

	return nil
}

// IncrementAttempts increments the attempts counter for the given keys.
func (s *Store) IncrementAttempts(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// Use a transaction for consistency and to speed up the operation
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	query := `UPDATE replication_keys SET attempts = attempts + 1 WHERE key IN (?)`
	query, args, err := sqlx.In(query, keys)
	if err != nil {
		return err
	}

	// Rebind is necessary to match placeholders style to the SQL driver in use
	// (it converts "?" placeholders into the correct placeholder style, e.g., "$1, $2" for PostgreSQL)
	query = tx.Rebind(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		// rollback the transaction in case of error
		_ = tx.Rollback()
		return err
	}

	// commit the transaction
	return tx.Commit()
}
