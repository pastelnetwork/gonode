package supernode

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	sc "github.com/pastelnetwork/gonode/supernode/services/common"
)

// this implements SN's GRPC methods that are called by another SNs during Cascade Registration
// meaning - these methods implements server side of SN to SN GRPC communication

// RegisterCascade represents grpc service for registration Sense tickets.
type RegisterCascade struct {
	pb.UnimplementedRegisterCascadeServer

	*common.RegisterCascade
}

// Session implements supernode.RegisterSenseServer.Session()
func (service *RegisterCascade) Session(stream pb.RegisterCascade_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *cascaderegister.CascadeRegistrationTask
	isTaskNew := false

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewSenseRegistrationTask()
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

// SendCascadeTicketSignature implements supernode.RegisterCascadeServer.SendCascadeTicketSignature()
func (service *RegisterCascade) SendCascadeTicketSignature(ctx context.Context, req *pb.SendTicketSignatureRequest) (*pb.SendTicketSignatureReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendCascadeTicketSignature request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.AddPeerTicketSignature(req.NodeID, req.Signature, sc.StatusImageProbed); err != nil {
		return nil, errors.Errorf("add peer signature %w", err)
	}

	return &pb.SendTicketSignatureReply{}, nil
}

// Desc returns a description of the service.
func (service *RegisterCascade) Desc() *grpc.ServiceDesc {
	return &pb.RegisterCascade_ServiceDesc
}

// NewRegisterCascade returns a new RegisterCascade instance.
func NewRegisterCascade(service *cascaderegister.CascadeRegistrationService) *RegisterCascade {
	return &RegisterCascade{
		RegisterCascade: common.NewRegisterCascade(service),
	}
}
