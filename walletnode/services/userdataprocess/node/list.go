package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/userdata"
)

// List represents multiple Node.
type List []*Node

// Add adds a new node to the list.
func (nodes *List) Add(node *Node) {
	*nodes = append(*nodes, node)
}

// Activate marks all nodes as activated.
// Since any node can be present in the same time in several List and Node is a pointer, this is reflected in all lists.
func (nodes *List) Activate() {
	for _, node := range *nodes {
		node.activated = true
	}
}

// DisconnectInactive disconnects nodes which were not marked as activated.
func (nodes *List) DisconnectInactive() {
	for _, node := range *nodes {
		if node.ConnectionInterface != nil && !node.activated {
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
func (nodes *List) FindByPastelID(id string) *Node {
	for _, node := range *nodes {
		if node.pastelID == id {
			return node
		}
	}
	return nil
}

// SendUserdata sends the userdata to supernodes for verification and store in MDL
func (nodes *List) SendUserdata(ctx context.Context, req *userdata.ProcessRequestSigned) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			res, err := node.SendUserdata(ctx, req)
			node.Result = res
			return err
		})
	}
	return group.Wait()
}

// ReceiveUserdata retrieve the userdata from Metadata layer
func (nodes *List) ReceiveUserdata(ctx context.Context, pasteluserid string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			res, err := node.ReceiveUserdata(ctx, pasteluserid)
			node.ResultGet = res
			return err
		})
	}
	return group.Wait()
}

// SetPrimary promotes a supernode to primary role \
func (node *Node) SetPrimary(primary bool) {
	node.isPrimary = primary
}

// IsPrimary returns true if this node has been promoted to primary in meshNode session
func (node *Node) IsPrimary() bool {
	return node.isPrimary
}
