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
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registerArtowrk struct {
	pb.WalletNode_RegisterArtowrkServer
	task *artworkregister.Task

	fileWriter *bufio.Writer
}

func (stream *registerArtowrk) handshake(ctx context.Context, req *pb.RegisterArtworkRequest_HandshakeRequest) error {
	if err := stream.task.Handshake(ctx, req.ConnID, req.IsPrimary); err != nil {
		return err
	}

	resp := &pb.RegisterArtworkReply{
		Replies: &pb.RegisterArtworkReply_Handshake{
			Handshake: &pb.RegisterArtworkReply_HandshakeReply{
				Error: stream.newError(),
			},
		},
	}
	return stream.send(ctx, resp)
}

func (stream *registerArtowrk) acceptedNodes(ctx context.Context, _ *pb.RegisterArtworkRequest_AcceptedNodesRequest) error {
	nodes, err := stream.task.AcceptedNodes(ctx)
	if err != nil {
		return err
	}

	var peers []*pb.RegisterArtworkReply_AcceptedNodesReply_Peer
	for _, node := range nodes {
		peers = append(peers, &pb.RegisterArtworkReply_AcceptedNodesReply_Peer{
			NodeKey: node.Key,
		})
	}

	resp := &pb.RegisterArtworkReply{
		Replies: &pb.RegisterArtworkReply_AcceptedNodes{
			AcceptedNodes: &pb.RegisterArtworkReply_AcceptedNodesReply{
				Peers: peers,
				Error: stream.newError(),
			},
		},
	}
	return stream.send(ctx, resp)
}

func (stream *registerArtowrk) connectTo(ctx context.Context, req *pb.RegisterArtworkRequest_ConnectToRequest) error {
	if err := stream.task.ConnectTo(ctx, req.NodeKey); err != nil {
		return err
	}

	resp := &pb.RegisterArtworkReply{
		Replies: &pb.RegisterArtworkReply_ConnectTo{
			ConnectTo: &pb.RegisterArtworkReply_ConnectToReply{
				Error: stream.newError(),
			},
		},
	}
	return stream.send(ctx, resp)
}

func (stream *registerArtowrk) uploadImage(ctx context.Context, req *pb.RegisterArtworkRequest_UploadImageRequest) error {

	if _, err := stream.fileWriter.Write(req.Payload); err != nil {
		return errors.Errorf("failed to write: %w", err)
	}

	// resp := &pb.RegisterArtworkReply{
	// 	Replies: &pb.RegisterArtworkReply_ConnectTo{
	// 		ConnectTo: &pb.RegisterArtworkReply_ConnectToReply{
	// 			Error: stream.newError(),
	// 		},
	// 	},
	// }

	// log.WithContext(ctx).WithField("resp", resp.String()).Debugf("Sending")
	// if err := stream.Send(resp); err != nil {
	// 	return errors.New(err)
	// }
	return nil
}

func (stream *registerArtowrk) newError() *pb.RegisterArtworkReply_Error {
	return &pb.RegisterArtworkReply_Error{
		Status: pb.RegisterArtworkReply_Error_OK,
	}
}

func (stream *registerArtowrk) send(ctx context.Context, resp *pb.RegisterArtworkReply) error {
	log.WithContext(ctx).WithField("resp", resp.String()).Debugf("Sending")

	if err := stream.Send(resp); err != nil {
		return errors.New(err)
	}
	return nil
}

func (stream *registerArtowrk) handle(ctx context.Context, service *artworkregister.Service, workDir string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream.task = service.NewTask(ctx)
	defer stream.task.Cancel()
	go func() {
		<-stream.task.Done()
		cancel()
	}()

	fileID, _ := random.String(16, random.Base62Chars)
	filename := filepath.Join(workDir, fileID)
	file, err := os.Create(filename)
	if err != nil {
		return errors.Errorf("failed to open file %q: %w", filename, err)
	}
	defer func() {
		os.Remove(filename)
		log.WithContext(ctx).Debugf("Removed temp file %a", filename)
	}()
	defer file.Close()
	log.WithContext(ctx).Debugf("Created temp file %a for uploading image", filename)

	stream.fileWriter = bufio.NewWriter(file)

	reqCh := make(chan *pb.RegisterArtworkRequest)
	errCh := make(chan error)
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				if status.Code(err) == codes.Canceled {
					err = errors.New("connection closed")
				}
				errCh <- errors.New(err)
				return
			}
			reqCh <- req
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		case req := <-reqCh:
			switch req.GetRequests().(type) {
			case *pb.RegisterArtworkRequest_UploadImage:
			default:
				log.WithContext(ctx).WithField("req", req.String()).Debugf("Receiving")
			}

			switch req.GetRequests().(type) {
			case *pb.RegisterArtworkRequest_Handshake:
				if err := stream.handshake(ctx, req.GetHandshake()); err != nil {
					return err
				}

			case *pb.RegisterArtworkRequest_AcceptedNodes:
				if err := stream.acceptedNodes(ctx, req.GetAcceptedNodes()); err != nil {
					return err
				}

			case *pb.RegisterArtworkRequest_ConnectTo:
				if err := stream.connectTo(ctx, req.GetConnectTo()); err != nil {
					return err
				}

			case *pb.RegisterArtworkRequest_UploadImage:
				if err := stream.uploadImage(ctx, req.GetUploadImage()); err != nil {
					return err
				}

			default:
				return errors.New("unsupported call")
			}
		}
	}
}

func newRegisterArtowrk(stream pb.WalletNode_RegisterArtowrkServer) *registerArtowrk {
	return &registerArtowrk{
		WalletNode_RegisterArtowrkServer: stream,
	}
}
