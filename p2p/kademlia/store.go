package kademlia

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
)

// Store is the interface for implementing the storage mechanism for the DHT.
type Store interface {
	// Store a key/value pair for the local node with the replication and expiration
	Store(ctx context.Context, key []byte, data []byte, replication time.Time, expiration time.Time) error

	// Retrieve the local key/value from store
	Retrieve(ctx context.Context, key []byte) ([]byte, error)

	// Delete a key/value pair from the store
	Delete(ctx context.Context, key []byte)

	// Keys returns all the keys from the store
	Keys(ctx context.Context) [][]byte

	// KeysForReplication returns the keys of all data to be replicated across the network
	KeysForReplication(ctx context.Context) [][]byte

	// ExpireKeys expires all key/values
	ExpireKeys(ctx context.Context)
}

// MemoryStore is a simple in-memory key/value store used for unit testing
type MemoryStore struct {
	mutex        sync.Mutex
	self         *Node
	data         map[string][]byte
	replications map[string]time.Time
	expirations  map[string]time.Time
}

// KeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (ms *MemoryStore) KeysForReplication(_ context.Context) [][]byte {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	var keys [][]byte
	for k := range ms.data {
		if time.Now().After(ms.replications[k]) {
			keys = append(keys, []byte(k))
		}
	}
	return keys
}

// ExpireKeys should expire all key/values due for expiration
func (ms *MemoryStore) ExpireKeys(_ context.Context) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	for k, v := range ms.expirations {
		if time.Now().After(v) {
			delete(ms.replications, k)
			delete(ms.expirations, k)
			delete(ms.data, k)
		}
	}
}

// Init initializes the Store
func NewMemStore() *MemoryStore {
	return &MemoryStore{
		data:         make(map[string][]byte),
		replications: make(map[string]time.Time),
		expirations:  make(map[string]time.Time),
	}
}

func (ms *MemoryStore) SetSelf(self *Node) {
	ms.self = self
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (ms *MemoryStore) Store(ctx context.Context, key []byte, data []byte, replication time.Time, expiration time.Time) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.replications[string(key)] = replication
	ms.expirations[string(key)] = expiration
	ms.data[string(key)] = data

	log.WithContext(ctx).Debugf("id: %s, store key: %s, data: %v", ms.self.String(), hex.EncodeToString(key), hex.EncodeToString(data))
	return nil
}

// Retrieve will return the local key/value if it exists
func (ms *MemoryStore) Retrieve(ctx context.Context, key []byte) ([]byte, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	data, ok := ms.data[string(key)]
	if !ok {
		return nil, nil
	}

	log.WithContext(ctx).Debugf("id: %s, retrieve key: %s, data: %v", ms.self.String(), hex.EncodeToString(key), hex.EncodeToString(data))
	return data, nil
}

// Delete deletes a key/value pair from the MemoryStore
func (ms *MemoryStore) Delete(_ context.Context, key []byte) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	delete(ms.replications, string(key))
	delete(ms.expirations, string(key))
	delete(ms.data, string(key))
}

// Keys returns all the keys from the Store
func (ms *MemoryStore) Keys(_ context.Context) [][]byte {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	var keys [][]byte
	for k := range ms.data {
		keys = append(keys, []byte(k))
	}
	return keys
}
