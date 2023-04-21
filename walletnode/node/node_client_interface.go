//go:generate mockery --name=ClientInterface
//go:generate mockery --name=ConnectionInterface
//go:generate mockery --name=RegisterSenseInterface
//go:generate mockery --name=RegisterCascadeInterface
//go:generate mockery --name=RegisterCollectionInterface
//go:generate mockery --name=RegisterNftInterface
//go:generate mockery --name=DownloadNftInterface
//go:generate mockery --name=ProcessUserdataInterface

package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/storage/files"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/common/types"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// ClientInterface represents a base connection interface.
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
	// RegisterNft returns a new RegisterNft stream.
	RegisterNft() RegisterNftInterface
	// DownloadNft returns a new DownloadNft stream.
	DownloadNft() DownloadNftInterface
	// RegisterSense returns new RegisterSense stream
	RegisterSense() RegisterSenseInterface
	// RegisterCascade returns new RegisterCascade stream
	RegisterCascade() RegisterCascadeInterface
	// RegisterCollection returns new RegisterCollection stream
	RegisterCollection() RegisterCollectionInterface
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
	// GetDupeDetectionDBHash requests hash of dd database
	GetDupeDetectionDBHash(ctx context.Context) (hash string, err error)
}

// RealNodeMaker interface to make concrete node types
type RealNodeMaker interface {
	MakeNode(conn ConnectionInterface) SuperNodeAPIInterface
}

// RegisterSenseInterface contains methods for sense register
type RegisterSenseInterface interface {
	SuperNodeAPIInterface

	// SendRegMetadata send metadata of registration to SNs for next steps
	SendRegMetadata(ctx context.Context, regMetadata *types.ActionRegMetadata) error
	// ProbeImage uploads image to supernode.
	ProbeImage(ctx context.Context, image *files.File) ([]byte, bool, string, bool, error)
	// SendSignedTicket send a reg-nft ticket signed by cNode to SuperNode
	SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, ddFpFile []byte) (string, error)
	// SendActionAct send action act to SNs for next steps
	SendActionAct(ctx context.Context, actionRegTxid string) error
}

// RegisterCascadeInterface contains methods for sense register
type RegisterCascadeInterface interface {
	SuperNodeAPIInterface

	// SendRegMetadata send metadata of registration to SNs for next steps
	SendRegMetadata(ctx context.Context, regMetadata *types.ActionRegMetadata) error
	// UploadAsset uploads asset to supernode for storing
	UploadAsset(ctx context.Context, asset *files.File) error
	// SendSignedTicket send a reg-nft ticket signed by cNode to SuperNode
	SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, rqIdsFile []byte, encoderParams rqnode.EncoderParameters) (string, error)
	// SendActionAct send action act to SNs for next steps
	SendActionAct(ctx context.Context, actionRegTxid string) error
}

// RegisterCollectionInterface contains methods for collection register
type RegisterCollectionInterface interface {
	SuperNodeAPIInterface

	// SendTicketForSignature send a collection ticket to be signed by other SNs
	SendTicketForSignature(ctx context.Context, ticket string, signature []byte, label string, fee int64) ([]byte, error)
}

// RegisterNftInterface contains methods for registering nft.
type RegisterNftInterface interface {
	SuperNodeAPIInterface

	// SendRegMetadata send metadata of registration to SNs for next steps
	SendRegMetadata(ctx context.Context, regMetadata *types.NftRegMetadata) error
	// ProbeImage uploads image to supernode.
	ProbeImage(ctx context.Context, image *files.File) ([]byte, bool, bool, error)
	// UploadImageWithThumbnail uploads the image with pqsignature and its thumbnail to supernodes
	UploadImageWithThumbnail(ctx context.Context, image *files.File, thumbnail files.ThumbnailCoordinate) ([]byte, []byte, []byte, error)
	// SendSignedTicket send a reg-nft ticket signed by cNode to SuperNode
	SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, label string, rqIdsFile []byte, ddFpFile []byte, encoderParams rqnode.EncoderParameters) (int64, error)
	// SendPreBurntFeeTxid send TxId of the transaction in which 10% of registration fee is preburned
	SendPreBurntFeeTxid(ctx context.Context, txid string) (string, error)
}

// DownloadNftInterface contains methods for downloading nft.
type DownloadNftInterface interface {
	SuperNodeAPIInterface

	// Download sends image downloading request to supernode.
	Download(ctx context.Context, txid, timestamp, signature, ttxid, ttype string) ([]byte, error)
	DownloadThumbnail(ctx context.Context, txid string, numNails int) (files map[int][]byte, err error)
	DownloadDDAndFingerprints(ctx context.Context, txid string) (file []byte, err error)
}

// ProcessUserdataInterface contains methods for processing userdata.
type ProcessUserdataInterface interface {
	SuperNodeAPIInterface

	// SendUserdata send user specified data (with other generated info like signature, previous block hash, timestamp,...) to supernode.
	SendUserdata(ctx context.Context, request *userdata.ProcessRequestSigned) (*userdata.ProcessResult, error)
	// RetrieveUserdata get user specified data from supernode
	RetrieveUserdata(ctx context.Context, userpastelid string) (*userdata.ProcessRequest, error)
}
