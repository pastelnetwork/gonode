//go:generate mockery --name=DDServerClient

package ddclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

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
	ImageRarenessScore(ctx context.Context, img []byte, format string, blockHash string, blockHeight string, timestamp string, pastelID string,
		supernode1 string, supernode2 string, supernode3 string, openAPIRequest bool, groupID string, collectionName string) (*pastel.DDAndFingerprints, error)
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
func (ddClient *ddServerClientImpl) callImageRarenessScore(ctx context.Context, client pb.DupeDetectionServerClient, img []byte, format string,
	blockHash string, blockHeight string, timestamp string, pastelID string, supernode1 string, supernode2 string, supernode3 string,
	openAPIRequest bool, groupID string, collectionName string) (*pastel.DDAndFingerprints, error) {
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
		PastelIdOfRegisteringSupernode_1:      supernode1,
		PastelIdOfRegisteringSupernode_2:      supernode2,
		PastelIdOfRegisteringSupernode_3:      supernode3,
		IsPastelOpenapiRequest:                openAPIRequest,
		OpenApiGroupIdString:                  groupID,
		CollectionNameString:                  collectionName,
	}

	// remove file after use
	defer os.Remove(inputPath)

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	res, err := client.ImageRarenessScore(reqCtx, &req)
	if err != nil {
		return nil, errors.Errorf("Error calling image rareness score: %w", err)
	}

	output := &pastel.DDAndFingerprints{
		BlockHash:   res.PastelBlockHashWhenRequestSubmitted,
		BlockHeight: res.PastelBlockHeightWhenRequestSubmitted,

		TimestampOfRequest: res.UtcTimestampWhenRequestSubmitted,
		SubmitterPastelID:  res.PastelIdOfSubmitter,
		SN1PastelID:        res.PastelIdOfRegisteringSupernode_1,
		SN2PastelID:        res.PastelIdOfRegisteringSupernode_2,
		SN3PastelID:        res.PastelIdOfRegisteringSupernode_3,

		IsOpenAPIRequest:           res.IsPastelOpenapiRequest,
		ImageFilePath:              res.ImageFilePath,
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

		CollectionNameString:                       res.CollectionNameString,
		OpenAPIGroupIDString:                       res.OpenApiGroupIdString,
		GroupRarenessScore:                         res.GroupRarenessScore,
		CandidateImageThumbnailWebpAsBase64String:  res.CandidateImageThumbnailWebpAsBase64String,
		DoesNotImpactTheFollowingCollectionStrings: res.DoesNotImpactTheFollowingCollectionStrings,
		IsInvalidSenseRequest:                      res.IsInvalidSenseRequest,
		InvalidSenseRequestReason:                  res.InvalidSenseRequestReason,
		SimilarityScoreToFirstEntryInCollection:    res.SimilarityScoreToFirstEntryInCollection,
		CPProbability:                              res.CpProbability,
		ChildProbability:                           res.ChildProbability,
	}

	log.WithContext(ctx).WithField("Rareness Score Response", output).Info("Image rareness score response from dd-server")

	return output, output.Validate()
}

// ImageRarenessScore call ddserver to calculate scores
func (ddClient *ddServerClientImpl) ImageRarenessScore(ctx context.Context, img []byte, format string, blockHash string, blockHeight string,
	timestamp string, pastelID string, supernode1 string, supernode2 string, supernode3 string, openAPIRequest bool,
	groupID string, collectionName string) (*pastel.DDAndFingerprints, error) {
	ctx = ddClient.contextWithLogPrefix(ctx)

	baseClient := NewClient()
	ddServerAddress := fmt.Sprint(ddClient.config.Host, ":", ddClient.config.Port)
	conn, err := baseClient.Connect(ctx, ddServerAddress)
	if err != nil {
		return nil, errors.Errorf("connect to dd-server  %w", err)
	}

	defer conn.Close()
	client := pb.NewDupeDetectionServerClient(conn)

	for {
		select {
		case <-ctx.Done():
			return nil, errors.Errorf("dd server request cancelled because context done: %w", ctx.Err())
		case <-time.After(15 * time.Second):
			// Create a context with timeout
			reqCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			stats, err := client.GetStatus(reqCtx, &pb.GetStatusRequest{})
			if err != nil {
				return nil, errors.Errorf("get status: %w", err)
			}

			if stats == nil {
				return nil, errors.Errorf("stats is nil")
			}

			taskCount := stats.GetTaskCount()
			if taskCount == nil {
				log.WithContext(ctx).Error("task count is nil")
				return ddClient.callImageRarenessScore(ctx, client, img, format, blockHash, blockHeight, timestamp, pastelID, supernode1, supernode2,
					supernode3, openAPIRequest, groupID, collectionName)
			}

			log.WithContext(ctx).WithField("executing", taskCount.GetExecuting()).WithField("waiting", taskCount.GetWaitingInQueue()).WithField("max-concurrent", taskCount.GetMaxConcurrent()).
				WithField("succeeded", taskCount.GetSucceeded()).WithField("failed", taskCount.GetFailed()).WithField("cancelled", taskCount.GetCancelled()).
				Info("dd-server task stats")

			taskMetrics := stats.GetTaskMetrics()
			if taskMetrics != nil {
				log.WithContext(ctx).WithField("avg-execution-time", taskMetrics.GetAverageTaskExecutionTimeSecs()).WithField("avg-queue-time", taskMetrics.GetAverageTaskWaitTimeSecs()).
					WithField("avg-vr-memory", taskMetrics.GetAverageTaskVirtualMemoryUsageBytes()).WithField("avg-rss-memory", taskMetrics.GetAverageTaskRssMemoryUsageBytes()).
					WithField("max-rss-memory", taskMetrics.GetPeakTaskRssMemoryUsageBytes()).WithField("max-vr-memory", taskMetrics.GetPeakTaskVmsMemoryUsageBytes()).
					WithField("max-task-wait-time", taskMetrics.GetMaxTaskWaitTimeSecs()).Info("dd-server task metrics")
			}

			if taskCount.GetWaitingInQueue() >= 2 {
				continue
			} else {
				return ddClient.callImageRarenessScore(ctx, client, img, format, blockHash, blockHeight, timestamp, pastelID, supernode1, supernode2,
					supernode3, openAPIRequest, groupID, collectionName)
			}
		case <-time.After(15 * time.Minute):
			return nil, errors.Errorf("dd-server request cancelled because of timeout")
		}
	}
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
