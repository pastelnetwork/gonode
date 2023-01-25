package pastel

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
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
	if lhs == nil || rhs == nil {
		return errors.New("nil input")
	}

	// BlockHash
	if !strings.EqualFold(lhs.BlockHash, rhs.BlockHash) {
		return errors.Errorf("block hash not matched: lhs(%s) != rhs(%s)", lhs.BlockHash, rhs.BlockHash)
	}

	// BlockHeight
	if !strings.EqualFold(lhs.BlockHeight, rhs.BlockHeight) {
		return errors.Errorf("block height not matched: lhs(%s) != rhs(%s)", lhs.BlockHeight, rhs.BlockHeight)
	}

	// Timestamp
	if !strings.EqualFold(lhs.TimestampOfRequest, rhs.TimestampOfRequest) {
		return errors.Errorf("timestamp not matched: lhs(%s) != rhs(%s)", lhs.TimestampOfRequest, rhs.TimestampOfRequest)
	}

	// SubmitterPastelID
	if !strings.EqualFold(lhs.SubmitterPastelID, rhs.SubmitterPastelID) {
		return errors.Errorf("SubmitterPastelID not matched: lhs(%s) != rhs(%s)", lhs.SubmitterPastelID, rhs.SubmitterPastelID)
	}

	// SN1PastelID
	if !strings.EqualFold(lhs.SN1PastelID, rhs.SN1PastelID) {
		return errors.Errorf("SN1PastelID not matched: lhs(%s) != rhs(%s)", lhs.SN1PastelID, rhs.SN1PastelID)
	}
	// SN2PastelID
	if !strings.EqualFold(lhs.SN2PastelID, rhs.SN2PastelID) {
		return errors.Errorf("SN2PastelID not matched: lhs(%s) != rhs(%s)", lhs.SN2PastelID, rhs.SN2PastelID)
	}
	// SN3PastelID
	if !strings.EqualFold(lhs.SN3PastelID, rhs.SN3PastelID) {
		return errors.Errorf("SN3PastelID not matched: lhs(%s) != rhs(%s)", lhs.SN3PastelID, rhs.SN3PastelID)
	}

	// IsOpenAPIRequest
	if lhs.IsOpenAPIRequest != rhs.IsOpenAPIRequest {
		return errors.Errorf("IsOpenAPIRequest not matched: lhs(%t) != rhs(%t)", lhs.IsOpenAPIRequest, rhs.IsOpenAPIRequest)
	}

	// OpenAPISubsetID
	if !strings.EqualFold(lhs.OpenAPISubsetID, rhs.OpenAPISubsetID) {
		return errors.Errorf("OpenAPISubsetID not matched: lhs(%s) != rhs(%s)", lhs.OpenAPISubsetID, rhs.OpenAPISubsetID)
	}

	// DupeDetectionSystemVersion
	if !compareDupeDetectionSystemVersionNumber(lhs.DupeDetectionSystemVersion, rhs.DupeDetectionSystemVersion) {
		return errors.Errorf("dupe_detection_system_version not matched: lhs(%s) != rhs(%s)", lhs.DupeDetectionSystemVersion, rhs.DupeDetectionSystemVersion)
	}

	// IsLikelyDupe
	if lhs.IsLikelyDupe != rhs.IsLikelyDupe {
		return errors.Errorf("is_likely_dupe not matched: lhs(%t) != rhs(%t)", lhs.IsLikelyDupe, rhs.IsLikelyDupe)
	}

	// IsRareOnInternet
	if lhs.IsRareOnInternet != rhs.IsRareOnInternet {
		return errors.Errorf("is_rare_on_internet not matched: lhs(%t) != rhs(%t)", lhs.IsRareOnInternet, rhs.IsRareOnInternet)
	}

	// OverallRarenessScore
	if !compareFloat(lhs.OverallRarenessScore, rhs.OverallRarenessScore) {
		return errors.Errorf("OverallRarenessScore do not match: lhs(%f) != rhs(%f)", lhs.OverallRarenessScore, rhs.OverallRarenessScore)
	}

	// PctOfTop10MostSimilarWithDupeProbAbove25pct
	if !compareFloat(lhs.PctOfTop10MostSimilarWithDupeProbAbove25pct, rhs.PctOfTop10MostSimilarWithDupeProbAbove25pct) {
		return errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove25pct do not match: lhs(%f) != rhs(%f)", lhs.PctOfTop10MostSimilarWithDupeProbAbove25pct, rhs.PctOfTop10MostSimilarWithDupeProbAbove25pct)
	}
	// PctOfTop10MostSimilarWithDupeProbAbove33pct
	if !compareFloat(lhs.PctOfTop10MostSimilarWithDupeProbAbove33pct, rhs.PctOfTop10MostSimilarWithDupeProbAbove33pct) {
		return errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove33pct do not match: lhs(%f) != rhs(%f)", lhs.PctOfTop10MostSimilarWithDupeProbAbove33pct, rhs.PctOfTop10MostSimilarWithDupeProbAbove33pct)
	}
	// PctOfTop10MostSimilarWithDupeProbAbove50pct
	if !compareFloat(lhs.PctOfTop10MostSimilarWithDupeProbAbove50pct, rhs.PctOfTop10MostSimilarWithDupeProbAbove50pct) {
		return errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove50pct do not match: lhs(%f) != rhs(%f)", lhs.PctOfTop10MostSimilarWithDupeProbAbove50pct, rhs.PctOfTop10MostSimilarWithDupeProbAbove50pct)
	}

	// RarenessScoresTableJSONCompressedB64
	if !strings.EqualFold(lhs.RarenessScoresTableJSONCompressedB64, rhs.RarenessScoresTableJSONCompressedB64) {
		return errors.Errorf("RarenessScoresTableJSONCompressedB64 not matched: lhs(%s) != rhs(%s)", lhs.RarenessScoresTableJSONCompressedB64, rhs.RarenessScoresTableJSONCompressedB64)
	}

	//InternetRareness
	if err := CompareInternetRareness(lhs.InternetRareness, rhs.InternetRareness); err != nil {
		return err
	}

	//open_nsfw_score
	if !compareFloat(lhs.OpenNSFWScore, rhs.OpenNSFWScore) {
		return errors.Errorf("open_nsfw_score do not match: lhs(%f) != rhs(%f)", lhs.OpenNSFWScore, rhs.OpenNSFWScore)
	}

	//alternative_nsfw_scores
	if err := CompareAlternativeNSFWScore(lhs.AlternativeNSFWScores, rhs.AlternativeNSFWScores); err != nil {
		return errors.Errorf("alternative nsfw score do not match: %w", err)
	}

	// ImageFingerprintOfCandidateImageFile
	if !compareFloats(lhs.ImageFingerprintOfCandidateImageFile, rhs.ImageFingerprintOfCandidateImageFile) {
		return errors.New("image_fingerprint_of_candidate_image_file nsfw score do not match")
	}

	//hash_of_candidate_image_file
	if !reflect.DeepEqual(lhs.HashOfCandidateImageFile, rhs.HashOfCandidateImageFile) {
		return errors.New("image hash do not match")
	}

	return nil
}

