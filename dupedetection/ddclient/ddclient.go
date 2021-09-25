//go:generate mockery --name=DDServerClient

package ddclient

import (
	"context"
	"fmt"
	"io/ioutil"
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
	DupeDetectionSystemVer  string              `json:"dupe_detection_system_version"`
	ImageHash               []byte              `json:"hash_of_candidate_image_file"`
	PastelRarenessScore     float32             `json:"overall_average_rareness_score"`
	IsRareOnInternet        bool                `json:"is_rare_on_internet"`
	MatchesFoundOnFirstPage uint32              `json:"matches_found_on_first_page"`
	NumberOfResultPages     uint32              `json:"number_of_pages_of_results"`
	FirstMatchURL           string              `json:"url_of_first_match_in_page"`
	OpenNSFWScore           float32             `json:"open_nsfw_score"`
	AlternateNSFWScores     AlternateNSFWScores `json:"alternative_nsfw_scores"`
	ImageHashes             ImageHashes         `json:"image_hashes"`
	Fingerprints            []float32           `json:"image_fingerprint_of_candidate_image_file"`
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
	return ioutil.WriteFile(path, data, 0777)
}

func createInputDDFile(base string, data []byte, format string) (string, error) {
	fileName := randID() + "." + format
	filePath := filepath.Join(base, fileName)

	err := writeFile(filePath, data)

	if err != nil {
		return "", errors.Errorf("failed to write data to %s: %w", filePath, err)
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
		return nil, errors.Errorf("failed to create input file: %w", err)
	}

	req := pb.RarenessScoreRequest{
		Path: inputPath,
	}

	res, err := client.ImageRarenessScore(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("failed to send request: %w", err)
	}

	/*
		DupeDetectionSystemVersion           string         `protobuf:"bytes,1,opt,name=dupe_detection_system_version,json=dupeDetectionSystemVersion,proto3" json:"dupe_detection_system_version,omitempty"`
		HashOfCandidateImageFile             []byte         `protobuf:"bytes,2,opt,name=hash_of_candidate_image_file,json=hashOfCandidateImageFile,proto3" json:"hash_of_candidate_image_file,omitempty"`
		IsLikelyDupe                         bool           `protobuf:"varint,3,opt,name=is_likely_dupe,json=isLikelyDupe,proto3" json:"is_likely_dupe,omitempty"`
		OverallAverageRarenessScore          float32        `protobuf:"fixed32,4,opt,name=overall_average_rareness_score,json=overallAverageRarenessScore,proto3" json:"overall_average_rareness_score,omitempty"`
		IsRareOnInternet                     bool           `protobuf:"varint,5,opt,name=is_rare_on_internet,json=isRareOnInternet,proto3" json:"is_rare_on_internet,omitempty"`
		MatchesFoundOnFirstPage              uint32         `protobuf:"varint,6,opt,name=matches_found_on_first_page,json=matchesFoundOnFirstPage,proto3" json:"matches_found_on_first_page,omitempty"`
		NumberOfPagesOfResults               uint32         `protobuf:"varint,7,opt,name=number_of_pages_of_results,json=numberOfPagesOfResults,proto3" json:"number_of_pages_of_results,omitempty"`
		UrlOfFirstMatchInPage                string         `protobuf:"bytes,8,opt,name=url_of_first_match_in_page,json=urlOfFirstMatchInPage,proto3" json:"url_of_first_match_in_page,omitempty"`
		OpenNsfwScore                        float32        `protobuf:"fixed32,9,opt,name=open_nsfw_score,json=openNsfwScore,proto3" json:"open_nsfw_score,omitempty"`
		AlternativeNsfwScores                *AltNsfwScores `protobuf:"bytes,10,opt,name=alternative_nsfw_scores,json=alternativeNsfwScores,proto3" json:"alternative_nsfw_scores,omitempty"`
		ImageHashes                          *ImageHashes   `protobuf:"bytes,11,opt,name=image_hashes,json=imageHashes,proto3" json:"image_hashes,omitempty"`
		ImageFingerprintOfCandidateImageFile []float32
	*/
	output := &DupeDetection{
		DupeDetectionSystemVer:  res.DupeDetectionSystemVersion,
		ImageHash:               res.HashOfCandidateImageFile,
		PastelRarenessScore:     res.OverallAverageRarenessScore,
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
		Fingerprints: res.ImageFingerprintOfCandidateImageFile,
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
		return nil, errors.Errorf("failed to connect to dd-server  %w", err)
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
