package nftregister

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

// RegisterNftNodeMaker makes concrete instance of NftRegistrationNode
type RegisterNftNodeMaker struct {
	node.NodeMaker
}

// MakeNode makes concrete instance of NftRegistrationNode
func (maker RegisterNftNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &NftRegistrationNode{RegisterNftInterface: conn.RegisterNft()}
}

// NftRegistrationNode represent supernode connection.
type NftRegistrationNode struct {
	node.RegisterNftInterface
}
