//go:generate mockery --name=Client

package p2p

import "context"

// Client exposes the interfaces for p2p service
type Client interface {
	// Retrieve retrieve data from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the data
	// - localOnly will
	Retrieve(ctx context.Context, key string, localOnly ...bool) ([]byte, error)
	// Store store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	Store(ctx context.Context, data []byte, typ int) (string, error)

	// StoreBatch will store a batch of values with their SHA256 hash as the key
	StoreBatch(ctx context.Context, values [][]byte, typ int) error

	// Delete a key, value
	Delete(ctx context.Context, key string) error

	// Stats return status of p2p
	Stats(ctx context.Context) (map[string]interface{}, error)

	// NClosestNodes return n closest supernodes to a given key string (NB full node string formatting)
	NClosestNodes(ctx context.Context, n int, key string, ignores ...string) []string

	// LocalStore store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	LocalStore(ctx context.Context, key string, data []byte) (string, error)

	// DisableKey adds key to disabled keys list - It takes in a B58 encoded SHA-256 hash
	DisableKey(ctx context.Context, b58EncodedHash string) error

	// EnableKey removes key from disabled list - It takes in a B58 encoded SHA-256 hash
	EnableKey(ctx context.Context, b58EncodedHash string) error
}
