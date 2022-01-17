package node

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type NftRegisterNode struct {
	node.BaseNode
	node.RegisterNftInterface
	common.FingerprintsHandler

	// thumbnail hash
	previewHash         []byte
	mediumThumbnailHash []byte
	smallThumbnailHash  []byte

	registrationFee int64
	regNFTTxid      string
}

// Connect connects to supernode.
func (node *NftRegisterNode) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	conn, err := node.BaseNode.Connect(ctx, timeout, secInfo)
	if err != nil {
		return err
	}
	node.RegisterNftInterface = conn.RegisterArtwork()
	return nil
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address string, pastelID string) *NftRegisterNode {
	return &NftRegisterNode{
		BaseNode: *node.NewBaseNode(client, address, pastelID),
	}
}
