package node

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type NftDownloadNode struct {
	node.BaseNode
	node.DownloadNftInterface

	done bool
	file []byte
}

// Connect connects to supernode.
func (node *NftDownloadNode) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	conn, err := node.BaseNode.Connect(ctx, timeout, secInfo)
	if err != nil {
		return err
	}
	node.DownloadNftInterface = conn.DownloadArtwork()
	return nil
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address, pastelID string) *NftDownloadNode {
	return &NftDownloadNode{
		BaseNode: *node.NewBaseNode(client, address, pastelID),
	}
}
