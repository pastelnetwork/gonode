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
)

const (
	logPrefix = "ddClient"
)

// DupeDetection is the dupe detection result which will be sent to the caller
type DupeDetection struct {
	DupeDetectionSystemVer                       string              `json:"dupe_detection_system_version"`
	ImageHash                                    string              `json:"hash_of_candidate_image_file"`
	PastelRarenessScore                          float32             `json:"overall_average_rareness_score"`
	IsLikelyDupe                                 bool                `json:"is_likely_dupe"`
	IsRareOnInternet                             bool                `json:"is_rare_on_internet"`
	MatchesFoundOnFirstPage                      uint32              `json:"matches_found_on_first_page"`
	NumberOfResultPages                          uint32              `json:"number_of_pages_of_results"`
	FirstMatchURL                                string              `json:"url_of_first_match_in_page"`
	OpenNSFWScore                                float32             `json:"open_nsfw_score"`
	AlternateNSFWScores                          AlternateNSFWScores `json:"alternative_nsfw_scores"`
	ImageHashes                                  ImageHashes         `json:"image_hashes"`
	Fingerprints                                 []float32           `json:"image_fingerprint_of_candidate_image_file"`
	PerceptualHashOverlapCount                   uint32              `json:"perceptual_hash_overlap_count"`
	NumberOfFingerprintsRequiringFurtherTesting1 uint32              `json:"number_of_fingerprints_requiring_further_testing_1"`
	NumberOfFingerprintsRequiringFurtherTesting2 uint32              `json:"number_of_fingerprints_requiring_further_testing_2"`
	NumberOfFingerprintsRequiringFurtherTesting3 uint32              `json:"number_of_fingerprints_requiring_further_testing_3"`
	NumberOfFingerprintsRequiringFurtherTesting4 uint32              `json:"number_of_fingerprints_requiring_further_testing_4"`
	NumberOfFingerprintsRequiringFurtherTesting5 uint32              `json:"number_of_fingerprints_requiring_further_testing_5"`
	NumberOfFingerprintsRequiringFurtherTesting6 uint32              `json:"number_of_fingerprints_requiring_further_testing_6"`
	NumberOfFingerprintsOfSuspectedDupes         uint32              `json:"number_of_fingerprints_of_suspected_dupes"`
	PearsonMax                                   float32             `json:"pearson_max"`
	SpearmanMax                                  float32             `json:"spearman_max"`
	KendallMax                                   float32             `json:"kendall_max"`
	HoeffdingMax                                 float32             `json:"hoeffding_max"`
	MutualInformationMax                         float32             `json:"mutual_information_max"`
	HsicMax                                      float32             `json:"hsic_max"`
	XgbimportanceMax                             float32             `json:"xgbimportance_max"`
	PearsonTop1BpsPercentile                     float32             `json:"pearson_top_1_bps_percentile"`
	SpearmanTop1BpsPercentile                    float32             `json:"spearman_top_1_bps_percentile"`
	KendallTop1BpsPercentile                     float32             `json:"kendall_top_1_bps_percentile"`
	HoeffdingTop10BpsPercentile                  float32             `json:"hoeffding_top_10_bps_percentile"`
	MutualInformationTop100BpsPercentile         float32             `json:"mutual_information_top_100_bps_percentile"`
	HsicTop100BpsPercentile                      float32             `json:"hsic_top_100_bps_percentile"`
	XgbimportanceTop100BpsPercentile             float32             `json:"xgbimportance_top_100_bps_percentile"`
	CombinedRarenessScore                        float32             `json:"combined_rareness_score"`
	XgboostPredictedRarenessScore                float32             `json:"xgboost_predicted_rareness_score"`
	NnPredictedRarenessScore                     float32             `json:"nn_predicted_rareness_score"`
	OverallAverageRarenessScore                  float32             `json:"overall_average_rareness_score"`
	NumberOfPagesOfResults                       uint32              `json:"number_of_pages_of_results"`
	UrlOfFirstMatchInPage                        string              `json:"url_of_first_match_in_page"`
}

