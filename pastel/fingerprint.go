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

	// Block
	if !strings.EqualFold(lhs.Block, rhs.Block) {
		return errors.Errorf("block not matched: lhs(%s) != rhs(%s)", lhs.Block, rhs.Block)
	}

	// Principal
	if !strings.EqualFold(lhs.Principal, rhs.Principal) {
		return errors.Errorf("principal not matched: lhs(%s) != rhs(%s)", lhs.Principal, rhs.Principal)
	}

	// DupeDetectionSystemVersion
	if !strings.EqualFold(lhs.DupeDetectionSystemVersion, rhs.DupeDetectionSystemVersion) {
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

	//RarenessScores
	if err := CompareRarenessScores(lhs.RarenessScores, rhs.RarenessScores); err != nil {
		return err
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

	//alternative_nsfw_scores
	if err := CompareFingerprintsStat(lhs.FingerprintsStat, rhs.FingerprintsStat); err != nil {
		return err
	}

	//hash_of_candidate_image_file
	if !reflect.DeepEqual(lhs.HashOfCandidateImageFile, rhs.HashOfCandidateImageFile) {
		return errors.New("image hash do not match")
	}

	//perceptual_image_hashes
	if err := ComparePerceptualImageHashes(lhs.PerceptualImageHashes, rhs.PerceptualImageHashes); err != nil {
		return errors.Errorf("image hashes do not match: %w", err)
	}

	//perceptual_hash_overlap_count
	if lhs.PerceptualHashOverlapCount != rhs.PerceptualHashOverlapCount {
		return errors.Errorf("perceptual_hash_overlap_count not matched: lhs(%d) != rhs(%d)", lhs.PerceptualHashOverlapCount, rhs.PerceptualHashOverlapCount)
	}

	//maxes
	if err := CompareMaxes(lhs.Maxes, rhs.Maxes); err != nil {
		return err
	}

	//percentile
	return ComparePercentile(lhs.Percentile, rhs.Percentile)
}

// CompareRarenessScores  return nil if two RarenessScores are equal
func CompareRarenessScores(lhs *RarenessScores, rhs *RarenessScores) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil rareness_scores")
	}

	if !compareFloat(lhs.CombinedRarenessScore, rhs.CombinedRarenessScore) {
		return errors.Errorf("combined_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.CombinedRarenessScore, rhs.CombinedRarenessScore)
	}

	if !compareFloat(lhs.XgboostPredictedRarenessScore, rhs.XgboostPredictedRarenessScore) {
		return errors.Errorf("xgboost_predicted_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.XgboostPredictedRarenessScore, rhs.XgboostPredictedRarenessScore)
	}

	if !compareFloat(lhs.NnPredictedRarenessScore, rhs.NnPredictedRarenessScore) {
		return errors.Errorf("nn_predicted_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.NnPredictedRarenessScore, rhs.NnPredictedRarenessScore)
	}

	if !compareFloat(lhs.OverallAverageRarenessScore, rhs.OverallAverageRarenessScore) {
		return errors.Errorf("overall_average_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.OverallAverageRarenessScore, rhs.OverallAverageRarenessScore)
	}

	return nil
}

// CompareInternetRareness return nil if two InternetRareness are equal
func CompareInternetRareness(lhs *InternetRareness, rhs *InternetRareness) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil internet_rareness")
	}

	if lhs.MatchesFoundOnFirstPage != rhs.MatchesFoundOnFirstPage {
		return errors.Errorf("matches_found_on_first_page not matched: lhs(%d) != rhs(%d)", lhs.MatchesFoundOnFirstPage, rhs.MatchesFoundOnFirstPage)
	}

	if lhs.NumberOfPagesOfResults != rhs.NumberOfPagesOfResults {
		return errors.Errorf("number_of_pages_of_results not matched: lhs(%d) != rhs(%d)", lhs.NumberOfPagesOfResults, rhs.NumberOfPagesOfResults)
	}

	if !strings.EqualFold(lhs.URLOfFirstMatchInPage, rhs.URLOfFirstMatchInPage) {
		return errors.Errorf("url_of_first_match_in_page not matched: lhs(%s) != rhs(%s)", lhs.URLOfFirstMatchInPage, rhs.URLOfFirstMatchInPage)
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

// CompareFingerprintsStat compares two FingerprintsStat
func CompareFingerprintsStat(lhs *FingerprintsStat, rhs *FingerprintsStat) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil fingerprints_stat")
	}

	if lhs.NumberOfFingerprintsRequiringFurtherTesting1 != rhs.NumberOfFingerprintsRequiringFurtherTesting1 {
		return errors.Errorf("number_of_fingerprints_requiring_further_testing_1 not matched: lhs(%d) != rhs(%d)", lhs.NumberOfFingerprintsRequiringFurtherTesting1, rhs.NumberOfFingerprintsRequiringFurtherTesting1)
	}

	if lhs.NumberOfFingerprintsRequiringFurtherTesting2 != rhs.NumberOfFingerprintsRequiringFurtherTesting2 {
		return errors.Errorf("number_of_fingerprints_requiring_further_testing_2 not matched: lhs(%d) != rhs(%d)", lhs.NumberOfFingerprintsRequiringFurtherTesting2, rhs.NumberOfFingerprintsRequiringFurtherTesting2)
	}

	if lhs.NumberOfFingerprintsRequiringFurtherTesting3 != rhs.NumberOfFingerprintsRequiringFurtherTesting3 {
		return errors.Errorf("number_of_fingerprints_requiring_further_testing_3 not matched: lhs(%d) != rhs(%d)", lhs.NumberOfFingerprintsRequiringFurtherTesting3, rhs.NumberOfFingerprintsRequiringFurtherTesting3)
	}

	if lhs.NumberOfFingerprintsRequiringFurtherTesting4 != rhs.NumberOfFingerprintsRequiringFurtherTesting4 {
		return errors.Errorf("number_of_fingerprints_requiring_further_testing_4 not matched: lhs(%d) != rhs(%d)", lhs.NumberOfFingerprintsRequiringFurtherTesting4, rhs.NumberOfFingerprintsRequiringFurtherTesting4)
	}

	if lhs.NumberOfFingerprintsRequiringFurtherTesting5 != rhs.NumberOfFingerprintsRequiringFurtherTesting5 {
		return errors.Errorf("number_of_fingerprints_requiring_further_testing_5 not matched: lhs(%d) != rhs(%d)", lhs.NumberOfFingerprintsRequiringFurtherTesting5, rhs.NumberOfFingerprintsRequiringFurtherTesting5)
	}

	if lhs.NumberOfFingerprintsRequiringFurtherTesting6 != rhs.NumberOfFingerprintsRequiringFurtherTesting6 {
		return errors.Errorf("number_of_fingerprints_requiring_further_testing_6 not matched: lhs(%d) != rhs(%d)", lhs.NumberOfFingerprintsRequiringFurtherTesting6, rhs.NumberOfFingerprintsRequiringFurtherTesting6)
	}

	if lhs.NumberOfFingerprintsOfSuspectedDupes != rhs.NumberOfFingerprintsOfSuspectedDupes {
		return errors.Errorf("number_of_fingerprints_of_suspected_dupes not matched: lhs(%d) != rhs(%d)", lhs.NumberOfFingerprintsOfSuspectedDupes, rhs.NumberOfFingerprintsOfSuspectedDupes)
	}
	return nil
}

// ComparePerceptualImageHashes return nil if two PerceptualImageHashes are equal
func ComparePerceptualImageHashes(lhs *PerceptualImageHashes, rhs *PerceptualImageHashes) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil perceptual_image_hashes")
	}

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

