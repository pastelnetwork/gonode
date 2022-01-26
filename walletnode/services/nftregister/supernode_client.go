package nftregister

import "github.com/pastelnetwork/gonode/walletnode/node"

type RegisterNftNodeMaker struct {
	node.NodeMaker
}

func (maker RegisterNftNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftRegisterNode{RegisterNftInterface: conn.RegisterArtwork()}
}

// Node represent supernode connection.
type NftRegisterNode struct {
	node.RegisterNftInterface
}
