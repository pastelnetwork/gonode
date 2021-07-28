package node

import (
	"bytes"
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
			node.fingerprint = res.FingerprintData
			node.rarenessScore = res.RarenessScore
			node.nSFWScore = res.NSFWScore
			node.seenScore = res.SeenScore
			return err
		})
	}
	return group.Wait()
}

// MatchFingerprintAndScores matches fingerprints.
func (nodes List) MatchFingerprintAndScores() error {
	node := nodes[0]

	for i := 1; i < len(nodes); i++ {
		if !bytes.Equal(node.fingerprint, nodes[i].fingerprint) {
			return errors.Errorf("fingerprints of nodes %q and %q didn't match", node.String(), nodes[i].String())
		}
		if node.rarenessScore != nodes[i].rarenessScore {
			return errors.Errorf("rareness score of nodes %q and %q didn't match", node.String(), nodes[i].String())
		}
		if node.nSFWScore != nodes[i].nSFWScore {
			return errors.Errorf("NSFWS score of nodes %q and %q didn't match", node.String(), nodes[i].String())
		}
		if node.seenScore != nodes[i].seenScore {
			return errors.Errorf("seen score of nodes %q and %q didn't match", node.String(), nodes[i].String())
		}
	}
	return nil
}

// UploadImageWithThumbnail uploads the image with pqsignatured appended and thumbnail's coordinate to super nodes
func (nodes *List) UploadImageWithThumbnail(ctx context.Context, file *artwork.File, thumbnail artwork.ThumbnailCoordinate) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			hash1, hash2, hash3, err := node.UploadImageWithThumbnail(ctx, file, thumbnail)
			if err != nil {
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

// UploadSignedTicket uploads regart ticket and its signature to super nodes
func (nodes *List) UploadSignedTicket(ctx context.Context, ticket []byte, signature []byte, rqids map[string][]byte, encoderParams rqnode.EncoderParameters) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range *nodes {
		node := node
		// TODO: Fix this when method to generate key1 and key2 are finalized
		key1 := "key1-" + uuid.New().String()
		key2 := "key2-" + uuid.New().String()
		group.Go(func() error {
			fee, err := node.SendSignedTicket(ctx, ticket, signature, key1, key2, rqids, encoderParams)
			if err != nil {
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
				return err
			}
			if !node.IsPrimary() && ticketTxid != "" {
				return errors.Errorf("receive response %s from secondary node", ticketTxid)
			}

			node.regArtTxid = ticketTxid
			return nil
		})
	}
	return group.Wait()
}

// RegistrationFee returns registration fee of the first node
func (nodes *List) RegistrationFee() int64 {
	return (*nodes)[0].registrationFee
}

// Fingerprint returns fingerprint of the first node.
func (nodes *List) Fingerprint() []byte {
	return (*nodes)[0].fingerprint
}

// RarenessScore returns rarenessScore score
func (nodes *List) RarenessScore() int {
	return (*nodes)[0].rarenessScore
}

// NSFWScore returns NSFWScore score
func (nodes *List) NSFWScore() int {
	return (*nodes)[0].nSFWScore
}

// SeenScore returns SeenScore score
func (nodes *List) SeenScore() int {
	return (*nodes)[0].seenScore
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

// SetPrimary promotes a supernode to primary role which handle the write to Kamedila
func (node *Node) SetPrimary(primary bool) {
	node.isPrimary = primary
}

// IsPrimary returns true if this node has been promoted to primary in meshNode session
func (node *Node) IsPrimary() bool {
	return node.isPrimary
}

// RegArtTicketID return txid of RegArt ticket
func (nodes *List) RegArtTicketID() string {
	for i := range *nodes {
		if (*nodes)[i].isPrimary {
			return (*nodes)[i].regArtTxid
		}
	}
	return string("")
}
