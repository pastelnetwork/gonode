//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=RegisterArtwork
//go:generate mockery --name=ProcessUserdata

package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
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
	// StorageChallenge returns a new StorageChallenge stream.
	StorageChallenge() StorageChallenge
}

// RegisterArtwork represents an interaction stream with supernodes for registering artwork.
type RegisterArtwork interface {
	// SessID returns the taskID received from the server during the handshake.
	SessID() (taskID string)
	// Session sets up an initial connection with primary supernode, by telling sessID and its own nodeID.
	Session(ctx context.Context, nodeID, sessID string) (err error)
	// SendSignedDDAndFingerprints send compressedDDAndFingerprints from fromNodeID to target SN
	SendSignedDDAndFingerprints(ctx context.Context, sessionID string, fromNodeID string, compressedDDAndFingerprints []byte) error
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

// StorageChallenge represents an interaction stream with supernodes for sending userdata.
type StorageChallenge pb.StorageChallengeClient
