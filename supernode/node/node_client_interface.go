//go:generate mockery --name=ClientInterface
//go:generate mockery --name=ConnectionInterface
//go:generate mockery --name=RegisterNftInterface
//go:generate mockery --name=ProcessUserdataInterface
//go:generate mockery --name=RegisterSenseInterface
//go:generate mockery --name=RegisterCascadeInterface
//go:generate mockery --name=RegisterCollectionInterface
//go:generate mockery --name=StorageChallengeInterface
//go:generate mockery --name=SelfHealingChallengeInterface

package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

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
	// RegisterNft returns a new RegisterNft stream.
	RegisterNft() RegisterNftInterface
	// RegisterSense returns a new RegisterSense stream
	RegisterSense() RegisterSenseInterface
	// RegisterCascade returns a new RegisterCascade stream
	RegisterCascade() RegisterCascadeInterface
	//StorageChallenge returns a new StorageChallenge stream
	StorageChallenge() StorageChallengeInterface
	//SelfHealingChallenge returns a new SelfHealingChallenge stream
	SelfHealingChallenge() SelfHealingChallengeInterface
	// RegisterCollection returns a new RegisterCollection stream
	RegisterCollection() RegisterCollectionInterface
}

// SuperNodePeerAPIInterface base interface for other Node API interfaces
type SuperNodePeerAPIInterface interface {
	// SessID returns the taskID received from the server during the handshake.
	SessID() (taskID string)
	// Session sets up an initial connection with primary supernode, by telling sessID and its own nodeID.
	Session(ctx context.Context, nodeID, sessID string) (err error)
}

// revive:disable:exported

// NodeMaker interface to make concrete node types
type NodeMaker interface {
	MakeNode(conn ConnectionInterface) SuperNodePeerAPIInterface
}

// revive:enable:exported

// RegisterNftInterface represents an interaction stream with supernodes for registering Nft.
type RegisterNftInterface interface {
	SuperNodePeerAPIInterface

	// SendSignedDDAndFingerprints send compressedDDAndFingerprints from fromNodeID to target SN
	SendSignedDDAndFingerprints(ctx context.Context, sessionID string, fromNodeID string, compressedDDAndFingerprints []byte) error
	// Send signature of ticket to primary supernode
	SendNftTicketSignature(ctx context.Context, nodeID string, signature []byte) error
}

// RegisterSenseInterface represents an interaction stream with supernodes for registering sense.
type RegisterSenseInterface interface {
	SuperNodePeerAPIInterface

	// SendSignedDDAndFingerprints send compressedDDAndFingerprints from fromNodeID to target SN
	SendSignedDDAndFingerprints(ctx context.Context, sessionID string, fromNodeID string, compressedDDAndFingerprints []byte) error
	// SendSenseTicketSignature sends signature of ticket to primary supernode
	SendSenseTicketSignature(ctx context.Context, nodeID string, signature []byte) error
}

// RegisterCollectionInterface represents an interaction stream with supernodes for registering collection.
type RegisterCollectionInterface interface {
	SuperNodePeerAPIInterface

	// SendCollectionTicketSignature sends signature of ticket to primary supernode
	SendCollectionTicketSignature(ctx context.Context, nodeID string, signature []byte) error
}

// RegisterCascadeInterface represents an interaction stream with supernodes for registering sense.
type RegisterCascadeInterface interface {
	SuperNodePeerAPIInterface

	// Send signature of ticket to primary supernode
	SendCascadeTicketSignature(ctx context.Context, nodeID string, signature []byte) error
}

// StorageChallengeInterface represents an interaction stream with supernodes for storage challenge communications
type StorageChallengeInterface interface {
	SuperNodePeerAPIInterface

	ProcessStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error

	VerifyStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) (*pb.StorageChallengeData, error)
}

// SelfHealingChallengeInterface represents an interaction stream with supernodes for self-healing challenge communications
type SelfHealingChallengeInterface interface {
	SuperNodePeerAPIInterface

	ProcessSelfHealingChallenge(ctx context.Context, challengeMessage *pb.SelfHealingData) error

	VerifySelfHealingChallenge(ctx context.Context, challengeMessage *pb.SelfHealingData) (*pb.SelfHealingData, error)
}

// ProcessUserdataInterface represents an interaction stream with supernodes for sending userdata.
type ProcessUserdataInterface interface {
	SuperNodePeerAPIInterface

	// Send userdata to primary supernode
	SendUserdataToPrimary(ctx context.Context, dataSigned userdata.SuperNodeRequest) (userdata.SuperNodeReply, error)
	// Send userdata to supernode with leader rqlite
	SendUserdataToLeader(ctx context.Context, finalUserdata userdata.ProcessRequestSigned) (userdata.SuperNodeReply, error)
}
