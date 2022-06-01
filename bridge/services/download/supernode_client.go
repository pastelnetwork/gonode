package download

import "github.com/pastelnetwork/gonode/bridge/node"

// DownloadingNodeMaker makes class DownloadNft for SuperNodeAPIInterface
type DownloadingNodeMaker struct {
	node.RealNodeMaker
}

// MakeNode makes class DownloadNft for SuperNodeAPIInterface
func (maker DownloadingNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &DownloadNode{DownloadDataInterface: conn.DownloadData()}
}

// DownloadNode represent supernode connection.
type DownloadNode struct {
	node.DownloadDataInterface
}
