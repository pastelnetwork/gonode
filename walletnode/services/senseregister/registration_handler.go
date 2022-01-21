package senseregister

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
)

type SenseRegisterHandler struct {
	meshHandler *mixins.MeshHandler
	fpHandler   *mixins.FingerprintsHandler
}

func NewSenseRegisterHandler(meshHandler *mixins.MeshHandler, fpHandler *mixins.FingerprintsHandler) *SenseRegisterHandler {
	return &SenseRegisterHandler{
		meshHandler: meshHandler,
		fpHandler:   fpHandler,
	}
}

// SendRegMetadata send metadata
func (h *SenseRegisterHandler) SendRegMetadata(ctx context.Context, regMetadata *types.ActionRegMetadata) error {
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.meshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(SenseRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegisterNode", someNode.String())
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

// ProbeImage sends the image to supernodes for image analysis, such as fingerprint, rareness score, NSWF.
func (h *SenseRegisterHandler) ProbeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	h.fpHandler.Clear()

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.meshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(SenseRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegisterNode", someNode.String())
		}
		group.Go(func() (err error) {
			compress, stateOk, err := senseRegNode.ProbeImage(ctx, file)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			someNode.SetRemoteState(stateOk)
			if !stateOk {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("probe image failed")
				return errors.Errorf("remote node %s: indicated processing error", someNode.String())
			}

			fingerprintAndScores, fingerprintAndScoresBytes, signature, err := pastel.ExtractCompressSignedDDAndFingerprints(compress)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", someNode).Error("extract compressed signed DDAandFingerprints failed")
				return errors.Errorf("node %s: extract failed: %w", someNode.String(), err)
			}
			h.fpHandler.AddNew(fingerprintAndScores, fingerprintAndScoresBytes, signature, someNode.PastelID())

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}

	if err := h.fpHandler.Match(ctx); err != nil {
		log.WithContext(ctx).WithError(err).WithField("filename", fileName).Error("probe image failed")
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}

	return nil
}

// UploadSignedTicket uploads sense ticket  and its signature to super nodes
func (h *SenseRegisterHandler) UploadSignedTicket(ctx context.Context, ticket []byte, signature []byte) (string, error) {
	var regSenseTxid string
	ddFpFile := h.fpHandler.DDAndFpFile
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.meshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(SenseRegisterNode)
		if !ok {
			//TODO: use assert here
			return "", errors.Errorf("node %s is not SenseRegisterNode", someNode.String())
		}
		group.Go(func() error {
			ticketTxid, err := senseRegNode.SendSignedTicket(ctx, ticket, signature, ddFpFile)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("send signed ticket failed")
				return err
			}
			if !someNode.IsPrimary() && ticketTxid != "" {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, someNode.PastelID())
			}
			if someNode.IsPrimary() {
				if ticketTxid == "" {
					return errors.Errorf("primary node - %s, returned empty txid", someNode.PastelID())
				}
				regSenseTxid = ticketTxid
			}
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return "", errors.Errorf("uploading ticket has failed: %w", err)
	}

	return regSenseTxid, nil
}

// UploadActionAct uploads action act to primary node
func (h *SenseRegisterHandler) UploadActionAct(ctx context.Context, activateTxID string) error {
	group, _ := errgroup.WithContext(ctx)

	for _, someNode := range h.meshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(SenseRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegisterNode", someNode.String())
		}
		if someNode.IsPrimary() {
			group.Go(func() error {
				return senseRegNode.SendActionAct(ctx, activateTxID)
			})
		}
	}
	return group.Wait()
}
