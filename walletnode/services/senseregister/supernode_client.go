package senseregister

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// RegisterSenseNodeMaker makes class RegisterSense for SuperNodeAPIInterface
type RegisterSenseNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class RegisterSense for SuperNodeAPIInterface
func (maker RegisterSenseNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &SenseRegistrationNode{RegisterSenseInterface: conn.RegisterSense()}
}

// SenseRegistrationNode represent supernode connection.
type SenseRegistrationNode struct {
	node.RegisterSenseInterface
}
