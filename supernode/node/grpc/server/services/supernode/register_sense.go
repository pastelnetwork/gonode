package supernode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/senseregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterSense represents grpc service for registration Sense tickets.
type RegisterSense struct {
	pb.UnimplementedRegisterSenseServer

	*common.RegisterSense
}

// Session implements supernode.RegisterSenseServer.Session()
func (service *RegisterSense) Session(stream pb.RegisterSense_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *senseregister.SenseRegistrationTask
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
func (service *RegisterSense) SendSignedDDAndFingerprints(ctx context.Context, req *pb.SendSignedDDAndFingerprintsRequest) (*pb.SendSignedDDAndFingerprintsReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendSignedDDAndFingerprints request")
	task := service.Task(req.SessID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", req.SessID)
	}

	return &pb.SendSignedDDAndFingerprintsReply{}, task.AddSignedDDAndFingerprints(req.NodeID, req.ZstdCompressedFingerprint)

}

// SendSenseTicketSignature implements supernode.RegisterSenseServer.SendSenseTicketSignature()
func (service *RegisterSense) SendSenseTicketSignature(ctx context.Context, req *pb.SendNftTicketSignatureRequest) (*pb.SendNftTicketSignatureReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendSenseTicketSignature request")
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
func (service *RegisterSense) Desc() *grpc.ServiceDesc {
	return &pb.RegisterSense_ServiceDesc
}

// NewRegisterSense returns a new RegisterSense instance.
func NewRegisterSense(service *senseregister.SenseRegistrationService) *RegisterSense {
	return &RegisterSense{
		RegisterSense: common.NewRegisterSense(service),
	}
}
