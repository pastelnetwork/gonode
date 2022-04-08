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
		ImageFilepath:                         inputPath,
		PastelBlockHashWhenRequestSubmitted:   blockHash,
		PastelBlockHeightWhenRequestSubmitted: blockHeight,
		UtcTimestampWhenRequestSubmitted:      timestamp,
		PastelIdOfSubmitter:                   pastelID,
		PastelIdOfRegisteringSupernode_1:      sn_1,
		PastelIdOfRegisteringSupernode_2:      sn_2,
		PastelIdOfRegisteringSupernode_3:      sn_3,
		IsPastelOpenapiRequest:                openapi_request,
		OpenApiSubsetIdString:                 open_api_subset_id,
	}

	// remove file after use
	defer os.Remove(inputPath)

	// 2/9/2022, we don't have this implemented yet, so this would fail
	res, err := client.ImageRarenessScore(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("send request: %w", err)
	}

	output := &pastel.DDAndFingerprints{
		BlockHash:   res.PastelBlockHashWhenRequestSubmitted,
		BlockHeight: res.PastelBlockHeightWhenRequestSubmitted,

		TimestampOfRequest: res.UtcTimestampWhenRequestSubmitted,
		SubmitterPastelID:  res.PastelIdOfSubmitter,
		SN1PastelID:        res.PastelIdOfSubmitter,
		SN2PastelID:        res.PastelIdOfRegisteringSupernode_2,
		SN3PastelID:        res.PastelIdOfRegisteringSupernode_3,

		IsOpenAPIRequest: res.IsPastelOpenapiRequest,

		OpenAPISubsetID: res.OpenApiSubsetIdString,

		DupeDetectionSystemVersion: res.DupeDetectionSystemVersion,

		IsLikelyDupe:         res.IsLikelyDupe,
		IsRareOnInternet:     res.IsRareOnInternet,
		OverallRarenessScore: res.OverallRarenessScore,

		PctOfTop10MostSimilarWithDupeProbAbove25pct: res.PctOfTop_10MostSimilarWithDupeProbAbove_25Pct,
		PctOfTop10MostSimilarWithDupeProbAbove33pct: res.PctOfTop_10MostSimilarWithDupeProbAbove_33Pct,
		PctOfTop10MostSimilarWithDupeProbAbove50pct: res.PctOfTop_10MostSimilarWithDupeProbAbove_50Pct,

		RarenessScoresTableJSONCompressedB64: res.RarenessScoresTableJsonCompressedB64,

		InternetRareness: &pastel.InternetRareness{
			RareOnInternetSummaryTableAsJSONCompressedB64:    res.InternetRareness.RareOnInternetSummaryTableAsJsonCompressedB64,
			RareOnInternetGraphJSONCompressedB64:             res.InternetRareness.RareOnInternetGraphJsonCompressedB64,
			AlternativeRareOnInternetDictAsJSONCompressedB64: res.InternetRareness.AlternativeRareOnInternetDictAsJsonCompressedB64,
			MinNumberOfExactMatchesInPage:                    res.InternetRareness.MinNumberOfExactMatchesInPage,
			EarliestAvailableDateOfInternetResults:           res.InternetRareness.EarliestAvailableDateOfInternetResults,
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

		HashOfCandidateImageFile: res.HashOfCandidateImageFile,
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
