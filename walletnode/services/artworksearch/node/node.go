package node

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type NftSearchNodeClient struct {
	common.SuperNodeClient
	node.DownloadNftInterface

	fingerprint []byte
}

// Connect connects to supernode.
func (node *NftSearchNodeClient) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	err := node.SuperNodeClient.Connect(ctx, timeout, secInfo)
	if err != nil {
		return err
	}
	node.DownloadNftInterface = node.DownloadArtwork()
	return nil
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address, pastelID string) *NftSearchNodeClient {
	return &NftSearchNodeClient{
		SuperNodeClient: *common.NewSuperNode(client, address, pastelID),
	}
}
