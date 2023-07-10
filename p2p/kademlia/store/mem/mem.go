package mem

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

// Store is a simple in-memory key/value store used for unit testing
type Store struct {
	mutex             sync.RWMutex
	data              map[string][]byte
	replications      map[string]time.Time
	replicateInterval time.Duration
	//republish time.Duration
}

// GetKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (s *Store) GetKeysForReplication(_ context.Context, _ time.Time) [][]byte {
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

// UpdateKeyReplication updates the replication status of the key
func (s *Store) UpdateKeyReplication(_ context.Context, _ []byte) error {
	return nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (s *Store) Store(_ context.Context, key []byte, value []byte, _ int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.replications[string(key)] = time.Now().Add(s.replicateInterval).UTC()
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

// Close the store
func (s *Store) Close(_ context.Context) {
}

// Stats returns stats of store
func (s *Store) Stats(_ context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	return stats, nil
}

// Count the records in store
func (s *Store) Count(_ context.Context /*, type RecordType*/) (int, error) {
	return len(s.data), nil
}

// DeleteAll the records in store
func (s *Store) DeleteAll(_ context.Context /*, type RecordType*/) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = make(map[string][]byte)
	s.replications = make(map[string]time.Time)
	return nil
}

// StoreBatch stores a batch of key/value pairs for the local node with the given
func (s *Store) StoreBatch(_ context.Context, values [][]byte, _ int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, value := range values {
		s.data[string(value)] = value
	}
	return nil
}

// RetrieveWithType will return the local key/value if it exists
func (s *Store) RetrieveWithType(_ context.Context, key []byte) ([]byte, int, error) {
	return []byte{}, 0, nil
}

// NewStore returns a new memory store
func NewStore() *Store {
	return &Store{
		data:              make(map[string][]byte),
		replications:      make(map[string]time.Time),
		replicateInterval: time.Second * 3600,
	}
}

// GetAllReplicationInfo returns all records in replication table
func (s *Store) GetAllReplicationInfo(_ context.Context) ([]domain.NodeReplicationInfo, error) {
	return nil, nil
}

// UpdateReplicationInfo updates replication info
func (s *Store) UpdateReplicationInfo(_ context.Context, _ domain.NodeReplicationInfo) error {
	return nil
}

// AddReplicationInfo adds replication info
func (s *Store) AddReplicationInfo(_ context.Context, _ domain.NodeReplicationInfo) error {
	return nil
}
