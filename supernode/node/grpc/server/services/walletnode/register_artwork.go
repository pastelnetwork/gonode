package walletnode

import (
	"bufio"
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterArtwork represents grpc service for registration artwork.
type RegisterArtwork struct {
	pb.UnimplementedRegisterArtworkServer

	*common.RegisterArtwork
}

// Session implements walletnode.RegisterArtworkServer.Session()
func (service *RegisterArtwork) Session(stream pb.RegisterArtwork_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *artworkregister.Task

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewTask()
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("failed to receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if err := task.Session(ctx, req.IsPrimary); err != nil {
		return err
	}

	resp := &pb.SessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("failed to send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Session response")

	for {
		if _, err := stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			switch status.Code(err) {
			case codes.Canceled, codes.Unavailable:
				return nil
			}
			return errors.Errorf("handshake stream closed: %w", err)
		}
	}
}

// AcceptedNodes implements walletnode.RegisterArtworkServer.AcceptedNodes()
func (service *RegisterArtwork) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := task.AcceptedNodes(ctx)
	if err != nil {
		return nil, err
	}

	var peers []*pb.AcceptedNodesReply_Peer
	for _, node := range nodes {
		peers = append(peers, &pb.AcceptedNodesReply_Peer{
			NodeID: node.ID,
		})
	}

	resp := &pb.AcceptedNodesReply{
		Peers: peers,
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")
	return resp, nil
}

// ConnectTo implements walletnode.RegisterArtworkServer.ConnectTo()
func (service *RegisterArtwork) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.ConnectTo(ctx, req.NodeID, req.SessID); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")
	return resp, nil
}

// ProbeImage implements walletnode.RegisterArtworkServer.ProbeImage()
func (service *RegisterArtwork) ProbeImage(stream pb.RegisterArtwork_ProbeImageServer) error {
	ctx := stream.Context()

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}

	image := service.Storage.NewFile()
	file, err := image.Create()
	if err != nil {
		return errors.Errorf("failed to open file %q: %w", file.Name(), err)
	}
	defer file.Close()
	log.WithContext(ctx).WithField("filename", file.Name()).Debugf("ProbeImage request")

	wr := bufio.NewWriter(file)
	defer wr.Flush()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			if status.Code(err) == codes.Canceled {
				return errors.New("connection closed")
			}
			return errors.Errorf("failed to receive ProbeImage: %w", err)
		}

		if _, err := wr.Write(req.Payload); err != nil {
			return errors.Errorf("failed to write to file %q: %w", file.Name(), err)
		}
	}

	fingerprint, err := task.ProbeImage(ctx, image)
	if err != nil {
		return err
	}
	if fingerprint == nil {
		return errors.New("could not compute fingerprint data")
	}

	resp := &pb.ProbeImageReply{
		Fingerprint: fingerprint,
	}
	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("failed to send ProbeImage response: %w", err)
	}
	log.WithContext(ctx).WithField("fingerprintLenght", len(resp.Fingerprint)).Debugf("ProbeImage response")
	return nil
}

// UploadImageWithThumbnail implements walletnode.RegisterArtwork.UploadImageWithThumbnail
func (service *RegisterArtwork) UploadImage(stream pb.RegisterArtwork_UploadImageServer) error {
	ctx := stream.Context()

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}

	image := service.Storage.NewFile()
	imageFile, err := image.Create()
	if err != nil {
		return errors.Errorf("failed to open image file %q: %w", imageFile.Name(), err)
	}
	// defer imageFile.Close()
	log.WithContext(ctx).WithField("filename", imageFile.Name()).Debugf("UploadImageWithThumbnail request")

	imageWriter := bufio.NewWriter(imageFile)
	// defer imageWriter.Flush()

	thumbnail := artwork.ImageThumbnail{}
	imageSize := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			if status.Code(err) == codes.Canceled {
				return errors.New("connection closed")
			}
			return errors.Errorf("failed to receive UploadImageWithThumbnail: %w", err)
		}

		if imagePiece := req.GetImagePiece(); imagePiece != nil {
			n, err := imageWriter.Write(imagePiece)
			imageSize += n
			if err != nil {
				return errors.Errorf("failed to write to file %q: %w", imageFile.Name(), err)
			}
		} else {
			if cordinates := req.GetThumbnail(); cordinates != nil {
				thumbnail.TopLeftX = cordinates.TopLeftX
				thumbnail.TopLeftY = cordinates.TopLeftY
				thumbnail.BottomRightX = cordinates.BottomRightX
				thumbnail.BottomRightY = cordinates.BottomRightY
			}
		}
	}

	log.WithContext(ctx).Debugf("Receive %d\n", imageSize)
	decoder := qrsignature.New()
	if err := image.Decode(decoder); err != nil {
		return errors.Errorf("failed to decode image %w", err).WithField("FileName", imageFile.Name())
	}

	thumbnailHash, err := task.UploadImageWithThumbnail(ctx, image, thumbnail)
	if err != nil {
		return err
	}
	if thumbnailHash == nil {
		return errors.New("could not compute thumbnailHash data")
	}

	resp := &pb.UploadImageReply{
		ThumbnailHash: thumbnailHash,
	}
	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("failed to send UploadImageAndThumbnail response: %w", err)
	}
	log.WithContext(ctx).WithField("thumbnailHashLenght", len(resp.ThumbnailHash)).Debugf("UploadImageAndThumbnail response")
	return nil
}

// Desc returns a description of the service.
func (service *RegisterArtwork) Desc() *grpc.ServiceDesc {
	return &pb.RegisterArtwork_ServiceDesc
}

// NewRegisterArtwork returns a new RegisterArtwork instance.
func NewRegisterArtwork(service *artworkregister.Service) *RegisterArtwork {
	return &RegisterArtwork{
		RegisterArtwork: common.NewRegisterArtwork(service),
	}
}
