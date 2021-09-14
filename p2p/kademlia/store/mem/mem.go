package mem

import (
	"context"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/log"
)

// Store is a simple in-memory key/value store used for unit testing
type Store struct {
	mutex        sync.RWMutex
	data         map[string][]byte
	replications map[string]time.Time
}

// KeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) KeysForReplication(_ context.Context) [][]byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var keys [][]byte
	for k := range s.data {
		if time.Now().After(s.replications[k]) {
			keys = append(keys, []byte(k))
		}
	}
	return keys
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (s *Store) Store(ctx context.Context, key []byte, value []byte, replication time.Time) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.replications[string(key)] = replication
	s.data[string(key)] = value

	log.WithContext(ctx).Debugf("store key: %s", base58.Encode(key))
	return nil
}

// Retrieve will return the local key/value if it exists
func (s *Store) Retrieve(ctx context.Context, key []byte) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[string(key)]
	if !ok {
		return nil, nil
	}

	log.WithContext(ctx).Debugf("retrieve key: %s", base58.Encode(key))
	return value, nil
}

// Delete a key/value pair from the Store
func (s *Store) Delete(_ context.Context, key []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.replications, string(key))
	delete(s.data, string(key))
}

// Keys returns all the keys from the Store
func (s *Store) Keys(_ context.Context) [][]byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var keys [][]byte
	for k := range s.data {
		keys = append(keys, []byte(k))
	}
	return keys
}

// Close the store
func (s *Store) Close(_ context.Context) {
}

// Stats returns stats of store
func (s *Store) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	stats["record_count"] = len(s.Keys(ctx))
	return stats, nil
}

// NewStore returns a new memory store
func NewStore() *Store {
	return &Store{
		data:         make(map[string][]byte),
		replications: make(map[string]time.Time),
	}
}
