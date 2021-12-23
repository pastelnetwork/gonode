package node

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Node represent supernode connection.
type Node struct {
	mtx *sync.RWMutex

	node.Client
	node.Connection
	node.SenseRegister

	isPrimary bool
	activated bool

	address  string
	pastelID string
}

// String returns node address.
func (node *Node) String() string {
	return node.address
}

// PastelID returns pastelID
func (node *Node) PastelID() string {
	return node.pastelID
}

// IsPrimary returns true if node is primary
func (node *Node) IsPrimary() bool {

	return node.isPrimary
}

// SetPrimary promotes a supernode to primary role which handle the write to Kamedila
func (node *Node) SetPrimary(primary bool) {
	node.isPrimary = primary
}

// Connect connects to supernode.
func (node *Node) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if node.Connection != nil {
		return nil
	}

	connCtx, connCancel := context.WithTimeout(ctx, timeout)
	defer connCancel()

	conn, err := node.Client.Connect(connCtx, node.address, secInfo)
	if err != nil {
		return err
	}

	node.Connection = conn
	return nil
}

// NewNOde returns new Node instance
func NewNode(client node.Client, address, pastelID string) *Node {
	return &Node{
		mtx:      &sync.RWMutex{},
		Client:   client,
		address:  address,
		pastelID: pastelID,
	}
}
