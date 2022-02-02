package walletnode

import (
	"context"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	"google.golang.org/grpc"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
)

// this implements SN's GRPC methods that are called by WNs during Sense Registration
// meaning - these methods implements server side of WN to SN GRPC communication

// RegisterSense represents grpc service for registration Sense.
type RegisterCascade struct {
	pb.UnimplementedRegisterCascadeServer

	*common.RegisterCascade
}

// Session implements walletnode.RegisterSenseServer.Session()
func (service *RegisterCascade) Session(stream pb.RegisterCascade_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *cascaderegister.CascadeRegistrationTask

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewTask()
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

	if err := task.NetworkHandler.Session(ctx, req.IsPrimary); err != nil {
		return err
	}

	resp := &pb.SessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("Session response")

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

// Desc returns a description of the service.
func (service *RegisterCascade) Desc() *grpc.ServiceDesc {
	return &pb.RegisterCascade_ServiceDesc
}

// NewRegisterSense returns a new RegisterSense instance.
func NewRegisterCascade(service *cascaderegister.CascadeRegistrationService) *RegisterCascade {
	return &RegisterCascade{
		RegisterCascade: common.NewRegisterCascade(service),
	}
}
