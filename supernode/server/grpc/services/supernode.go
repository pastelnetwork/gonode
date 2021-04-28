package services

import (
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto"
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
	var (
		ctx  = stream.Context()
		task = new(artworkregister.Task)
	)

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			if status.Code(err) == codes.Canceled {
				return errors.New("connection closed")
			}
			return err
		}

		switch req := req.GetTestOneof().(type) {
		case *pb.SuperNodeRegisterArtworkRequest_Hello:
			if task != nil {
				return errors.New("task is already registered")
			}

			task = service.artworkRegister.Task(req.Hello.ConnID)
			if task == nil {
				return errors.Errorf("connID %q not found", req.Hello.ConnID)
			}

			if err := task.RegisterSecondaryNode(ctx, req.Hello.SecondaryNodeKey); err != nil {
				return err
			}

		default:
			return errors.New("unsupported call")
		}
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
