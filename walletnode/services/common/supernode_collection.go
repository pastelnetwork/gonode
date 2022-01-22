package common

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type SuperNodeList []*SuperNodeClient

// AddNewNode created and adds a new node to the list.
func (nodes *SuperNodeList) AddNewNode(client node.ClientInterface, address string, pastelID string, nodeMaker node.NodeMaker) {
	someNode := NewSuperNode(client, address, pastelID, nodeMaker)
	*nodes = append(*nodes, someNode)
}

// Add adds a new node to the list.
func (nodes *SuperNodeList) Add(node *SuperNodeClient) {
	*nodes = append(*nodes, node)
}

// Activate marks all nodes as activated.
// Since any node can be present in the same time in several List and Node is a pointer, this is reflected in all lists.
func (nodes *SuperNodeList) Activate() {
	for _, node := range *nodes {
		node.SetActive(true)
	}
}

// DisconnectAll disconnects all nodes
func (nodes *SuperNodeList) DisconnectAll() {
	for _, node := range *nodes {
		node.RLock()
		defer node.RUnlock()

		if node.ConnectionInterface != nil {
			node.ConnectionInterface.Close()
		}
	}
}

// WaitConnClose waits for the connection closing by any supernodes.
func (nodes *SuperNodeList) WaitConnClose(ctx context.Context, done <-chan struct{}) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			case <-node.ConnectionInterface.Done():
				return errors.Errorf("%q unexpectedly closed the connection", node)
			case <-done:
				return nil
			}
		})
	}

	return group.Wait()
}

// FindByPastelID returns node by its patstelID.
func (nodes *SuperNodeList) FindByPastelID(id string) *SuperNodeClient {
	for _, someNode := range *nodes {
		if someNode.PastelID() == id {
			return someNode
		}
	}
	return nil
}
