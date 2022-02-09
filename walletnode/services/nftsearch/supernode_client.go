package nftsearch

import "github.com/pastelnetwork/gonode/walletnode/node"

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
