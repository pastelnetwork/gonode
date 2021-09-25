package grpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/dupedetection"
	"github.com/pastelnetwork/gonode/dupedetection/node"
)

const (
	inputEncodeFileName = "input.data"
	symbolIDFileSubDir  = "meta"
	symbolFileSubdir    = "symbols"
)

type dupedetectionImpl struct {
	conn   *clientConn
	client pb.DupeDetectionServerClient
	config *node.DDServerConfig
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

// Encode data, and return a list of identifier of symbols
func (service *dupedetectionImpl) ImageRarenessScore(ctx context.Context, img []byte, format string) (*node.DupeDetection, error) {
	if img == nil {
		return nil, errors.Errorf("invalid data")
	}

	ctx = service.contextWithLogPrefix(ctx)

	inputPath, err := createInputDDFile(service.config.DDFilesDir, img, format)
	if err != nil {
		return nil, errors.Errorf("failed to create input file: %w", err)
	}

	req := pb.RarenessScoreRequest{
		Path: inputPath,
	}

	res, err := service.client.ImageRarenessScore(ctx, &req)
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
	output := &node.DupeDetection{
		DupeDetectionSystemVer:  res.DupeDetectionSystemVersion,
		ImageHash:               res.HashOfCandidateImageFile,
		PastelRarenessScore:     res.OverallAverageRarenessScore,
		IsRareOnInternet:        res.IsRareOnInternet,
		MatchesFoundOnFirstPage: res.MatchesFoundOnFirstPage,
		NumberOfResultPages:     res.NumberOfPagesOfResults,
		FirstMatchURL:           res.UrlOfFirstMatchInPage,
		OpenNSFWScore:           res.OpenNsfwScore,
		AlternateNSFWScores: node.AlternateNSFWScores{
			Drawings: res.AlternativeNsfwScores.Drawings,
			Hentai:   res.AlternativeNsfwScores.Hentai,
			Neutral:  res.AlternativeNsfwScores.Neutral,
			Porn:     res.AlternativeNsfwScores.Porn,
			Sexy:     res.AlternativeNsfwScores.Sexy,
		},
		ImageHashes: node.ImageHashes{
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

func (service *dupedetectionImpl) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newDupedetectionImpl(conn *clientConn, config *node.DDServerConfig) node.Dupedetection {
	return &dupedetectionImpl{
		conn:   conn,
		client: pb.NewDupeDetectionServerClient(conn),
		config: config,
	}
}
