package kademlia

import (
	"crypto/sha1"
	"sync"
	"time"
	        _ "github.com/mattn/go-sqlite3"

)

// SQLiteStore is a simple in-memory key/value store used for unit testing, and
// the CLI example
type SQLiteStore struct {
	db			 *sql.DB
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (ss *SQLiteStore) GetAllKeysForReplication() [][]byte {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	var keys [][]byte
	for k := range ms.data {
		if time.Now().After(ms.replicateMap[k]) {
			keys = append(keys, []byte(k))
		}
	}
	return keys
}

// ExpireKeys should expire all key/values due for expiration.
func (ss *SQLiteStore) ExpireKeys() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	for k, v := range ms.expireMap {
		if time.Now().After(v) {
			delete(ms.replicateMap, k)
			delete(ms.expireMap, k)
			delete(ms.data, k)
		}
	}
}

// Init initializes the Store
func (ss *SQLiteStore) Init() error {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS `keys` (
        `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
        `key` VARCHAR(64) NULL,
        `replication` DATE NULL
        `expiration` DATE NULL"
    );
   	if err != nil {
		return err
	}

	ss.db = db
	return nil
}


// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (ss *SQLiteStore) Store(key []byte, data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	return nil
}


// Retrieve will return the local key/value if it exists
func (ss *SQLiteStore) Retrieve(key []byte) (data []byte, found bool) {
	return []byte, false
}

// Delete deletes a key/value pair from the MemoryStore
func (ss *SQLiteStore) Delete(key []byte) {
}