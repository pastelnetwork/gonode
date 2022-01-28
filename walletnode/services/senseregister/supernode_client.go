package senseregister

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type RegisterSenseNodeMaker struct {
	node.NodeMaker
}

func (maker RegisterSenseNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &SenseRegistrationNode{RegisterSenseInterface: conn.RegisterSense()}
}

// Node represent supernode connection.
type SenseRegistrationNode struct {
	node.RegisterSenseInterface
}
