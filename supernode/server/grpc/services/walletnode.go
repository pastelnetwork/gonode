package services

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/server/grpc/log"
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
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	task := service.artworkRegister.NewTask(ctx)

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
		case *pb.WalletNodeRegisterArtworkRequest_AcceptSecondaryNodes:
			log.WithContext(ctx).Debugf("Request AcceptSecondaryNodes")

			nodes, err := task.AcceptSecondaryNodes(ctx, req.AcceptSecondaryNodes.ConnID)
			if err != nil {
				return err
			}
			defer task.Cancel()

			var peers []*pb.WalletNodeRegisterArtworkReply_AcceptSecondaryNodesReply_Peer
			for _, node := range nodes {
				peers = append(peers, &pb.WalletNodeRegisterArtworkReply_AcceptSecondaryNodesReply_Peer{
					SecondaryKey: node.Key,
				})
			}

			res := &pb.WalletNodeRegisterArtworkReply{
				TestOneof: &pb.WalletNodeRegisterArtworkReply_AcceptSecondaryNodes{
					AcceptSecondaryNodes: &pb.WalletNodeRegisterArtworkReply_AcceptSecondaryNodesReply{
						Peers: peers,
					},
				},
			}

			if err := stream.Send(res); err != nil {
				return errors.New(err)
			}

		case *pb.WalletNodeRegisterArtworkRequest_ConnectToPrimaryNode:
			log.WithContext(ctx).Debugf("Request ConnectToPrimaryNode")

			if err := task.ConnectToPrimaryNode(ctx, req.ConnectToPrimaryNode.ConnID, req.ConnectToPrimaryNode.PrimaryNodeKey); err != nil {
				return err
			}
			defer task.Cancel()

			res := &pb.WalletNodeRegisterArtworkReply{
				TestOneof: &pb.WalletNodeRegisterArtworkReply_ConnectToPrimaryNode{},
			}

			if err := stream.Send(res); err != nil {
				return errors.New(err)
			}

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
