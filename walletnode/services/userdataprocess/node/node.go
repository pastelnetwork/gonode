package node

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type Node struct {
	node.ClientInterface
	node.ConnectionInterface
	node.ProcessUserdataInterface
	activated bool
	address   string
	pastelID  string
	Result    *userdata.ProcessResult
	ResultGet *userdata.ProcessRequest
	isPrimary bool
}

func (node *Node) String() string {
	return node.address
}

// PastelID returns pastelID
func (node *Node) PastelID() string {
	return node.pastelID
}

// Connect connects to supernode.
func (node *Node) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	if node.ConnectionInterface != nil {
		return nil
	}

	connCtx, connCancel := context.WithTimeout(ctx, timeout)
	defer connCancel()

	conn, err := node.ClientInterface.Connect(connCtx, node.address, secInfo)
	if err != nil {
		return err
	}
	node.ConnectionInterface = conn
	node.ProcessUserdataInterface = conn.ProcessUserdata()
	return nil
}

// NewNode returns a new Node instance.
func NewNode(client node.ClientInterface, address, pastelID string) *Node {
	return &Node{
		ClientInterface: client,
		address:         address,
		pastelID:        pastelID,
	}
}
