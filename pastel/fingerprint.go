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
	if !compareFloatWithPrecision(lhs.OpenNSFWScore, rhs.OpenNSFWScore, 7.0) {
		return errors.Errorf("open_nsfw_score do not match: lhs(%f) != rhs(%f)", lhs.OpenNSFWScore, rhs.OpenNSFWScore)
	}

	//alternative_nsfw_scores
	if err := CompareAlternativeNSFWScore(lhs.AlternativeNSFWScores, rhs.AlternativeNSFWScores); err != nil {
		return errors.Errorf("alternative nsfw score do not match: %w", err)
	}

	// ImageFingerprintOfCandidateImageFile
	if !compareFloatsWithPrecision(lhs.ImageFingerprintOfCandidateImageFile, rhs.ImageFingerprintOfCandidateImageFile, 5.0) {
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
	if err := ComparePercentile(lhs.Percentile, rhs.Percentile); err != nil {
		return err
	}

	return nil
}

// CompareRarenessScores  return nil if two RarenessScores are equal
func CompareRarenessScores(lhs *RarenessScores, rhs *RarenessScores) error {
	if !compareFloatWithPrecision(lhs.CombinedRarenessScore, rhs.CombinedRarenessScore, 5.0) {
		return errors.Errorf("combined_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.CombinedRarenessScore, rhs.CombinedRarenessScore)
	}

	if !compareFloatWithPrecision(lhs.XgboostPredictedRarenessScore, rhs.XgboostPredictedRarenessScore, 5.0) {
		return errors.Errorf("xgboost_predicted_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.XgboostPredictedRarenessScore, rhs.XgboostPredictedRarenessScore)
	}

	if !compareFloatWithPrecision(lhs.NnPredictedRarenessScore, rhs.NnPredictedRarenessScore, 5.0) {
		return errors.Errorf("nn_predicted_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.NnPredictedRarenessScore, rhs.NnPredictedRarenessScore)
	}

	if !compareFloatWithPrecision(lhs.OverallAverageRarenessScore, rhs.OverallAverageRarenessScore, 5.0) {
		return errors.Errorf("overall_average_rareness_score not matched: lhs(%f) != rhs(%f)", lhs.OverallAverageRarenessScore, rhs.OverallAverageRarenessScore)
	}

	return nil
}

// CompareInternetRareness return nil if two InternetRareness are equal
func CompareInternetRareness(lhs *InternetRareness, rhs *InternetRareness) error {
	if lhs.MatchesFoundOnFirstPage != rhs.MatchesFoundOnFirstPage {
		return errors.Errorf("matches_found_on_first_page not matched: lhs(%d) != rhs(%d)", lhs.MatchesFoundOnFirstPage, rhs.MatchesFoundOnFirstPage)
	}

	if lhs.NumberOfPagesOfResults != rhs.NumberOfPagesOfResults {
		return errors.Errorf("number_of_pages_of_results not matched: lhs(%d) != rhs(%d)", lhs.NumberOfPagesOfResults, rhs.NumberOfPagesOfResults)
	}

	if !strings.EqualFold(lhs.UrlOfFirstMatchInPage, rhs.UrlOfFirstMatchInPage) {
		return errors.Errorf("url_of_first_match_in_page not matched: lhs(%s) != rhs(%s)", lhs.UrlOfFirstMatchInPage, rhs.UrlOfFirstMatchInPage)
	}
	return nil
}

// CompareAlternativeNSFWScore return nil if two AlternativeNSFWScore are equal
func CompareAlternativeNSFWScore(lhs *AlternativeNSFWScores, rhs *AlternativeNSFWScores) error {

	if !compareFloatWithPrecision(lhs.Drawings, rhs.Drawings, 5.0) {
		return errors.Errorf("drawings score not matched: lhs(%f) != rhs(%f)", lhs.Drawings, rhs.Drawings)
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

// CompareFingerprintsStat compares two FingerprintsStat
func CompareFingerprintsStat(lhs *FingerprintsStat, rhs *FingerprintsStat) error {
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

	if !compareFloatWithPrecision(lhs.PearsonMax, rhs.PearsonMax, 5.0) {
		return errors.Errorf("pearson_max not matched: lhs(%f) != rhs(%f)", lhs.PearsonMax, rhs.PearsonMax)
	}
	if !compareFloatWithPrecision(lhs.SpearmanMax, rhs.SpearmanMax, 5.0) {
		return errors.Errorf("spearman_max not matched: lhs(%f) != rhs(%f)", lhs.SpearmanMax, rhs.SpearmanMax)
	}
	if !compareFloatWithPrecision(lhs.KendallMax, rhs.KendallMax, 5.0) {
		return errors.Errorf("kendall_max not matched: lhs(%f) != rhs(%f)", lhs.KendallMax, rhs.KendallMax)
	}
	if !compareFloatWithPrecision(lhs.HoeffdingMax, rhs.HoeffdingMax, 5.0) {
		return errors.Errorf("hoeffding_maxnot matched: lhs(%f) != rhs(%f)", lhs.HoeffdingMax, rhs.HoeffdingMax)
	}
	if !compareFloatWithPrecision(lhs.MutualInformationMax, rhs.MutualInformationMax, 5.0) {
		return errors.Errorf("mutual_information_max not matched: lhs(%f) != rhs(%f)", lhs.MutualInformationMax, rhs.MutualInformationMax)
	}
	if !compareFloatWithPrecision(lhs.HsicMax, rhs.HsicMax, 5.0) {
		return errors.Errorf("hsic_max not matched: lhs(%f) != rhs(%f)", lhs.HsicMax, rhs.HsicMax)
	}
	if !compareFloatWithPrecision(lhs.XgbimportanceMax, rhs.XgbimportanceMax, 5.0) {
		return errors.Errorf("xgbimportance_max not matched: lhs(%f) != rhs(%f)", lhs.XgbimportanceMax, rhs.XgbimportanceMax)
	}
	return nil
}

// ComparePercentile return nil if two Percentile are equal
func ComparePercentile(lhs *Percentile, rhs *Percentile) error {

	if !compareFloatWithPrecision(lhs.PearsonTop1BpsPercentile, rhs.PearsonTop1BpsPercentile, 5.0) {
		return errors.Errorf("pearson_top_1_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.PearsonTop1BpsPercentile, rhs.PearsonTop1BpsPercentile)
	}
	if !compareFloatWithPrecision(lhs.SpearmanTop1BpsPercentile, rhs.SpearmanTop1BpsPercentile, 5.0) {
		return errors.Errorf("spearman_top_1_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.SpearmanTop1BpsPercentile, rhs.SpearmanTop1BpsPercentile)
	}
	if !compareFloatWithPrecision(lhs.KendallTop1BpsPercentile, rhs.KendallTop1BpsPercentile, 5.0) {
		return errors.Errorf("kendall_top_1_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.KendallTop1BpsPercentile, rhs.KendallTop1BpsPercentile)
	}
	if !compareFloatWithPrecision(lhs.HoeffdingTop10BpsPercentile, rhs.HoeffdingTop10BpsPercentile, 5.0) {
		return errors.Errorf("hoeffding_top_10_bps_percentile matched: lhs(%f) != rhs(%f)", lhs.HoeffdingTop10BpsPercentile, rhs.HoeffdingTop10BpsPercentile)
	}
	if !compareFloatWithPrecision(lhs.MutualInformationTop100BpsPercentile, rhs.MutualInformationTop100BpsPercentile, 5.0) {
		return errors.Errorf("mutual_information_top_100_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.MutualInformationTop100BpsPercentile, rhs.MutualInformationTop100BpsPercentile)
	}
	if !compareFloatWithPrecision(lhs.HsicTop100BpsPercentile, rhs.HsicTop100BpsPercentile, 5.0) {
		return errors.Errorf("hsic_top_100_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.HsicTop100BpsPercentile, rhs.HsicTop100BpsPercentile)
	}
	if !compareFloatWithPrecision(lhs.XgbimportanceTop100BpsPercentile, rhs.XgbimportanceTop100BpsPercentile, 5.0) {
		return errors.Errorf("xgbimportance_top_100_bps_percentile not matched: lhs(%f) != rhs(%f)", lhs.XgbimportanceTop100BpsPercentile, rhs.XgbimportanceTop100BpsPercentile)
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

func compareFloatsWithPrecision(l []float32, r []float32, prec float64) bool {
	if len(l) != len(r) {
		return false
	}

	for i := 0; i < len(l); i = i + 1 {
		if !compareFloatWithPrecision(l[i], r[i], prec) {
			return false
		}
	}

	return true
}
