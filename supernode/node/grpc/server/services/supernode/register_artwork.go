package supernode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/nftregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterNft represents grpc service for NFT registration.
type RegisterNft struct {
	pb.UnimplementedRegisterNftServer

	*common.RegisterNft
}

// Session implements supernode.RegisterNftServer.Session()
func (service *RegisterNft) Session(stream pb.RegisterNft_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *nftregister.NftRegistrationTask
	isTaskNew := false

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewTask()
		isTaskNew = true
	}
	go func() {
		<-task.Done()
		cancel()
	}()

	if isTaskNew {
		defer task.Cancel()
	}

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("receive handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if err := task.NetworkHandler.SessionNode(ctx, req.NodeID); err != nil {
		return err
	}

	if !isTaskNew {
		defer task.Cancel()
	}

	resp := &pb.SessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Session response")

	for {
		if _, err := stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			switch status.Code(err) {
			case codes.Canceled, codes.Unavailable:
				return nil
			}
			return errors.Errorf("handshake stream closed: %w", err)
		}
	}
}

// SendSignedDDAndFingerprints sends signed dd and fp
func (service *RegisterNft) SendSignedDDAndFingerprints(ctx context.Context, req *pb.SendSignedDDAndFingerprintsRequest) (*pb.SendSignedDDAndFingerprintsReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendSignedDDAndFingerprints request")
	task := service.Task(req.SessID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", req.SessID)
	}

	return &pb.SendSignedDDAndFingerprintsReply{}, task.AddSignedDDAndFingerprints(req.NodeID, req.ZstdCompressedFingerprint)

}

// SendNftTicketSignature implements supernode.RegisterNftServer.SendNftTicketSignature()
func (service *RegisterNft) SendNftTicketSignature(ctx context.Context, req *pb.SendNftTicketSignatureRequest) (*pb.SendNftTicketSignatureReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendNftTicketSignature request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.AddPeerTicketSignature(req.NodeID, req.Signature); err != nil {
		return nil, errors.Errorf("add peer signature %w", err)
	}

	return &pb.SendNftTicketSignatureReply{}, nil
}

// Desc returns a description of the service.
func (service *RegisterNft) Desc() *grpc.ServiceDesc {
	return &pb.RegisterNft_ServiceDesc
}

// NewRegisterNft returns a new RegisterNft instance.
func NewRegisterNft(service *nftregister.NftRegistrationService) *RegisterNft {
	return &RegisterNft{
		RegisterNft: common.NewRegisterNft(service),
	}
}
