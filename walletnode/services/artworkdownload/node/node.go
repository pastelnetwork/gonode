package node

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type NftDownloadNodeClient struct {
	common.SuperNodeClient
	node.DownloadNftInterface

	done bool
	file []byte
}

// Connect connects to supernode.
func (node *NftDownloadNodeClient) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	err := node.SuperNodeClient.Connect(ctx, timeout, secInfo)
	if err != nil {
		return err
	}
	node.DownloadNftInterface = node.DownloadArtwork()
	return nil
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address, pastelID string) *NftDownloadNodeClient {
	return &NftDownloadNodeClient{
		SuperNodeClient: *common.NewSuperNode(client, address, pastelID),
	}
}
