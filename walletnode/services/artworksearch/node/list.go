package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
)

// List represents multiple Node.
type List []*NftSearchNode

// Add adds a new node to the list.
func (nodes *List) Add(node *NftSearchNode) {
	*nodes = append(*nodes, node)
}

// Activate marks all nodes as activated.
// Since any node can be present in the same time in several List and Node is a pointer, this is reflected in all lists.
func (nodes *List) Activate() {
	for _, node := range *nodes {
		node.SetActive(true)
	}
}

// DisconnectInactive disconnects nodes which were not marked as activated.
func (nodes *List) DisconnectInactive() {
	for _, node := range *nodes {
		if node.ConnectionInterface != nil && !node.IsActive() {
			node.ConnectionInterface.Close()
		}
	}
}

// WaitConnClose waits for the connection closing by any supernodes.
func (nodes *List) WaitConnClose(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			case <-node.ConnectionInterface.Done():
				return errors.Errorf("%q unexpectedly closed the connection", node)
			}
		})
	}

	return group.Wait()
}

// FindByPastelID returns node by its patstelID.
func (nodes List) FindByPastelID(id string) *NftSearchNode {
	for _, node := range nodes {
		if node.PastelID() == id {
			return node
		}
	}
	return nil
}

// Fingerprint returns fingerprint of the first node.
func (nodes List) Fingerprint() []byte {
	return nodes[0].fingerprint
}
