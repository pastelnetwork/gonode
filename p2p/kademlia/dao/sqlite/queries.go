package sqlite

import (
	"context"
	"database/sql"
	"time"

	// add sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/p2p/kademlia/crypto"
)

// Migrate runs migrations such as creating `keys` table.
func Migrate(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	query := "CREATE TABLE IF NOT EXISTS `keys` (`uid` INTEGER PRIMARY KEY AUTOINCREMENT, `key` VARCHAR(64) NULL, `data` BLOB NULL, `replication` DATE NULL, `expiration` DATE NULL)"

	if _, err := db.ExecContext(ctx, query); err != nil {
		return errors.Errorf("could not create \"keys\" table: %w", err).WithField("query", query)
	}
	return nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func Store(ctx context.Context, db *sql.DB, data []byte, replication, expiration time.Time) (int64, error) {
	key := crypto.GetKey(data)
	query := "INSERT INTO keys(key, data, replication, expiration) VALUES (?, ?, ?, ?)"
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return 0, errors.Errorf("could not prepare statement: %w", err).WithField("query", query)
	}

	res, err := stmt.ExecContext(ctx, string(key), data, replication, expiration)
	if err != nil {
		return 0, errors.Errorf("could not insert keys: %w", err)
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := stmt.Close(); err != nil {
		return 0, errors.Errorf("could not close statement: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Errorf("could not get the number of affected rows: %w", err)
	}

	return rows, nil
}

// Retrieve will return the local key/value if it exists
func Retrieve(ctx context.Context, db *sql.DB, key []byte) ([]byte, error) {
	var data []byte
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	query := "SELECT data FROM keys WHERE key=? LIMIT 1"
	rows, err := db.QueryContext(ctx, query, string(key))
	if err != nil {
		return []byte(""), errors.Errorf("could not retrieve data for key=%q: %w", string(key), err).WithField("query", query)
	}

	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			return []byte(""), errors.Errorf("could not scan columns: %w", err)
		}
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		return []byte(""), errors.Errorf("could not close the Rows: %w", err)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return []byte(""), errors.Errorf("could not scann rows: %w", err)
	}

	return data, nil
}

// ExpireKeys should expire all key/values due for expiration.
func ExpireKeys(ctx context.Context, db *sql.DB) error {
	query := "DELETE FROM keys WHERE expiration <= ?"
	_, err := db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return errors.Errorf("could not delete expired keys: %w", err).WithField("query", query)
	}

	return nil
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func GetAllKeysForReplication(ctx context.Context, db *sql.DB) ([][]byte, error) {
	var closerr error
	var keys [][]byte

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	query := "SELECT key FROM keys WHERE replication <= ?"
	rows, err := db.QueryContext(ctx, query, time.Now())
	if err != nil {
		return [][]byte{}, errors.Errorf("could not retrieve keys: %w", err).WithField("query", query)
	}

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			return [][]byte{}, errors.Errorf("could not scan columns: %w", err)
		}
		keys = append(keys, []byte(key))
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		return [][]byte{}, errors.Errorf("could not close the Rows: %w", err)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return [][]byte{}, errors.Errorf("could not scann rows: %w", err)
	}

	return keys, closerr
}

// Remove deletes a key/value pair from the Key
func Remove(ctx context.Context, db *sql.DB, key []byte) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()
	query := "DELETE FROM keys WHERE key=?"
	_, err := db.ExecContext(ctx, query, string(key))
	if err != nil {
		return errors.Errorf("could not delete key=%q: %w", string(key), err).WithField("query", query)
	}

	return nil
}
