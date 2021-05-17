package kademlia

import (
	"time"
)

// Store is the interface for implementing the storage mechanism for the
// DHT.
type Store interface {
	// Store should store a key/value pair for the local node with the
	// given replication and expiration times.
	Store(key []byte, data []byte, replication time.Time, expiration time.Time, publisher bool) error

	// Retrieve should return the local key/value if it exists.
	Retrieve(key []byte) (data []byte, found bool)

	// Delete should delete a key/value pair from the Store
	Delete(key []byte) error

	// Init initializes the Store
	Init() error

	// GetAllKeysForReplication should return the keys of all data to be
	// replicated across the network. Typically all data should be
	// replicated every tReplicate seconds.
	GetAllKeysForReplication() ([][]byte, error)

	// ExpireKeys should expire all key/values due for expiration.
	ExpireKeys() error

	// GetKey returns the key for data
	GetKey(data []byte) ([]byte, error)
}
