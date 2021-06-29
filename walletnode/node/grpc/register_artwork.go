package grpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	uploadImageBufferSize = 32 * 1024
)

type registerArtwork struct {
	conn   *clientConn
	client pb.RegisterArtworkClient

	sessID string
}

func (service *registerArtwork) SessID() string {
	return service.sessID
}

// Session implements node.RegisterArtwork.Session()
func (service *registerArtwork) Session(ctx context.Context, isPrimary bool) error {
	ctx = service.contextWithLogPrefix(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("failed to open Health stream: %w", err)
	}

	req := &pb.SessionRequest{
		IsPrimary: isPrimary,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if err := stream.Send(req); err != nil {
		return errors.Errorf("failed to send Session request: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		switch status.Code(err) {
		case codes.Canceled, codes.Unavailable:
			return nil
		}
		return errors.Errorf("failed to receive Session response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Session response")
	service.sessID = resp.SessID

	go func() {
		defer service.conn.Close()
		for {
			if _, err := stream.Recv(); err != nil {
				return
			}
		}
	}()

	return nil
}

// AcceptedNodes implements node.RegisterArtwork.AcceptedNodes()
func (service *registerArtwork) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.AcceptedNodesRequest{}
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")

	resp, err := service.client.AcceptedNodes(ctx, req)
	if err != nil {
		return nil, errors.Errorf("failed to request to accepted secondary nodes: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")

	var ids []string
	for _, peer := range resp.Peers {
		ids = append(ids, peer.NodeID)
	}
	return ids, nil
}

// ConnectTo implements node.RegisterArtwork.ConnectTo()
func (service *registerArtwork) ConnectTo(ctx context.Context, nodeID, sessID string) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.ConnectToRequest{
		NodeID: nodeID,
		SessID: sessID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return errors.Errorf("failed to request to connect to primary node: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")

	return nil
}

// ProbeImage implements node.RegisterArtwork.ProbeImage()
func (service *registerArtwork) ProbeImage(ctx context.Context, image *artwork.File) ([]byte, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.ProbeImage(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := image.Open()
	if err != nil {
		return nil, errors.Errorf("failed to open file %q: %w", file.Name(), err)
	}
	defer file.Close()

	buffer := make([]byte, uploadImageBufferSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		req := &pb.ProbeImageRequest{
			Payload: buffer[:n],
		}
		if err := stream.Send(req); err != nil {
			return nil, errors.Errorf("failed to send image data: %w", err).WithField("reqID", service.conn.id)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, errors.Errorf("failed to receive send image response: %w", err)
	}
	log.WithContext(ctx).WithField("fingerprintLenght", len(resp.Fingerprint)).Debugf("ProbeImage response")

	return resp.Fingerprint, nil
}

// UploadImageWithThumbnail implements node.RegisterArtwork.UploadImageWithThumbnail()
func (service *registerArtwork) UploadImageWithThumbnail(ctx context.Context, image *artwork.File, thumbnail artwork.ImageThumbnail) ([]byte, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	log.WithContext(ctx).Debug("Start upload image and thumbnail to node")
	stream, err := service.client.UploadImage(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := image.Open()
	if err != nil {
		return nil, errors.Errorf("failed to open file %q: %w", file.Name(), err)
	}
	defer file.Close()

	buffer := make([]byte, uploadImageBufferSize)
	lastPiece := false
	payloadSize := 0
	for {
		n, err := file.Read(buffer)
		payloadSize += n
		if err != nil && err == io.EOF {
			log.WithContext(ctx).WithField("Filename", file.Name()).Debugf("EOF")
			lastPiece = true
			break
		} else if err != nil {
			log.WithContext(ctx).WithField("Filename", file.Name()).Debugf("Error %w", err)
			return nil, err
		}

		req := &pb.UploadImageRequest{
			Payload: &pb.UploadImageRequest_ImagePiece{
				ImagePiece: buffer[:n],
			},
		}

		if err := stream.Send(req); err != nil {
			return nil, errors.Errorf("failed to send image data: %w", err).WithField("ReqID", service.conn.id)
		}
	}

	log.WithContext(ctx).Debugf("Encoded Image Size :%d\n", payloadSize)
	if !lastPiece {
		return nil, errors.Errorf("failed to read all image data")
	}

	file.Seek(0, io.SeekStart)
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, errors.Errorf("failed to compute artwork hash %w", err)
	}
	hash := hasher.Sum(nil)
	log.WithContext(ctx).WithField("Filename", file.Name()).Debugf("hash: %s", base64.URLEncoding.EncodeToString(hash))

	thumnailReq := &pb.UploadImageRequest{
		Payload: &pb.UploadImageRequest_MetaData_{
			MetaData: &pb.UploadImageRequest_MetaData{
				Hash:   hash[:],
				Size:   int64(payloadSize),
				Format: image.Format().String(),
				Thumbnail: &pb.UploadImageRequest_Coordinate{
					TopLeftX:     thumbnail.TopLeftX,
					TopLeftY:     thumbnail.TopLeftY,
					BottomRightX: thumbnail.BottomRightX,
					BottomRightY: thumbnail.BottomRightY,
				},
			},
		},
	}

	if err := stream.Send(thumnailReq); err != nil {
		return nil, errors.Errorf("failed to send image thumbnail: %w", err).WithField("ReqID", service.conn.id)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, errors.Errorf("failed to receive upload image response: %w", err)
	}
	log.WithContext(ctx).WithField("thumbnailHashLength", len(resp.ThumbnailHash)).Debugf("UploadImageWithThumbnail response")

	return resp.ThumbnailHash, nil
}

func (service *registerArtwork) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerArtwork) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRegisterArtwork(conn *clientConn) node.RegisterArtwork {
	return &registerArtwork{
		conn:   conn,
		client: pb.NewRegisterArtworkClient(conn),
	}
}
