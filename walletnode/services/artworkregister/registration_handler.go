package artworkregister

import (
	"context"
	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
)

type NftRegisterHandler struct {
	meshHandler *mixins.MeshHandler
	fpHandler   *mixins.FingerprintsHandler
	rqHandler   *mixins.RQHandler
}

func NewNftRegisterHandler(meshHandler *mixins.MeshHandler, fpHandler *mixins.FingerprintsHandler) *NftRegisterHandler {
	return &NftRegisterHandler{
		meshHandler: meshHandler,
		fpHandler:   fpHandler,
	}
}

// SendRegMetadata send metadata
func (h *NftRegisterHandler) SendRegMetadata(ctx context.Context, regMetadata *types.NftRegMetadata) error {
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.meshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() (err error) {
			err = nftRegNode.SendRegMetadata(ctx, regMetadata)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

func (h *NftRegisterHandler) ProbeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	h.fpHandler.Clear()

	// Send image to supernodes for probing.
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.meshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() (err error) {
			compress, stateOk, err := nftRegNode.ProbeImage(ctx, file)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			someNode.SetRemoteState(stateOk)
			if !stateOk {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("probe image failed")
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

// UploadSignedTicket uploads regart ticket and its signature to super nodes
func (h *NftRegisterHandler) UploadSignedTicket(ctx context.Context, ticket []byte, signature []byte /*, rqidsFile []byte, encoderParams rqnode.EncoderParameters*/) error {
	ddFpFile := h.fpHandler.DDAndFpFile
	rqidsFile := h.rqHandler.RQIdsFile
	encoderParams := h.rqHandler.EncoderParams

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range h.meshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		// TODO: Fix this when method to generate key1 and key2 are finalized
		key1 := "key1-" + uuid.New().String()
		key2 := "key2-" + uuid.New().String()
		group.Go(func() error {
			fee, err := nftRegNode.SendSignedTicket(ctx, ticket, signature, key1, key2, rqidsFile, ddFpFile, encoderParams)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("send signed ticket failed")
				return err
			}
			nftRegNode.registrationFee = fee
			return nil
		})
	}
	return group.Wait()
}
