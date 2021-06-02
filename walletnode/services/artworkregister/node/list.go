package node

import (
	"bytes"
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
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
		if node.Connection != nil && !node.activated {
			node.Connection.Close()
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
				return errors.Errorf("%q unexpectedly closed the connection", node)
			}
		})
	}

	return group.Wait()
}

// FindByPastelID returns node by its patstelID.
func (nodes List) FindByPastelID(id string) *Node {
	for _, node := range nodes {
		if node.pastelID == id {
			return node
		}
	}
	return nil
}

// ProbeImage sends the image to supernodes for image analysis, such as fingerprint, raraness score, NSWF.
func (nodes *List) ProbeImage(ctx context.Context, file *artwork.File) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			node.fingerprint, err = node.ProbeImage(ctx, file)
			return err
		})
	}
	return group.Wait()
}

// MatchFingerprints matches fingerprints.
func (nodes *List) MatchFingerprints() error {
	node := (*nodes)[0]

	for i := 1; i < len(*nodes); i++ {
		if !bytes.Equal(node.fingerprint, (*nodes)[i].fingerprint) {
			return errors.Errorf("fingerprints of nodes %q and %q didn't match", node.String(), (*nodes)[i].String())
		}
	}
	return nil
}
