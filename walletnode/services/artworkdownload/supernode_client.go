package artworkdownload

import "github.com/pastelnetwork/gonode/walletnode/node"

type NftDownloadNodeMaker struct {
	node.NodeMaker
}

func (maker NftDownloadNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &NftDownload{DownloadNftInterface: conn.DownloadArtwork()}
}

// Node represent supernode connection.
type NftDownload struct {
	node.DownloadNftInterface
}
