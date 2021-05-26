package sqlite

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/crypto"
)

var (
	logPrefix = "p2p"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	_, err := db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS `keys` (`uid` INTEGER PRIMARY KEY AUTOINCREMENT, `key` VARCHAR(64) NULL, `data` BLOB NULL, `replication` DATE NULL, `expiration` DATE NULL)")
	if err != nil {
		log.WithContext(ctx).Infof("Error %s when creating keys table", err)
		return err
	}

	return nil
}

func Store(ctx context.Context, db *sql.DB, data []byte, replication, expiration time.Time) (int64, error) {
	key := crypto.GetKey(data)
	query := "INSERT INTO keys(key, data, replication, expiration) VALUES (?, ?, ?, ?)"

	ctx = log.ContextWithPrefix(ctx, logPrefix)

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.WithContext(ctx).Infof("Error %s when preparing SQL statement", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, string(key), data, replication, expiration)
	if err != nil {
		log.WithContext(ctx).Infof("Error %s when inserting row into keys table", err)
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.WithContext(ctx).Infof("Error %s when finding rows affected", err)
		return 0, err
	}

	return rows, nil
}

func Retrieve(ctx context.Context, db *sql.DB, key []byte) ([]byte, error) {
	var data []byte

	ctx = log.ContextWithPrefix(ctx, logPrefix)

	rows, err := db.QueryContext(ctx, "SELECT data FROM keys WHERE key=? LIMIT 1", string(key))
	if err != nil {
		log.WithContext(ctx).Infof("Error %s while selecting key=%s", err, string(key))
		return []byte(""), err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			log.WithContext(ctx).Infof("Error %s when scanning rows", err)
			return []byte(""), err
		}
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		log.WithContext(ctx).Infof("Error %s while closing rows", err)
		return []byte(""), err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.WithContext(ctx).Infof("Error %s after scanning rows", err)
		return []byte(""), err
	}

	return data, nil
}

func ExpireKeys(ctx context.Context, db *sql.DB) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	_, err := db.ExecContext(ctx, "DELETE FROM keys WHERE expiration > TIME('now')")
	if err != nil {
		log.WithContext(ctx).Infof("Error %s while keys with due expiration", err)
		return err
	}

	return nil
}

func GetAllKeysForReplication(ctx context.Context, db *sql.DB) ([][]byte, error) {
	var keys [][]byte
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	rows, err := db.QueryContext(ctx, "SELECT key FROM keys WHERE replication > TIME('now')")
	if err != nil {
		log.WithContext(ctx).Infof("Error %s while selecting keys for replication", err)
		return [][]byte{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.WithContext(ctx).Infof("Error %s while scanning rows", err)
			return [][]byte{}, err
		}
		keys = append(keys, []byte(key))
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		log.WithContext(ctx).Infof("Error %s while closing rows", err)
		return [][]byte{}, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.WithContext(ctx).Infof("Error %s after scanning rows", err)
		return [][]byte{}, err
	}

	return keys, nil
}

func Remove(ctx context.Context, db *sql.DB, key []byte) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	_, err := db.ExecContext(ctx, "DELETE FROM keys WHERE key=?", string(key))
	if err != nil {
		log.WithContext(ctx).Infof("Error %s while deleting key=%s", err, string(key))
		return err
	}

	return nil
}
