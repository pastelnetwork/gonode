//go:generate mockery --name=DDServerClient

package ddclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/dupedetection"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix = "ddClient"
)

// DDServerClient contains methods for request services from dd-server service.
type DDServerClient interface {
	// ImageRarenessScore returns rareness score of image
	ImageRarenessScore(ctx context.Context, img []byte, format string, blockHash string, blockHeight string, timestamp string, pastelID string, sn_1 string, sn_2 string, sn_3 string, openapi_request bool, open_api_subset_id string) (*pastel.DDAndFingerprints, error)
}

type ddServerClientImpl struct {
	config *Config
}

func randID() string {
	id := uuid.New().String()
	return id[0:8]
}

func writeFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0750)
}

func createInputDDFile(base string, data []byte, format string) (string, error) {
	fileName := randID() + "." + format
	filePath := filepath.Join(base, fileName)

	err := writeFile(filePath, data)

	if err != nil {
		return "", errors.Errorf("write data to %s: %w", filePath, err)
	}

	return filePath, nil
}

// call ddserver's ImageRarenessScore service
func (ddClient *ddServerClientImpl) callImageRarenessScore(ctx context.Context, client pb.DupeDetectionServerClient, img []byte, format string, blockHash string, blockHeight string, timestamp string, pastelID string, sn_1 string, sn_2 string, sn_3 string, openapi_request bool, open_api_subset_id string) (*pastel.DDAndFingerprints, error) {
	if img == nil {
		return nil, errors.Errorf("invalid data")
	}

	inputPath, err := createInputDDFile(ddClient.config.DDFilesDir, img, format)
	if err != nil {
		return nil, errors.Errorf("create input file: %w", err)
	}

	req := pb.RarenessScoreRequest{
		Path:                   inputPath,
		BlockHash:              blockHash,
		BlockHeight:            blockHeight,
		Timestamp:              timestamp,
		PastelId:               pastelID,
		RegisteringSn1:         sn_1,
		RegisteringSn2:         sn_2,
		RegisteringSn3:         sn_3,
		IsPastelOpenapiRequest: openapi_request,
		OpenApiSubsetIdString:  open_api_subset_id,
	}

	// remove file after use
	defer os.Remove(inputPath)

	// 2/9/2022, we don't have this implemented yet, so this would fail
	res, err := client.ImageRarenessScore(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("send request: %w", err)
	}

	/*
		{
		  "block": string           // block_hash from NFT ticket
		  "principal": string       // PastelID of the author from NFT ticket

		  "dupe_detection_system_version": string,

		  "is_likely_dupe": bool,
		  "is_rare_on_internet": bool,

		  "rareness_scores": {
		      "combined_rareness_score": float,                      // 0 to 1
		      "xgboost_predicted_rareness_score": float,            // 0 to 1
		      "nn_predicted_rareness_score": float,                  // 0 to 1
		      "overall_average_rareness_score": float,              // 0 to 1
		  },

		  "internet_rareness": {
		      "matches_found_on_first_page": uint32,
		      "number_of_pages_of_results": uint32,
		      "url_of_first_match_in_page": string,
		  },

		  "open_nsfw_score": float,                                 // 0 to 1
		  "alternative_nsfw_scores": {
		      "drawings": float,                                    // 0 to 1
		      "hentai": float,                                      // 0 to 1
		      "neutral": float,                                      // 0 to 1
		      "porn": float,                                        // 0 to 1
		      "sexy": float,                                        // 0 to 1
		  },

		  "image_fingerprint_of_candidate_image_file": [float],           // array with fingerprints values
		  "fingerprints_stat": {
		      "number_of_fingerprints_requiring_further_testing_1": uint32,
		      "number_of_fingerprints_requiring_further_testing_2": uint32,
		      "number_of_fingerprints_requiring_further_testing_3": uint32,
		      "number_of_fingerprints_requiring_further_testing_4": uint32,
		      "number_of_fingerprints_requiring_further_testing_5": uint32,
		      "number_of_fingerprints_requiring_further_testing_6":  uint32,
		      "number_of_fingerprints_of_suspected_dupes": uint32,
		  },

		  "hash_of_candidate_image_file": string,
		  "perceptual_image_hashes": {
		      "pdq_hash": string,
		      "perceptual_hash": string,
		      "average_hash": string,
		      "difference_hash": string,
		      "neuralhash_hash": string,
		  },
		  "perceptual_hash_overlap_count": uint32,

		  "maxes": {
		      "pearson_max": float,                                  // 0 to 1
		      "spearman_max": float,                                // 0 to 1
		      "kendall_max": float,                                  // 0 to 1
		      "hoeffding_max": float,                                // 0 to 1
		      "mutual_information_max": float,                      // 0 to 1
		      "hsic_max": float,                                    // 0 to 1
		      "xgbimportance_max": float,                            // 0 to 1
		  },

		  "percentile": {
		      "pearson_top_1_bps_percentile": float,                // 0 to 1
		      "spearman_top_1_bps_percentile": float,                // 0 to 1
		      "kendall_top_1_bps_percentile": float,                // 0 to 1
		      "hoeffding_top_10_bps_percentile": float,              // 0 to 1
		      "mutual_information_top_100_bps_percentile": float,    // 0 to 1
		      "hsic_top_100_bps_percentile": float,                  // 0 to 1
		      "xgbimportance_top_100_bps_percentile": float,        // 0 to 1
		  },
		}
	*/
	output := &pastel.DDAndFingerprints{
		Block:                      res.Block,
		Principal:                  res.Principal,
		DupeDetectionSystemVersion: res.DupeDetectionSystemVersion,

		IsLikelyDupe:     res.IsLikelyDupe,
		IsRareOnInternet: res.IsRareOnInternet,

		RarenessScores: &pastel.RarenessScores{
			CombinedRarenessScore:         res.RarenessScores.CombinedRarenessScore,
			XgboostPredictedRarenessScore: res.RarenessScores.XgboostPredictedRarenessScore,
			NnPredictedRarenessScore:      res.RarenessScores.NnPredictedRarenessScore,
			OverallAverageRarenessScore:   res.RarenessScores.OverallAverageRarenessScore,
		},
		InternetRareness: &pastel.InternetRareness{
			MatchesFoundOnFirstPage: res.InternetRareness.MatchesFoundOnFirstPage,
			NumberOfPagesOfResults:  res.InternetRareness.NumberOfPagesOfResults,
			URLOfFirstMatchInPage:   res.InternetRareness.UrlOfFirstMatchInPage,
		},

		OpenNSFWScore: res.OpenNsfwScore,
		AlternativeNSFWScores: &pastel.AlternativeNSFWScores{
			Drawings: res.AlternativeNsfwScores.Drawings,
			Hentai:   res.AlternativeNsfwScores.Hentai,
			Neutral:  res.AlternativeNsfwScores.Neutral,
			Porn:     res.AlternativeNsfwScores.Porn,
			Sexy:     res.AlternativeNsfwScores.Sexy,
		},

		ImageFingerprintOfCandidateImageFile: res.ImageFingerprintOfCandidateImageFile,
		FingerprintsStat: &pastel.FingerprintsStat{
			NumberOfFingerprintsRequiringFurtherTesting1: res.FingerprintsStat.NumberOfFingerprintsRequiringFurtherTesting_1,
			NumberOfFingerprintsRequiringFurtherTesting2: res.FingerprintsStat.NumberOfFingerprintsRequiringFurtherTesting_2,
			NumberOfFingerprintsRequiringFurtherTesting3: res.FingerprintsStat.NumberOfFingerprintsRequiringFurtherTesting_3,
			NumberOfFingerprintsRequiringFurtherTesting4: res.FingerprintsStat.NumberOfFingerprintsRequiringFurtherTesting_4,
			NumberOfFingerprintsRequiringFurtherTesting5: res.FingerprintsStat.NumberOfFingerprintsRequiringFurtherTesting_5,
			NumberOfFingerprintsRequiringFurtherTesting6: res.FingerprintsStat.NumberOfFingerprintsRequiringFurtherTesting_6,
			NumberOfFingerprintsOfSuspectedDupes:         res.FingerprintsStat.NumberOfFingerprintsOfSuspectedDupes,
		},

		HashOfCandidateImageFile: res.HashOfCandidateImageFile,
		PerceptualImageHashes: &pastel.PerceptualImageHashes{
			PDQHash:        res.PerceptualImageHashes.PdqHash,
			PerceptualHash: res.PerceptualImageHashes.PerceptualHash,
			AverageHash:    res.PerceptualImageHashes.AverageHash,
			DifferenceHash: res.PerceptualImageHashes.DifferenceHash,
			NeuralHash:     res.PerceptualImageHashes.NeuralhashHash,
		},
		PerceptualHashOverlapCount: res.PerceptualHashOverlapCount,

		Maxes: &pastel.Maxes{
			PearsonMax:           res.Maxes.PearsonMax,
			SpearmanMax:          res.Maxes.SpearmanMax,
			KendallMax:           res.Maxes.KendallMax,
			HoeffdingMax:         res.Maxes.HoeffdingMax,
			MutualInformationMax: res.Maxes.MutualInformationMax,
			HsicMax:              res.Maxes.HsicMax,
			XgbimportanceMax:     res.Maxes.XgbimportanceMax,
		},
		Percentile: &pastel.Percentile{
			PearsonTop1BpsPercentile:             res.Percentile.PearsonTop_1BpsPercentile,
			SpearmanTop1BpsPercentile:            res.Percentile.SpearmanTop_1BpsPercentile,
			KendallTop1BpsPercentile:             res.Percentile.KendallTop_1BpsPercentile,
			HoeffdingTop10BpsPercentile:          res.Percentile.HoeffdingTop_10BpsPercentile,
			MutualInformationTop100BpsPercentile: res.Percentile.MutualInformationTop_100BpsPercentile,
			HsicTop100BpsPercentile:              res.Percentile.HsicTop_100BpsPercentile,
			XgbimportanceTop100BpsPercentile:     res.Percentile.XgbimportanceTop_100BpsPercentile,
		},
	}

	return output, nil
}

// ImageRarenessScore call ddserver to calculate scores
func (ddClient *ddServerClientImpl) ImageRarenessScore(ctx context.Context, img []byte, format string, blockHash string, blockHeight string, timestamp string, pastelID string, sn_1 string, sn_2 string, sn_3 string, openapi_request bool, open_api_subset_id string) (*pastel.DDAndFingerprints, error) {
	ctx = ddClient.contextWithLogPrefix(ctx)

	baseClient := NewClient()
	ddServerAddress := fmt.Sprint(ddClient.config.Host, ":", ddClient.config.Port)
	conn, err := baseClient.Connect(ctx, ddServerAddress)
	if err != nil {
		return nil, errors.Errorf("connect to dd-server  %w", err)
	}

	defer conn.Close()
	client := pb.NewDupeDetectionServerClient(conn)

	return ddClient.callImageRarenessScore(ctx, client, img, format, blockHash, blockHeight, timestamp, pastelID, sn_1, sn_2, sn_3, openapi_request, open_api_subset_id)
}

func (ddClient *ddServerClientImpl) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, logPrefix)
}

// NewDDServerClient returns a DDServerClient
func NewDDServerClient(config *Config) DDServerClient {
	return &ddServerClientImpl{
		config: config,
	}
}
