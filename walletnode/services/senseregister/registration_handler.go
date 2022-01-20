package senseregister

import (
	"context"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
)

type SenseRegisterHandler struct {
	*mixins.MeshHandler
	*mixins.FingerprintsHandler
	isValidBurnTxID bool
}

// ValidBurnTxID returns whether the burn txid is valid at ALL SNs
func (h *SenseRegisterHandler) ValidBurnTxID() bool {
	for _, someNode := range h.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(SenseRegisterNode)
		if !ok {
			return false
		}
		if !senseRegNode.isValidBurnTxID {
			return false
		}
	}
	return true
}

// SendRegMetadata send metadata
func (h *SenseRegisterHandler) SendRegMetadata(ctx context.Context, regMetadata *types.ActionRegMetadata) error {
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(SenseRegisterNode)
		if !ok {
			return nil
		}
		group.Go(func() (err error) {
			err = senseRegNode.SendRegMetadata(ctx, regMetadata)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

// ProbeImage sends the image to supernodes for image analysis, such as fingerprint, raraness score, NSWF.
func (h *SenseRegisterHandler) ProbeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	h.FingerprintsHandler.Clear()

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(SenseRegisterNode)
		if !ok {
			return nil
		}
		group.Go(func() (err error) {
			compress, isValidBurnTxID, err := senseRegNode.ProbeImage(ctx, file)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			senseRegNode.isValidBurnTxID = isValidBurnTxID

			fingerprintAndScores, fingerprintAndScoresBytes, signature, err := pastel.ExtractCompressSignedDDAndFingerprints(compress)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("extract compressed signed DDAandFingerprints failed")
				return errors.Errorf("node %s: extract failed: %w", someNode.String(), err)
			}
			h.FingerprintsHandler.AddNew(fingerprintAndScores, fingerprintAndScoresBytes, signature)

			return nil
		})
	}
	group.Wait()

	h.FingerprintsHandler.ProcessFingerprints(ctx, h.Nodes)
}

// MatchFingerprintAndScores matches fingerprints.
func (h *SenseRegisterHandler) MatchFingerprintAndScores() error {
	node := h.Nodes[0]
	for i := 1; i < len(h.Nodes); i++ {
		if err := pastel.CompareFingerPrintAndScore(node.FingerprintAndScores, (*nodes)[i].FingerprintAndScores); err != nil {
			return errors.Errorf("node[%s] and node[%s] not matched: %w", node.PastelID(), (*nodes)[i].PastelID(), err)
		}
	}

	return nil
}

// FingerAndScores returns fingerprint of the image and dupedetection scores
func (h *SenseRegisterHandler) FingerAndScores() *pastel.DDAndFingerprints {
	return h.Nodes[0].FingerprintAndScores
}

// UploadSignedTicket uploads regart ticket and its signature to super nodes
func (h *SenseRegisterHandler) UploadSignedTicket(ctx context.Context, ticket []byte, signature []byte, ddFpFile []byte) error {
	group, _ := errgroup.WithContext(ctx)
	for _, node := range h.Nodes {
		node := node
		group.Go(func() error {
			ticketTxid, err := node.SendSignedTicket(ctx, ticket, signature, ddFpFile)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", node).Error("send signed ticket failed")
				return err
			}
			if !node.IsPrimary() && ticketTxid != "" {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, node.PastelID())
			}
			node.regActionTxid = ticketTxid
			return nil
		})
	}
	return group.Wait()
}

// UploadActionAct uploads action act to primary node
func (h *SenseRegisterHandler) UploadActionAct(ctx context.Context) error {
	group, _ := errgroup.WithContext(ctx)

	for _, node := range h.Nodes {
		node := node
		group.Go(func() error {
			if node.IsPrimary() {
				return node.SendActionAct(ctx, node.regActionTxid)
			}

			return nil
		})
	}

	return group.Wait()
}

// CompressedFingerAndScores returns compressed fingerprint and other scores
func (h *SenseRegisterHandler) CompressedFingerAndScores() *pastel.DDAndFingerprints {
	return h.Nodes[0].FingerprintAndScores
}

// RegActionTicketID return txid of reg action ticket
func (h *SenseRegisterHandler) RegActionTicketID() string {
	for i := range h.Nodes {
		if h.Nodes[i].IsPrimary() {
			return h.Nodes[i].regActionTxid
		}
	}
	return string("")
}
