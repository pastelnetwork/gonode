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
	RegisterArtowrk(taskID string) RegisterArtowrk
}

// RegisterArtowrk represents an interaction stream with supernodes for registering artwork.
type RegisterArtowrk interface {
	// TaskID returns the taskID received from the server during the handshake.
	TaskID() (taskID string)
	// Handshake secondary sets up an initial connection with primary supernode, by telling connection id and its own node key.
	Handshake(ctx context.Context, nodeKey string) (err error)
}
