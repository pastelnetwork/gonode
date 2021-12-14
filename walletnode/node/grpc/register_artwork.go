package grpc

import (
	"context"
	"encoding/base64"

	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
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

// AcceptedNodes implements node.RegisterArtwork.AcceptedNodes()
func (service *registerArtwork) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.AcceptedNodesRequest{}
	log.WithContext(ctx).WithField("req", req).Debug("AcceptedNodes request")

	resp, err := service.client.AcceptedNodes(ctx, req)
	if err != nil {
		return nil, errors.Errorf("request to accepted secondary nodes: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("AcceptedNodes response")

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
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return err
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")

	return nil
}

// ProbeImage implements node.RegisterArtwork.ProbeImage()
func (service *registerArtwork) ProbeImage(ctx context.Context, image *artwork.File) (*pastel.DDAndFingerprints, []byte, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.ProbeImage(ctx)
	if err != nil {
		return nil, nil, errors.Errorf("open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := image.Open()
	if err != nil {
		return nil, nil, errors.Errorf("open file %q: %w", file.Name(), err)
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
			return nil, nil, errors.Errorf("send image data: %w", err)
		}
	}

	// -------------------- WIP: PSL-142 -----------------------------
	// stream.CloseAndRecv() is supposed to return  a string response
	var ddDataSignature []byte
	/*resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, errors.Errorf("receive image response: %w", err)
	}
	ddDataSignature = resp.ddDataSignature
	*/

	return pastel.GetDDFingerprintsAndSigFromProbeImageReply(ddDataSignature)
}

// UploadImageWithThumbnail implements node.RegisterArtwork.UploadImageWithThumbnail()
func (service *registerArtwork) UploadImageWithThumbnail(ctx context.Context, image *artwork.File, thumbnail artwork.ThumbnailCoordinate) ([]byte, []byte, []byte, error) {
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
		return nil, nil, nil, errors.Errorf("compute artwork hash:%w", err)
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

func (service *registerArtwork) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerArtwork) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

// SendSignedTicket
func (service *registerArtwork) SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, key1 string, key2 string, rqids []byte, ddFp []byte, encoderParams rqnode.EncoderParameters) (int64, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := pb.SendSignedNFTTicketRequest{
		NftTicket:        ticket,
		CreatorSignature: signature,
		Key1:             key1,
		Key2:             key2,
		EncodeParameters: &pb.EncoderParameters{
			Oti: encoderParams.Oti,
		},

		// ----- WIP: PSL-142-------------Waiting for contract update
		//EncodeFiles: rqids,
	}

	rsp, err := service.client.SendSignedNFTTicket(ctx, &req)
	if err != nil {
		return -1, err
	}

	return rsp.RegistrationFee, nil
}

func (service *registerArtwork) SendPreBurntFeeTxid(ctx context.Context, txid string) (string, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	log.WithContext(ctx).Debug("send burned txid to super node")
	req := pb.SendPreBurntFeeTxidRequest{
		Txid: txid,
	}

	rsp, err := service.client.SendPreBurntFeeTxid(ctx, &req)
	if err != nil {
		return "", errors.Errorf("send burned txid to super node: %w", err)
	}

	// TODO: response from sending preburned TxId should be the TxId of RegActTicket
	return rsp.NFTRegTxid, nil
}

func newRegisterArtwork(conn *clientConn) node.RegisterArtwork {
	return &registerArtwork{
		conn:   conn,
		client: pb.NewRegisterArtworkClient(conn),
	}
}
