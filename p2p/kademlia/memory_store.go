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
func (ms *MemoryStore) GetAllKeysForReplication() [][]byte {
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
func (ms *MemoryStore) ExpireKeys() {
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
func (ms *MemoryStore) Init() {
	ms.data = make(map[string][]byte)
	ms.mutex = &sync.Mutex{}
	ms.replicateMap = make(map[string]time.Time)
	ms.expireMap = make(map[string]time.Time)
}

// GetKey returns the key for data
func (ms *MemoryStore) GetKey(data []byte) []byte {
	sha := sha1.Sum(data)
	return sha[:]
}