// CompareInternetRareness return nil if two InternetRareness are equal
func CompareInternetRareness(lhs *InternetRareness, rhs *InternetRareness) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil internet_rareness")
	}

	if !strings.EqualFold(lhs.RareOnInternetSummaryTableAsJSONCompressedB64, rhs.RareOnInternetSummaryTableAsJSONCompressedB64) {
		return errors.Errorf("RareOnInternetSummaryTableAsJSONCompressedB64 not matched: lhs(%s) != rhs(%s)", lhs.RareOnInternetSummaryTableAsJSONCompressedB64, rhs.RareOnInternetSummaryTableAsJSONCompressedB64)
	}

	if !strings.EqualFold(lhs.RareOnInternetGraphJSONCompressedB64, rhs.RareOnInternetGraphJSONCompressedB64) {
		return errors.Errorf("RareOnInternetGraphJSONCompressedB64 not matched: lhs(%s) != rhs(%s)", lhs.RareOnInternetGraphJSONCompressedB64, rhs.RareOnInternetGraphJSONCompressedB64)
	}

	if !strings.EqualFold(lhs.AlternativeRareOnInternetDictAsJSONCompressedB64, rhs.AlternativeRareOnInternetDictAsJSONCompressedB64) {
		return errors.Errorf("AlternativeRareOnInternetDictAsJSONCompressedB64 not matched: lhs(%s) != rhs(%s)", lhs.AlternativeRareOnInternetDictAsJSONCompressedB64, rhs.AlternativeRareOnInternetDictAsJSONCompressedB64)
	}

	if lhs.MinNumberOfExactMatchesInPage != rhs.MinNumberOfExactMatchesInPage {
		return errors.Errorf("MinNumberOfExactMatchesInPage not matched: lhs(%d) != rhs(%d)", lhs.MinNumberOfExactMatchesInPage, rhs.MinNumberOfExactMatchesInPage)
	}

	if !strings.EqualFold(lhs.EarliestAvailableDateOfInternetResults, rhs.EarliestAvailableDateOfInternetResults) {
		return errors.Errorf("EarliestAvailableDateOfInternetResults not matched: lhs(%s) != rhs(%s)", lhs.EarliestAvailableDateOfInternetResults, rhs.EarliestAvailableDateOfInternetResults)
	}
	return nil
}

