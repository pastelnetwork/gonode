//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=RaptorQ

package node

import (
	"context"
	//"github.com/pastelnetwork/gonode/common/service/artwork"
)

// Client represents a base connection interface.
type Client interface {
	// Connect connects to the server at the given address.
	Connect(ctx context.Context, address string) (Connection, error)
}

// Connection represents a client connection
type Connection interface {
	// Close closes connection.
	Close() error
	// Done returns a channel that's closed when connection is shutdown.
	Done() <-chan struct{}
	// RaptorQ returns a new RaptorQ stream.
	RaptorQ() RaptorQ
}

// RegisterArtwork contains methods for registering artwork.
type RaptorQ interface {
}
