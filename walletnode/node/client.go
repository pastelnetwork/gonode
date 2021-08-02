//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=RegisterArtwork
//go:generate mockery --name=DownloadArtwork

package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/service/artwork"
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
	// DownloadArtwork returns a new DownloadArtwork stream.
	DownloadArtwork() DownloadArtwork
}

// RegisterArtwork contains methods for registering artwork.
type RegisterArtwork interface {
	// SessID returns the sessID received from the server during the handshake.
	SessID() (sessID string)
	// Session sets up an initial connection with supernode, with given supernode mode primary/secondary.
	Session(ctx context.Context, IsPrimary bool) (err error)
	// AcceptedNodes requests information about connected secondary nodes.
	AcceptedNodes(ctx context.Context) (pastelIDs []string, err error)
	// ConnectTo commands to connect to the primary node, where nodeKey is primary key.
	ConnectTo(ctx context.Context, nodeKey, sessID string) error
	// ProbeImage uploads image to supernode.
	ProbeImage(ctx context.Context, image *artwork.File) (fingerprintData []byte, err error)
}

// DownloadArtwork contains methods for downloading artwork.
type DownloadArtwork interface {
	// Download sends image downloading request to supernode.
	Download(ctx context.Context, txid, timestamp, signature, ttxid string) (file []byte, err error)
}
