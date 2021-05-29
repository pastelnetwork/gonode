package p2p

import (
	"context"

	"github.com/pastelnetwork/gonode/common/service/artwork"
)

// Client represents an interaction with a distributed hash table.
type Client interface {
	// Get retrieves data from the networking using key. Key is the base58 encoded
	// identifier of the data.
	Get(ctx context.Context, key string, file *artwork.File) (found bool, err error)

	// Store stores data on the network. This will trigger an iterateStore message.
	// The base58 encoded identifier will be returned if the store is successful.
	Store(ctx context.Context, file *artwork.File) (id string, err error)
}
