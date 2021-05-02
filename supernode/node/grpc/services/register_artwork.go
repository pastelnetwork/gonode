package services

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registerArtowrk struct {
	pb.WalletNode_RegisterArtowrkServer
	task *artworkregister.Task

	reqCh chan *pb.RegisterArtworkRequest
	errCh chan error
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

	log.WithContext(ctx).WithField("resp", resp.String()).Debugf("Sending")
	if err := stream.Send(resp); err != nil {
		return errors.New(err)
	}
	return nil
}

func (stream *registerArtowrk) secondaryNodes(ctx context.Context, req *pb.RegisterArtworkRequest_ConnectedNodesRequest) error {
	nodes, err := stream.task.PrimaryWaitSecondary(ctx)
	if err != nil {
		return err
	}

	var peers []*pb.RegisterArtworkReply_ConnectedNodesReply_Peer
	for _, node := range nodes {
		peers = append(peers, &pb.RegisterArtworkReply_ConnectedNodesReply_Peer{
			NodeKey: node.Key,
		})
	}

	resp := &pb.RegisterArtworkReply{
		Replies: &pb.RegisterArtworkReply_ConnectedNodes{
			ConnectedNodes: &pb.RegisterArtworkReply_ConnectedNodesReply{
				Peers: peers,
				Error: stream.newError(),
			},
		},
	}

	log.WithContext(ctx).WithField("resp", resp.String()).Debugf("Sending")
	if err := stream.Send(resp); err != nil {
		return errors.New(err)
	}
	return nil
}

func (stream *registerArtowrk) connectToPrimary(ctx context.Context, req *pb.RegisterArtworkRequest_ConnectToRequest) error {
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

	log.WithContext(ctx).WithField("resp", resp.String()).Debugf("Sending")
	if err := stream.Send(resp); err != nil {
		return errors.New(err)
	}
	return nil
}

func (stream *registerArtowrk) newError() *pb.RegisterArtworkReply_Error {
	return &pb.RegisterArtworkReply_Error{
		Status: pb.RegisterArtworkReply_Error_OK,
	}
}

func (stream *registerArtowrk) start(ctx context.Context, service *artworkregister.Service) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream.task = service.NewTask(ctx)
	defer stream.task.Cancel()
	go func() {
		<-stream.task.Done()
		cancel()
	}()

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
				stream.errCh <- errors.New(err)
				return
			}
			stream.reqCh <- req
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-stream.errCh:
			return err
		case req := <-stream.reqCh:
			log.WithContext(ctx).WithField("req", req.String()).Debugf("Receiving")

			switch req.GetRequests().(type) {
			case *pb.RegisterArtworkRequest_Handshake:
				if err := stream.handshake(ctx, req.GetHandshake()); err != nil {
					return err
				}

			case *pb.RegisterArtworkRequest_ConnectedNodes:
				if err := stream.secondaryNodes(ctx, req.GetConnectedNodes()); err != nil {
					return err
				}

			case *pb.RegisterArtworkRequest_ConnectTo:
				if err := stream.connectToPrimary(ctx, req.GetConnectTo()); err != nil {
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
		reqCh:                            make(chan *pb.RegisterArtworkRequest),
		errCh:                            make(chan error),
	}
}
