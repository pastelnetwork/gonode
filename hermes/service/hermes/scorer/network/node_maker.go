package network

import "github.com/pastelnetwork/gonode/hermes/service/node"

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

// NftSearchingNodeMaker makes class DownloadNft for SuperNodeAPIInterface
type NftSearchingNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class DownloadNft for SuperNodeAPIInterface
func (maker NftSearchingNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftSearchingNode{DownloadNftInterface: conn.DownloadNft()}
}

// NftSearchingNode represent supernode connection.
type NftSearchingNode struct {
	node.DownloadNftInterface
}
