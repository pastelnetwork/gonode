//go:generate mockery --name=Client

package p2p

import (
	"context"
	"time"
)

// Client exposes the interfaces for p2p service
type Client interface {
	// Retrieve retrieve data from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the data
	// - localOnly will
	Retrieve(ctx context.Context, key string, localOnly ...bool) ([]byte, error)

	// BatchRetrieve retrieve data from the kademlia network by keys
	// reqCount is the minimum number of keys that are actually required by the caller
	// to successfully perform the reuquired operation
	BatchRetrieve(ctx context.Context, keys []string, reqCount int, txID string, localOnly ...bool) (map[string][]byte, error)
	// Store store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	Store(ctx context.Context, data []byte, typ int) (string, error)

	// StoreBatch will store a batch of values with their SHA256 hash as the key
	StoreBatch(ctx context.Context, values [][]byte, typ int, taskID string) error

	// Delete a key, value
	Delete(ctx context.Context, key string) error

	// Stats return status of p2p
	Stats(ctx context.Context) (map[string]interface{}, error)

	// NClosestNodes return n closest supernodes to a given key string (NB full node string formatting)
	NClosestNodes(ctx context.Context, n int, key string, ignores ...string) []string

	//NClosestNodesWithIncludingNodeList return n closest nodes to a given key with including node list
	NClosestNodesWithIncludingNodeList(ctx context.Context, n int, key string, ignores, nodesToInclude []string) []string

	// LocalStore store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	LocalStore(ctx context.Context, key string, data []byte) (string, error)

	// DisableKey adds key to disabled keys list - It takes in a B58 encoded SHA-256 hash
	DisableKey(ctx context.Context, b58EncodedHash string) error

	// EnableKey removes key from disabled list - It takes in a B58 encoded SHA-256 hash
	EnableKey(ctx context.Context, b58EncodedHash string) error

	// GetLocalKeys returns a list of all keys stored locally
	GetLocalKeys(ctx context.Context, from *time.Time, to time.Time) ([]string, error)
}
