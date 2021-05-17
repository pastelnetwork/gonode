package kademlia

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func migrate(ctx context.Context, db *sql.DB) error {
	_, err = db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS `keys` (`uid` INTEGER PRIMARY KEY AUTOINCREMENT, `key` VARCHAR(64) NULL, `data` BLOB NULL, `replication` DATE NULL, `expiration` DATE NULL")
	if err != nil {
		log.Printf("Error %s when creating keys table", err)
		return err
	}
}

func store(db *sql.DB, key []byte, data []byte, replication, expiration time.Time) error {
	query := "INSERT INTO keys(key, data, replication, expiration) VALUES (?, ?, ?, ?)"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, string(key), data, replication, expiration)
	if err != nil {
		log.Printf("Error %s when inserting row into keys table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}

	return nil
}

func retrieve(db *sql.DB, key []byte) ([]byte, error) {
	var data []byte
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT data FROM keys WHERE key=? LIMIT 1", string(key))
	if err != nil {
		log.Printf("Error %s while selecting key=%s", err, string(key))
		return []byte{}, err
	}
	defer rows.Close()

	if err := rows.Scan(&data); err != nil {
		log.Printf("Error %s when scanning rows", err)
		return []byte{}, err
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		log.Printf("Error %s while closing rows", err)
		return []byte{}, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Printf("Error %s after scanning rows", rerr)
		return []byte{}, err
	}

	return data, nil
}

func expireKeys(db *sql.DB) error {
	var data []byte
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "UPDATE FROM keys where expiration > ?")
	if err != nil {
		log.Printf("Error %s while keys with due expiration", err)
		return err
	}

	return nil
}

func getAllKeysForReplication(db *sql.DB) ([][]byte, error) {
	var keys [][]byte
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT key FROM keys WHERE expiration > ?", time.Now())
	if err != nil {
		log.Printf("Error %s while selecting keys for replication", err)
		return [][]byte{}, err
	}
	defer rows.Close()

	if err := rows.Scan(&data); err != nil {
		log.Printf("Error %s when scanning rows", err)
		return [][]byte{}, err
	}

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		names = append(keys, key)
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err := rows.Close(); err != nil {
		log.Printf("Error %s while closing rows", err)
		return [][]byte{}, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Printf("Error %s after scanning rows", rerr)
		return [][]byte{}, err
	}
}
