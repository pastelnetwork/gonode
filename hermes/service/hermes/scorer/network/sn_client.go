package network

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/hermes/service/node"
)

// SuperNodeClient represents base SN client
type SuperNodeClient struct {
	node.SNClientInterface
	node.RealNodeMaker
	node.ConnectionInterface
	node.SuperNodeAPIInterface

	mtx *sync.RWMutex

	address  string
	pastelID string
	txid     string

	isPrimary bool
	activated bool

	isRemoteState bool
	idx           int
}

// String returns node as string (address)
func (node *SuperNodeClient) String() string {
	return node.address
}

// PastelID returns pastelID
func (node *SuperNodeClient) PastelID() string {
	return node.pastelID
}

// Idx returns outpoint idx
func (node *SuperNodeClient) Idx() int {
	return node.idx
}

// TxID returns txid
func (node *SuperNodeClient) TxID() string {
	return node.txid
}

// Address returns address of node
func (node *SuperNodeClient) Address() string {
	return node.address
}

// SetPrimary promotes a supernode to primary role which handle the writes to Kademila
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

// IsActive returns true if nodes is active
func (node *SuperNodeClient) IsActive() bool {
	return node.activated
}

// SetRemoteState set state returned by nodes
func (node *SuperNodeClient) SetRemoteState(remote bool) {
	node.isRemoteState = remote
}

// IsRemoteState returns true if remote nodes processing status is ok
func (node *SuperNodeClient) IsRemoteState() bool {
	return node.isRemoteState
}

// RLock set nodes active or not
func (node *SuperNodeClient) RLock() {
	node.mtx.RLock()
}

// RUnlock set nodes active or not
func (node *SuperNodeClient) RUnlock() {
	node.mtx.RUnlock()
}

// Connect connects to supernode.
func (node *SuperNodeClient) Connect(ctx context.Context, timeout time.Duration, _ *alts.SecInfo) error {
	node.mtx.Lock()
	defer node.mtx.Unlock()

	if node.ConnectionInterface != nil {
		return nil
	}

	connCtx, connCancel := context.WithTimeout(ctx, timeout)
	defer connCancel()

	conn, err := node.SNClientInterface.Connect(connCtx, node.address)
	if err != nil {
		return err
	}

	node.ConnectionInterface = conn
	node.SuperNodeAPIInterface = node.MakeNode(conn)

	return nil
}

// NewSuperNode returns a new Node instance.
func NewSuperNode(client node.SNClientInterface, txid string, address string, pastelID string, idx int,
	nodeMaker node.RealNodeMaker) *SuperNodeClient {
	return &SuperNodeClient{
		SNClientInterface: client,
		RealNodeMaker:     nodeMaker,
		txid:              txid,
		mtx:               &sync.RWMutex{},
		address:           address,
		pastelID:          pastelID,
		idx:               idx,
	}
}
