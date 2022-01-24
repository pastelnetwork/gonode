package mixins

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

type hashes struct {
	previewHash         []byte
	mediumThumbnailHash []byte
	smallThumbnailHash  []byte
	pastelID            string
}

type NftImageHandler struct {
	task          *common.WalletNodeTask
	pastelHandler *PastelHandler

	ImageEncodedWithFingerprints *files.File

	received            []*hashes
	PreviewHash         []byte
	MediumThumbnailHash []byte
	SmallThumbnailHash  []byte
}

// AddNew adds fingerprints info to 'received' array
func (h *NftImageHandler) AddHashes(hashes *hashes) {
	h.received = append(h.received, hashes)
}

// AddNew creates and adds fingerprints info to 'received' array
func (h *NftImageHandler) AddNewHashes(preview []byte, mediumThumbnail []byte, smallThumbnail []byte, pastelID string) {
	hashes := &hashes{preview, mediumThumbnail, smallThumbnail, pastelID}
	h.received = append(h.received, hashes)
}

// ClearHashes clears stored SNs fingerprints info
func (h *NftImageHandler) ClearHashes() {
	h.received = nil
}

// EncodeFingerprintIntoImage - signs fingerprints and encode them into image
func (h *NftImageHandler) CreateCopyWithEncodedFingerprint(ctx context.Context,
	creatorPatelID string, creatorPassphrase string,
	fingerprints Fingerprints,
	image *files.File,
) error {
	img, err := image.Copy()
	if err != nil {
		return errors.Errorf("copy image to encode: %w", err)
	}
	log.WithContext(ctx).WithField("FileName", img.Name()).Debug("image with embedded fingerprints")

	ticket, err := h.pastelHandler.PastelClient.FindTicketByID(ctx, creatorPatelID)
	if err != nil {
		return errors.Errorf("find register ticket of artist pastel id(%s):%w", creatorPatelID, err)
	}

	ed448PubKey, err := common.GetPubKey(creatorPatelID)
	if err != nil {
		return fmt.Errorf("encodeFingerprint: %v", err)
	}

	pqPubKey, err := common.GetPubKey(ticket.PqKey)
	if err != nil {
		return fmt.Errorf("encodeFingerprint: %v", err)
	}

	// Sign fingerprint
	ed448Signature, err := h.pastelHandler.PastelClient.Sign(ctx,
		fingerprints.fingerprintAndScoresBytes,
		creatorPatelID,
		creatorPassphrase,
		pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign fingerprint: %w", err)
	}

	pqSignature, err := h.pastelHandler.PastelClient.Sign(ctx,
		fingerprints.fingerprintAndScoresBytes,
		creatorPatelID,
		creatorPassphrase,
		pastel.SignAlgorithmLegRoast)
	if err != nil {
		return errors.Errorf("sign fingerprint with legroats: %w", err)
	}

	// Encode data to the image.
	encSig := qrsignature.New(
		qrsignature.Fingerprint(fingerprints.fingerprintAndScoresBytes),
		qrsignature.PostQuantumSignature(pqSignature),
		qrsignature.PostQuantumPubKey(pqPubKey),
		qrsignature.Ed448Signature(ed448Signature),
		qrsignature.Ed448PubKey(ed448PubKey),
	)
	if err := img.Encode(encSig); err != nil {
		return errors.Errorf("encode fingerprint into image: %w", err)
	}

	h.ImageEncodedWithFingerprints = img

	return nil
}

// MatchThumbnailHashes matches thumbnail's hashes recevied from super nodes
func (h *NftImageHandler) MatchThumbnailHashes() error {
	if len(h.received) != 3 {
		return errors.Errorf("wrong number of hashes - %d", len(h.received))
	}

	first := h.received[0]
	for _, some := range h.received[1:] {
		if !bytes.Equal(first.previewHash, some.previewHash) {
			return errors.Errorf("hash of preview thumbnail of nodes %q and %q didn't match", first.pastelID, some.pastelID)
		}
		if !bytes.Equal(first.mediumThumbnailHash, some.mediumThumbnailHash) {
			return errors.Errorf("hash of medium thumbnail of nodes %q and %q didn't match", first.pastelID, some.pastelID)
		}
		if !bytes.Equal(first.smallThumbnailHash, some.smallThumbnailHash) {
			return errors.Errorf("hash of small thumbnail of nodes %q and %q didn't match", first.pastelID, some.pastelID)
		}
	}
	h.PreviewHash = first.previewHash
	h.MediumThumbnailHash = first.mediumThumbnailHash
	h.SmallThumbnailHash = first.smallThumbnailHash
	return nil
}

// GetHash returns hash of the image
func (h *NftImageHandler) GetHash() ([]byte, error) {
	imgBytes, err := h.ImageEncodedWithFingerprints.Bytes()
	if err != nil {
		return nil, errors.Errorf("convert image to byte stream %w", err)
	}
	dataHash, err := utils.Sha3256hash(imgBytes)
	if err != nil {
		return nil, errors.Errorf("hash encoded image: %w", err)
	}
	return dataHash, nil
}

// NewImageHandler create new Image Handler
func NewImageHandler(task *common.WalletNodeTask, pastelHandler *PastelHandler) *NftImageHandler {
	return &NftImageHandler{task: task, pastelHandler: pastelHandler}
}

func (h NftImageHandler) IsEmpty() bool {
	return h.PreviewHash == nil || len(h.PreviewHash) == 0 ||
		h.MediumThumbnailHash == nil || len(h.MediumThumbnailHash) == 0 ||
		h.SmallThumbnailHash == nil || len(h.SmallThumbnailHash) == 0
}
