//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=ProcessUserdata

package node

import (
	"context"
)

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
	// TODO: Define function to send data to multiple Super Node
}
