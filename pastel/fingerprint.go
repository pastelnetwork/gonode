package pastel

import (
	"bytes"
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
)

// Fingerprint uniquely identify an image
type Fingerprint []float32

func fingerprintBaseTypeSize() int {
	return int(unsafe.Sizeof(float32(0)))
}

// FingerprintFromBytes deserialize a slice of bytes into Fingerprint
func FingerprintFromBytes(data []byte) (Fingerprint, error) {
	typeSize := fingerprintBaseTypeSize()
	if len(data)%typeSize != 0 {
		return nil, errors.Errorf("invalid data length %d, length should be multiple of sizeof(float64)", len(data))
	}
	fg := make([]float32, len(data)/typeSize)

	for i := range fg {
		bits := binary.LittleEndian.Uint32(data[i*typeSize : (i+1)*typeSize])
		fg[i] = math.Float32frombits(bits)
	}
	return fg, nil
}

// Bytes serialize a Fingerprint into slice of byte
func (fg Fingerprint) Bytes() []byte {
	output := new(bytes.Buffer)
	_ = binary.Write(output, binary.LittleEndian, fg)
	return output.Bytes()
}

// CompareFingerPrintAndScore returns nil if two FingerAndScores are equal
func CompareFingerPrintAndScore(lhs *FingerAndScores, rhs *FingerAndScores) error {
	lhsFingerprint, err := zstd.Decompress(nil, lhs.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("failed to decompress lhs fingerprint: %w", err)
	}

	rhsFingerpint, err := zstd.Decompress(nil, rhs.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("failed to decompress rhs fingerprint: %w", err)
	}

	if !bytes.Equal(lhsFingerprint, rhsFingerpint) {
		return errors.Errorf("fingerprint not matched")
	}

	if lhs.OverallAverageRarenessScore != rhs.OverallAverageRarenessScore {
		return errors.Errorf("overall average rareness score not matched: lhs(%f) != rhs(%f)", lhs.OverallAverageRarenessScore, rhs.OverallAverageRarenessScore)
	}

	if lhs.IsLikelyDupe != rhs.IsLikelyDupe {
		return errors.Errorf("is likely dupe score not matched: lhs(%t) != rhs(%t)", lhs.IsLikelyDupe, rhs.IsLikelyDupe)
	}

	if lhs.IsRareOnInternet != rhs.IsRareOnInternet {
		return errors.Errorf("is rare on internet score not matched: lhs(%t) != rhs(%t)", lhs.IsRareOnInternet, rhs.IsRareOnInternet)
	}

	if lhs.OpenNSFWScore != rhs.OpenNSFWScore {
		return errors.Errorf("open nsfw score not matched: lhs(%f) != rhs(%f)", lhs.OpenNSFWScore, rhs.OpenNSFWScore)
	}

	if err := CompareAlternativeNSFWScore(&lhs.AlternativeNSFWScore, &rhs.AlternativeNSFWScore); err != nil {
		return errors.Errorf("alternative nsfw score not matched: %w", err)
	}
	return nil
}

// CompareAlternativeNSFWScore return nil if two AlternativeNSFWScore are equal
func CompareAlternativeNSFWScore(lhs *AlternativeNSFWScore, rhs *AlternativeNSFWScore) error {
	if lhs.Drawing != rhs.Drawing {
		return errors.Errorf("drawing score not matched: lhs(%f) != rhs(%f)", lhs.Drawing, rhs.Drawing)
	}
	if lhs.Hentai != rhs.Hentai {
		return errors.Errorf("hentai score not matched: lhs(%f) != rhs(%f)", lhs.Hentai, rhs.Hentai)
	}
	if lhs.Neutral != rhs.Neutral {
		return errors.Errorf("neutral score not matched: lhs(%f) != rhs(%f)", lhs.Neutral, rhs.Neutral)
	}
	if lhs.Porn != rhs.Porn {
		return errors.Errorf("porn score not matched: lhs(%f) != rhs(%f)", lhs.Porn, rhs.Porn)
	}
	if lhs.Sexy != rhs.Sexy {
		return errors.Errorf("sexy score not matched: lhs(%f) != rhs(%f)", lhs.Sexy, rhs.Sexy)
	}
	return nil
}
