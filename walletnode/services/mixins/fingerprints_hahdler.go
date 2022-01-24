package mixins

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"math/rand"
)

type Fingerprints struct {
	fingerprintAndScores      *pastel.DDAndFingerprints
	fingerprintAndScoresBytes []byte // JSON bytes of FingerprintAndScores
	signature                 []byte
	pastelID                  string
}

type FingerprintsHandler struct {
	task          *common.WalletNodeTask
	pastelHandler *PastelHandler

	received []*Fingerprints

	//final combined data
	FinalFingerprints    Fingerprints
	DDAndFingerprintsIDs []string
	DDAndFingerprintsIc  uint32
	DDAndFpFile          []byte
	SNsSignatures        [][]byte
}

// NewFingerprintsHandler create new Fingerprints Handler
func NewFingerprintsHandler(task *common.WalletNodeTask, pastelHandler *PastelHandler) *FingerprintsHandler {
	return &FingerprintsHandler{task: task, pastelHandler: pastelHandler}
}

// AddNew adds fingerprints info to 'received' array
func (h *FingerprintsHandler) Add(fingerprints *Fingerprints) {
	h.received = append(h.received, fingerprints)
}

// AddNew creates and adds fingerprints info to 'received' array
func (h *FingerprintsHandler) AddNew(ddf *pastel.DDAndFingerprints, bytes []byte, signature []byte, pastelid string) {
	fingerprints := &Fingerprints{ddf, bytes, signature, pastelid}
	h.received = append(h.received, fingerprints)
}

// Clear clears stored SNs fingerprints info
func (h *FingerprintsHandler) Clear() {
	h.received = nil
	h.SNsSignatures = nil
}

// Match inter checks fingerprints data from 3 SNs
func (h *FingerprintsHandler) Match(ctx context.Context) error {
	if len(h.received) != 3 {
		return errors.Errorf("wrong number of fingerprints - %d", len(h.received))
	}

	// Match signatures received from supernodes.
	for _, someNode := range h.received {
		// Validate signatures received from supernodes.
		verified, err := h.pastelHandler.VerifySignature(ctx,
			someNode.fingerprintAndScoresBytes,
			string(someNode.signature),
			someNode.pastelID,
			pastel.SignAlgorithmED448)
		if err != nil {
			return errors.Errorf("probeImage: pastelClient.Verify %w", err)
		}

		if !verified {
			h.task.UpdateStatus(common.StatusErrorSignaturesNotMatch)
			return errors.Errorf("node[%s] signature doesn't match", someNode.pastelID)
		}

		h.SNsSignatures = append(h.SNsSignatures, someNode.signature)
	}

	// Match fingerprints received from supernodes.
	if err := h.matchFingerprintAndScores(); err != nil {
		h.task.UpdateStatus(common.StatusErrorFingerprintsNotMatch)
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}

	h.FinalFingerprints.fingerprintAndScores = h.received[0].fingerprintAndScores
	h.task.UpdateStatus(common.StatusImageProbed)

	return nil
}

// MatchFingerprintAndScores matches fingerprints.
func (h *FingerprintsHandler) matchFingerprintAndScores() error {
	for _, someNode := range h.received[1:] {
		if err := pastel.CompareFingerPrintAndScore(h.received[0].fingerprintAndScores, someNode.fingerprintAndScores); err != nil {
			return errors.Errorf("node[%s] and node[%s] not matched: %w", h.received[0].pastelID, someNode.pastelID, err)
		}
	}
	return nil
}

// GenerateDDAndFingerprintsIDs generates redundant IDs and assigns to task.redundantIDs
func (h *FingerprintsHandler) GenerateDDAndFingerprintsIDs(_ context.Context, max uint32) error {
	if len(h.SNsSignatures) != 3 {
		return errors.Errorf("wrong number of signature for fingerprints - %d", len(h.SNsSignatures))
	}

	ddDataJSON, err := json.Marshal(h.FinalFingerprints.fingerprintAndScores)
	if err != nil {
		return errors.Errorf("failed to marshal dd-data: %w", err)
	}

	ddEncoded := utils.B64Encode(ddDataJSON)

	var buffer bytes.Buffer
	buffer.Write(ddEncoded)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(h.SNsSignatures[0])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(h.SNsSignatures[1])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(h.SNsSignatures[2])
	ddFpFile := buffer.Bytes()

	ddAndFingerprintsIc := rand.Uint32()
	ids, _, err := pastel.GetIDFiles(ddFpFile, ddAndFingerprintsIc, max)
	if err != nil {
		return errors.Errorf("get ID Files: %w", err)
	}

	comp, err := zstd.CompressLevel(nil, ddFpFile, 22)
	if err != nil {
		return errors.Errorf("compress: %w", err)
	}
	ddFile := utils.B64Encode(comp)

	h.DDAndFingerprintsIDs = ids
	h.DDAndFpFile = ddFile
	h.DDAndFingerprintsIc = ddAndFingerprintsIc

	return nil
}

func (h FingerprintsHandler) IsEmpty() bool {
	return h.DDAndFingerprintsIDs == nil || len(h.DDAndFingerprintsIDs) == 0 ||
		h.DDAndFpFile == nil || len(h.DDAndFpFile) == 0 ||
		h.SNsSignatures == nil || len(h.SNsSignatures) == 0
}
