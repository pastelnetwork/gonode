package kademlia

import (
	"context"
	"time"
)

// Store is the interface for implementing the storage mechanism for the DHT
type Store interface {
	// Store a key/value pair for the local node with the replication
	Store(ctx context.Context, key []byte, data []byte, replication time.Time) error

	// Retrieve the local key/value from store
	Retrieve(ctx context.Context, key []byte) ([]byte, error)

	// Delete a key/value pair from the store
	Delete(ctx context.Context, key []byte)

	// Keys returns all the keys from the store
	Keys(ctx context.Context) [][]byte

	// KeysForReplication returns the keys of all data to be replicated across the network
	KeysForReplication(ctx context.Context) [][]byte

	// Stats returns stats of store
	Stats(ctx context.Context) (map[string]interface{}, error)

	// Close the store
	Close(ctx context.Context)
}
