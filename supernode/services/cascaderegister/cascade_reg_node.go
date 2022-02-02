package cascaderegister

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

type RegisterCascadeNodeMaker struct {
	node.NodeMaker
}

func (maker RegisterCascadeNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &CascadeRegistrationNode{RegisterCascadeInterface: conn.RegisterCascade()}
}

// Node represent supernode connection.
type CascadeRegistrationNode struct {
	node.RegisterCascadeInterface
}
