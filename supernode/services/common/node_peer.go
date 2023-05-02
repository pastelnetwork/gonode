package common

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/supernode/node"
)

const (
	defaultConnectToNodeTimeout = time.Second * 20
)

// SuperNodePeer represents a single supernode
type SuperNodePeer struct {
	node.ClientInterface
	node.NodeMaker
	node.ConnectionInterface
	node.SuperNodePeerAPIInterface

	ID      string
	Address string
}

// Connect connects to grpc Server and setup pointer to concrete client wrapper
func (node *SuperNodePeer) Connect(ctx context.Context) error {
	connCtx, connCancel := context.WithTimeout(ctx, defaultConnectToNodeTimeout)
	defer connCancel()

	conn, err := node.ClientInterface.Connect(connCtx, node.Address)
	if err != nil {
		return err
	}

	node.ConnectionInterface = conn
	node.SuperNodePeerAPIInterface = node.MakeNode(conn)
	return nil
}

// NewSuperNode returns a new Node instance.
func NewSuperNode(client node.ClientInterface, address string, pastelID string,
	nodeMaker node.NodeMaker) *SuperNodePeer {
	return &SuperNodePeer{
		ClientInterface: client,
		NodeMaker:       nodeMaker,
		Address:         address,
		ID:              pastelID,
	}
}

// SuperNodePeerList represents muptiple SenseRegistrationNodes
type SuperNodePeerList []*SuperNodePeer

// Add adds a new node to the list
func (list *SuperNodePeerList) Add(node *SuperNodePeer) {
	*list = append(*list, node)
}

// ByID returns a node from the list by the given id.
func (list *SuperNodePeerList) ByID(id string) *SuperNodePeer {
	for _, someNode := range *list {
		if someNode.ID == id {
			return someNode
		}
	}
	return nil
}

// Remove removes a node from the list by the given id.
func (list *SuperNodePeerList) Remove(id string) {
	for i, someNode := range *list {
		if someNode.ID == id {
			if i+1 < len(*list) {
				*list = append((*list)[:i], (*list)[i+1:]...)
			} else {
				*list = (*list)[:i]
			}
			break
		}
	}
}
