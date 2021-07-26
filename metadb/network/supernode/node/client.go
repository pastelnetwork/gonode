//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=ProcessUserdata

package node

import (
	"context"
	"github.com/pastelnetwork/gonode/common/service/userdata"
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
	// ProcessUserdata returns a new ProcessUserdata stream.
	ProcessUserdata() ProcessUserdata
}

// ProcessUserdata represents an interaction stream with supernodes for registering artwork.
type ProcessUserdata interface {
	// SessID returns the taskID received from the server during the handshake.
	SessID() (taskID string)
	// Session sets up an initial connection with primary supernode, by telling sessID and its own nodeID.
	Session(ctx context.Context, nodeID, sessID string) (err error)
	// Send userdata to primary supernode
	SendUserdataToPrimary(ctx context.Context, dataSigned userdata.SuperNodeRequest) (userdata.SuperNodeReply, error)
}
