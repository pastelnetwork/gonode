package mixins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"math/rand"
)

type Fingerprints struct {
	FingerprintAndScores      *pastel.DDAndFingerprints
	FingerprintAndScoresBytes []byte // JSON bytes of FingerprintAndScores
	Signature                 []byte
}

type FingerprintsHandler struct {
	received []*Fingerprints
	common   Fingerprints
}

func (h *FingerprintsHandler) Add(fingerprints *Fingerprints) {
	h.received = append(h.received, fingerprints)
}

func (h *FingerprintsHandler) AddNew(ddf *pastel.DDAndFingerprints, bytes []byte, signature []byte) {
	fingerprints := &Fingerprints{ddf, bytes, signature}
	h.received = append(h.received, fingerprints)
}

func (h *FingerprintsHandler) Clear() {
	h.received = nil
}

// GenerateDDAndFingerprintsIDs generates redundant IDs and assigns to task.redundantIDs
func (h *FingerprintsHandler) GenerateDDAndFingerprintsIDs(ctx context.Context, signatures [][]byte, ddAndFingerprintsIDs *[]string, ddAndFpFile *[]byte, max uint32) (uint32, error) {
	if len(signatures) != 3 {
		return 0, errors.Errorf("wrong number of signature for fingerprints - %d", len(signatures))
	}

	ddDataJSON, err := json.Marshal(h.common.FingerprintAndScores)
	if err != nil {
		return 0, errors.Errorf("failed to marshal dd-data: %w", err)
	}

	ddEncoded := utils.B64Encode(ddDataJSON)

	var buffer bytes.Buffer
	buffer.Write(ddEncoded)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(signatures[0])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(signatures[1])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(signatures[2])
	ddFpFile := buffer.Bytes()

	ddAndFingerprintsIc := rand.Uint32()
	ids, _, err := pastel.GetIDFiles(ddFpFile, ddAndFingerprintsIc, max)
	if err != nil {
		return 0, fmt.Errorf("get ID Files: %w", err)
	}

	comp, err := zstd.CompressLevel(nil, ddFpFile, 22)
	if err != nil {
		return 0, errors.Errorf("compress: %w", err)
	}
	ddFile := utils.B64Encode(comp)

	ddAndFingerprintsIDs = &ids
	ddAndFpFile = &ddFile

	return ddAndFingerprintsIc, nil
}
