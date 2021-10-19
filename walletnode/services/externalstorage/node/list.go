package node

import (
	"context"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
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
		node.mtx.RLock()
		defer node.mtx.RUnlock()

		if node.Connection != nil && !node.activated {
			node.Connection.Close()
		}
	}
}

// WaitConnClose waits for the connection closing by any supernodes.
func (nodes *List) WaitConnClose(ctx context.Context, done <-chan struct{}) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			case <-node.Connection.Done():
				return errors.Errorf("%q unexpectedly closed the connection", node)
			case <-done:
				return nil
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

// UploadImage uploads the image with pqsignatured appended to super nodes
func (nodes *List) UploadImage(ctx context.Context, file *artwork.File) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			return node.UploadImage(ctx, file)
		})
	}
	return group.Wait()
}

// UploadSignedTicket uploads regart ticket and its signature to super nodes
func (nodes *List) UploadSignedTicket(ctx context.Context, ticket []byte, signature []byte, rqids map[string][]byte, encoderParams rqnode.EncoderParameters) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		// TODO: Fix this when method to generate key1 and key2 are finalized
		key1 := "key1-" + uuid.New().String()
		key2 := "key2-" + uuid.New().String()
		group.Go(func() error {
			fee, err := node.SendSignedExternalStorageTicket(ctx, ticket, signature, key1, key2, rqids, encoderParams)
			if err != nil {
				return err
			}
			node.storageFee = fee
			return nil
		})
	}
	return group.Wait()
}

// MatchStorageFee matches storage fee returned by super nodes
func (nodes *List) MatchStorageFee() error {
	node := (*nodes)[0]
	for i := 1; i < len(*nodes); i++ {
		if node.storageFee != (*nodes)[i].storageFee {
			return errors.Errorf("storage fee of nodes %q and %q didn't match", node.String(), (*nodes)[i].String())
		}
	}
	return nil
}

// SendPreBurnedFeeTxID send txid of transaction in which 10% of storage fee is preburnt
func (nodes *List) SendPreBurnedFeeTxID(ctx context.Context, txid string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			ticketTxID, err := node.SendPreBurnedFeeExternalStorageTxID(ctx, txid)
			if err != nil {
				return err
			}
			if !node.IsPrimary() && ticketTxID != "" {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxID, node.pastelID)
			}

			node.regExternalStorageTxID = ticketTxID
			return nil
		})
	}
	return group.Wait()
}

// StorageFee returns storage fee of the first node
func (nodes *List) StorageFee() int64 {
	return (*nodes)[0].storageFee
}

// RegEDDTicketID return txid of RegEDD ticket
func (nodes *List) RegEDDTicketID() string {
	for i := range *nodes {
		if (*nodes)[i].isPrimary {
			return (*nodes)[i].regExternalStorageTxID
		}
	}
	return string("")
}
