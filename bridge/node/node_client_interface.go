//go:generate mockery --name=ClientInterface
//go:generate mockery --name=ConnectionInterface
//go:generate mockery --name=RegisterSenseInterface
//go:generate mockery --name=RegisterCascadeInterface
//go:generate mockery --name=RegisterNftInterface
//go:generate mockery --name=DownloadNftInterface
//go:generate mockery --name=ProcessUserdataInterface

package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/types"
)

// SNClientInterface represents a base connection interface.
type SNClientInterface interface {
	// Connect connects to the server at the given address.
	Connect(ctx context.Context, address string, secInfo *alts.SecInfo) (ConnectionInterface, error)
}

// ClientInterface represents a base connection interface.
type ClientInterface interface {
	// Connect connects to the server at the given address.
	Connect(ctx context.Context, address string) (ConnectionInterface, error)
}

// ConnectionInterface represents a client connection
type ConnectionInterface interface {
	// Close closes connection.
	Close() error
	// Done returns a channel that's closed when connection is shutdown.
	Done() <-chan struct{}

	// DownloadData returns a new DownloadData stream.
	DownloadData() DownloadDataInterface
}

// SuperNodeAPIInterface base API interface
type SuperNodeAPIInterface interface {
	// SessID returns the sessID received from the server during the handshake.
	SessID() (sessID string)
	// Session sets up an initial connection with supernode, with given supernode mode primary/secondary.
	Session(ctx context.Context, IsPrimary bool) (err error)
	// AcceptedNodes requests information about connected secondary nodes.
	AcceptedNodes(ctx context.Context) (pastelIDs []string, err error)
	// ConnectTo commands to connect to the primary node
	ConnectTo(ctx context.Context, primaryNode types.MeshedSuperNode) error
	// MeshNodes send to supernode all info of nodes are meshed together (include the received supernode)
	MeshNodes(ctx context.Context, meshedNodes []types.MeshedSuperNode) error
}

// RealNodeMaker interface to make concrete node types
type RealNodeMaker interface {
	MakeNode(conn ConnectionInterface) SuperNodeAPIInterface
}

// DownloadDataInterface contains methods for downloading data.
type DownloadDataInterface interface {
	SuperNodeAPIInterface

	DownloadThumbnail(ctx context.Context, txid string, numNails int) (files map[int][]byte, err error)
	DownloadDDAndFingerprints(ctx context.Context, txid string) (file []byte, err error)
}
