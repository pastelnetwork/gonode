package mock

import (
	"encoding/json"
)

func getDDServerResponse() ([]byte, error) {
	reply := &ImageRarenessScoreReply{
		Block:                                "160647",
		DupeDetectionSystemVersion:           "1.0",
		IsLikelyDupe:                         false,
		IsRareOnInternet:                     true,
		OpenNsfwScore:                        0.4,
		ImageFingerprintOfCandidateImageFile: []float32{3.2, 2.5, 6.7, 0.4},
		PerceptualHashOverlapCount:           1,
		RarenessScores: &RarenessScores{
			CombinedRarenessScore:         0.443,
			XgboostPredictedRarenessScore: 0.532,
			NnPredictedRarenessScore:      0.344,
			OverallAverageRarenessScore:   0.429,
		},
		AlternativeNsfwScores: &AltNsfwScores{
			Drawings: 0.34,
			Hentai:   0.23,
			Neutral:  0.89,
			Porn:     0.13,
			Sexy:     0.18,
		},
		InternetRareness: &InternetRareness{
			MatchesFoundOnFirstPage: 0,
			NumberOfPagesOfResults:  0,
			UrlOfFirstMatchInPage:   "",
		},
		PerceptualImageHashes: &PerceptualImageHashes{
			PdqHash:        "abc",
			PerceptualHash: "bcd",
			AverageHash:    "tyf",
			DifferenceHash: "dsc",
			NeuralhashHash: "uip",
		},
		Maxes: &Maxes{
			PearsonMax:           0.94,
			SpearmanMax:          0.23,
			KendallMax:           0.33,
			HoeffdingMax:         0.65,
			MutualInformationMax: 0.22,
			HsicMax:              0.13,
			XgbimportanceMax:     0.67,
		},
		Percentile: &Percentile{
			PearsonTop_1BpsPercentile:             78.2,
			SpearmanTop_1BpsPercentile:            22.3,
			KendallTop_1BpsPercentile:             76.3,
			HoeffdingTop_10BpsPercentile:          45.3,
			MutualInformationTop_100BpsPercentile: 33.23,
			HsicTop_100BpsPercentile:              87.33,
			XgbimportanceTop_100BpsPercentile:     20.3,
		},
		FingerprintsStat: &FingerprintsStat{
			NumberOfFingerprintsRequiringFurtherTesting_1: 0,
			NumberOfFingerprintsRequiringFurtherTesting_2: 0,
			NumberOfFingerprintsRequiringFurtherTesting_3: 0,
			NumberOfFingerprintsRequiringFurtherTesting_4: 0,
			NumberOfFingerprintsRequiringFurtherTesting_5: 0,
			NumberOfFingerprintsRequiringFurtherTesting_6: 0,
			NumberOfFingerprintsOfSuspectedDupes:          0,
		},
	}

	return json.Marshal(reply)
}

type ImageRarenessScoreReply struct {
	Block                                string                 `json:"block,omitempty"`
	Principal                            string                 `json:"principal,omitempty"`
	DupeDetectionSystemVersion           string                 `json:"dupe_detection_system_version,omitempty"`
	IsLikelyDupe                         bool                   `json:"is_likely_dupe,omitempty"`
	IsRareOnInternet                     bool                   `json:"is_rare_on_internet,omitempty"`
	RarenessScores                       *RarenessScores        `json:"rareness_scores,omitempty"`
	InternetRareness                     *InternetRareness      `json:"internet_rareness,omitempty"`
	OpenNsfwScore                        float32                `json:"open_nsfw_score,omitempty"`
	AlternativeNsfwScores                *AltNsfwScores         `json:"alternative_nsfw_scores,omitempty"`
	ImageFingerprintOfCandidateImageFile []float32              `json:"image_fingerprint_of_candidate_image_file,omitempty"`
	FingerprintsStat                     *FingerprintsStat      `json:"fingerprints_stat,omitempty"`
	HashOfCandidateImageFile             string                 `json:"hash_of_candidate_image_file,omitempty"`
	PerceptualImageHashes                *PerceptualImageHashes `json:"perceptual_image_hashes,omitempty"`
	PerceptualHashOverlapCount           uint32                 `json:"perceptual_hash_overlap_count,omitempty"`
	Maxes                                *Maxes                 `json:"maxes,omitempty"`
	Percentile                           *Percentile            `json:"percentile,omitempty"`
}

