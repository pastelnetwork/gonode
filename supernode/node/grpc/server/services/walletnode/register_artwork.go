package walletnode

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
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
	workDir string
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
		task = service.NewTask(ctx)
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
		SessID: task.ID,
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

// UploadImage implements walletnode.RegisterArtworkServer.UploadImage()
func (service *RegisterArtwork) UploadImage(stream pb.RegisterArtwork_UploadImageServer) error {
	ctx := stream.Context()

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}

	fileID, _ := random.String(16, random.Base62Chars)
	filename := filepath.Join(service.workDir, fileID)

	file, err := os.Create(filename)
	if err != nil {
		return errors.Errorf("failed to open file %q: %w", filename, err)
	}

	defer file.Close()
	log.WithContext(ctx).Debugf("Created temp file %q for uploading image", filename)

	wr := bufio.NewWriter(file)

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			if status.Code(err) == codes.Canceled {
				return errors.New("connection closed")
			}
			return errors.Errorf("failed to receive UploadImage: %w", err)
		}

		if _, err := wr.Write(req.Payload); err != nil {
			return errors.Errorf("failed to write to file %q: %w", filename, err)
		}
	}

	if err := task.UploadImage(ctx, filename); err != nil {
		return err
	}

	resp := &pb.UploadImageReply{}
	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("failed to send UploadImage response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("UploadImage response")
	return nil
}

// Desc returns a description of the service.
func (service *RegisterArtwork) Desc() *grpc.ServiceDesc {
	return &pb.RegisterArtwork_ServiceDesc
}

// NewRegisterArtwork returns a new RegisterArtwork instance.
func NewRegisterArtwork(service *artworkregister.Service, workDir string) *RegisterArtwork {
	return &RegisterArtwork{
		RegisterArtwork: common.NewRegisterArtwork(service),
		workDir:         workDir,
	}
}
