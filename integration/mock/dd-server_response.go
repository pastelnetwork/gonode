package mock

import (
	"time"

	json "github.com/json-iterator/go"
)

func getDDServerResponse() ([]byte, error) {
	reply := &ImageRarenessScoreReply{
		BlockHash:   "39720f0d10c7a877f2bd6edcacbfe18e1333fe77a3032f6d394a1df742f42b37",
		BlockHeight: "160647",

		TimestampOfRequest: time.Now().UTC().Format("YYYY-MM-DD hh:mm:ss"),
		SubmitterPastelID:  "SubmitterPastelID",
		SN1PastelID:        "Supernode1PastelID",
		SN2PastelID:        "Supernode2PastelID",
		SN3PastelID:        "Supernode3PastelID",

		IsOpenAPIRequest: false,

		DupeDetectionSystemVersion: "1.0",
		IsLikelyDupe:               false,
		IsRareOnInternet:           true,
		OverallRarenessScore:       0.9,

		PctOfTop10MostSimilarWithDupeProbAbove25pct: 0.3,
		PctOfTop10MostSimilarWithDupeProbAbove33pct: 0.2,
		PctOfTop10MostSimilarWithDupeProbAbove50pct: 0.1,

		RarenessScoresTableJSONCompressedB64: "eyJQYXN0ZWwgcmFyZW5lc3Mgc2NvcmUgdGFibGUiOiJ0YWJsZSJ9Cg==",

		InternetRareness: &InternetRareness{
			RareOnInternetSummaryTableAsJSONCompressedB64:    "eyJQYXN0ZWwgcmFyZSBvbiBpbnRlcm5ldCBzdW1tYXJ5IHRhYmxlIjoidGFibGUifQo=",
			RareOnInternetGraphJSONCompressedB64:             "eyJQYXN0ZWwgcmFyZSBvbiBpbnRlcm5ldCBncmFwaCI6ImdyYXBoIn0K",
			AlternativeRareOnInternetDictAsJSONCompressedB64: "eyJQYXN0ZWwgYWx0ZXJuYXRpdmUgcmFyZSBvbiBpbnRlcm5ldCBkaWN0IjoiZGljdCJ9Cg==",
			MinNumberOfExactMatchesInPage:                    3,
			EarliestAvailableDateOfInternetResults:           time.Now().UTC().Format("2022-01-01 hh:mm:ss"),
		},
		OpenNSFWScore: 0.4,

		AlternativeNSFWScores: &AltNsfwScores{
			Drawings: 0.1,
			Hentai:   0.2,
			Neutral:  0.3,
			Porn:     0.4,
			Sexy:     0.5,
		},

		ImageFingerprintOfCandidateImageFile: []float64{3.2, 2.5, 6.7, 0.4},

		HashOfCandidateImageFile: "53a66ab645b9f56b3087760185035abe122aeb808fdad50a037acb150b5528bd",
	}

	return json.Marshal(reply)
}

type ImageRarenessScoreReply struct {
	BlockHash   string `json:"pastel_block_hash_when_request_submitted"`
	BlockHeight string `json:"pastel_block_height_when_request_submitted"`

	TimestampOfRequest string `json:"utc_timestamp_when_request_submitted"`
	SubmitterPastelID  string `json:"pastel_id_of_submitter"`
	SN1PastelID        string `json:"pastel_id_of_registering_supernode_1"`
	SN2PastelID        string `json:"pastel_id_of_registering_supernode_2"`
	SN3PastelID        string `json:"pastel_id_of_registering_supernode_3"`

	IsOpenAPIRequest bool `json:"is_pastel_openapi_request"`

	DupeDetectionSystemVersion string `json:"dupe_detection_system_version"`

	IsLikelyDupe         bool    `json:"is_likely_dupe"`
	IsRareOnInternet     bool    `json:"is_rare_on_internet"`
	OverallRarenessScore float32 `json:"overall_rareness_score "`

	PctOfTop10MostSimilarWithDupeProbAbove25pct float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_25pct"`
	PctOfTop10MostSimilarWithDupeProbAbove33pct float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_33pct"`
	PctOfTop10MostSimilarWithDupeProbAbove50pct float32 `json:"pct_of_top_10_most_similar_with_dupe_prob_above_50pct"`

	RarenessScoresTableJSONCompressedB64 string `json:"rareness_scores_table_json_compressed_b64"`

	InternetRareness *InternetRareness `json:"internet_rareness"`

	OpenNSFWScore float32 `json:"open_nsfw_score"`

	AlternativeNSFWScores *AltNsfwScores `json:"alternative_nsfw_scores"`

	ImageFingerprintOfCandidateImageFile []float64 `json:"image_fingerprint_of_candidate_image_file"`

	HashOfCandidateImageFile string `json:"hash_of_candidate_image_file"`
}

