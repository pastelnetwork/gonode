package nftdownload

import "github.com/pastelnetwork/gonode/walletnode/node"

// NftDownloadingNodeMaker makes class DownloadNft for SuperNodeAPIInterface
type NftDownloadingNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class DownloadNft for SuperNodeAPIInterface
func (maker NftDownloadingNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftDownloadingNode{DownloadNftInterface: conn.DownloadNft()}
}

// NftDownloadingNode represent supernode connection.
type NftDownloadingNode struct {
	node.DownloadNftInterface
}
