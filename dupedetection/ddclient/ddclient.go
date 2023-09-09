//go:generate mockery --name=DDServerClient

package ddclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	json "github.com/json-iterator/go"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/dupedetection"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/dupedetection"
	"github.com/pastelnetwork/gonode/pastel"
	"google.golang.org/grpc"
)

const (
	logPrefix                   = "ddClient"
	defaultDDServiceCallTimeout = 26 * time.Minute
)

// DDServerClient contains methods for request services from dd-server service.
type DDServerClient interface {
	// ImageRarenessScore returns rareness score of image
	ImageRarenessScore(ctx context.Context, img []byte, format string, blockHash string, blockHeight string, timestamp string, pastelID string,
		supernode1 string, supernode2 string, supernode3 string, openAPIRequest bool, groupID string, collectionName string) (*pastel.DDAndFingerprints, error)

	//GetStats returns the stats from dd-server
	GetStats(ctx context.Context) (dupedetection.DDServerStats, error)
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

	// Limits the dial timeout, prevent got stuck too long
	dialCtx, dcancel := context.WithTimeout(ctx, defaultDDServiceCallTimeout)
	defer dcancel()

	res, err := client.ImageRarenessScore(dialCtx, &req, grpc.MaxCallRecvMsgSize(35000000))
	if err != nil {
		return nil, errors.Errorf("Error calling image rareness score, dd-serivce returned error: %w", err)
	}

	bytes, _ := json.Marshal(res)
	sizeInKBObj := float64(len(bytes)) / 1024

	log.WithContext(ctx).WithField("size", sizeInKBObj).Infof("Size of the object: %.2f KB", sizeInKBObj)

	sizeInKB := len(res.RarenessScoresTableJsonCompressedB64) / 1024
	if sizeInKB < 50 {
		log.WithContext(ctx).WithField("size", sizeInKB).Info("Size of the rareness object")
		//return nil, errors.Errorf("Error calling image rareness score - response rcvd had very low size: obj: %.2f  rareness: %.2f", sizeInKBObj, sizeInKB)
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

	//log.WithContext(ctx).WithField("Rareness Score Response", output).Info("Image rareness score response from dd-server")

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

	startTime := time.Now()
	for {
		// Check if the context is done
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("original context cancelled: %w", ctx.Err())
		default:
			// Continue if not done
		}

		reqCtx, reqCancel := context.WithTimeout(ctx, 1*time.Minute)
		stats, err := client.GetStatus(reqCtx, &pb.GetStatusRequest{})
		reqCancel() // cancel context after use

		if err != nil {
			return nil, errors.Errorf("retry: get status: %w", err)
		} else if stats == nil {
			return nil, errors.Errorf("retry: get status: stats is nil")
		}

		taskCount := stats.GetTaskCount()
		if taskCount == nil {
			return nil, errors.New("retry: task count is nil")
		}

		log.WithContext(reqCtx).WithField("executing", taskCount.GetExecuting()).WithField("waiting", taskCount.GetWaitingInQueue()).WithField("max-concurrent", taskCount.GetMaxConcurrent()).
			WithField("succeeded", taskCount.GetSucceeded()).WithField("failed", taskCount.GetFailed()).WithField("cancelled", taskCount.GetCancelled()).
			Info("dd-server task stats")

		if taskCount.GetWaitingInQueue() == 0 {
			break
		}

		if time.Since(startTime).Minutes() > 25 {
			return nil, errors.New("retry: waiting count did not reach 0 after 25 minutes")
		}

		// Wait for a minute before checking the status again
		time.Sleep(1 * time.Minute)
	}

	log.WithContext(ctx).Info("calling image rareness score now")
	return ddClient.callImageRarenessScore(ctx, client, img, format, blockHash, blockHeight, timestamp, pastelID, supernode1, supernode2,
		supernode3, openAPIRequest, groupID, collectionName)
}

// GetStats return the stats from the dd-server
func (ddClient ddServerClientImpl) GetStats(ctx context.Context) (stats dupedetection.DDServerStats, err error) {
	ctx = ddClient.contextWithLogPrefix(ctx)

	baseClient := NewClient()
	ddServerAddress := fmt.Sprint(ddClient.config.Host, ":", ddClient.config.Port)
	conn, err := baseClient.Connect(ctx, ddServerAddress)
	if err != nil {
		log.WithContext(ctx).WithError(err).Info("unable to connect to dd-server")
		return stats, errors.Errorf("connect to dd-server  %w", err)
	}

	defer conn.Close()
	client := pb.NewDupeDetectionServerClient(conn)
	log.WithContext(ctx).Info("connection established, sending request to dd-server for stats from SN")

	resp, err := client.GetStatus(ctx, &pb.GetStatusRequest{})
	if err != nil {
		log.WithContext(ctx).WithError(err).Info("error getting stats from sn - dd-server")
		return stats, errors.Errorf("error getting dd server stats from the dd-server")
	}
	log.WithContext(ctx).WithField("stats", resp).Info("stats received")

	stats = dupedetection.DDServerStats{
		MaxConcurrent:  resp.TaskCount.GetMaxConcurrent(),
		Executing:      resp.TaskCount.GetExecuting(),
		WaitingInQueue: resp.TaskCount.GetWaitingInQueue(),
	}

	return stats, nil
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
