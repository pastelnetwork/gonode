package kademlia

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

// Store is the interface for implementing the storage mechanism for the DHT
type Store interface {
	// Store a key/value pair for the local node with the replication
	Store(ctx context.Context, key []byte, data []byte, typ int, isOriginal bool) error

	// Retrieve the local key/value from store
	Retrieve(ctx context.Context, key []byte) ([]byte, error)

	// RetrieveWithType gets data with type
	RetrieveWithType(_ context.Context, key []byte) ([]byte, int, error)

	// Delete a key/value pair from the store
	Delete(ctx context.Context, key []byte)

	// KeysForReplication returns the keys of all data to be replicated across the network
	GetKeysForReplication(ctx context.Context, from time.Time, to time.Time) [][]byte

	// Stats returns stats of store
	Stats(ctx context.Context) (map[string]interface{}, error)

	// Close the store
	Close(ctx context.Context)

	// Count the records in store
	Count(ctx context.Context /*, type RecordType*/) (int, error)

	// DeleteAll the records in store
	DeleteAll(ctx context.Context /*, type RecordType*/) error

	// UpdateKeyReplication updates the replication status of the key
	UpdateKeyReplication(ctx context.Context, key []byte) error

	// StoreBatch stores a batch of key/value pairs for the local node with the replication
	StoreBatch(ctx context.Context, values [][]byte, typ int, isOriginal bool) error
	// GetAllReplicationInfo returns all records in replication table
	GetAllReplicationInfo(ctx context.Context) ([]domain.NodeReplicationInfo, error)

	// UpdateReplicationInfo updates replication info
	UpdateReplicationInfo(ctx context.Context, rep domain.NodeReplicationInfo) error

	// AddReplicationInfo adds replication info
	AddReplicationInfo(ctx context.Context, rep domain.NodeReplicationInfo) error

	//  GetOwnCreatedAt ...
	GetOwnCreatedAt(ctx context.Context) (time.Time, error)

	// StoreBatchRepKeys ...
	StoreBatchRepKeys(values [][]byte, id string, ip string, port int) error

	// GetAllToDoRepKeys gets all keys that need to be replicated
	GetAllToDoRepKeys() (retKeys domain.ToRepKeys, err error)

	// DeleteRepKey deletes a key from the replication table
	DeleteRepKey(key []byte) error
}
