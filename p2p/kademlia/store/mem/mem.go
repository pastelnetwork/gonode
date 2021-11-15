package mem

import (
	"bytes"
	"context"
	"errors"
	"sort"
	"sync"
	"time"

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
func (s *Store) Store(_ context.Context, key []byte, value []byte, replication time.Time) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.replications[string(key)] = replication
	s.data[string(key)] = value

	return nil
}

// Retrieve will return the local key/value if it exists
func (s *Store) Retrieve(_ context.Context, key []byte) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[string(key)]
	if !ok {
		return nil, errors.New("not found")
	}

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
// return all keys of limit is -1
func (s *Store) Keys(_ context.Context, offset int, limit int) [][]byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if offset < 0 {
		offset = 0
	}

	if limit == -1 {
		limit = len(s.data)
	}

	if offset >= len(s.data) {
		return [][]byte{}
	}

	var keys [][]byte
	for k := range s.data {
		keys = append(keys, []byte(k))
	}

	// Sort keys
	sort.SliceStable(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})

	if offset+limit > len(keys) {
		return keys[offset:]
	}

	return keys[offset : offset+limit]
}

// Close the store
func (s *Store) Close(_ context.Context) {
}

// Stats returns stats of store
func (s *Store) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	stats["record_count"] = len(s.Keys(ctx, 0, -1))
	return stats, nil
}

// NewStore returns a new memory store
func NewStore() *Store {
	return &Store{
		data:         make(map[string][]byte),
		replications: make(map[string]time.Time),
	}
}

// InitCleanup is suppossed to initiate grabage cleanup
// not applicable on this test implementation
func (s *Store) InitCleanup(ctx context.Context, _ time.Duration) {
	log.WithContext(ctx).Error("s.InitCleanup not implemented")
}
