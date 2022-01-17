package node

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type NftSearchNode struct {
	node.BaseNode
	node.DownloadNftInterface

	fingerprint []byte
}

// Connect connects to supernode.
func (node *NftSearchNode) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	conn, err := node.BaseNode.Connect(ctx, timeout, secInfo)
	if err != nil {
		return err
	}
	node.DownloadNftInterface = conn.DownloadArtwork()
	return nil
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address, pastelID string) *NftSearchNode {
	return &NftSearchNode{
		BaseNode: *node.NewBaseNode(client, address, pastelID),
	}
}
