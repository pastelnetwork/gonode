//go:generate mockery --name=Client

package p2p

import (
	"context"

	"github.com/pastelnetwork/gonode/p2p/kademlia"
)

// Client exposes the interfaces for p2p service
type Client interface {
	// Retrieve retrieve data from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the data
	Retrieve(ctx context.Context, key string, localOnly ...bool) ([]byte, error)
	// Store store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	Store(ctx context.Context, data []byte) (string, error)

	// Delete a key, value
	Delete(ctx context.Context, key string) error

	// Stats return status of p2p
	Stats(ctx context.Context) (map[string]interface{}, error)

	// NClosestNodes return n closest masternode to a given string
	NClosestNodes(ctx context.Context, n int, key string) []*kademlia.Node
}
