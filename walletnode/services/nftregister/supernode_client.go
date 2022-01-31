package nftregister

import "github.com/pastelnetwork/gonode/walletnode/node"

// RegisterNftNodeMaker makes class ProcessUserdata for RegisterNft
type RegisterNftNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class ProcessUserdata for RegisterNft
func (maker RegisterNftNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftRegistrationNode{RegisterNftInterface: conn.RegisterNft()}
}

// NftRegistrationNode represent supernode connection.
type NftRegistrationNode struct {
	node.RegisterNftInterface
}
