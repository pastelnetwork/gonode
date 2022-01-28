package nftsearch

import "github.com/pastelnetwork/gonode/walletnode/node"

type NftSearchNodeMaker struct {
	node.NodeMaker
}

func (maker NftSearchNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftSearchNode{DownloadNftInterface: conn.DownloadNft()}
}

// Node represent supernode connection.
type NftSearchNode struct {
	node.DownloadNftInterface
}
