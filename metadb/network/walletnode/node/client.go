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

// ProcessUserdata contains methods for processing userdata.
type ProcessUserdata interface {
	// SessID returns the sessID received from the server during the handshake.
	SessID() (sessID string)
	// Session sets up an initial connection with supernode, with given supernode mode primary/secondary.
	Session(ctx context.Context, IsPrimary bool) (err error)
	// AcceptedNodes requests information about connected secondary nodes.
	AcceptedNodes(ctx context.Context) (pastelIDs []string, err error)
	// ConnectTo commands to connect to the primary node, where nodeKey is primary key.
	ConnectTo(ctx context.Context, nodeKey, sessID string) error
	// SendUserdata send user specified data (with other generated info like signature, previous block hash, timestamp,...) to supernode.
	SendUserdata(ctx context.Context, request *userdata.ProcessRequestSigned) (result *userdata.ProcessResult, err error)
	// ReceiveUserdata get user specified data from supernode
	ReceiveUserdata(ctx context.Context, userpastelid string) (result *userdata.ProcessRequest, err error)
}
