package senseregister

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

type RegisterSenseNodeMaker struct {
	node.NodeMaker
}

func (maker RegisterSenseNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &SenseRegistrationNode{RegisterSenseInterface: conn.RegisterSense()}
}

// Node represent supernode connection.
type SenseRegistrationNode struct {
	node.RegisterSenseInterface
}
