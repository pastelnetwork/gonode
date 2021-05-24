package sqlite

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore is a simple in-memory key/value store used for unit testing, and
// the CLI example
type Key struct {
	db *sql.DB
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (k *Key) GetAllKeysForReplication() ([][]byte, error) {
	return GetAllKeysForReplication(k.db)
}

// ExpireKeys should expire all key/values due for expiration.
func (k *Key) ExpireKeys() error {
	return ExpireKeys(k.db)
}

// Init initializes the Store
func (k *Key) Init() error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	defer db.Close()

	if err := Migrate(k.db); err != nil {
		return err
	}

	k.db = db
	return nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (k *Key) Store(data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	_, err := Store(k.db, data, replication, expiration)
	return err
}

// Retrieve will return the local key/value if it exists
func (k *Key) Retrieve(key []byte) (data []byte, found bool) {
	data, err := Retrieve(k.db, key)
	if err != nil {
		return data, false
	}

	return data, true
}

// Delete deletes a key/value pair from the MemoryStore
func (k *Key) Delete(key []byte) error {
	return Remove(k.db, key)
}
