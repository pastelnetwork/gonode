package node

import (
	"bytes"
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
)

// List represents multiple Node.
type List []*Node

// Add adds a new node to the list.
func (nodes *List) Add(node *Node) {
	*nodes = append(*nodes, node)
}

// Active marks all nodes as activated.
// Since any node can be present in the same time in several List and Node is a pointer, this is reflected in all lists.
func (nodes *List) Active() List {
	activeNodes := List{}
	for _, node := range *nodes {
		if node.activated {
			activeNodes = append(activeNodes, node)
		}
	}
	return activeNodes
}

// DisconnectInactive disconnects nodes which were not marked as activated.
func (nodes *List) DisconnectInactive() {
	for _, node := range *nodes {
		if node.Connection != nil && !node.activated {
			node.Connection.Close()
		}
	}
}

// Disconnect disconnects all nodes.
func (nodes *List) Disconnect() {
	for _, node := range *nodes {
		if node.Connection != nil {
			node.Connection.Close()
			node.done = true
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
			case <-node.Connection.Done():
				if node.done {
					return nil
				}
				return errors.Errorf("%q unexpectedly closed the connection", node)
			}
		})
	}

	return group.Wait()
}

// MatchFiles matches files.
func (nodes List) MatchFiles() error {
	node := nodes[0]

	for i := 1; i < len(nodes); i++ {
		if !bytes.Equal(node.file, nodes[i].file) {
			return errors.Errorf("file of nodes %q and %q didn't match", node.String(), nodes[i].String())
		}
	}
	return nil
}

// File returns downloaded file.
func (nodes List) File() []byte {
	return nodes[0].file
}

// Download download image from supernodes.
func (nodes *List) Download(ctx context.Context, txid, timestamp, signature, ttxid string, downloadErrs error) error {
	group, _ := errgroup.WithContext(ctx)
	errChan := make(chan error, len(*nodes))

	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			node.file, err = node.Download(ctx, txid, timestamp, signature, ttxid)
			if err != nil {
				errChan <- err
			} else {
				node.activated = true
			}
			return err
		})
	}
	err := group.Wait()

	close(errChan)

	for err := range errChan {
		downloadErrs = errors.Append(downloadErrs, err)
	}
	return err
}
