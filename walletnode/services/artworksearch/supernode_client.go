package artworksearch

import "github.com/pastelnetwork/gonode/walletnode/node"

type SearchNftNodeMaker struct {
	node.NodeMaker
}

func (maker SearchNftNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftSearchNode{DownloadNftInterface: conn.DownloadArtwork()}
}

// Node represent supernode connection.
type NftSearchNode struct {
	node.DownloadNftInterface
}
