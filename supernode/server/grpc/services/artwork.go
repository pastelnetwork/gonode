package services

import (
	"fmt"
	"io"

	"github.com/pastelnetwork/go-commons/errors"
	pb "github.com/pastelnetwork/supernode-proto"
	"github.com/pastelnetwork/supernode/services/artworkregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Artwork represents grpc service
type Artwork struct {
	pb.UnimplementedArtworkServer

	artworkRegister *artworkregister.Service
}

// SuperNode is responsible for communication between supernodes.
func (service *Artwork) SuperNode(stream pb.Artwork_SuperNodeServer) error {
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
		case *pb.SuperNodeRequest_Hello:
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

// Register is responsible for communication between walletnote and supernodes.
func (service *Artwork) Register(stream pb.Artwork_RegisterServer) error {
	var (
		ctx  = stream.Context()
		task *artworkregister.Task
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
		case *pb.RegisterRequest_PrimaryNode:
			if task != nil {
				return errors.New("task is already registered")
			}

			task = service.artworkRegister.NewTask(ctx, req.PrimaryNode.ConnID)
			defer task.Cancel()

			if err := task.AcceptSecondaryNodes(ctx); err != nil {
				return err
			}

			fmt.Println("primary", req.PrimaryNode.ConnID)
		case *pb.RegisterRequest_SecondaryNode:
			if task != nil {
				return errors.New("task is already registered")
			}

			task = service.artworkRegister.NewTask(ctx, req.SecondaryNode.ConnID)
			defer task.Cancel()

			if err := task.ConnectPrimaryNode(ctx, req.SecondaryNode.PrimaryNodeKey); err != nil {
				return err
			}

			fmt.Println("secondary", req.SecondaryNode.ConnID)
		default:
			return errors.New("unsupported call")
		}
	}
}

// Desc returns a description of the service.
func (service *Artwork) Desc() *grpc.ServiceDesc {
	return &pb.Artwork_ServiceDesc
}

// NewArtwork returns a new Artwork instance.
func NewArtwork(artworkRegister *artworkregister.Service) *Artwork {
	return &Artwork{
		artworkRegister: artworkRegister,
	}
}
