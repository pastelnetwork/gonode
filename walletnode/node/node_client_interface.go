//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=RegisterArtwork
//go:generate mockery --name=DownloadArtwork
//go:generate mockery --name=ProcessUserdata
//go:generate mockery --name=RegisterSense

package node

import (
	"context"
	"github.com/pastelnetwork/gonode/common/storage/files"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/common/types"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// Client represents a base connection interface.
type ClientInterface interface {
	// Connect connects to the server at the given address.
	Connect(ctx context.Context, address string, secInfo *alts.SecInfo) (ConnectionInterface, error)
}

// ConnectionInterface represents a client connection
type ConnectionInterface interface {
	// Close closes connection.
	Close() error
	// Done returns a channel that's closed when connection is shutdown.
	Done() <-chan struct{}
	// RegisterArtwork returns a new RegisterArtwork stream.
	RegisterArtwork() RegisterNftInterface
	// DownloadArtwork returns a new DownloadArtwork stream.
	DownloadArtwork() DownloadNftInterface
	// ProcessUserdata returns a new ProcessUserdata stream.
	ProcessUserdata() ProcessUserdataInterface
	// RegisterSense returns new RegisterSense stream
	RegisterSense() RegisterSenseInterface
}

// RegisterSense contains methods for sense register
type RegisterSenseInterface interface {
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
	// SendRegMetadata send metadata of registration to SNs for next steps
	SendRegMetadata(ctx context.Context, regMetadata *types.ActionRegMetadata) error
	// ProbeImage uploads image to supernode.
	ProbeImage(ctx context.Context, image *files.File) (compressedDDAndFingerprints []byte, isValidBurnTxID bool, err error)
	// SendSignedTicket send a reg-art ticket signed by cNode to SuperNode
	SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, ddFpFile []byte) (string, error)
	// SendActionAct send action act to SNs for next steps
	SendActionAct(ctx context.Context, actionRegTxid string) error
}

// RegisterArtwork contains methods for registering artwork.
type RegisterNftInterface interface {
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
	// SendRegMetadata send metadata of registration to SNs for next steps
	SendRegMetadata(ctx context.Context, regMetadata *types.NftRegMetadata) error
	// ProbeImage uploads image to supernode.
	ProbeImage(ctx context.Context, image *files.File) (compressedDDAndFingerprints []byte, err error)
	// UploadImageImageWithThumbnail uploads the image with pqsignature and its thumbnail to supernodes
	UploadImageWithThumbnail(ctx context.Context, image *files.File, thumbnail files.ThumbnailCoordinate) (previewHash []byte, mediumThumbnailHash []byte, smallThumbnailHash []byte, err error)
	// SendSignedTicket send a reg-art ticket signed by cNode to SuperNode
	SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, key1 string, key2 string, rqdisFile []byte, ddFpFile []byte, encoderParams rqnode.EncoderParameters) (int64, error)
	// SendPreBurnedFreeTxId send TxId of the transaction in which 10% of registration fee is preburned
	SendPreBurntFeeTxid(ctx context.Context, txid string) (string, error)
}

// DownloadArtwork contains methods for downloading artwork.
type DownloadNftInterface interface {
	// Download sends image downloading request to supernode.
	Download(ctx context.Context, txid, timestamp, signature, ttxid string) (file []byte, err error)
	DownloadThumbnail(ctx context.Context, key []byte) (file []byte, err error)
}

// ProcessUserdata contains methods for processing userdata.
type ProcessUserdataInterface interface {
	// SessID returns the sessID received from the server during the handshake.
	SessID() (sessID string)
	// Session sets up an initial connection with supernode, with given supernode mode primary/secondary.
	Session(ctx context.Context, IsPrimary bool) (err error)
	// AcceptedNodes requests information about connected secondary nodes.
	AcceptedNodes(ctx context.Context) (pastelIDs []string, err error)
	// ConnectTo commands to connect to the primary node, where nodeKey is primary key.
	ConnectTo(ctx context.Context, primaryNode types.MeshedSuperNode) error
	// SendUserdata send user specified data (with other generated info like signature, previous block hash, timestamp,...) to supernode.
	SendUserdata(ctx context.Context, request *userdata.ProcessRequestSigned) (result *userdata.ProcessResult, err error)
	// ReceiveUserdata get user specified data from supernode
	ReceiveUserdata(ctx context.Context, userpastelid string) (result *userdata.ProcessRequest, err error)
}
