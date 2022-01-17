package node

import (
	"context"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"sync"
	"time"
)

type BaseNode struct {
	mtx *sync.RWMutex
	ClientInterface
	ConnectionInterface

	isPrimary bool
	activated bool
	address   string
	pastelID  string
}

// String returns node as string (address)
func (node *BaseNode) String() string {
	return node.address
}

// Address returns address of node
func (node *BaseNode) Address() string {
	return node.address
}

// PastelID returns pastelID
func (node *BaseNode) PastelID() string {
	return node.pastelID
}

// PastelID returns pastelID
func (node *BaseNode) SetPastelID(pastelID string) {
	node.pastelID = pastelID
}

// Connect connects to supernode.
func (node *BaseNode) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) (ConnectionInterface, error) {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if node.ConnectionInterface != nil {
		return node.ConnectionInterface, nil
	}

	connCtx, connCancel := context.WithTimeout(ctx, timeout)
	defer connCancel()

	conn, err := node.ClientInterface.Connect(connCtx, node.address, secInfo)
	if err != nil {
		return nil, err
	}
	node.ConnectionInterface = conn
	return node.ConnectionInterface, nil
}

// SetPrimary promotes a supernode to primary role which handle the write to Kamedila
func (node *BaseNode) SetPrimary(primary bool) {
	node.isPrimary = primary
}

// IsPrimary returns true if this node has been promoted to primary in meshNode session
func (node *BaseNode) IsPrimary() bool {
	return node.isPrimary
}

// SetActive set nodes active or not
func (node *BaseNode) SetActive(active bool) {
	node.activated = active
}

// SetActive set nodes active or not
func (node *BaseNode) IsActive() bool {
	return node.activated
}

// RLock set nodes active or not
func (node *BaseNode) RLock() {
	node.mtx.RLock()
}

// RUnLock set nodes active or not
func (node *BaseNode) RUnlock() {
	node.mtx.RUnlock()
}

// NewBaseNode returns a new Node instance.
func NewBaseNode(client ClientInterface, address string, pastelID string) *BaseNode {
	return &BaseNode{
		ClientInterface: client,
		address:         address,
		pastelID:        pastelID,
		mtx:             &sync.RWMutex{},
	}
}
