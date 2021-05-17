package kademlia

import (
	"crypto/sha1"
	"sync"
	"time"
)

// MemoryStore is a simple in-memory key/value store used for unit testing, and
// the CLI example
type MemoryStore struct {
	mutex        *sync.Mutex
	data         map[string][]byte
	replicateMap map[string]time.Time
	expireMap    map[string]time.Time
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (ms *MemoryStore) GetAllKeysForReplication() ([][]byte, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	var keys [][]byte
	for k := range ms.data {
		if time.Now().After(ms.replicateMap[k]) {
			keys = append(keys, []byte(k))
		}
	}
	return keys, nil
}

// ExpireKeys should expire all key/values due for expiration.
func (ms *MemoryStore) ExpireKeys() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	for k, v := range ms.expireMap {
		if time.Now().After(v) {
			delete(ms.replicateMap, k)
			delete(ms.expireMap, k)
			delete(ms.data, k)
		}
	}
	return nil
}

// Init initializes the Store
func (ms *MemoryStore) Init() error {
	ms.data = make(map[string][]byte)
	ms.mutex = &sync.Mutex{}
	ms.replicateMap = make(map[string]time.Time)
	ms.expireMap = make(map[string]time.Time)
	return nil
}

// GetKey returns the key for data
func (ms *MemoryStore) GetKey(data []byte) ([]byte, error) {
	sha := sha1.Sum(data)
	return sha[:], nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (ms *MemoryStore) Store(key []byte, data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.replicateMap[string(key)] = replication
	ms.expireMap[string(key)] = expiration
	ms.data[string(key)] = data
	return nil
}

// Retrieve will return the local key/value if it exists
func (ms *MemoryStore) Retrieve(key []byte) (data []byte, found bool) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	data, found = ms.data[string(key)]
	return data, found
}

// Delete deletes a key/value pair from the MemoryStore
func (ms *MemoryStore) Delete(key []byte) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	delete(ms.replicateMap, string(key))
	delete(ms.expireMap, string(key))
	delete(ms.data, string(key))
	return nil
}
