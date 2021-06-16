package memory

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
)

// Store is a simple in-memory key/value store used for unit testing
type Store struct {
	mutex        sync.Mutex
	data         map[string][]byte
	replications map[string]time.Time
	expirations  map[string]time.Time
}

// KeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) KeysForReplication(_ context.Context) [][]byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var keys [][]byte
	for k := range s.data {
		if time.Now().After(s.replications[k]) {
			keys = append(keys, []byte(k))
		}
	}
	return keys
}

// ExpireKeys should expire all key/values due for expiration
func (s *Store) ExpireKeys(_ context.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for k, v := range s.expirations {
		if time.Now().After(v) {
			delete(s.replications, k)
			delete(s.expirations, k)
			delete(s.data, k)
		}
	}
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (s *Store) Store(ctx context.Context, key []byte, data []byte, replication time.Time, expiration time.Time) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.replications[string(key)] = replication
	s.expirations[string(key)] = expiration
	s.data[string(key)] = data

	log.WithContext(ctx).Debugf("store key: %s, data: %v", hex.EncodeToString(key), hex.EncodeToString(data))
	return nil
}

// Retrieve will return the local key/value if it exists
func (s *Store) Retrieve(ctx context.Context, key []byte) ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, ok := s.data[string(key)]
	if !ok {
		return nil, nil
	}

	log.WithContext(ctx).Debugf("retrieve key: %s, data: %v", hex.EncodeToString(key), hex.EncodeToString(data))
	return data, nil
}

// Delete a key/value pair from the Store
func (s *Store) Delete(_ context.Context, key []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.replications, string(key))
	delete(s.expirations, string(key))
	delete(s.data, string(key))
}

// Keys returns all the keys from the Store
func (s *Store) Keys(_ context.Context) [][]byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var keys [][]byte
	for k := range s.data {
		keys = append(keys, []byte(k))
	}
	return keys
}

// Close the store
func (s *Store) Close(_ context.Context) {
}

// NewStore returns a new memory store
func NewStore() *Store {
	return &Store{
		data:         make(map[string][]byte),
		replications: make(map[string]time.Time),
		expirations:  make(map[string]time.Time),
	}
}
