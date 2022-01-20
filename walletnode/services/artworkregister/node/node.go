package node

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type NftRegisterNodeClient struct {
	common.SuperNodeClient
	node.RegisterNftInterface
	mixins.FingerprintsHandler

	// thumbnail hash
	previewHash         []byte
	mediumThumbnailHash []byte
	smallThumbnailHash  []byte

	registrationFee int64
	regNFTTxid      string
}

// Connect connects to supernode.
func (node *NftRegisterNodeClient) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	err := node.SuperNodeClient.Connect(ctx, timeout, secInfo)
	if err != nil {
		return err
	}
	node.RegisterNftInterface = node.RegisterArtwork()
	return nil
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address string, pastelID string) *NftRegisterNodeClient {
	return &NftRegisterNodeClient{
		SuperNodeClient: *common.NewSuperNode(client, address, pastelID),
	}
}
