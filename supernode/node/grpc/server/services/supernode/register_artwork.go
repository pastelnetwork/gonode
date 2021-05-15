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
	"google.golang.org/grpc/status"
)

// RegisterArtowrk represents grpc service for registration artowrk.
type RegisterArtowrk struct {
	pb.UnimplementedRegisterArtowrkServer

	*common.RegisterArtowrk
}

func (service *RegisterArtowrk) HealthCheck(stream pb.RegisterArtowrk_HealthCheckServer) error {
	ctx := stream.Context()

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}
	defer task.Cancel()

	go func() {
		defer task.Cancel()

		for {
			_, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.WithContext(ctx).Debug("Stream closed by peer")
				}

				switch status.Code(err) {
				case codes.Canceled, codes.Unavailable:
					log.WithContext(ctx).WithError(err).Debug("Stream closed")
				default:
					log.WithContext(ctx).WithError(err).Error("Stream closed")
				}
				return
			}
		}
	}()

	<-task.Done()
	return nil
}

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
