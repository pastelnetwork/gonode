package collectionregister

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// RegisterCollectionNodeMaker makes class RegisterCollection for SuperNodeAPIInterface
type RegisterCollectionNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class RegisterCascade for SuperNodeAPIInterface
func (maker RegisterCollectionNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &CollectionRegistrationNode{RegisterCollectionInterface: conn.RegisterCollection()}
}

// CollectionRegistrationNode represent supernode connection.
type CollectionRegistrationNode struct {
	node.RegisterCollectionInterface
}
