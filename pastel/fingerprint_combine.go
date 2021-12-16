package pastel

import (
	"reflect"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
)

func copyFingerPrintAndScores(origin *DDAndFingerprints) *DDAndFingerprints {
	copiedRarenessScores := *origin.RarenessScores
	copiedInternetRareness := *origin.InternetRareness
	copiedAlternativeNSFWScores := *origin.AlternativeNSFWScores
	copiedFingerprintsStat := *origin.FingerprintsStat
	copiedPerceptualImageHashes := *origin.PerceptualImageHashes
	copiedPercentile := *origin.Percentile
	copiedMaxes := *origin.Maxes
	copied := &DDAndFingerprints{
		Block:                      origin.Block,
		Principal:                  origin.Principal,
		DupeDetectionSystemVersion: origin.DupeDetectionSystemVersion,

		IsLikelyDupe:     origin.IsLikelyDupe,
		IsRareOnInternet: origin.IsRareOnInternet,

		RarenessScores:   &copiedRarenessScores,
		InternetRareness: &copiedInternetRareness,

		OpenNSFWScore:         origin.OpenNSFWScore, // Update later
		AlternativeNSFWScores: &copiedAlternativeNSFWScores,

		ImageFingerprintOfCandidateImageFile: origin.ImageFingerprintOfCandidateImageFile,
		FingerprintsStat:                     &copiedFingerprintsStat,

		HashOfCandidateImageFile:   origin.HashOfCandidateImageFile,
		PerceptualImageHashes:      &copiedPerceptualImageHashes,
		PerceptualHashOverlapCount: origin.PerceptualHashOverlapCount,

		Maxes:      &copiedMaxes,
		Percentile: &copiedPercentile,
	}
	return copied
}

