package sqlite

import (
	"context"
	"database/sql"
	"time"

	// add sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/crypto"

	"github.com/pastelnetwork/gonode/common/errors"
)

// Migrate runs migrations such as creating `keys` table.
func Migrate(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	_, err := db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS `keys` (`uid` INTEGER PRIMARY KEY AUTOINCREMENT, `key` VARCHAR(64) NULL, `data` BLOB NULL, `replication` DATE NULL, `expiration` DATE NULL)")
	if err != nil {
		return errors.Errorf("Error %s when creating keys table", err)
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
		return 0, errors.Errorf("failed to prepare SQL statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, string(key), data, replication, expiration)
	if err != nil {
		return 0, errors.Errorf("Error %s when inserting row into keys table", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Errorf("Error %s when finding rows affected", err)
	}

	return rows, nil
}

// Retrieve will return the local key/value if it exists
func Retrieve(ctx context.Context, db *sql.DB, key []byte) ([]byte, error) {
	var data []byte
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT data FROM keys WHERE key=? LIMIT 1", string(key))
	if err != nil {
		log.WithContext(ctx).Fatalf("Error %s while selecting key=%s", err, string(key))
		return []byte(""), errors.New(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			log.WithContext(ctx).Fatalf("Error %s when scanning rows", err)
			return []byte(""), errors.New(err)
		}
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		log.WithContext(ctx).Fatalf("Error %s while closing rows", err)
		return []byte(""), errors.New(err)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.WithContext(ctx).Fatalf("Error %s after scanning rows", err)
		return []byte(""), errors.New(err)
	}

	return data, nil
}

// ExpireKeys should expire all key/values due for expiration.
func ExpireKeys(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "DELETE FROM keys WHERE expiration <= ?", time.Now())
	if err != nil {
		return errors.Errorf("Error %s while keys with due expiration", err)
	}

	return nil
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func GetAllKeysForReplication(ctx context.Context, db *sql.DB) ([][]byte, error) {
	var keys [][]byte

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT key FROM keys WHERE replication <= ?", time.Now())
	if err != nil {
		log.WithContext(ctx).Fatalf("Error %s while selecting keys for replication", err)
		return [][]byte{}, errors.New(err)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			return [][]byte{}, errors.Errorf("Error %s while scanning rows", err)
		}
		keys = append(keys, []byte(key))
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		return [][]byte{}, errors.Errorf("Error %s while closing rows", err)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return [][]byte{}, errors.Errorf("Error %s after scanning rows", err)
	}

	return keys, nil
}

// Remove deletes a key/value pair from the Key
func Remove(ctx context.Context, db *sql.DB, key []byte) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()
	_, err := db.ExecContext(ctx, "DELETE FROM keys WHERE key=?", string(key))
	if err != nil {
		return errors.Errorf("Error %s while deleting key=%s", err, string(key))
	}

	return nil
}
