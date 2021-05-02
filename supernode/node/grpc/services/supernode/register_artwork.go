package supernode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registerArtowrk struct {
	pb.SuperNode_RegisterArtowrkServer
	task *artworkregister.Task

	reqCh chan *pb.RegisterArtworkRequest
	errCh chan error
}

func (stream *registerArtowrk) handshake(ctx context.Context, req *pb.RegisterArtworkRequest_HandshakeRequest) error {
	if err := stream.task.HandshakeNode(ctx, req.NodeKey); err != nil {
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

func (stream *registerArtowrk) newError() *pb.RegisterArtworkReply_Error {
	return &pb.RegisterArtworkReply_Error{
		Status: pb.RegisterArtworkReply_Error_OK,
	}
}

func (stream *registerArtowrk) start(ctx context.Context, service *artworkregister.Service) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
				req := req.GetHandshake()

				if stream.task != nil {
					return errors.New("task is already registered")
				}

				stream.task = service.TaskByConnID(req.ConnID)
				if stream.task == nil {
					return errors.Errorf("connID %q not found", req.ConnID)
				}
				defer stream.task.Cancel()
				go func() {
					<-stream.task.Done()
					cancel()
				}()

				if err := stream.handshake(ctx, req); err != nil {
					return err
				}

			default:
				return errors.New("unsupported call")
			}
		}
	}
}

func newRegisterArtowrk(stream pb.SuperNode_RegisterArtowrkServer) *registerArtowrk {
	return &registerArtowrk{
		SuperNode_RegisterArtowrkServer: stream,
		reqCh:                           make(chan *pb.RegisterArtworkRequest),
		errCh:                           make(chan error),
	}
}
