package node

import (
	"context"
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
	// RegisterArtowrk returns a new RegisterArtowrk stream.
	RegisterArtowrk(ctx context.Context) (RegisterArtowrk, error)
}

// RegisterArtowrk represents an interaction stream with supernodes for registering artwork.
type RegisterArtowrk interface {
	// Handshake sets up an initial connection with supernode, by telling connection id and supernode mode, primary or secondary.
	Handshake(ctx context.Context, connID string, IsPrimary bool) error
	// AcceptedNodes requests information about connected secondary nodes.
	AcceptedNodes(ctx context.Context) (pastelIDs []string, err error)
	// ConnectTo commands to connect to the primary node, where nodeKey is primary key.
	ConnectTo(ctx context.Context, nodeKey string) error
}
