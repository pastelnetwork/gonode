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
)

// RegisterArtowrk represents grpc service for registration artowrk.
type RegisterArtowrk struct {
	pb.UnimplementedRegisterArtowrkServer

	*common.RegisterArtowrk
}

func (service *RegisterArtowrk) Handshake(stream pb.RegisterArtowrk_HandshakeServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	errCh := make(chan error)
	reqCh := make(chan *pb.HandshakeRequest)

	go func() {
		defer cancel()
		for {
			req, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			reqCh <- req
		}
	}()

	select {
	case err := <-errCh:
		if err == io.EOF {
			return nil
		}
		return errors.Errorf("failed to receive Handshake request: %w", err)
	case req := <-reqCh:
		log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

		task, err := service.Task(ctx)
		if err != nil {
			return err
		}
		defer task.Cancel()
		go func() {
			<-task.Done()
			cancel()
		}()

		if err := task.HandshakeNode(ctx, req.NodeKey); err != nil {
			return err
		}

		resp := &pb.HandshakeReply{}
		if err := stream.Send(resp); err != nil {
			return errors.Errorf("failed to send Handshake response: %w", err)
		}
		log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")

		<-ctx.Done()
	}
	return nil
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
