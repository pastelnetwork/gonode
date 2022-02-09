package cascaderegister

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

// RegisterCascadeNodeMaker makes concrete instance of CascadeRegistrationNode
type RegisterCascadeNodeMaker struct {
	node.NodeMaker
}

// MakeNode makes concrete instance of CascadeRegistrationNode
func (maker RegisterCascadeNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &CascadeRegistrationNode{RegisterCascadeInterface: conn.RegisterCascade()}
}

// CascadeRegistrationNode represent supernode connection.
type CascadeRegistrationNode struct {
	node.RegisterCascadeInterface
}
