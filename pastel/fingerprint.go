package pastel

import (
	"bytes"
	"encoding/binary"
	"math"
	"strings"
	"unsafe"

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

// CompareFingerPrintAndScore returns nil if two DDAndFingerprints are equal
func CompareFingerPrintAndScore(lhs *DDAndFingerprints, rhs *DDAndFingerprints) error {
	// -------------------WIP: PSL-142------------------------------------------------
	/*//hash_of_candidate_image_file
	if !reflect.DeepEqual(lhs.HashOfCandidateImageFile, rhs.HashOfCandidateImageFile) {
		return errors.New("image hash do not match")
	}*/

	//is_likely_dupe
	if lhs.IsLikelyDupe != rhs.IsLikelyDupe {
		return errors.Errorf("is likely dupe score do not match: lhs(%t) != rhs(%t)", lhs.IsLikelyDupe, rhs.IsLikelyDupe)
	}

	//overall_average_rareness_score
	if !compareFloatWithPrecision(lhs.Score.OverallAverageRarenessScore, rhs.Score.OverallAverageRarenessScore, 7.0) {
		return errors.Errorf("overall average rareness score do not match: lhs(%f) != rhs(%f)", lhs.Score.OverallAverageRarenessScore, rhs.Score.OverallAverageRarenessScore)
	}

	//is_rare_on_internet
	if lhs.IsRareOnInternet != rhs.IsRareOnInternet {
		return errors.Errorf("is rare on internet score do not match: lhs(%t) != rhs(%t)", lhs.IsRareOnInternet, rhs.IsRareOnInternet)
	}

	//open_nsfw_score
	if !compareFloatWithPrecision(lhs.OpenNSFWScore, rhs.OpenNSFWScore, 7.0) {
		return errors.Errorf("open nsfw score do not match: lhs(%f) != rhs(%f)", lhs.OpenNSFWScore, rhs.OpenNSFWScore)
	}

	//alternative_nsfw_scores
	if err := CompareAlternativeNSFWScore(lhs.AlternateNSFWScores, rhs.AlternateNSFWScores); err != nil {
		return errors.Errorf("alternative nsfw score do not match: %w", err)
	}

	//image_hashes
	if err := CompareImageHashes(lhs.PerceptualImageHashes, rhs.PerceptualImageHashes); err != nil {
		return errors.Errorf("image hashes do not match: %w", err)
	}

	// -------------------WIP: PSL-142------------------------------------------------

	/*lhsFingerprint, err := zstd.Decompress(nil, lhs.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("decompress lhs fingerprint: %w", err)
	}

	rhsFingerprint, err := zstd.Decompress(nil, rhs.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("decompress rhs fingerprint: %w", err)
	}

	if !bytes.Equal(lhsFingerprint, rhsFingerprint) {

		lfg, err := FingerprintFromBytes(lhsFingerprint)
		if err != nil {
			return errors.New("fingerprints corrupted")
		}
		rfg, err := FingerprintFromBytes(lhsFingerprint)
		if err != nil {
			return errors.New("fingerprints corrupted")
		}
		if len(lfg) != len(rfg) {
			return errors.New("fingerprints do not match")
		}
		for i := range lfg {
			if !compareFloatWithPrecision(lfg[i], rfg[i], 4.0) {
				return errors.Errorf("fingerprints do not match: lfg(%f) != rfg(%f)", lfg[i], rfg[i])
			}
		}
	}
	*/

	return nil
}

// CompareAlternativeNSFWScore return nil if two AlternativeNSFWScore are equal
func CompareAlternativeNSFWScore(lhs *AlternativeNSFWScore, rhs *AlternativeNSFWScore) error {

	if !compareFloatWithPrecision(lhs.Drawing, rhs.Drawing, 5.0) {
		return errors.Errorf("drawing score not matched: lhs(%f) != rhs(%f)", lhs.Drawing, rhs.Drawing)
	}
	if !compareFloatWithPrecision(lhs.Hentai, rhs.Hentai, 5.0) {
		return errors.Errorf("hentai score not matched: lhs(%f) != rhs(%f)", lhs.Hentai, rhs.Hentai)
	}
	if !compareFloatWithPrecision(lhs.Neutral, rhs.Neutral, 5.0) {
		return errors.Errorf("neutral score not matched: lhs(%f) != rhs(%f)", lhs.Neutral, rhs.Neutral)
	}
	if !compareFloatWithPrecision(lhs.Porn, rhs.Porn, 5.0) {
		return errors.Errorf("porn score not matched: lhs(%f) != rhs(%f)", lhs.Porn, rhs.Porn)
	}
	if !compareFloatWithPrecision(lhs.Sexy, rhs.Sexy, 5.0) {
		return errors.Errorf("sexy score not matched: lhs(%f) != rhs(%f)", lhs.Sexy, rhs.Sexy)
	}
	return nil
}

// CompareImageHashes return nil if two PerceptualImageHashes are equal
func CompareImageHashes(lhs *PerceptualImageHashes, rhs *PerceptualImageHashes) error {
	if !strings.EqualFold(lhs.PDQHash, rhs.PDQHash) {
		return errors.Errorf("pdq_hash not matched: lhs(%s) != rhs(%s)", lhs.PDQHash, rhs.PDQHash)
	}
	if !strings.EqualFold(lhs.PerceptualHash, rhs.PerceptualHash) {
		return errors.Errorf("perceptual_hash not matched: lhs(%s) != rhs(%s)", lhs.PerceptualHash, rhs.PerceptualHash)
	}
	if !strings.EqualFold(lhs.AverageHash, rhs.AverageHash) {
		return errors.Errorf("average_hash not matched: lhs(%s) != rhs(%s)", lhs.AverageHash, rhs.AverageHash)
	}
	if !strings.EqualFold(lhs.DifferenceHash, rhs.DifferenceHash) {
		return errors.Errorf("difference_hash not matched: lhs(%s) != rhs(%s)", lhs.DifferenceHash, rhs.DifferenceHash)
	}
	if !strings.EqualFold(lhs.NeuralHash, rhs.NeuralHash) {
		return errors.Errorf("neuralhash_hash not matched: lhs(%s) != rhs(%s)", lhs.NeuralHash, rhs.NeuralHash)
	}
	return nil
}

func compareFloatWithPrecision(l float32, r float32, prec float64) bool {
	if l == r {
		return true
	}
	multiplier := math.Pow(10, prec)
	return math.Round(float64(l)*multiplier)/multiplier == math.Round(float64(l)*multiplier)/multiplier
}
