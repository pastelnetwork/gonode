package dao

import (
	"context"
	"time"
)

// Key is the interface for implementing the storage mechanism for the
// DHT.
type Key interface {
	// Store should store a key/value pair for the local node with the
	// given replication and expiration times.
	Store(ctx context.Context, data []byte, replication time.Time, expiration time.Time, publisher bool) error

	// Retrieve should return the local key/value if it exists.
	Retrieve(ctx context.Context, key []byte) (data []byte, found bool)

	// Delete should delete a key/value pair from the Store
	Delete(ctx context.Context, key []byte) error

	// Init initializes the Store
	Init(ctx context.Context) error

	// GetAllKeysForReplication should return the keys of all data to be
	// replicated across the network. Typically all data should be
	// replicated every tReplicate seconds.
	GetAllKeysForReplication(ctx context.Context) ([][]byte, error)

	// ExpireKeys should expire all key/values due for expiration.
	ExpireKeys(ctx context.Context) error
}
