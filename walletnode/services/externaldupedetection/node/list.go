package node

import (
	"context"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/pastel"
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

// ProbeImage sends the image to supernodes for image analysis, such as fingerprint, raraness score, NSWF.
func (nodes *List) ProbeImage(ctx context.Context, file *artwork.File) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			res, err := node.ProbeImage(ctx, file)
			if err != nil {
				return errors.Errorf("failed to probe image: %w", err)
			}
			node.fingerprintAndScores = res
			return nil
		})
	}
	return group.Wait()
}

// MatchFingerprintAndScores matches fingerprints.
func (nodes *List) MatchFingerprintAndScores() error {
	node := (*nodes)[0]
	for i := 1; i < len(*nodes); i++ {
		if err := pastel.CompareFingerPrintAndScore(node.fingerprintAndScores, (*nodes)[i].fingerprintAndScores); err != nil {
			return errors.Errorf("fingerprint or score of node[%s] and node[%s] not matched: %w", node.PastelID(), (*nodes)[i].PastelID(), err)
		}
	}
	return nil
}

// FingerAndScores returns fingerprint of the image and dupedetection scores
func (nodes *List) FingerAndScores() *pastel.FingerAndScores {
	return (*nodes)[0].fingerprintAndScores
}

// UploadImage uploads the image with pqsignatured appended to super nodes
func (nodes *List) UploadImage(ctx context.Context, file *artwork.File) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {

			if err := node.UploadImage(ctx, file); err != nil {
				return err
			}
			return nil
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
			fee, err := node.SendSignedEDDTicket(ctx, ticket, signature, key1, key2, rqids, encoderParams)
			if err != nil {
				return err
			}
			node.registrationFee = fee
			return nil
		})
	}
	return group.Wait()
}

// MatchDetectionFee matches registration fee returned by super nodes
func (nodes *List) MatchDetectionFee() error {
	node := (*nodes)[0]
	for i := 1; i < len(*nodes); i++ {
		if node.registrationFee != (*nodes)[i].registrationFee {
			return errors.Errorf("registration fee of nodes %q and %q didn't match", node.String(), (*nodes)[i].String())
		}
	}
	return nil
}

// SendPreBurnedFeeTxID send txid of transaction in which 10% of registration fee is preburnt
func (nodes *List) SendPreBurnedFeeTxID(ctx context.Context, txid string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			ticketTxid, err := node.SendPreBurnedFeeEDDTxID(ctx, txid)
			if err != nil {
				return err
			}
			if !node.IsPrimary() && ticketTxid != "" {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, node.pastelID)
			}

			node.regEDDTxid = ticketTxid
			return nil
		})
	}
	return group.Wait()
}

// DetectionFee returns registration fee of the first node
func (nodes *List) DetectionFee() int64 {
	return (*nodes)[0].registrationFee
}

// CompressedFingerAndScores returns compressed fingerprint and other scores
func (nodes *List) CompressedFingerAndScores() *pastel.FingerAndScores {
	return (*nodes)[0].fingerprintAndScores
}

// RegEDDTicketID return txid of RegEDD ticket
func (nodes *List) RegEDDTicketID() string {
	for i := range *nodes {
		if (*nodes)[i].isPrimary {
			return (*nodes)[i].regEDDTxid
		}
	}
	return string("")
}
