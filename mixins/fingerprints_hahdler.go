package mixins

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"math/rand"
)

type Fingerprints struct {
	fingerprintAndScores      *pastel.DDAndFingerprints
	FingerprintAndScoresBytes []byte // JSON bytes of fingerprintAndScores
	signature                 []byte
	pastelID                  string
}

type FingerprintsHandler struct {
	pastelHandler *PastelHandler

	// array of fingerprints usually from different sources for analysis
	fingerprints []*Fingerprints

	//final combined data
	FinalFingerprints    Fingerprints
	DDAndFingerprintsIDs []string
	DDAndFingerprintsIc  uint32
	DDAndFpFile          []byte
	SNsSignatures        [][]byte
}

// NewFingerprintsHandler create new Fingerprints Handler
func NewFingerprintsHandler(pastelHandler *PastelHandler) *FingerprintsHandler {
	return &FingerprintsHandler{pastelHandler: pastelHandler}
}

// AddNew adds fingerprints info to 'fingerprints' array
func (h *FingerprintsHandler) Add(fingerprints *Fingerprints) {
	h.fingerprints = append(h.fingerprints, fingerprints)
}

// AddNew creates and adds fingerprints info to 'fingerprints' array
func (h *FingerprintsHandler) AddNew(ddf *pastel.DDAndFingerprints, bytes []byte, signature []byte, pastelid string) {
	fingerprints := &Fingerprints{ddf, bytes, signature, pastelid}
	h.fingerprints = append(h.fingerprints, fingerprints)
}

// Clear clears stored SNs fingerprints info
func (h *FingerprintsHandler) Clear() {
	h.fingerprints = nil
	h.SNsSignatures = nil
}

// Match inter checks fingerprints data from 3 SNs
func (h *FingerprintsHandler) Match(ctx context.Context) error {
	if len(h.fingerprints) != 3 {
		return errors.Errorf("wrong number of fingerprints - %d", len(h.fingerprints))
	}

	// Match signatures received from supernodes.
	for _, someNode := range h.fingerprints {
		// Validate signatures received from supernodes.
		verified, err := h.pastelHandler.VerifySignature(ctx,
			someNode.FingerprintAndScoresBytes,
			string(someNode.signature),
			someNode.pastelID,
			pastel.SignAlgorithmED448)
		if err != nil {
			return errors.Errorf("probeImage: pastelClient.Verify %w", err)
		}

		if !verified {
			return errors.Errorf("node[%s] signature doesn't match", someNode.pastelID)
		}

		h.SNsSignatures = append(h.SNsSignatures, someNode.signature)
	}

	// Match fingerprints received from supernodes.
	if err := h.matchFingerprintAndScores(); err != nil {
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}

	h.FinalFingerprints.fingerprintAndScores = h.fingerprints[0].fingerprintAndScores

	return nil
}

// MatchFingerprintAndScores matches fingerprints.
func (h *FingerprintsHandler) matchFingerprintAndScores() error {
	for _, someNode := range h.fingerprints[1:] {
		if err := pastel.CompareFingerPrintAndScore(h.fingerprints[0].fingerprintAndScores, someNode.fingerprintAndScores); err != nil {
			return errors.Errorf("node[%s] and node[%s] not matched: %w", h.fingerprints[0].pastelID, someNode.pastelID, err)
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
