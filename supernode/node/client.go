//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=RegisterArtwork
//go:generate mockery --name=ProcessUserdata
//go:generate mockery --name=ExternalDupeDetection

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
	// RegisterArtwork returns a new RegisterArtwork stream.
	RegisterArtwork() RegisterArtwork
	// ProcessUserdata returns a new ProcessUserdata stream.
	ProcessUserdata() ProcessUserdata
	// ExternalDupeDetection returns a new ExternalDupeDetection stream.
	ExternalDupeDetection() ExternalDupeDetection
}

// RegisterArtwork represents an interaction stream with supernodes for registering artwork.
type RegisterArtwork interface {
	// SessID returns the taskID received from the server during the handshake.
	SessID() (taskID string)
	// Session sets up an initial connection with primary supernode, by telling sessID and its own nodeID.
	Session(ctx context.Context, nodeID, sessID string) (err error)
	// Send signature of ticket to primary supernode
	SendArtTicketSignature(ctx context.Context, nodeID string, signature []byte) error
}

// ProcessUserdata represents an interaction stream with supernodes for sending userdata.
type ProcessUserdata interface {
	// SessID returns the taskID received from the server during the handshake.
	SessID() (taskID string)
	// Session sets up an initial connection with primary supernode, by telling sessID and its own nodeID.
	Session(ctx context.Context, nodeID, sessID string) (err error)
	// Send userdata to primary supernode
	SendUserdataToPrimary(ctx context.Context, dataSigned userdata.SuperNodeRequest) (userdata.SuperNodeReply, error)
	// Send userdata to supernode with leader rqlite
	SendUserdataToLeader(ctx context.Context, finalUserdata userdata.ProcessRequestSigned) (userdata.SuperNodeReply, error)
}

// ExternalDupeDetection  represents an interaction stream with supernodes for sending external dupe detection.
type ExternalDupeDetection interface {
	// SessID returns the taskID received from the server during the handshake.
	SessID() (taskID string)
	// Session sets up an initial connection with primary supernode, by telling sessID and its own nodeID.
	Session(ctx context.Context, nodeID, sessID string) (err error)
	// SendEDDTicketSignature Send signature of ticket to primary supernode
	SendEDDTicketSignature(ctx context.Context, nodeID string, signature []byte) error
}
