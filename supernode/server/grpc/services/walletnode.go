package services

import (
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WalletNode represents grpc service
type WalletNode struct {
	pb.UnimplementedWalletNodeServer

	artworkRegister *artworkregister.Service
}

// RegisterArtowrk is responsible for communication between walletnote and supernodes.
func (service *WalletNode) RegisterArtowrk(stream pb.WalletNode_RegisterArtowrkServer) error {
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
		case *pb.WalletNodeRegisterArtworkRequest_PrimaryNode:
			if task != nil {
				return errors.New("task is already registered")
			}

			task = service.artworkRegister.NewTask(ctx, req.PrimaryNode.ConnID)
			defer task.Cancel()

			if err := task.AcceptSecondaryNodes(ctx); err != nil {
				return err
			}

			fmt.Println("primary", req.PrimaryNode.ConnID)
		case *pb.WalletNodeRegisterArtworkRequest_SecondaryNode:
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
func (service *WalletNode) Desc() *grpc.ServiceDesc {
	return &pb.WalletNode_ServiceDesc
}

// NewWalletNode returns a new WalletNode instance.
func NewWalletNode(artworkRegister *artworkregister.Service) *WalletNode {
	return &WalletNode{
		artworkRegister: artworkRegister,
	}
}
