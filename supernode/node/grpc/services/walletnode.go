package services

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
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
				req := req.GetHandshake()

				log.WithContext(ctx).WithField("connID", req.ConnID).WithField("isPrimary", req.IsPrimary).Debugf("Request Handshake")

				if err := task.Handshake(task.Context(ctx), req.ConnID, req.IsPrimary); err != nil {
					return err
				}

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

			case *pb.RegisterArtworkRequest_PrimaryAcceptSecondary:
				log.WithContext(ctx).Debugf("Request PrimaryAcceptSecondary")

				nodes, err := task.PrimaryWaitSecondary(task.Context(ctx))
				if err != nil {
					return err
				}
				defer task.Cancel()
				go func() {
					<-task.Done()
					cancel()
				}()

				var peers []*pb.RegisterArtworkReply_PrimaryAcceptSecondaryReply_Peer
				for _, node := range nodes {
					peers = append(peers, &pb.RegisterArtworkReply_PrimaryAcceptSecondaryReply_Peer{
						NodeKey: node.Key,
					})
				}

				repl := &pb.RegisterArtworkReply{
					Replies: &pb.RegisterArtworkReply_PrimayAcceptSecondary{
						PrimayAcceptSecondary: &pb.RegisterArtworkReply_PrimaryAcceptSecondaryReply{
							Peers: peers,
							Error: service.NewEmptyError(),
						},
					},
				}
				if err := stream.Send(repl); err != nil {
					return errors.New(err)
				}

			case *pb.RegisterArtworkRequest_SecondaryConnectToPrimary:
				req := req.GetSecondaryConnectToPrimary()

				log.WithContext(ctx).Debugf("Request SecondaryConnectToPrimary")

				if err := task.SecondaryConnectToPrimary(task.Context(ctx), req.NodeKey); err != nil {
					return err
				}
				defer task.Cancel()
				go func() {
					<-task.Done()
					cancel()
				}()

				repl := &pb.RegisterArtworkReply{
					Replies: &pb.RegisterArtworkReply_SecondaryConnectToPrimary{
						SecondaryConnectToPrimary: &pb.RegisterArtworkReply_SecondaryConnectToPrimaryReply{
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

func (service *WalletNode) NewEmptyError() *pb.RegisterArtworkReply_Error {
	return &pb.RegisterArtworkReply_Error{
		Status: pb.RegisterArtworkReply_Error_OK,
	}
}

func (service *WalletNode) NewError(err error) *pb.RegisterArtworkReply_Error {
	return &pb.RegisterArtworkReply_Error{
		Status: pb.RegisterArtworkReply_Error_ERR,
		ErrMsg: err.Error(),
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
