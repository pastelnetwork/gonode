//go:generate mockery --name=Client

package p2p

import "context"

// Client exposes the interfaces for p2p service
type Client interface {
	// Retrieve data from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the data
	Retrieve(ctx context.Context, key string) ([]byte, error)

	// Store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	Store(ctx context.Context, data []byte) (string, error)

	// Stats return status of p2p
	Stats(ctx context.Context) (map[string]interface{}, error)
}
