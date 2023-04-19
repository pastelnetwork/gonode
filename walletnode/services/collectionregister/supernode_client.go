package collectionregister

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// RegisterCollectionNodeMaker makes class RegisterCollection for SuperNodeAPIInterface
type RegisterCollectionNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class RegisterCascade for SuperNodeAPIInterface
func (maker RegisterCollectionNodeMaker) MakeNode(_ node.ConnectionInterface) node.SuperNodeAPIInterface {
	return nil
}

// CollectionRegistrationNode represent supernode connection.
type CollectionRegistrationNode struct {
	node.RegisterCascadeInterface
}
