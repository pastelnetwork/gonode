package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
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

// DisconnectAll disconnects all nodes
func (nodes *List) DisconnectAll() {
	for _, node := range *nodes {
		node.mtx.RLock()
		defer node.mtx.RUnlock()

		if node.Connection != nil {
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

// SendRegMetadata send metadata
func (nodes *List) SendRegMetadata(ctx context.Context, regMetadata *types.NftRegMetadata) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			err = node.SendRegMetadata(ctx, regMetadata)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", node.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

// ProbeImage sends the image to supernodes for image analysis, such as fingerprint, raraness score, NSWF.
func (nodes *List) ProbeImage(ctx context.Context, file *artwork.File) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() (err error) {
			compress, isValidBurnTxID, err := node.ProbeImage(ctx, file)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", node.String(), err)
			}

			node.isValidBurnTxID = isValidBurnTxID
			node.FingerprintAndScores, node.FingerprintAndScoresBytes, node.Signature, err = pastel.ExtractCompressSignedDDAndFingerprints(compress)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("extract compressed signed DDAandFingerprints failed")
				return errors.Errorf("node %s: extract failed: %w", node.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

// MatchFingerprintAndScores matches fingerprints.
func (nodes *List) MatchFingerprintAndScores() error {
	node := (*nodes)[0]
	for i := 1; i < len(*nodes); i++ {
		if err := pastel.CompareFingerPrintAndScore(node.FingerprintAndScores, (*nodes)[i].FingerprintAndScores); err != nil {
			return errors.Errorf("node[%s] and node[%s] not matched: %w", node.PastelID(), (*nodes)[i].PastelID(), err)
		}
	}

	return nil
}

// ValidBurnTxID returns whether the burn txid is valid
func (nodes *List) ValidBurnTxID() bool {
	for _, node := range *nodes {
		if !node.isValidBurnTxID {
			return false
		}
	}

	return true
}

// FingerAndScores returns fingerprint of the image and dupedetection scores
func (nodes *List) FingerAndScores() *pastel.DDAndFingerprints {
	return (*nodes)[0].FingerprintAndScores
}

// UploadSignedTicket uploads regart ticket and its signature to super nodes
func (nodes *List) UploadSignedTicket(ctx context.Context, ticket []byte, signature []byte, ddFpFile []byte) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			fee, err := node.SendSignedTicket(ctx, ticket, signature, ddFpFile)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("send signed ticket failed")
				return err
			}
			node.registrationFee = fee
			return nil
		})
	}
	return group.Wait()
}

// MatchRegistrationFee matches registration fee returned by super nodes
func (nodes *List) MatchRegistrationFee() error {
	node := (*nodes)[0]
	for i := 1; i < len(*nodes); i++ {
		if node.registrationFee != (*nodes)[i].registrationFee {
			return errors.Errorf("registration fee of nodes %q and %q didn't match", node.String(), (*nodes)[i].String())
		}
	}
	return nil
}

// SendPreBurntFeeTxid send txid of transaction in which 10% of registration fee is preburnt
func (nodes *List) SendPreBurntFeeTxid(ctx context.Context, txid string) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			ticketTxid, err := node.SendPreBurntFeeTxid(ctx, txid)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("send pre-burnt fee txid failed")
				return err
			}
			if !node.IsPrimary() && ticketTxid != "" {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, node.pastelID)
			}

			node.regNFTTxid = ticketTxid
			return nil
		})
	}
	return group.Wait()
}

// RegistrationFee returns registration fee of the first node
func (nodes *List) RegistrationFee() int64 {
	return (*nodes)[0].registrationFee
}

// CompressedFingerAndScores returns compressed fingerprint and other scores
func (nodes *List) CompressedFingerAndScores() *pastel.DDAndFingerprints {
	return (*nodes)[0].FingerprintAndScores
}

// PreviewHash returns the hash of the preview thumbnail calculated by the first node
func (nodes *List) PreviewHash() []byte {
	return (*nodes)[0].previewHash
}

// MediumThumbnailHash returns the hash of the medium thumbnail calculated by the first node
func (nodes *List) MediumThumbnailHash() []byte {
	return (*nodes)[0].previewHash
}

// SmallThumbnailHash returns the hash of the small thumbnail calculated by the first node
func (nodes *List) SmallThumbnailHash() []byte {
	return (*nodes)[0].smallThumbnailHash
}

// RegArtTicketID return txid of RegArt ticket
func (nodes *List) RegArtTicketID() string {
	for i := range *nodes {
		if (*nodes)[i].isPrimary {
			return (*nodes)[i].regNFTTxid
		}
	}
	return string("")
}
