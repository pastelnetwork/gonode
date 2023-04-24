package collectionregister

import (
	"github.com/pastelnetwork/gonode/supernode/node"
)

// RegisterCollectionNodeMaker makes concrete instance of CollectionRegistrationNode
type RegisterCollectionNodeMaker struct {
	node.NodeMaker
}

// MakeNode makes concrete instance of CollectionRegistrationNode
func (maker RegisterCollectionNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodePeerAPIInterface {
	return &CollectionRegistrationNode{RegisterCollectionInterface: conn.RegisterCollection()}
}

// CollectionRegistrationNode represent supernode connection.
type CollectionRegistrationNode struct {
	node.RegisterCollectionInterface
}
