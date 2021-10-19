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
	node.ExternalStorage
	node.Connection

	isPrimary bool
	activated bool

	storageFee             int64
	regExternalStorageTxID string

	address  string
	pastelID string
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
	node.ExternalStorage = conn.ExternalStorage()
	return nil
}

// SetPrimary promotes a supernode to primary role which handle the write to Kamedila
func (node *Node) SetPrimary(primary bool) {
	node.isPrimary = primary
}

// IsPrimary returns true if this node has been promoted to primary in meshNode session
func (node *Node) IsPrimary() bool {
	return node.isPrimary
}

// Address returns address of node
func (node *Node) Address() string {
	return node.address
}

// NewNode returns a new Node instance.
func NewNode(client node.Client, address, pastelID string) *Node {
	return &Node{
		Client:   client,
		address:  address,
		pastelID: pastelID,
		mtx:      &sync.RWMutex{},
	}
}