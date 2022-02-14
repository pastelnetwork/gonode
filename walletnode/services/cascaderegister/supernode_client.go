package cascaderegister

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// RegisterCascadeNodeMaker makes class RegisterCascade for SuperNodeAPIInterface
type RegisterCascadeNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class RegisterCascade for SuperNodeAPIInterface
func (maker RegisterCascadeNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &CascadeRegistrationNode{RegisterCascadeInterface: conn.RegisterCascade()}
}

// CascadeRegistrationNode represent supernode connection.
type CascadeRegistrationNode struct {
	node.RegisterCascadeInterface
}