// CompareMaxes return nil if two Maxes are equal
func CompareMaxes(lhs *Maxes, rhs *Maxes) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil maxes")
	}

	if !compareFloat(lhs.PearsonMax, rhs.PearsonMax) {
		return errors.Errorf("pearson_max not matched: lhs(%f) != rhs(%f)", lhs.PearsonMax, rhs.PearsonMax)
	}
	if !compareFloat(lhs.SpearmanMax, rhs.SpearmanMax) {
		return errors.Errorf("spearman_max not matched: lhs(%f) != rhs(%f)", lhs.SpearmanMax, rhs.SpearmanMax)
	}
	if !compareFloat(lhs.KendallMax, rhs.KendallMax) {
		return errors.Errorf("kendall_max not matched: lhs(%f) != rhs(%f)", lhs.KendallMax, rhs.KendallMax)
	}
	if !compareFloat(lhs.HoeffdingMax, rhs.HoeffdingMax) {
		return errors.Errorf("hoeffding_maxnot matched: lhs(%f) != rhs(%f)", lhs.HoeffdingMax, rhs.HoeffdingMax)
	}
	if !compareFloat(lhs.MutualInformationMax, rhs.MutualInformationMax) {
		return errors.Errorf("mutual_information_max not matched: lhs(%f) != rhs(%f)", lhs.MutualInformationMax, rhs.MutualInformationMax)
	}
	if !compareFloat(lhs.HsicMax, rhs.HsicMax) {
		return errors.Errorf("hsic_max not matched: lhs(%f) != rhs(%f)", lhs.HsicMax, rhs.HsicMax)
	}
	if !compareFloat(lhs.XgbimportanceMax, rhs.XgbimportanceMax) {
		return errors.Errorf("xgbimportance_max not matched: lhs(%f) != rhs(%f)", lhs.XgbimportanceMax, rhs.XgbimportanceMax)
	}
	return nil
}