// type RarenessScores struct {
// 	CombinedRarenessScore         float32 `json:"combined_rareness_score,omitempty"`
// 	XgboostPredictedRarenessScore float32 `json:"xgboost_predicted_rareness_score,omitempty"`
// 	NnPredictedRarenessScore      float32 `json:"nn_predicted_rareness_score,omitempty"`
// 	OverallAverageRarenessScore   float32 `json:"overall_average_rareness_score,omitempty"`
// }

// type InternetRareness struct {
// 	MatchesFoundOnFirstPage uint32 `json:"matches_found_on_first_page,omitempty"`
// 	NumberOfPagesOfResults  uint32 `json:"number_of_pages_of_results,omitempty"`
// 	URLOfFirstMatchInPage   string `json:"url_of_first_match_in_page,omitempty"`
// }

type AltNsfwScores struct {
	Drawings float32 `json:"drawings,omitempty"`
	Hentai   float32 `json:"hentai,omitempty"`
	Neutral  float32 `json:"neutral,omitempty"`
	Porn     float32 `json:"porn,omitempty"`
	Sexy     float32 `json:"sexy,omitempty"`
}

type InternetRareness struct {
	RareOnInternetSummaryTableAsJSONCompressedB64    string `json:"rare_on_internet_summary_table_as_json_compressed_b64"`
	RareOnInternetGraphJSONCompressedB64             string `json:"rare_on_internet_graph_json_compressed_b64"`
	AlternativeRareOnInternetDictAsJSONCompressedB64 string `json:"alternative_rare_on_internet_dict_as_json_compressed_b64"`
	MinNumberOfExactMatchesInPage                    uint32 `json:"min_number_of_exact_matches_in_page"`
	EarliestAvailableDateOfInternetResults           string `json:"earliest_available_date_of_internet_results"`
}

// type PerceptualImageHashes struct {
// 	PDQHash        string `json:"pdq_hash,omitempty"`
// 	PerceptualHash string `json:"perceptual_hash,omitempty"`
// 	AverageHash    string `json:"average_hash,omitempty"`
// 	DifferenceHash string `json:"difference_hash,omitempty"`
// 	NeuralHash     string `json:"neuralhash_hash,omitempty"`
// }

// type FingerprintsStat struct {
// 	NumberOfFingerprintsRequiringFurtherTesting1 uint32 `json:"number_of_fingerprints_requiring_further_testing_1,omitempty"`
// 	NumberOfFingerprintsRequiringFurtherTesting2 uint32 `json:"number_of_fingerprints_requiring_further_testing_2,omitempty"`
// 	NumberOfFingerprintsRequiringFurtherTesting3 uint32 `json:"number_of_fingerprints_requiring_further_testing_3,omitempty"`
// 	NumberOfFingerprintsRequiringFurtherTesting4 uint32 `json:"number_of_fingerprints_requiring_further_testing_4,omitempty"`
// 	NumberOfFingerprintsRequiringFurtherTesting5 uint32 `json:"number_of_fingerprints_requiring_further_testing_5,omitempty"`
// 	NumberOfFingerprintsRequiringFurtherTesting6 uint32 `json:"number_of_fingerprints_requiring_further_testing_6,omitempty"`
// 	NumberOfFingerprintsOfSuspectedDupes         uint32 `json:"number_of_fingerprints_of_suspected_dupes,omitempty"`
// }

// type Maxes struct {
// 	PearsonMax           float32 `json:"pearson_max,omitempty"`
// 	SpearmanMax          float32 `json:"spearman_max,omitempty"`
// 	KendallMax           float32 `json:"kendall_max,omitempty"`
// 	HoeffdingMax         float32 `json:"hoeffding_max,omitempty"`
// 	MutualInformationMax float32 `json:"mutual_information_max,omitempty"`
// 	HsicMax              float32 `json:"hsic_max,omitempty"`
// 	XgbimportanceMax     float32 `json:"xgbimportance_max,omitempty"`
// }

// type Percentile struct {
// 	PearsonTop1BpsPercentile             float32 `json:"pearson_top_1_bps_percentile,omitempty"`
// 	SpearmanTop1BpsPercentile            float32 `json:"spearman_top_1_bps_percentile,omitempty"`
// 	KendallTop1BpsPercentile             float32 `json:"kendall_top_1_bps_percentile,omitempty"`
// 	HoeffdingTop10BpsPercentile          float32 `json:"hoeffding_top_10_bps_percentile,omitempty"`
// 	MutualInformationTop100BpsPercentile float32 `json:"mutual_information_top_100_bps_percentile,omitempty"`
// 	HsicTop100BpsPercentile              float32 `json:"hsic_top_100_bps_percentile,omitempty"`
// 	XgbimportanceTop100BpsPercentile     float32 `json:"xgbimportance_top_100_bps_percentile,omitempty"`
// }
