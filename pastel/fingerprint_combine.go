package pastel

import (
	"reflect"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
)

func copyFingerPrintAndScores(origin *DDAndFingerprints) *DDAndFingerprints {
	copiedInternetRareness := *origin.InternetRareness
	copiedAlternativeNSFWScores := *origin.AlternativeNSFWScores
	copied := &DDAndFingerprints{
		BlockHash:                  origin.BlockHash,
		BlockHeight:                origin.BlockHeight,
		TimestampOfRequest:         origin.TimestampOfRequest,
		SubmitterPastelID:          origin.SubmitterPastelID,
		SN1PastelID:                origin.SN1PastelID,
		SN2PastelID:                origin.SN2PastelID,
		SN3PastelID:                origin.SN3PastelID,
		IsOpenAPIRequest:           origin.IsOpenAPIRequest,
		OpenAPISubsetID:            origin.OpenAPISubsetID,
		DupeDetectionSystemVersion: origin.DupeDetectionSystemVersion,
		IsLikelyDupe:               origin.IsLikelyDupe,
		IsRareOnInternet:           origin.IsRareOnInternet,
		OverallRarenessScore:       origin.OverallRarenessScore,

		PctOfTop10MostSimilarWithDupeProbAbove25pct: origin.PctOfTop10MostSimilarWithDupeProbAbove25pct,
		PctOfTop10MostSimilarWithDupeProbAbove33pct: origin.PctOfTop10MostSimilarWithDupeProbAbove33pct,
		PctOfTop10MostSimilarWithDupeProbAbove50pct: origin.PctOfTop10MostSimilarWithDupeProbAbove50pct,

		RarenessScoresTableJSONCompressedB64: origin.RarenessScoresTableJSONCompressedB64,

		InternetRareness: &copiedInternetRareness,

		OpenNSFWScore:         origin.OpenNSFWScore, // Update later
		AlternativeNSFWScores: &copiedAlternativeNSFWScores,

		ImageFingerprintOfCandidateImageFile: origin.ImageFingerprintOfCandidateImageFile,

		HashOfCandidateImageFile: origin.HashOfCandidateImageFile,

		CollectionNameString:                       origin.CollectionNameString,
		OpenAPIGroupIDString:                       origin.OpenAPIGroupIDString,
		GroupRarenessScore:                         origin.GroupRarenessScore,
		CandidateImageThumbnailWebpAsBase64String:  origin.CandidateImageThumbnailWebpAsBase64String,
		DoesNotImpactTheFollowingCollectionStrings: origin.DoesNotImpactTheFollowingCollectionStrings,
		IsInvalidSenseRequest:                      origin.IsInvalidSenseRequest,
		InvalidSenseRequestReason:                  origin.InvalidSenseRequestReason,
		SimilarityScoreToFirstEntryInCollection:    origin.SimilarityScoreToFirstEntryInCollection,
		CPProbability:                              origin.CPProbability,
		ChildProbability:                           origin.ChildProbability,
		ImageFilePath:                              origin.ImageFilePath,
	}
	return copied
}