// ComparePercentile return nil if two Percentile are equal
func ComparePercentile(lhs *Percentile, rhs *Percentile) error {
	if lhs == nil || rhs == nil {
		return errors.New("nil percentile")
	}

	if !compareFloat(lhs.PearsonTop1BpsPercentile, rhs.PearsonTop1BpsPercentile) {
		return errors.Errorf("pearson_top_1_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.PearsonTop1BpsPercentile, rhs.PearsonTop1BpsPercentile)
	}
	if !compareFloat(lhs.SpearmanTop1BpsPercentile, rhs.SpearmanTop1BpsPercentile) {
		return errors.Errorf("spearman_top_1_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.SpearmanTop1BpsPercentile, rhs.SpearmanTop1BpsPercentile)
	}
	if !compareFloat(lhs.KendallTop1BpsPercentile, rhs.KendallTop1BpsPercentile) {
		return errors.Errorf("kendall_top_1_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.KendallTop1BpsPercentile, rhs.KendallTop1BpsPercentile)
	}
	if !compareFloat(lhs.HoeffdingTop10BpsPercentile, rhs.HoeffdingTop10BpsPercentile) {
		return errors.Errorf("hoeffding_top_10_bps_percentile matched: lhs(%f) != rhs(%f)", lhs.HoeffdingTop10BpsPercentile, rhs.HoeffdingTop10BpsPercentile)
	}
	if !compareFloat(lhs.MutualInformationTop100BpsPercentile, rhs.MutualInformationTop100BpsPercentile) {
		return errors.Errorf("mutual_information_top_100_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.MutualInformationTop100BpsPercentile, rhs.MutualInformationTop100BpsPercentile)
	}
	if !compareFloat(lhs.HsicTop100BpsPercentile, rhs.HsicTop100BpsPercentile) {
		return errors.Errorf("hsic_top_100_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.HsicTop100BpsPercentile, rhs.HsicTop100BpsPercentile)
	}
	if !compareFloat(lhs.XgbimportanceTop100BpsPercentile, rhs.XgbimportanceTop100BpsPercentile) {
		return errors.Errorf("xgbimportance_top_100_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.XgbimportanceTop100BpsPercentile, rhs.XgbimportanceTop100BpsPercentile)
	}
	return nil
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
