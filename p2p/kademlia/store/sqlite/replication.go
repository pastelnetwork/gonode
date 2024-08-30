package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
		LastSeen:         n.LastSeen,
	}
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
func (s *Store) GetAllToDoRepKeys(minAttempts, maxAttempts int) (domain.ToRepKeys, error) {
	var keys []domain.ToRepKey
	if err := s.db.Select(&keys, "SELECT * FROM replication_keys where attempts <= ? AND attempts >= ? LIMIT 10000", maxAttempts, minAttempts); err != nil {
		return nil, fmt.Errorf("error reading all keys from database: %w", err)
	}

	return domain.ToRepKeys(keys), nil
}

// DeleteRepKey will delete a key from the replication_keys table
func (s *Store) DeleteRepKey(hkey string) error {
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
func (s *Store) StoreBatchRepKeys(values []string, id string, ip string, port int) error {
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
			r := domain.ToRepKey{Key: values[i], UpdatedAt: now, ID: id, IP: ip, Port: port}

			// Execute the insert statement
			_, err = stmt.Exec(r)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("cannot insert or update replication key record with key %s: %w", values[i], err)
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
func (s *Store) GetKeysForReplication(ctx context.Context, from time.Time, to time.Time) domain.KeysWithTimestamp {
	var results []domain.KeyWithTimestamp
	query := `SELECT key, createdAt FROM data WHERE createdAt > ? AND createdAt < ? ORDER BY createdAt ASC`
	if err := s.db.Select(&results, query, from, to); err != nil {
		log.P2P().WithError(err).WithContext(ctx).Errorf("failed to get records for replication between %s and %s", from, to)
		return nil
	}

	return domain.KeysWithTimestamp(results)
}

// GetAllReplicationInfo returns all records in replication table
func (s *Store) GetAllReplicationInfo(_ context.Context) ([]domain.NodeReplicationInfo, error) {
	r := []nodeReplicationInfo{}
	err := s.db.Select(&r, `SELECT * FROM replication_info 
	ORDER BY 
		CASE WHEN lastReplicatedAt IS NULL THEN 0 ELSE 1 END,
		lastReplicatedAt
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get replication info: %w", err)
	}

	list := make([]domain.NodeReplicationInfo, len(r))
	for i := range r {
		list[i] = r[i].toDomain()
	}

	return list, nil
}

// RecordExists checks if a record with the given nodeID exists in the database.
func (s *Store) RecordExists(nodeID string) (bool, error) {
	const query = "SELECT 1 FROM replication_info WHERE id = ? LIMIT 1"

	var exists int
	err := s.db.QueryRow(query, nodeID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			// No matching record found
			return false, nil
		}
		// Other database error
		return false, fmt.Errorf("error checking record existence: %v", err)
	}

	return exists == 1, nil
}

// UpdateReplicationInfo updates replication info
func (s *Store) UpdateReplicationInfo(_ context.Context, rep domain.NodeReplicationInfo) error {
	if rep.IP == "" || len(rep.ID) == 0 || rep.Port == 0 {
		return fmt.Errorf("invalid replication info: %v", rep)
	}

	_, err := s.db.Exec(`UPDATE replication_info SET ip = ?, is_active = ?, is_adjusted = ?, lastReplicatedAt = ?, updatedAt =?, port = ?, last_seen = ? WHERE id = ?`,
		rep.IP, rep.Active, rep.IsAdjusted, rep.LastReplicatedAt, rep.UpdatedAt, rep.Port, string(rep.ID), rep.LastSeen)
	if err != nil {
		return fmt.Errorf("failed to update replication info: %v", err)
	}

	return err
}

// AddReplicationInfo adds replication info
func (s *Store) AddReplicationInfo(_ context.Context, rep domain.NodeReplicationInfo) error {
	if rep.IP == "" || len(rep.ID) == 0 || rep.Port == 0 {
		return fmt.Errorf("invalid replication info: %v", rep)
	}

	_, err := s.db.Exec(`INSERT INTO replication_info(id, ip, is_active, is_adjusted, lastReplicatedAt, updatedAt, port, last_seen) values(?, ?, ?, ?, ?, ?, ?, ?)`,
		string(rep.ID), rep.IP, rep.Active, rep.IsAdjusted, rep.LastReplicatedAt, rep.UpdatedAt, rep.Port, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to insert replicate record: %v", err)
	}

	return err
}

// UpdateLastSeen updates last seen
func (s *Store) UpdateLastSeen(_ context.Context, id string) error {
	_, err := s.db.Exec(`UPDATE replication_info SET last_seen = ? WHERE id = ?`, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update last_seen of node: %s: - err: %v", id, err)
	}

	return err
}

// UpdateLastReplicated updates replication info last replicated
func (s *Store) UpdateLastReplicated(_ context.Context, id string, t time.Time) error {
	_, err := s.db.Exec(`UPDATE replication_info SET lastReplicatedAt = ?, updatedAt = ? WHERE id = ?`, t, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update last replicated: %v", err)
	}

	return err
}

// UpdateIsAdjusted updates adjusted
func (s *Store) UpdateIsAdjusted(_ context.Context, id string, isAdjusted bool) error {
	_, err := s.db.Exec(`UPDATE replication_info SET is_adjusted = ?, updatedAt = ? WHERE id = ?`, isAdjusted, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update is_adjusted of node: %s: - err: %v", id, err)
	}

	return err
}

// UpdateIsActive updates active
func (s *Store) UpdateIsActive(_ context.Context, id string, isActive bool, isAdjusted bool) error {
	_, err := s.db.Exec(`UPDATE replication_info SET is_active = ?, is_adjusted = ?, updatedAt = ? WHERE id = ?`, isActive, isAdjusted, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update is_active of node: %s: - err: %v", id, err)
	}

	return err
}

// RetrieveBatchNotExist returns a list of keys (hex-decoded) that do not exist in the table
func (s *Store) RetrieveBatchNotExist(ctx context.Context, keys []string, batchSize int) ([]string, error) {
	keyMap := make(map[string]bool)

	// Convert the keys to map them for easier lookups
	for i := 0; i < len(keys); i++ {
		keyMap[keys[i]] = true
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
			args[j] = key
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
	nonExistingKeys := make([]string, 0, len(keyMap))
	for key := range keyMap {
		nonExistingKeys = append(nonExistingKeys, key)

	}

	return nonExistingKeys, nil
}

func (s Store) RetrieveBatchValues(ctx context.Context, keys []string, getFromCloud bool) ([][]byte, int, error) {
	return retrieveBatchValues(ctx, s.db, keys, getFromCloud, s)
}

// RetrieveBatchValues returns a list of values  for the given keys (hex-encoded)
func retrieveBatchValues(ctx context.Context, db *sqlx.DB, keys []string, getFromCloud bool, s Store) ([][]byte, int, error) {
	placeholders := make([]string, len(keys))
	args := make([]interface{}, len(keys))
	keyToIndex := make(map[string]int)

	for i := 0; i < len(keys); i++ {
		placeholders[i] = "?"
		args[i] = keys[i]
		keyToIndex[keys[i]] = i
	}

	query := fmt.Sprintf(`SELECT key, data, is_on_cloud FROM data WHERE key IN (%s)`, strings.Join(placeholders, ","))
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve records: %w", err)
	}
	defer rows.Close()

	values := make([][]byte, len(keys))
	var cloudKeys []string
	keysFound := 0

	for rows.Next() {
		var key string
		var value []byte
		var is_on_cloud bool
		if err := rows.Scan(&key, &value, &is_on_cloud); err != nil {
			return nil, keysFound, fmt.Errorf("failed to scan key and value: %w", err)
		}

		if idx, found := keyToIndex[key]; found {
			values[idx] = value
			keysFound++

			if s.IsCloudBackupOn() && !s.migrationStore.isSyncInProgress {
				if len(value) == 0 && is_on_cloud {
					cloudKeys = append(cloudKeys, key)
				}
				PostAccessUpdate([]string{key})
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, keysFound, fmt.Errorf("rows processing error: %w", err)
	}

	if len(cloudKeys) > 0 && getFromCloud {
		// Fetch from cloud
		cloudValues, err := s.cloud.FetchBatch(cloudKeys)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("failed to fetch from cloud")
		}

		for key, value := range cloudValues {
			if idx, found := keyToIndex[key]; found {
				values[idx] = value
				keysFound++
			}
		}

		go func() {
			datList := make([][]byte, 0, len(cloudValues))
			for _, v := range cloudValues {
				datList = append(datList, v)
			}

			// Store the fetched data in the local store
			if err := s.StoreBatch(ctx, datList, 0, false); err != nil {
				log.WithError(err).Error("failed to store fetched data in local store")
			}
		}()
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