// CompareAlternativeNSFWScore return nil if two AlternativeNSFWScore are equal
func CompareAlternativeNSFWScore(lhs *AlternativeNSFWScores, rhs *AlternativeNSFWScores) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil alternative_nsfw_scores")
	}

	if !compareFloat(lhs.Drawings, rhs.Drawings) {
		return errors.Errorf("drawings score not matched: lhs(%f) != rhs(%f)", lhs.Drawings, rhs.Drawings)
	}
	if !compareFloat(lhs.Hentai, rhs.Hentai) {
		return errors.Errorf("hentai score not matched: lhs(%f) != rhs(%f)", lhs.Hentai, rhs.Hentai)
	}
	if !compareFloat(lhs.Neutral, rhs.Neutral) {
		return errors.Errorf("neutral score not matched: lhs(%f) != rhs(%f)", lhs.Neutral, rhs.Neutral)
	}
	if !compareFloat(lhs.Porn, rhs.Porn) {
		return errors.Errorf("porn score not matched: lhs(%f) != rhs(%f)", lhs.Porn, rhs.Porn)
	}
	if !compareFloat(lhs.Sexy, rhs.Sexy) {
		return errors.Errorf("sexy score not matched: lhs(%f) != rhs(%f)", lhs.Sexy, rhs.Sexy)
	}
	return nil
}

func compareDupeDetectionSystemVersionNumber(l, r string) bool {
	lhs := strings.Split(l, ".")
	rhs := strings.Split(r, ".")

	return lhs[0] == rhs[0]
}

func compareFloat(l float32, r float32) bool {
	if compareFloatWithPrecision(l, r, 5.0) {
		return true
	}

	return compareFloatWithPrecision(l, r, 3.0)
}

func compareFloatWithPrecision(l float32, r float32, prec float64) bool {
	if l == r {
		return true
	}
	multiplier := math.Pow(10, prec)
	return math.Round(float64(l)*multiplier)/multiplier == math.Round(float64(r)*multiplier)/multiplier
}

func compareFloats(l []float32, r []float32) bool {
	if len(l) != len(r) {
		return false
	}

	for i := 0; i < len(l); i = i + 1 {
		if !compareFloat(l[i], r[i]) {
			return false
		}
	}

	return true
}
