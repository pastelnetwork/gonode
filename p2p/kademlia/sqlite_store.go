package kademlia

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore is a simple in-memory key/value store used for unit testing, and
// the CLI example
type SQLiteStore struct {
	db *sql.DB
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (ss *SQLiteStore) GetAllKeysForReplication() ([][]byte, error) {
	return getAllKeysForReplication(ss.db)
}

// ExpireKeys should expire all key/values due for expiration.
func (ss *SQLiteStore) ExpireKeys() error {
	return expireKeys(ss.db)
}

// Init initializes the Store
func (ss *SQLiteStore) Init() error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	defer db.Close()

	if err := migrate(ss.db); err != nil {
		return err
	}

	ss.db = db
	return nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (ss *SQLiteStore) Store(key []byte, data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	return store(ss.db, key, data, replication, expiration)
}

// Retrieve will return the local key/value if it exists
func (ss *SQLiteStore) Retrieve(key []byte) (data []byte, found bool) {
	data, err := retrieve(ss.db, key)
	if err != nil {
		return data, false
	}

	return data, true
}

// Delete deletes a key/value pair from the MemoryStore
func (ss *SQLiteStore) Delete(key []byte) error {
	return nil
}
