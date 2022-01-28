package nftregister

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

type RegisterNftNodeMaker struct {
	node.NodeMaker
}

func (maker RegisterNftNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &NftRegistrationNode{RegisterNftInterface: conn.RegisterNft()}
}

// Node represent supernode connection.
type NftRegistrationNode struct {
	node.RegisterNftInterface
}
