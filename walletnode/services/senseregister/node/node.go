package node

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type SenseRegisterNode struct {
	node.BaseNode
	node.RegisterSenseInterface
	common.FingerprintsHandler

	regActionTxid   string
	isValidBurnTxID bool
}

// Connect connects to supernode.
func (node *SenseRegisterNode) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	conn, err := node.BaseNode.Connect(ctx, timeout, secInfo)
	if err != nil {
		return err
	}
	node.RegisterSenseInterface = conn.RegisterSense()
	return nil
}

// SetValidBurnTxID sets whether the burn txid is valid
func (node *SenseRegisterNode) SetValidBurnTxID(valid bool) {
	node.isValidBurnTxID = valid
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address string, pastelID string) *SenseRegisterNode {
	return &SenseRegisterNode{
		BaseNode: *node.NewBaseNode(client, address, pastelID),
	}
}
