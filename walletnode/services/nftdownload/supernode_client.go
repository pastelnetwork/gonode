package nftdownload

import "github.com/pastelnetwork/gonode/walletnode/node"

type NftDownloadNodeMaker struct {
	node.NodeMaker
}

func (maker NftDownloadNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftDownloadNode{DownloadNftInterface: conn.DownloadArtwork()}
}

// Node represent supernode connection.
type NftDownloadNode struct {
	node.DownloadNftInterface
}
