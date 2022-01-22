package node

import (
	"bytes"
	"context"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/pastel"
)

// List represents multiple Node.
type List []*NftRegisterNodeClient

// UploadImageWithThumbnail uploads the image with pqsignatured appended and thumbnail's coordinate to super nodes
func (nodes *List) UploadImageWithThumbnail(ctx context.Context, file *files.File, thumbnail files.ThumbnailCoordinate) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			hash1, hash2, hash3, err := node.UploadImageWithThumbnail(ctx, file, thumbnail)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("upload image with thumbnail failed")
				return err
			}
			node.previewHash, node.mediumThumbnailHash, node.smallThumbnailHash = hash1, hash2, hash3
			return nil
		})
	}
	return group.Wait()
}

// MatchThumbnailHashes matches thumbnail's hashes recevied from super nodes
func (nodes *List) MatchThumbnailHashes() error {
	node := (*nodes)[0]
	for i := 1; i < len(*nodes); i++ {
		if !bytes.Equal(node.previewHash, (*nodes)[i].previewHash) {
			return errors.Errorf("hash of preview thumbnail of nodes %q and %q didn't match", node.String(), (*nodes)[i].String())
		}
		if !bytes.Equal(node.mediumThumbnailHash, (*nodes)[i].mediumThumbnailHash) {
			return errors.Errorf("hash of medium thumbnail of nodes %q and %q didn't match", node.String(), (*nodes)[i].String())
		}
		if !bytes.Equal(node.smallThumbnailHash, (*nodes)[i].smallThumbnailHash) {
			return errors.Errorf("hash of small thumbnail of nodes %q and %q didn't match", node.String(), (*nodes)[i].String())
		}
	}
	return nil
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
