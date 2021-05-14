package supernode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	metadataConnID = "connID"
)

// registerArtowrk represents grpc service for registration artowrk.
type registerArtowrk struct {
	pb.UnimplementedRegisterArtowrkServer

	*artworkregister.Service
}

func (service *registerArtowrk) Handshake(stream pb.RegisterArtowrk_HandshakeServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}

				log.WithContext(ctx).WithError(err).Errorf("Handshake receving")
				return errors.Errorf("failed to receive Handshake: %w", err)
			}
			log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

			task, err := service.task(ctx)
			if err != nil {
				return err
			}
			defer task.Cancel()

			if err := task.HandshakeNode(ctx, req.NodeKey); err != nil {
				return err
			}

			resp := &pb.HandshakeReply{
				Error: &pb.Error{
					Status: pb.Error_OK,
				},
			}
			if err := stream.SendAndClose(resp); err != nil {
				log.WithContext(ctx).WithError(err).Errorf("Handshake sending")
				return errors.Errorf("failed to send Handshake response: %w", err)
			}
			log.WithContext(ctx).WithField("resp", resp).Debugf("HandshakeNodes response")
		}
	})

	return group.Wait()
}

func (service *registerArtowrk) task(ctx context.Context) (*artworkregister.Task, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("not found metadata")
	}

	mdVals := md.Get(metadataConnID)
	if len(mdVals) == 0 {
		return nil, errors.Errorf("not found %q in metadata", metadataConnID)
	}
	connID := mdVals[0]

	task := service.Task(connID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", connID)
	}
	return task, nil
}

// Desc returns a description of the service.
func (service *registerArtowrk) Desc() *grpc.ServiceDesc {
	return &pb.RegisterArtowrk_ServiceDesc
}

// NewRegisterArtowrk returns a new registerArtowrk instance.
func NewRegisterArtowrk(service *artworkregister.Service) pb.RegisterArtowrkServer {
	return &registerArtowrk{
		Service: service,
	}
}
