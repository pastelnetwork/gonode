package supernode

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/collectionregister"
)

// RegisterCollection represents grpc service for registration collection tickets.
type RegisterCollection struct {
	pb.UnimplementedRegisterCollectionServer

	*common.RegisterCollection
}

// Session implements supernode.RegisterSenseServer.Session()
func (service *RegisterCollection) Session(stream pb.RegisterCollection_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *collectionregister.CollectionRegistrationTask
	isTaskNew := false

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewCollectionRegistrationTask()
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

// SendCollectionTicketSignature implements supernode.RegisterCollectionServer.SendCollectionTicketSignature()
func (service *RegisterCollection) SendCollectionTicketSignature(ctx context.Context, req *pb.SendTicketSignatureRequest) (*pb.SendTicketSignatureReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendSenseTicketSignature request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.AddPeerCollectionTicketSignature(req.NodeID, req.Signature); err != nil {
		return nil, errors.Errorf("add peer signature %w", err)
	}

	return &pb.SendTicketSignatureReply{}, nil
}

// NewRegisterCollection returns a new RegisterCollection instance.
func NewRegisterCollection(service *collectionregister.CollectionRegistrationService) *RegisterCollection {
	return &RegisterCollection{
		RegisterCollection: common.NewRegisterCollection(service),
	}
}
