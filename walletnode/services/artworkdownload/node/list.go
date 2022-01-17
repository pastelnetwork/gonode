package node

import (
	"bytes"
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
)

// List represents multiple Node.
type List []*NftDownloadNode

// Add adds a new node to the list.
func (nodes *List) Add(node *NftDownloadNode) {
	*nodes = append(*nodes, node)
}

// Active marks all nodes as activated.
// Since any node can be present in the same time in several List and Node is a pointer, this is reflected in all lists.
func (nodes *List) Activate() List {
	activeNodes := List{}
	for _, node := range *nodes {
		if node.IsActive() {
			activeNodes = append(activeNodes, node)
		}
	}
	return activeNodes
}

// DisconnectInactive disconnects nodes which were not marked as activated.
func (nodes *List) DisconnectInactive() {
	for _, node := range *nodes {
		if node.ConnectionInterface != nil && !node.IsActive() {
			node.ConnectionInterface.Close()
		}
	}
}

// Disconnect disconnects all nodes.
func (nodes *List) Disconnect() {
	for _, node := range *nodes {
		if node.ConnectionInterface != nil {
			node.ConnectionInterface.Close()
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
			case <-node.ConnectionInterface.Done():
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
func (nodes *List) Download(ctx context.Context, txid, timestamp, signature, ttxid string, connectTimeout time.Duration, secInfo *alts.SecInfo) ([]error, error) {
	group, _ := errgroup.WithContext(ctx)
	errChan := make(chan error, len(*nodes))

	for _, node := range *nodes {
		subNode := node
		group.Go(func() error {
			var subErr error
			if subErr := subNode.Connect(ctx, connectTimeout, secInfo); subErr == nil {
				log.WithContext(ctx).WithField("address", subNode.String()).Info("Connected to supernode")
			} else {
				log.WithContext(ctx).WithError(subErr).WithField("address", subNode.String()).Error("Could not connect to supernode")
				errChan <- subErr
				return nil
			}

			subNode.file, subErr = subNode.Download(ctx, txid, timestamp, signature, ttxid)
			if subErr != nil {
				log.WithContext(ctx).WithField("address", subNode.String()).WithError(subErr).Error("Could not download from supernode")
				errChan <- subErr
			} else {
				log.WithContext(ctx).WithField("address", subNode.String()).Info("Downloaded from supernode")
				subNode.SetActive(true)
			}
			return nil
		})
	}
	err := group.Wait()

	close(errChan)

	downloadErrors := []error{}
	for subErr := range errChan {
		downloadErrors = append(downloadErrors, subErr)
	}

	return downloadErrors, err
}
