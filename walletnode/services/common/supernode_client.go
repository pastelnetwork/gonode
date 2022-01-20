package common

import (
	"context"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"sync"
	"time"
)

type SuperNodeClient struct {
	node.ClientInterface
	node.NodeMaker
	node.ConnectionInterface
	node.SuperNodeAPIInterface

	mtx *sync.RWMutex

	address  string
	pastelID string

	isPrimary bool
	activated bool
}

// String returns node as string (address)
func (node *SuperNodeClient) String() string {
	return node.address
}

// Address returns address of node
func (node *SuperNodeClient) Address() string {
	return node.address
}

// PastelID returns pastelID
func (node *SuperNodeClient) PastelID() string {
	return node.pastelID
}

// PastelID returns pastelID
func (node *SuperNodeClient) SetPastelID(pastelID string) {
	node.pastelID = pastelID
}

// Connect connects to supernode.
func (node *SuperNodeClient) Connect(ctx context.Context, timeout time.Duration, secInfo *alts.SecInfo) error {
	node.mtx.Lock()
	defer node.mtx.Unlock()

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
	node.SuperNodeAPIInterface = node.MakeNode(conn)
	return nil
}

// SetPrimary promotes a supernode to primary role which handle the writes to Kamedila
func (node *SuperNodeClient) SetPrimary(primary bool) {
	node.isPrimary = primary
}

// IsPrimary returns true if this node has been promoted to primary in meshNode session
func (node *SuperNodeClient) IsPrimary() bool {
	return node.isPrimary
}

// SetActive set nodes active or not
func (node *SuperNodeClient) SetActive(active bool) {
	node.activated = active
}

// SetActive set nodes active or not
func (node *SuperNodeClient) IsActive() bool {
	return node.activated
}

// RLock set nodes active or not
func (node *SuperNodeClient) RLock() {
	node.mtx.RLock()
}

// RUnLock set nodes active or not
func (node *SuperNodeClient) RUnlock() {
	node.mtx.RUnlock()
}

// NewSuperNode returns a new Node instance.
func NewSuperNode(client node.ClientInterface,
	address string, pastelID string,
	nodeMaker node.NodeMaker,
) *SuperNodeClient {
	return &SuperNodeClient{
		ClientInterface: client,
		NodeMaker:       nodeMaker,
		mtx:             &sync.RWMutex{},
		address:         address,
		pastelID:        pastelID,
	}
}
