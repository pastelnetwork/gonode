package grpc

import (
	"context"
	"encoding/base64"

	"github.com/pastelnetwork/gonode/common/storage/files"

	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	uploadImageBufferSize = 32 * 1024
)

type registerNft struct {
	conn   *clientConn
	client pb.RegisterNftClient

	sessID string
}

func (service *registerNft) SessID() string {
	return service.sessID
}

// Session implements node.RegisterNft.Session()
func (service *registerNft) Session(ctx context.Context, isPrimary bool) error {
	ctx = service.contextWithLogPrefix(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("open Session stream: %w", err)
	}

	req := &pb.SessionRequest{
		IsPrimary: isPrimary,
	}
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

	if err := stream.Send(req); err != nil {
		return errors.Errorf("send Session request: %w", err)
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
		return errors.Errorf("receive Session response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("Session response")
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

// AcceptedNodes implements node.RegisterNft.AcceptedNodes()
func (service *registerNft) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.AcceptedNodesRequest{}
	log.WithContext(ctx).WithField("req", req).Info("AcceptedNodes request")

	resp, err := service.client.AcceptedNodes(ctx, req)
	if err != nil {
		return nil, errors.Errorf("request to accepted secondary nodes: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Info("AcceptedNodes response")

	var ids []string
	for _, peer := range resp.Peers {
		ids = append(ids, peer.NodeID)
	}
	return ids, nil
}

// ConnectTo implements node.RegisterNft.ConnectTo()
func (service *registerNft) ConnectTo(ctx context.Context, primaryNode types.MeshedSuperNode) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.ConnectToRequest{
		NodeID: primaryNode.NodeID,
		SessID: primaryNode.SessID,
	}
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return err
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")

	return nil
}

// MeshNodes informs SNs which SNs are connected to do NFT request
func (service *registerNft) MeshNodes(ctx context.Context, meshedNodes []types.MeshedSuperNode) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	request := &pb.MeshNodesRequest{
		Nodes: []*pb.MeshNodesRequest_Node{},
	}

	for _, node := range meshedNodes {
		request.Nodes = append(request.Nodes, &pb.MeshNodesRequest_Node{
			SessID: node.SessID,
			NodeID: node.NodeID,
		})
	}

	_, err := service.client.MeshNodes(ctx, request)

	return err
}

// SendRegMetadata send metadata of registration to SNs for next steps
func (service *registerNft) SendRegMetadata(ctx context.Context, regMetadata *types.NftRegMetadata) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	request := &pb.SendRegMetadataRequest{
		BlockHash:       regMetadata.BlockHash,
		CreatorPastelID: regMetadata.CreatorPastelID,
		BlockHeight:     regMetadata.BlockHeight,
		Timestamp:       regMetadata.Timestamp,
	}

	_, err := service.client.SendRegMetadata(ctx, request)
	return err
}

// ProbeImage implements node.RegisterNft.ProbeImage()
// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Step 3, expect 4.B.5 in return
func (service *registerNft) ProbeImage(ctx context.Context, image *files.File) ([]byte, bool, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.ProbeImage(ctx)
	if err != nil {
		return nil, false, errors.Errorf("open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := image.Open()
	if err != nil {
		return nil, false, errors.Errorf("open file %q: %w", file.Name(), err)
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
			return nil, false, errors.Errorf("send image data: %w", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, false, errors.Errorf("receive image response: %w", err)
	}

	return resp.CompressedSignedDDAndFingerprints, true, nil
}

// UploadImageWithThumbnail implements node.RegisterNft.UploadImageWithThumbnail()
func (service *registerNft) UploadImageWithThumbnail(ctx context.Context, image *files.File, thumbnail files.ThumbnailCoordinate) ([]byte, []byte, []byte, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	log.WithContext(ctx).Debug("Start upload image and thumbnail to node")
	stream, err := service.client.UploadImage(ctx)
	if err != nil {
		return nil, nil, nil, errors.Errorf("open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := image.Open()
	if err != nil {
		return nil, nil, nil, errors.Errorf("open file %q: %w", file.Name(), err)
	}
	defer file.Close()

	buffer := make([]byte, uploadImageBufferSize)
	lastPiece := false
	payloadSize := 0
	for {
		n, err := file.Read(buffer)
		payloadSize += n
		if err != nil && err == io.EOF {
			log.WithContext(ctx).WithField("Filename", file.Name()).Debug("EOF")
			lastPiece = true
			break
		} else if err != nil {
			return nil, nil, nil, errors.Errorf("read file: %w", err)
		}

		req := &pb.UploadImageRequest{
			Payload: &pb.UploadImageRequest_ImagePiece{
				ImagePiece: buffer[:n],
			},
		}

		if err := stream.Send(req); err != nil {
			return nil, nil, nil, errors.Errorf("send image data: %w", err)
		}
	}

	log.WithContext(ctx).Debugf("Encoded Image Size :%d\n", payloadSize)
	if !lastPiece {
		return nil, nil, nil, errors.Errorf("read all image data failed")
	}

	file.Seek(0, io.SeekStart)
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, nil, nil, errors.Errorf("compute nft hash:%w", err)
	}
	hash := hasher.Sum(nil)
	log.WithContext(ctx).WithField("Filename", file.Name()).Debugf("hash: %s", base64.URLEncoding.EncodeToString(hash))

	thumbnailReq := &pb.UploadImageRequest{
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

	if err := stream.Send(thumbnailReq); err != nil {
		return nil, nil, nil, errors.Errorf("send image thumbnail: %w", err)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, nil, nil, errors.Errorf("receive upload image response: %w", err)
	}
	log.WithContext(ctx).Debugf("preview medium hash: %x", resp.PreviewThumbnailHash)
	log.WithContext(ctx).Debugf("medium thumbnail hash: %x", resp.MediumThumbnailHash)
	log.WithContext(ctx).Debugf("small thumbnail hash: %x", resp.SmallThumbnailHash)

	return resp.PreviewThumbnailHash, resp.MediumThumbnailHash, resp.SmallThumbnailHash, nil
}

func (service *registerNft) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerNft) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRegisterNft(conn *clientConn) node.RegisterNftInterface {
	return &registerNft{
		conn:   conn,
		client: pb.NewRegisterNftClient(conn),
	}
}