type RarenessScores struct {
	CombinedRarenessScore         float32 `json:"combined_rareness_score,omitempty"`
	XgboostPredictedRarenessScore float32 `json:"xgboost_predicted_rareness_score,omitempty"`
	NnPredictedRarenessScore      float32 `json:"nn_predicted_rareness_score,omitempty"`
	OverallAverageRarenessScore   float32 `json:"overall_average_rareness_score,omitempty"`
}

type InternetRareness struct {
	MatchesFoundOnFirstPage uint32 `json:"matches_found_on_first_page,omitempty"`
	NumberOfPagesOfResults  uint32 `json:"number_of_pages_of_results,omitempty"`
	UrlOfFirstMatchInPage   string `json:"url_of_first_match_in_page,omitempty"`
}

type AltNsfwScores struct {
	Drawings float32 `json:"drawings,omitempty"`
	Hentai   float32 `json:"hentai,omitempty"`
	Neutral  float32 `json:"neutral,omitempty"`
	Porn     float32 `json:"porn,omitempty"`
	Sexy     float32 `json:"sexy,omitempty"`
}

type PerceptualImageHashes struct {
	PdqHash        string `json:"pdq_hash,omitempty"`
	PerceptualHash string `json:"perceptual_hash,omitempty"`
	AverageHash    string `json:"average_hash,omitempty"`
	DifferenceHash string `json:"difference_hash,omitempty"`
	NeuralhashHash string `json:"neuralhash_hash,omitempty"`
}

type FingerprintsStat struct {
	NumberOfFingerprintsRequiringFurtherTesting_1 uint32 `json:"number_of_fingerprints_requiring_further_testing_1,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_2 uint32 `json:"number_of_fingerprints_requiring_further_testing_2,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_3 uint32 `json:"number_of_fingerprints_requiring_further_testing_3,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_4 uint32 `json:"number_of_fingerprints_requiring_further_testing_4,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_5 uint32 `json:"number_of_fingerprints_requiring_further_testing_5,omitempty"`
	NumberOfFingerprintsRequiringFurtherTesting_6 uint32 `json:"number_of_fingerprints_requiring_further_testing_6,omitempty"`
	NumberOfFingerprintsOfSuspectedDupes          uint32 `json:"number_of_fingerprints_of_suspected_dupes,omitempty"`
}

type Maxes struct {
	PearsonMax           float32 `json:"pearson_max,omitempty"`
	SpearmanMax          float32 `json:"spearman_max,omitempty"`
	KendallMax           float32 `json:"kendall_max,omitempty"`
	HoeffdingMax         float32 `json:"hoeffding_max,omitempty"`
	MutualInformationMax float32 `json:"mutual_information_max,omitempty"`
	HsicMax              float32 `json:"hsic_max,omitempty"`
	XgbimportanceMax     float32 `json:"xgbimportance_max,omitempty"`
}

type Percentile struct {
	PearsonTop_1BpsPercentile             float32 `json:"pearson_top_1_bps_percentile,omitempty"`
	SpearmanTop_1BpsPercentile            float32 `json:"spearman_top_1_bps_percentile,omitempty"`
	KendallTop_1BpsPercentile             float32 `json:"kendall_top_1_bps_percentile,omitempty"`
	HoeffdingTop_10BpsPercentile          float32 `json:"hoeffding_top_10_bps_percentile,omitempty"`
	MutualInformationTop_100BpsPercentile float32 `json:"mutual_information_top_100_bps_percentile,omitempty"`
	HsicTop_100BpsPercentile              float32 `json:"hsic_top_100_bps_percentile,omitempty"`
	XgbimportanceTop_100BpsPercentile     float32 `json:"xgbimportance_top_100_bps_percentile,omitempty"`
}
