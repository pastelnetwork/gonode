package common

import (
	"context"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"golang.org/x/crypto/sha3"
)

// BlockMessage is block message
type BlockMessage struct {
	Block int
}

const (
	//P2PConnectionCloseError represents the p2p connection close error
	P2PConnectionCloseError = "grpc: the client connection is closing"
	//P2PServiceNotRunningError represents the p2p service not running error
	P2PServiceNotRunningError = "p2p service is not running"
)

// CreateNewP2PConnection creates a new p2p connection using sn host and port
func CreateNewP2PConnection(host string, port int, sn node.SNClientInterface) (node.HermesP2PInterface, error) {
	snAddr := fmt.Sprintf("%s:%d", host, port)
	log.WithField("sn-addr", snAddr).Info("connecting with SN-Service")

	conn, err := sn.Connect(context.Background(), snAddr)
	if err != nil {
		return nil, errors.Errorf("unable to connect with SN service: %w", err)
	}
	log.Info("connection established with SN-Service")

	p2p := conn.HermesP2P()
	log.Info("connection established with SN-P2P-Service")

	return p2p, nil
}

// IsP2PConnectionCloseError checks if the passed errString is of gRPC connection closing
func IsP2PConnectionCloseError(errString string) bool {
	return strings.Contains(errString, P2PConnectionCloseError)
}

// IsP2PServiceNotRunningError checks if the passed errString is of p2p not running
func IsP2PServiceNotRunningError(errString string) bool {
	return strings.Contains(errString, P2PServiceNotRunningError)
}

// ThumbnailType is thumbnail type
type ThumbnailType int

// CreateAndHashThumbnail creates a thumbnail from the given image and returns its hash.
func CreateAndHashThumbnail(srcImg image.Image, _ ThumbnailType, rect *image.Rectangle, quality float32, targetFileSize int) ([]byte, error) {

	var thumbnailImg image.Image
	if rect != nil {
		thumbnailImg = imaging.Crop(srcImg, *rect)
	} else {
		thumbnailImg = srcImg
	}

	log.Debugf("Encode with target size %d and quality %f", targetFileSize, quality)
	encoderOptions, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return nil, errors.Errorf("create lossless encoder option %w", err)
	}
	encoderOptions.TargetSize = targetFileSize

	if err := webp.Encode(nil, thumbnailImg, encoderOptions); err != nil {
		return nil, errors.Errorf("encode to webp format: %w", err)
	}

	hasher := sha3.New256()
	if _, err := io.Copy(hasher, nil); err != nil {
		return nil, errors.Errorf("hash failed: %w", err)
	}

	return hasher.Sum(nil), nil
}
