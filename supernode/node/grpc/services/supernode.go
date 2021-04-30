package services

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/log"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SuperNoder represents grpc service
type SuperNoder struct {
	pb.UnimplementedSuperNodeServer

	artworkRegister *artworkregister.Service
}

// RegisterArtowrk is responsible for communication between supernodes.
func (service *SuperNoder) RegisterArtowrk(stream pb.SuperNode_RegisterArtowrkServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	task := &artworkregister.Task{}

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
			case *pb.RegisterArtworkRequest_Handshake:
				log.WithContext(ctx).Debugf("Request Handshake")

				if task != nil {
					return errors.New("task is already registered")
				}

				task = service.artworkRegister.TaskByConnID(req.GetHandshake().ConnID)
				if task == nil {
					return errors.Errorf("connID %q not found", req.GetHandshake().ConnID)
				}

				if err := task.PrimaryAcceptSecondary(ctx, req.GetHandshake().NodeKey); err != nil {
					return err
				}
				defer task.Cancel()
				go func() {
					<-task.Done()
					cancel()
				}()

				repl := &pb.RegisterArtworkReply{
					Replies: &pb.RegisterArtworkReply_Handshake{
						Handshake: &pb.RegisterArtworkReply_HandshakeReply{
							Error: service.NewEmptyError(),
						},
					},
				}
				if err := stream.Send(repl); err != nil {
					return errors.New(err)
				}

			default:
				return errors.New("unsupported call")
			}
		}
	}
}

func (service *SuperNoder) NewEmptyError() *pb.RegisterArtworkReply_Error {
	return &pb.RegisterArtworkReply_Error{
		Status: pb.RegisterArtworkReply_Error_OK,
	}
}

func (service *SuperNoder) NewError(err error) *pb.RegisterArtworkReply_Error {
	return &pb.RegisterArtworkReply_Error{
		Status: pb.RegisterArtworkReply_Error_ERR,
		ErrMsg: err.Error(),
	}
}

// Desc returns a description of the service.
func (service *SuperNoder) Desc() *grpc.ServiceDesc {
	return &pb.SuperNode_ServiceDesc
}

// NewSuperNode returns a new SuperNoder instance.
func NewSuperNode(artworkRegister *artworkregister.Service) *SuperNoder {
	return &SuperNoder{
		artworkRegister: artworkRegister,
	}
}
