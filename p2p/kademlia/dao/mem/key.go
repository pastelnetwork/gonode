package mem

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/p2p/kademlia/crypto"
)

// Key is a simple in-memory key/value store used for unit testing, and
// the CLI example
type Key struct {
	mutex        *sync.Mutex
	data         map[string][]byte
	replicateMap map[string]time.Time
	expireMap    map[string]time.Time
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (k *Key) GetAllKeysForReplication(ctx context.Context) ([][]byte, error) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	var keys [][]byte
	for kk := range k.data {
		if time.Now().After(k.replicateMap[kk]) {
			keys = append(keys, []byte(kk))
		}
	}
	return keys, nil
}

// ExpireKeys should expire all key/values due for expiration.
func (k *Key) ExpireKeys(ctx context.Context) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	for kk, v := range k.expireMap {
		if time.Now().After(v) {
			delete(k.replicateMap, kk)
			delete(k.expireMap, kk)
			delete(k.data, kk)
		}
	}
	return nil
}

// Init initializes the Store
func (k *Key) Init(ctx context.Context) error {
	k.data = make(map[string][]byte)
	k.mutex = &sync.Mutex{}
	k.replicateMap = make(map[string]time.Time)
	k.expireMap = make(map[string]time.Time)
	return nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (k *Key) Store(ctx context.Context, data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	key := crypto.GetKey(data)
	k.mutex.Lock()
	defer k.mutex.Unlock()
	k.replicateMap[string(key)] = replication
	k.expireMap[string(key)] = expiration
	k.data[string(key)] = data
	return nil
}

// Retrieve will return the local key/value if it exists
func (k *Key) Retrieve(ctx context.Context, key []byte) (data []byte, found bool) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	data, found = k.data[string(key)]
	return data, found
}

// Delete deletes a key/value pair from the MemoryStore
func (k *Key) Delete(ctx context.Context, key []byte) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	delete(k.replicateMap, string(key))
	delete(k.expireMap, string(key))
	delete(k.data, string(key))
	return nil
}
