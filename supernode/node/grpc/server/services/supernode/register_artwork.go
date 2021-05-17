package supernode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterArtowrk represents grpc service for registration artowrk.
type RegisterArtowrk struct {
	pb.UnimplementedRegisterArtowrkServer

	*common.RegisterArtowrk
}

// Handshake implements supernode.RegisterArtowrkServer.Handshake()
func (service *RegisterArtowrk) Handshake(stream pb.RegisterArtowrk_HandshakeServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *artworkregister.Task

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewTask(ctx)
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Handshake stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Handshake stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("failed to receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

	if err := task.HandshakeNode(ctx, req.NodeID); err != nil {
		return err
	}

	resp := &pb.HandshakeReply{
		SessID: task.ID,
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("failed to send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")

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
func (service *RegisterArtowrk) Desc() *grpc.ServiceDesc {
	return &pb.RegisterArtowrk_ServiceDesc
}

// NewRegisterArtowrk returns a new RegisterArtowrk instance.
func NewRegisterArtowrk(service *artworkregister.Service) *RegisterArtowrk {
	return &RegisterArtowrk{
		RegisterArtowrk: common.NewRegisterArtowrk(service),
	}
}
