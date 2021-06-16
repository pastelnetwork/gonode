package kademlia

import (
	"context"
	"time"
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

	// KeysForExpiration returns the keys of all data to be removed from local storage
	KeysForExpiration(ctx context.Context) [][]byte

	// Close the store
	Close(ctx context.Context)
}
