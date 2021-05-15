package supernode

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// RegisterArtowrk represents grpc service for registration artowrk.
type RegisterArtowrk struct {
	pb.UnimplementedRegisterArtowrkServer

	*common.RegisterArtowrk
}

// Health implements supernode.RegisterArtowrkServer.Health()
func (service *RegisterArtowrk) Health(stream pb.RegisterArtowrk_HealthServer) error {
	ctx := stream.Context()

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Helath stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Helath stream closed")

	go func() {
		defer task.Cancel()
		for {
			if _, err := stream.Recv(); err != nil {
				return
			}
		}
	}()

	<-task.Done()
	return nil
}

// Handshake implements supernode.RegisterArtowrkServer.Handshake()
func (service *RegisterArtowrk) Handshake(ctx context.Context, req *pb.HandshakeRequest) (*pb.HandshakeReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

	task := service.Task(req.ConnID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", req.ConnID)
	}

	if err := task.HandshakeNode(ctx, req.NodeID); err != nil {
		return nil, err
	}

	resp := &pb.HandshakeReply{}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")
	return resp, nil
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