// AlternateNSFWScores represents alternate NSFW scores in the output of dupe detection service
type AlternateNSFWScores struct {
	Drawings float32 `json:"drawings"`
	Hentai   float32 `json:"hentai"`
	Neutral  float32 `json:"neutral"`
	Porn     float32 `json:"porn"`
	Sexy     float32 `json:"sexy"`
}

// ImageHashes represents image hashes in the output of dupe detection service
type ImageHashes struct {
	PDQHash        string `json:"pdq_hash"`
	PerceptualHash string `json:"perceptual_hash"`
	AverageHash    string `json:"average_hash"`
	DifferenceHash string `json:"difference_hash"`
	NeuralHash     string `json:"neuralhash_hash"`
}

// DDServerClient contains methods for request services from dd-server service.
type DDServerClient interface {
	// ImageRarenessScore returns rareness score of image
	ImageRarenessScore(ctx context.Context, img []byte, format string) (*DupeDetection, error)
}

type ddServerClientImpl struct {
	config *Config
}

func randID() string {
	id := uuid.NewString()
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

// call ddsever's ImageRarenessScore service
func (ddClient *ddServerClientImpl) callImageRarenessScore(ctx context.Context, client pb.DupeDetectionServerClient, img []byte, format string) (*DupeDetection, error) {
	if img == nil {
		return nil, errors.Errorf("invalid data")
	}

	inputPath, err := createInputDDFile(ddClient.config.DDFilesDir, img, format)
	if err != nil {
		return nil, errors.Errorf("create input file: %w", err)
	}

	req := pb.RarenessScoreRequest{
		Path: inputPath,
	}

	// remove file after use
	defer os.Remove(inputPath)

	res, err := client.ImageRarenessScore(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("send request: %w", err)
	}

	/*
			DupeDetectionSystemVersion                    string                 `protobuf:"bytes,1,opt,name=dupe_detection_system_version,json=dupeDetectionSystemVersion,proto3" json:"dupe_detection_system_version,omitempty"`
		HashOfCandidateImageFile                      string                 `protobuf:"bytes,2,opt,name=hash_of_candidate_image_file,json=hashOfCandidateImageFile,proto3" json:"hash_of_candidate_image_file,omitempty"`
		IsLikelyDupe                                  bool                   `protobuf:"varint,3,opt,name=is_likely_dupe,json=isLikelyDupe,proto3" json:"is_likely_dupe,omitempty"`
		PerceptualHashOverlapCount                    uint32                 `protobuf:"varint,4,opt,name=perceptual_hash_overlap_count,json=perceptualHashOverlapCount,proto3" json:"perceptual_hash_overlap_count,omitempty"`
		NumberOfFingerprintsRequiringFurtherTesting_1 uint32                 `protobuf:"varint,5,opt,name=number_of_fingerprints_requiring_further_testing_1,json=numberOfFingerprintsRequiringFurtherTesting1,proto3" json:"number_of_fingerprints_requiring_further_testing_1,omitempty"`
		NumberOfFingerprintsRequiringFurtherTesting_2 uint32                 `protobuf:"varint,6,opt,name=number_of_fingerprints_requiring_further_testing_2,json=numberOfFingerprintsRequiringFurtherTesting2,proto3" json:"number_of_fingerprints_requiring_further_testing_2,omitempty"`
		NumberOfFingerprintsRequiringFurtherTesting_3 uint32                 `protobuf:"varint,7,opt,name=number_of_fingerprints_requiring_further_testing_3,json=numberOfFingerprintsRequiringFurtherTesting3,proto3" json:"number_of_fingerprints_requiring_further_testing_3,omitempty"`
		NumberOfFingerprintsRequiringFurtherTesting_4 uint32                 `protobuf:"varint,8,opt,name=number_of_fingerprints_requiring_further_testing_4,json=numberOfFingerprintsRequiringFurtherTesting4,proto3" json:"number_of_fingerprints_requiring_further_testing_4,omitempty"`
		NumberOfFingerprintsRequiringFurtherTesting_5 uint32                 `protobuf:"varint,9,opt,name=number_of_fingerprints_requiring_further_testing_5,json=numberOfFingerprintsRequiringFurtherTesting5,proto3" json:"number_of_fingerprints_requiring_further_testing_5,omitempty"`
		NumberOfFingerprintsRequiringFurtherTesting_6 uint32                 `protobuf:"varint,10,opt,name=number_of_fingerprints_requiring_further_testing_6,json=numberOfFingerprintsRequiringFurtherTesting6,proto3" json:"number_of_fingerprints_requiring_further_testing_6,omitempty"`
		NumberOfFingerprintsOfSuspectedDupes          uint32                 `protobuf:"varint,11,opt,name=number_of_fingerprints_of_suspected_dupes,json=numberOfFingerprintsOfSuspectedDupes,proto3" json:"number_of_fingerprints_of_suspected_dupes,omitempty"`
		PearsonMax                                    float32                `protobuf:"fixed32,12,opt,name=pearson_max,json=pearsonMax,proto3" json:"pearson_max,omitempty"`
		SpearmanMax                                   float32                `protobuf:"fixed32,13,opt,name=spearman_max,json=spearmanMax,proto3" json:"spearman_max,omitempty"`
		KendallMax                                    float32                `protobuf:"fixed32,14,opt,name=kendall_max,json=kendallMax,proto3" json:"kendall_max,omitempty"`
		HoeffdingMax                                  float32                `protobuf:"fixed32,15,opt,name=hoeffding_max,json=hoeffdingMax,proto3" json:"hoeffding_max,omitempty"`
		MutualInformationMax                          float32                `protobuf:"fixed32,16,opt,name=mutual_information_max,json=mutualInformationMax,proto3" json:"mutual_information_max,omitempty"`
		HsicMax                                       float32                `protobuf:"fixed32,17,opt,name=hsic_max,json=hsicMax,proto3" json:"hsic_max,omitempty"`
		XgbimportanceMax                              float32                `protobuf:"fixed32,18,opt,name=xgbimportance_max,json=xgbimportanceMax,proto3" json:"xgbimportance_max,omitempty"`
		PearsonTop_1BpsPercentile                     float32                `protobuf:"fixed32,19,opt,name=pearson_top_1_bps_percentile,json=pearsonTop1BpsPercentile,proto3" json:"pearson_top_1_bps_percentile,omitempty"`
		SpearmanTop_1BpsPercentile                    float32                `protobuf:"fixed32,20,opt,name=spearman_top_1_bps_percentile,json=spearmanTop1BpsPercentile,proto3" json:"spearman_top_1_bps_percentile,omitempty"`
		KendallTop_1BpsPercentile                     float32                `protobuf:"fixed32,21,opt,name=kendall_top_1_bps_percentile,json=kendallTop1BpsPercentile,proto3" json:"kendall_top_1_bps_percentile,omitempty"`
		HoeffdingTop_10BpsPercentile                  float32                `protobuf:"fixed32,22,opt,name=hoeffding_top_10_bps_percentile,json=hoeffdingTop10BpsPercentile,proto3" json:"hoeffding_top_10_bps_percentile,omitempty"`
		MutualInformationTop_100BpsPercentile         float32                `protobuf:"fixed32,23,opt,name=mutual_information_top_100_bps_percentile,json=mutualInformationTop100BpsPercentile,proto3" json:"mutual_information_top_100_bps_percentile,omitempty"`
		HsicTop_100BpsPercentile                      float32                `protobuf:"fixed32,24,opt,name=hsic_top_100_bps_percentile,json=hsicTop100BpsPercentile,proto3" json:"hsic_top_100_bps_percentile,omitempty"`
		XgbimportanceTop_100BpsPercentile             float32                `protobuf:"fixed32,25,opt,name=xgbimportance_top_100_bps_percentile,json=xgbimportanceTop100BpsPercentile,proto3" json:"xgbimportance_top_100_bps_percentile,omitempty"`
		CombinedRarenessScore                         float32                `protobuf:"fixed32,26,opt,name=combined_rareness_score,json=combinedRarenessScore,proto3" json:"combined_rareness_score,omitempty"`
		XgboostPredictedRarenessScore                 float32                `protobuf:"fixed32,27,opt,name=xgboost_predicted_rareness_score,json=xgboostPredictedRarenessScore,proto3" json:"xgboost_predicted_rareness_score,omitempty"`
		NnPredictedRarenessScore                      float32                `protobuf:"fixed32,28,opt,name=nn_predicted_rareness_score,json=nnPredictedRarenessScore,proto3" json:"nn_predicted_rareness_score,omitempty"`
		OverallAverageRarenessScore                   float32                `protobuf:"fixed32,29,opt,name=overall_average_rareness_score,json=overallAverageRarenessScore,proto3" json:"overall_average_rareness_score,omitempty"`
		IsRareOnInternet                              bool                   `protobuf:"varint,30,opt,name=is_rare_on_internet,json=isRareOnInternet,proto3" json:"is_rare_on_internet,omitempty"`
		MatchesFoundOnFirstPage                       uint32                 `protobuf:"varint,31,opt,name=matches_found_on_first_page,json=matchesFoundOnFirstPage,proto3" json:"matches_found_on_first_page,omitempty"`
		NumberOfPagesOfResults                        uint32                 `protobuf:"varint,32,opt,name=number_of_pages_of_results,json=numberOfPagesOfResults,proto3" json:"number_of_pages_of_results,omitempty"`
		UrlOfFirstMatchInPage                         string                 `protobuf:"bytes,33,opt,name=url_of_first_match_in_page,json=urlOfFirstMatchInPage,proto3" json:"url_of_first_match_in_page,omitempty"`
		OpenNsfwScore                                 float32                `protobuf:"fixed32,34,opt,name=open_nsfw_score,json=openNsfwScore,proto3" json:"open_nsfw_score,omitempty"`
		AlternativeNsfwScores                         *AltNsfwScores         `protobuf:"bytes,35,opt,name=alternative_nsfw_scores,json=alternativeNsfwScores,proto3" json:"alternative_nsfw_scores,omitempty"`
		ImageHashes                                   *PerceptualImageHashes `protobuf:"bytes,36,opt,name=image_hashes,json=imageHashes,proto3" json:"image_hashes,omitempty"`
		ImageFingerprintOfCandidateImageFile          []float32              `protobuf:"fixed32,37,rep,packed,name=image_fingerprint_of_candidate_image_file,json=imageFingerprintOfCandidateImageFile,proto3" json:"image_fingerprint_of_candidate_image_file,omitempty"`
	*/
	output := &DupeDetection{
		DupeDetectionSystemVer:  res.DupeDetectionSystemVersion,
		ImageHash:               res.HashOfCandidateImageFile,
		PastelRarenessScore:     res.OverallAverageRarenessScore,
		IsLikelyDupe:            res.IsLikelyDupe,
		IsRareOnInternet:        res.IsRareOnInternet,
		MatchesFoundOnFirstPage: res.MatchesFoundOnFirstPage,
		NumberOfResultPages:     res.NumberOfPagesOfResults,
		FirstMatchURL:           res.UrlOfFirstMatchInPage,
		OpenNSFWScore:           res.OpenNsfwScore,
		AlternateNSFWScores: AlternateNSFWScores{
			Drawings: res.AlternativeNsfwScores.Drawings,
			Hentai:   res.AlternativeNsfwScores.Hentai,
			Neutral:  res.AlternativeNsfwScores.Neutral,
			Porn:     res.AlternativeNsfwScores.Porn,
			Sexy:     res.AlternativeNsfwScores.Sexy,
		},
		ImageHashes: ImageHashes{
			PDQHash:        res.ImageHashes.PdqHash,
			PerceptualHash: res.ImageHashes.PerceptualHash,
			AverageHash:    res.ImageHashes.AverageHash,
			DifferenceHash: res.ImageHashes.DifferenceHash,
			NeuralHash:     res.ImageHashes.NeuralhashHash,
		},
		Fingerprints:                                 res.ImageFingerprintOfCandidateImageFile,
		PerceptualHashOverlapCount:                   res.PerceptualHashOverlapCount,
		NumberOfFingerprintsRequiringFurtherTesting1: res.NumberOfFingerprintsRequiringFurtherTesting_1,
		NumberOfFingerprintsRequiringFurtherTesting2: res.NumberOfFingerprintsRequiringFurtherTesting_2,
		NumberOfFingerprintsRequiringFurtherTesting3: res.NumberOfFingerprintsRequiringFurtherTesting_3,
		NumberOfFingerprintsRequiringFurtherTesting4: res.NumberOfFingerprintsRequiringFurtherTesting_4,
		NumberOfFingerprintsRequiringFurtherTesting5: res.NumberOfFingerprintsRequiringFurtherTesting_5,
		NumberOfFingerprintsRequiringFurtherTesting6: res.NumberOfFingerprintsRequiringFurtherTesting_6,
		NumberOfFingerprintsOfSuspectedDupes:         res.NumberOfFingerprintsOfSuspectedDupes,
		PearsonMax:                                   res.PearsonMax,
		SpearmanMax:                                  res.SpearmanMax,
		KendallMax:                                   res.KendallMax,
		HoeffdingMax:                                 res.HoeffdingMax,
		MutualInformationMax:                         res.MutualInformationMax,
		HsicMax:                                      res.HsicMax,
		XgbimportanceMax:                             res.XgbimportanceMax,
		PearsonTop1BpsPercentile:                     res.PearsonTop_1BpsPercentile,
		SpearmanTop1BpsPercentile:                    res.SpearmanTop_1BpsPercentile,
		KendallTop1BpsPercentile:                     res.KendallTop_1BpsPercentile,
		HoeffdingTop10BpsPercentile:                  res.HoeffdingTop_10BpsPercentile,
		MutualInformationTop100BpsPercentile:         res.MutualInformationTop_100BpsPercentile,
		HsicTop100BpsPercentile:                      res.HsicTop_100BpsPercentile,
		XgbimportanceTop100BpsPercentile:             res.XgbimportanceTop_100BpsPercentile,
		CombinedRarenessScore:                        res.CombinedRarenessScore,
		XgboostPredictedRarenessScore:                res.XgboostPredictedRarenessScore,
		NnPredictedRarenessScore:                     res.NnPredictedRarenessScore,
		OverallAverageRarenessScore:                  res.OverallAverageRarenessScore,
		NumberOfPagesOfResults:                       res.NumberOfPagesOfResults,
		UrlOfFirstMatchInPage:                        res.UrlOfFirstMatchInPage,
	}

	return output, nil
}

// ImageRarenessScore call ddserver to calculate scores
func (ddClient *ddServerClientImpl) ImageRarenessScore(ctx context.Context, img []byte, format string) (*DupeDetection, error) {
	ctx = ddClient.contextWithLogPrefix(ctx)

	baseClient := NewClient()
	ddServerAddress := fmt.Sprint(ddClient.config.Host, ":", ddClient.config.Port)
	conn, err := baseClient.Connect(ctx, ddServerAddress)
	if err != nil {
		return nil, errors.Errorf("connect to dd-server  %w", err)
	}

	defer conn.Close()
	client := pb.NewDupeDetectionServerClient(conn)

	return ddClient.callImageRarenessScore(ctx, client, img, format)
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
