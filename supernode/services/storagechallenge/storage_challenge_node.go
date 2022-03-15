package storagechallenge

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

// SCNodeMaker makes concrete instance of StorageChallengeNode
type SCNodeMaker struct {
	node.NodeMaker
}

// MakeNode makes concrete instance of StorageChallengeNode
func (maker SCNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &SCNode{StorageChallengeInterface: conn.StorageChallenge()}
}

// SCNode (for Storage Challenge Node) represent supernode connection.
type SCNode struct {
	node.StorageChallengeInterface
}