// CombineFingerPrintAndScores verifies erify dd_and_fingerprints data received from 2 other
//
//	SuperNodes and SNs own dd_and_fingerprints are the same or “close enough”
//
// Refer: https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration - 4.B.3
// The order of DDAndFingerprints should be higher rank first
func CombineFingerPrintAndScores(first *DDAndFingerprints, second *DDAndFingerprints, third *DDAndFingerprints) (*DDAndFingerprints, error) {
	if first == nil || second == nil || third == nil {
		return nil, errors.New("nil input")
	}

	finalResult := copyFingerPrintAndScores(first)

	// BlockHash
	if !strings.EqualFold(first.BlockHash, second.BlockHash) {
		return nil, errors.Errorf("BlockHash not matched: first(%s) != second(%s)", first.BlockHash, second.BlockHash)
	}

	if !strings.EqualFold(first.BlockHash, third.BlockHash) {
		return nil, errors.Errorf("BlockHash not matched: first(%s) != third(%s)", first.BlockHash, third.BlockHash)
	}

	// BlockHeight
	if !strings.EqualFold(first.BlockHeight, second.BlockHeight) {
		return nil, errors.Errorf("BlockHeight not matched: first(%s) != second(%s)", first.BlockHeight, second.BlockHeight)
	}

	if !strings.EqualFold(first.BlockHeight, third.BlockHeight) {
		return nil, errors.Errorf("BlockHeight not matched: first(%s) != third(%s)", first.BlockHeight, third.BlockHeight)
	}

	// TimestampOfRequest
	if !strings.EqualFold(first.TimestampOfRequest, second.TimestampOfRequest) {
		return nil, errors.Errorf("TimestampOfRequest not matched: first(%s) != second(%s)", first.TimestampOfRequest, second.TimestampOfRequest)
	}

	if !strings.EqualFold(first.TimestampOfRequest, third.TimestampOfRequest) {
		return nil, errors.Errorf("TimestampOfRequest not matched: first(%s) != third(%s)", first.TimestampOfRequest, third.TimestampOfRequest)
	}

	// SubmitterPastelID
	if !strings.EqualFold(first.SubmitterPastelID, second.SubmitterPastelID) {
		return nil, errors.Errorf("SubmitterPastelID not matched: first(%s) != second(%s)", first.SubmitterPastelID, second.SubmitterPastelID)
	}

	if !strings.EqualFold(first.SubmitterPastelID, third.SubmitterPastelID) {
		return nil, errors.Errorf("SubmitterPastelID not matched: first(%s) != third(%s)", first.SubmitterPastelID, third.SubmitterPastelID)
	}

	// SN1PastelID
	if !strings.EqualFold(first.SN1PastelID, second.SN1PastelID) {
		return nil, errors.Errorf("SN1PastelID not matched: first(%s) != second(%s)", first.SN1PastelID, second.SN1PastelID)
	}

	if !strings.EqualFold(first.SN1PastelID, third.SN1PastelID) {
		return nil, errors.Errorf("SN1PastelID not matched: first(%s) != third(%s)", first.SN1PastelID, third.SN1PastelID)
	}

	// SN2PastelID
	if !strings.EqualFold(first.SN2PastelID, second.SN2PastelID) {
		return nil, errors.Errorf("SN2PastelID not matched: first(%s) != second(%s)", first.SN2PastelID, second.SN2PastelID)
	}

	if !strings.EqualFold(first.SN2PastelID, third.SN2PastelID) {
		return nil, errors.Errorf("SN2PastelID not matched: first(%s) != third(%s)", first.SN2PastelID, third.SN2PastelID)
	}

	// SN3PastelID
	if !strings.EqualFold(first.SN3PastelID, second.SN3PastelID) {
		return nil, errors.Errorf("SN3PastelID not matched: first(%s) != second(%s)", first.SN3PastelID, second.SN3PastelID)
	}

	if !strings.EqualFold(first.SN3PastelID, third.SN3PastelID) {
		return nil, errors.Errorf("SN3PastelID not matched: first(%s) != third(%s)", first.SN3PastelID, third.SN3PastelID)
	}

	// IsOpenAPIRequest
	if first.IsOpenAPIRequest != second.IsOpenAPIRequest {
		return nil, errors.Errorf("IsOpenAPIRequest not matched: first(%t) != second(%t)", first.IsOpenAPIRequest, second.IsOpenAPIRequest)
	}
	if first.IsOpenAPIRequest != third.IsOpenAPIRequest {
		return nil, errors.Errorf("IsOpenAPIRequest not matched: first(%t) != third(%t)", first.IsOpenAPIRequest, third.IsOpenAPIRequest)
	}

	// OpenAPISubsetID
	if !strings.EqualFold(first.OpenAPISubsetID, second.OpenAPISubsetID) {
		return nil, errors.Errorf("OpenAPISubsetID not matched: first(%s) != second(%s)", first.OpenAPISubsetID, second.OpenAPISubsetID)
	}

	if !strings.EqualFold(first.OpenAPISubsetID, third.OpenAPISubsetID) {
		return nil, errors.Errorf("OpenAPISubsetID not matched: first(%s) != third(%s)", first.OpenAPISubsetID, third.OpenAPISubsetID)
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

	finalResult.IsRareOnInternet = combineBool(first.IsRareOnInternet, second.IsRareOnInternet, third.IsRareOnInternet)

	//OverallRarenessScore
	if !compareFloat(first.OverallRarenessScore, second.OverallRarenessScore) {
		return nil, errors.Errorf("OverallRarenessScore do not match: first(%f) != second(%f)", first.OverallRarenessScore, second.OverallRarenessScore)
	}
	if !compareFloat(first.OverallRarenessScore, third.OverallRarenessScore) {
		return nil, errors.Errorf("OverallRarenessScore do not match: first(%f) != third(%f)", first.OverallRarenessScore, third.OverallRarenessScore)
	}

	//PctOfTop10MostSimilarWithDupeProbAbove25pct
	if !compareFloat(first.PctOfTop10MostSimilarWithDupeProbAbove25pct, second.PctOfTop10MostSimilarWithDupeProbAbove25pct) {
		return nil, errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove25pct do not match: first(%f) != second(%f)", first.PctOfTop10MostSimilarWithDupeProbAbove25pct, second.PctOfTop10MostSimilarWithDupeProbAbove25pct)
	}
	if !compareFloat(first.PctOfTop10MostSimilarWithDupeProbAbove25pct, third.PctOfTop10MostSimilarWithDupeProbAbove25pct) {
		return nil, errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove25pct do not match: first(%f) != third(%f)", first.PctOfTop10MostSimilarWithDupeProbAbove25pct, third.PctOfTop10MostSimilarWithDupeProbAbove25pct)
	}

	//PctOfTop10MostSimilarWithDupeProbAbove33pct
	if !compareFloat(first.PctOfTop10MostSimilarWithDupeProbAbove33pct, second.PctOfTop10MostSimilarWithDupeProbAbove33pct) {
		return nil, errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove33pct do not match: first(%f) != second(%f)", first.PctOfTop10MostSimilarWithDupeProbAbove33pct, second.PctOfTop10MostSimilarWithDupeProbAbove33pct)
	}
	if !compareFloat(first.PctOfTop10MostSimilarWithDupeProbAbove33pct, third.PctOfTop10MostSimilarWithDupeProbAbove33pct) {
		return nil, errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove33pct do not match: first(%f) != third(%f)", first.PctOfTop10MostSimilarWithDupeProbAbove33pct, third.PctOfTop10MostSimilarWithDupeProbAbove33pct)
	}

	//PctOfTop10MostSimilarWithDupeProbAbove50pct
	if !compareFloat(first.PctOfTop10MostSimilarWithDupeProbAbove50pct, second.PctOfTop10MostSimilarWithDupeProbAbove50pct) {
		return nil, errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove50pct do not match: first(%f) != second(%f)", first.PctOfTop10MostSimilarWithDupeProbAbove50pct, second.PctOfTop10MostSimilarWithDupeProbAbove50pct)
	}
	if !compareFloat(first.PctOfTop10MostSimilarWithDupeProbAbove50pct, third.PctOfTop10MostSimilarWithDupeProbAbove50pct) {
		return nil, errors.Errorf("PctOfTop10MostSimilarWithDupeProbAbove50pct do not match: first(%f) != third(%f)", first.PctOfTop10MostSimilarWithDupeProbAbove50pct, third.PctOfTop10MostSimilarWithDupeProbAbove50pct)
	}

	// RarenessScoresTableJSONCompressedB64
	if first.RarenessScoresTableJSONCompressedB64 == "" {
		return nil, errors.New("rarenessScoresTableJSONCompressedB64 empty")
	}

	// RarenessScoresTableJSONCompressedB64
	if second.RarenessScoresTableJSONCompressedB64 == "" {
		return nil, errors.New("second rarenessScoresTableJSONCompressedB64 empty")
	}

	if third.RarenessScoresTableJSONCompressedB64 == "" {
		return nil, errors.Errorf("third rarenessScoresTableJSONCompressedB64 empty")
	}
	// //RarenessScores
	// if err := CompareRarenessScores(first.RarenessScores, second.RarenessScores); err != nil {
	// 	return nil, errors.Errorf("first and second RarenessScores do not match: %w", err)
	// }
	// if err := CompareRarenessScores(first.RarenessScores, third.RarenessScores); err != nil {
	// 	return nil, errors.Errorf("first and third RarenessScores do not match: %w", err)
	// }

	//InternetRareness
	// if err := CompareInternetRareness(first.InternetRareness, second.InternetRareness); err != nil {
	// 	return nil, errors.Errorf("first couple of InternetRareness do not match: %w", err)
	// }
	// if err := CompareInternetRareness(first.InternetRareness, third.InternetRareness); err != nil {
	// 	return nil, errors.Errorf("second couple of InternetRareness do not third: %w", err)
	// }
	finalResult.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64 = combineString(
		first.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64,
		second.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64,
		third.InternetRareness.RareOnInternetSummaryTableAsJSONCompressedB64,
	)
	finalResult.InternetRareness.RareOnInternetGraphJSONCompressedB64 = combineString(
		first.InternetRareness.RareOnInternetGraphJSONCompressedB64,
		second.InternetRareness.RareOnInternetGraphJSONCompressedB64,
		third.InternetRareness.RareOnInternetGraphJSONCompressedB64,
	)
	finalResult.InternetRareness.AlternativeRareOnInternetDictAsJSONCompressedB64 = combineString(
		first.InternetRareness.AlternativeRareOnInternetDictAsJSONCompressedB64,
		second.InternetRareness.AlternativeRareOnInternetDictAsJSONCompressedB64,
		third.InternetRareness.AlternativeRareOnInternetDictAsJSONCompressedB64,
	)
	finalResult.InternetRareness.MinNumberOfExactMatchesInPage = combineUint32(
		first.InternetRareness.MinNumberOfExactMatchesInPage,
		second.InternetRareness.MinNumberOfExactMatchesInPage,
		third.InternetRareness.MinNumberOfExactMatchesInPage,
	)
	finalResult.InternetRareness.EarliestAvailableDateOfInternetResults = combineString(
		first.InternetRareness.EarliestAvailableDateOfInternetResults,
		second.InternetRareness.EarliestAvailableDateOfInternetResults,
		third.InternetRareness.EarliestAvailableDateOfInternetResults,
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
	if !compareDoubles(first.ImageFingerprintOfCandidateImageFile, second.ImageFingerprintOfCandidateImageFile) {
		return nil, errors.New("first couple of image_fingerprint_of_candidate_image_file nsfw score do not match")
	}
	if !compareDoubles(first.ImageFingerprintOfCandidateImageFile, third.ImageFingerprintOfCandidateImageFile) {
		return nil, errors.New("second couple of image_fingerprint_of_candidate_image_file nsfw score do not match")
	}

	//hash_of_candidate_image_file
	if !reflect.DeepEqual(first.HashOfCandidateImageFile, second.HashOfCandidateImageFile) {
		return nil, errors.New("first couple of hash_of_candidate_image_filedo not match")
	}
	if !reflect.DeepEqual(first.HashOfCandidateImageFile, third.HashOfCandidateImageFile) {
		return nil, errors.New("second couple of hash_of_candidate_image_file do not match")
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
