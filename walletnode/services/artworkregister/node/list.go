package node

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

// List represents multiple Node.
type List []*NftRegisterNodeClient

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
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, node.PastelID())
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
		if (*nodes)[i].IsPrimary() {
			return (*nodes)[i].regNFTTxid
		}
	}
	return string("")
}
