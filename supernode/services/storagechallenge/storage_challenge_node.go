package storagechallenge

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

// RegisterSenseNodeMaker makes concrete instance of SenseRegistrationNode
type StorageChallengeNodeMaker struct {
	node.NodeMaker
}

// MakeNode makes concrete instance of SenseRegistrationNode
func (maker StorageChallengeNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &StorageChallengeNode{StorageChallengeInterface: conn.StorageChallenge()}
}

// SenseRegistrationNode represent supernode connection.
type StorageChallengeNode struct {
	node.StorageChallengeInterface
}