// CombineFingerPrintAndScores determines
// Refer: https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration - 4B
// The order of DDAndFingerprints should be higher rank first
func CombineFingerPrintAndScores(first *DDAndFingerprints, second *DDAndFingerprints, third *DDAndFingerprints) (*DDAndFingerprints, error) {
	if first == nil || second == nil || third == nil {
		return nil, errors.New("nil input")
	}

	finalResult := copyFingerPrintAndScores(first)

	// Block
	if !strings.EqualFold(first.Block, second.Block) {
		return nil, errors.Errorf("block not matched: first(%s) != second(%s)", first.Block, second.Block)
	}

	if !strings.EqualFold(first.Block, third.Block) {
		return nil, errors.Errorf("block not matched: first(%s) != third(%s)", first.Block, third.Block)
	}

	// Principal
	if !strings.EqualFold(first.Principal, second.Principal) {
		return nil, errors.Errorf("principal not matched: first(%s) != second(%s)", first.Principal, second.Principal)
	}
	if !strings.EqualFold(first.Principal, third.Principal) {
		return nil, errors.Errorf("principal not matched: first(%s) != third(%s)", first.Principal, third.Principal)
	}

	// DupeDetectionSystemVersion
	if !strings.EqualFold(first.DupeDetectionSystemVersion, second.DupeDetectionSystemVersion) {
		return nil, errors.Errorf("dupe_detection_system_version not matched: first(%s) != second(%s)", first.DupeDetectionSystemVersion, second.DupeDetectionSystemVersion)
	}

	if !strings.EqualFold(first.DupeDetectionSystemVersion, third.DupeDetectionSystemVersion) {
		return nil, errors.Errorf("dupe_detection_system_version not matched: first(%s) != third(%s)", first.DupeDetectionSystemVersion, third.DupeDetectionSystemVersion)
	}

	// IsLikelyDupe
	if first.IsLikelyDupe != second.IsLikelyDupe {
		return nil, errors.Errorf("is_likely_dupe not matched: first(%t) != second(%t)", first.IsLikelyDupe, second.IsLikelyDupe)
	}
	if first.IsLikelyDupe != third.IsLikelyDupe {
		return nil, errors.Errorf("is_likely_dupe not matched: first(%t) != third(%t)", first.IsLikelyDupe, third.IsLikelyDupe)
	}

	// IsRareOnInternet
	// if first.IsRareOnInternet != second.IsRareOnInternet {
	// 	return nil, errors.Errorf("is_rare_on_internet not matched: first(%t) != second(%t)", first.IsRareOnInternet, second.IsRareOnInternet)
	// }

	// if first.IsRareOnInternet != third.IsRareOnInternet {
	// 	return nil, errors.Errorf("is_rare_on_internet not matched: first(%t) != third(%t)", first.IsRareOnInternet, third.IsRareOnInternet)
	// }
	finalResult.IsRareOnInternet = combineBool(first.IsRareOnInternet, second.IsRareOnInternet, third.IsRareOnInternet)

	//RarenessScores
	if err := CompareRarenessScores(first.RarenessScores, second.RarenessScores); err != nil {
		return nil, errors.Errorf("first couple of RarenessScores do not match: %w", err)
	}
	if err := CompareRarenessScores(first.RarenessScores, third.RarenessScores); err != nil {
		return nil, errors.Errorf("second couple of RarenessScores do not third: %w", err)
	}

	//InternetRareness
	// if err := CompareInternetRareness(first.InternetRareness, second.InternetRareness); err != nil {
	// 	return nil, errors.Errorf("first couple of InternetRareness do not match: %w", err)
	// }
	// if err := CompareInternetRareness(first.InternetRareness, third.InternetRareness); err != nil {
	// 	return nil, errors.Errorf("second couple of InternetRareness do not third: %w", err)
	// }
	finalResult.InternetRareness.MatchesFoundOnFirstPage = combineUint32(
		first.InternetRareness.MatchesFoundOnFirstPage,
		second.InternetRareness.MatchesFoundOnFirstPage,
		third.InternetRareness.MatchesFoundOnFirstPage,
	)
	finalResult.InternetRareness.NumberOfPagesOfResults = combineUint32(
		first.InternetRareness.NumberOfPagesOfResults,
		second.InternetRareness.NumberOfPagesOfResults,
		third.InternetRareness.NumberOfPagesOfResults,
	)
	finalResult.InternetRareness.UrlOfFirstMatchInPage = combineString(
		first.InternetRareness.UrlOfFirstMatchInPage,
		second.InternetRareness.UrlOfFirstMatchInPage,
		third.InternetRareness.UrlOfFirstMatchInPage,
	)

	//open_nsfw_score
	if !compareFloat(first.OpenNSFWScore, second.OpenNSFWScore) {
		return nil, errors.Errorf("open_nsfw_score do not match: first(%f) != second(%f)", first.OpenNSFWScore, second.OpenNSFWScore)
	}
	if !compareFloat(first.OpenNSFWScore, third.OpenNSFWScore) {
		return nil, errors.Errorf("open_nsfw_score do not match: first(%f) != third(%f)", first.OpenNSFWScore, third.OpenNSFWScore)
	}

	//alternative_nsfw_scores
	if err := CompareAlternativeNSFWScore(first.AlternativeNSFWScores, second.AlternativeNSFWScores); err != nil {
		return nil, errors.Errorf("first couple of alternative_nsfw_scores do not match: %w", err)
	}
	if err := CompareAlternativeNSFWScore(first.AlternativeNSFWScores, third.AlternativeNSFWScores); err != nil {
		return nil, errors.Errorf("second couple of alternative_nsfw_scores do not third: %w", err)
	}

	// ImageFingerprintOfCandidateImageFile
	if !compareFloats(first.ImageFingerprintOfCandidateImageFile, second.ImageFingerprintOfCandidateImageFile) {
		return nil, errors.New("first couple of image_fingerprint_of_candidate_image_file nsfw score do not match")
	}
	if !compareFloats(first.ImageFingerprintOfCandidateImageFile, third.ImageFingerprintOfCandidateImageFile) {
		return nil, errors.New("second couple of image_fingerprint_of_candidate_image_file nsfw score do not match")
	}

	//alternative_nsfw_scores
	if err := CompareFingerprintsStat(first.FingerprintsStat, second.FingerprintsStat); err != nil {
		return nil, errors.Errorf("first couple of alternative_nsfw_scores do not match: %w", err)
	}
	if err := CompareFingerprintsStat(first.FingerprintsStat, third.FingerprintsStat); err != nil {
		return nil, errors.Errorf("second couple of alternative_nsfw_scores  do not match: %w", err)
	}

	//hash_of_candidate_image_file
	if !reflect.DeepEqual(first.HashOfCandidateImageFile, second.HashOfCandidateImageFile) {
		return nil, errors.New("first couple of hash_of_candidate_image_filedo not match")
	}
	if !reflect.DeepEqual(first.HashOfCandidateImageFile, third.HashOfCandidateImageFile) {
		return nil, errors.New("second couple of hash_of_candidate_image_file do not match")
	}

	//perceptual_image_hashes
	if err := ComparePerceptualImageHashes(first.PerceptualImageHashes, second.PerceptualImageHashes); err != nil {
		return nil, errors.Errorf("first couple of perceptual_image_hashes do not match: %w", err)
	}
	if err := ComparePerceptualImageHashes(first.PerceptualImageHashes, third.PerceptualImageHashes); err != nil {
		return nil, errors.Errorf("second couple of perceptual_image_hashes do not match: %w", err)
	}

	//perceptual_hash_overlap_count
	if first.PerceptualHashOverlapCount != second.PerceptualHashOverlapCount {
		return nil, errors.Errorf("perceptual_hash_overlap_count not matched: first(%d) != second(%d)", first.PerceptualHashOverlapCount, second.PerceptualHashOverlapCount)
	}
	if first.PerceptualHashOverlapCount != third.PerceptualHashOverlapCount {
		return nil, errors.Errorf("perceptual_hash_overlap_count not matched: first(%d) != third(%d)", first.PerceptualHashOverlapCount, second.PerceptualHashOverlapCount)
	}

	//maxes
	if err := CompareMaxes(first.Maxes, second.Maxes); err != nil {
		return nil, errors.Errorf("first maxes compare not match: %w", err)
	}
	if err := CompareMaxes(first.Maxes, third.Maxes); err != nil {
		return nil, errors.Errorf("second maxes compare not match: %w", err)
	}

	//percentile
	if err := ComparePercentile(first.Percentile, second.Percentile); err != nil {
		return nil, errors.Errorf("first couple of percentile compare not match: %w", err)
	}
	if err := ComparePercentile(first.Percentile, third.Percentile); err != nil {
		return nil, errors.Errorf("second couple of percentile compare not match: %w", err)
	}

	return finalResult, nil
}

func combineBool(a, b, c bool) bool {
	if b == c {
		return b
	}

	return a
}

func combineUint32(a, b, c uint32) uint32 {
	if b == c {
		return b
	}

	return a
}

func combineString(a, b, c string) string {
	if strings.EqualFold(a, b) {
		return a
	}

	if strings.EqualFold(a, c) {
		return a
	}

	if strings.EqualFold(b, c) {
		return b
	}

	return a
}
