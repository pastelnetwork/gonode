package senseregister

import (
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type SenseRegisterNodeMaker struct {
	node.NodeMaker
}

func (maker SenseRegisterNodeMaker) MakeNode(conn node.ConnectionInterface) node.SuperNodeAPIInterface {
	return &SenseRegisterNode{RegisterSenseInterface: conn.RegisterSense()}
}

// Node represent supernode connection.
type SenseRegisterNode struct {
	node.RegisterSenseInterface

	regActionTxid   string
	isValidBurnTxID bool
}

// SetValidBurnTxID sets whether the burn txid is valid
func (node *SenseRegisterNode) SetValidBurnTxID(valid bool) {
	node.isValidBurnTxID = valid
}

//// Connect connects to supernode.
//func (node *SenseRegisterNodeClient) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
//	err := node.SuperNodeClient.Connect(ctx, timeout, secInfo)
//	if err != nil {
//		return err
//	}
//	node.RegisterSenseInterface = node.RegisterSense()
//	return nil
//}
//
//// NewNode returns a new Node instance.
//func NewNode(client node.ClientInterface, address string, pastelID string) *SenseRegisterNodeClient {
//	return &SenseRegisterNodeClient{
//		SuperNodeClient: *common.NewSuperNode(client, address, pastelID),
//	}
//}
