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
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// registerArtowrk represents grpc service for registration artowrk.
type registerArtowrk struct {
	pb.UnimplementedRegisterArtowrkServer

	*common.RegisterArtowrk

	workDir string
}

func (service *registerArtowrk) Handshake(stream pb.RegisterArtowrk_HandshakeServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	task := service.NewTask(ctx)
	defer task.Cancel()

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}

				log.WithContext(ctx).WithError(err).Errorf("Handshake receving")
				return errors.Errorf("failed to receive Handshake: %w", err)
			}
			log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

			if err := task.Handshake(ctx, req.IsPrimary); err != nil {
				return err
			}

			resp := &pb.HandshakeReply{
				ConnID: task.ID,
				Error: &pb.Error{
					Status: pb.Error_OK,
				},
			}
			if err := stream.SendAndClose(resp); err != nil {
				log.WithContext(ctx).WithError(err).Errorf("Handshake sending")
				return errors.Errorf("failed to send Handshake response: %w", err)
			}
			log.WithContext(ctx).WithField("resp", resp).Debugf("HandshakeNodes response")
		}
	})

	return group.Wait()
}

func (service *registerArtowrk) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")
	task, err := service.Task(ctx)
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
			NodeKey: node.Key,
		})
	}

	resp := &pb.AcceptedNodesReply{
		Peers: peers,
		Error: &pb.Error{
			Status: pb.Error_OK,
		},
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")
	return resp, nil
}

func (service *registerArtowrk) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")
	task, err := service.Task(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.ConnectTo(ctx, req.ConnID, req.NodeKey); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{
		Error: &pb.Error{
			Status: pb.Error_OK,
		},
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")
	return resp, nil
}

func (service *registerArtowrk) SendImage(stream pb.RegisterArtowrk_SendImageServer) error {
	ctx := stream.Context()
	task, err := service.Task(ctx)
	if err != nil {
		return err
	}

	fileID, _ := random.String(16, random.Base62Chars)
	filename := filepath.Join(service.workDir, fileID)
	file, err := os.Create(filename)
	if err != nil {
		return errors.Errorf("failed to open file %q: %w", filename, err)
	}
	// defer func() {
	// 	os.Remove(filename)
	// 	log.WithContext(ctx).Debugf("Removed temp file %a", filename)
	// }()
	defer file.Close()
	log.WithContext(ctx).Debugf("Created temp file %a for uploading image", filename)

	wr := bufio.NewWriter(file)

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}

				log.WithContext(ctx).WithError(err).Errorf("SendImage receving")
				return errors.Errorf("failed to receive SendImage: %w", err)
			}

			if _, err := wr.Write(req.Payload); err != nil {
				return errors.Errorf("failed to write to file %q: %w", filename, err)
			}
		}
	})

	// TODO: pass filename to the task
	_ = task

	resp := &pb.SendImageReply{
		Error: &pb.Error{
			Status: pb.Error_OK,
		},
	}
	if err := stream.SendAndClose(resp); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("SendImage send error")
		return errors.New("failed to send SendImage response")
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("SendImage response")
	return nil
}

// Desc returns a description of the service.
func (service *registerArtowrk) Desc() *grpc.ServiceDesc {
	return &pb.RegisterArtowrk_ServiceDesc
}

// NewRegisterArtowrk returns a new registerArtowrk instance.
func NewRegisterArtowrk(service *artworkregister.Service, workDir string) pb.RegisterArtowrkServer {
	return &registerArtowrk{
		RegisterArtowrk: common.NewRegisterArtowrk(service),
		workDir:         workDir,
	}
}
