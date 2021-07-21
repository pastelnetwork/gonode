package node

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/service/userdata"
)

// Node represent supernode connection.
type Node struct {
	node.Client
	node.Connection
	node.ProcessUserdata
	activated     bool
	address  string
	pastelID string
	result *UserdataProcessResult
}

func (node *Node) String() string {
	return node.address
}

// PastelID returns pastelID
func (node *Node) PastelID() string {
	return node.pastelID
}

// Connect connects to supernode.
func (node *Node) Connect(ctx context.Context, timeout time.Duration) error {
	if node.Connection != nil {
		return nil
	}

	connCtx, connCancel := context.WithTimeout(ctx, timeout)
	defer connCancel()

	conn, err := node.Client.Connect(connCtx, node.address)
	if err != nil {
		return err
	}
	node.Connection = conn
	node.ProcessUserdata = conn.ProcessUserdata()
	return nil
}


// NewNode returns a new Node instance.
func NewNode(client node.Client, address, pastelID string) *Node {
	return &Node{
		Client:   client,
		address:  address,
		pastelID: pastelID,
	}
}
