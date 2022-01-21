package senseregister

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type SenseRegisterNodeMaker struct {
	node.NodeMaker
}

func (maker SenseRegisterNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &SenseRegisterNode{RegisterSenseInterface: conn.RegisterSense()}
}

// Node represent supernode connection.
type SenseRegisterNode struct {
	node.RegisterSenseInterface
}
