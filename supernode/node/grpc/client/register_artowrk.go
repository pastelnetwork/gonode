package client

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc/metadata"
)

type registerArtowrk struct {
	conn   *clientConn
	client pb.RegisterArtowrkClient
	sessID string
}

func (service *registerArtowrk) SessID() string {
	return service.sessID
}

func (service *registerArtowrk) healthCheck(ctx context.Context) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.Health(ctx)
	if err != nil {
		return errors.Errorf("failed to open Health stream: %w", err)
	}

	go func() {
		defer service.conn.Close()

		for {
			if _, err := stream.Recv(); err != nil {
				return
			}
		}
	}()

	return nil
}

func (service *registerArtowrk) Handshake(ctx context.Context, nodeID, sessID string) error {
	ctx = service.contextWithLogPrefix(ctx)

	req := &pb.HandshakeRequest{
		NodeID: nodeID,
		SessID: sessID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

	resp, err := service.client.Handshake(ctx, req)
	if err != nil {
		return errors.Errorf("failed to reqeust Handshake: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")
	service.sessID = sessID

	return service.healthCheck(ctx)
}

func (service *registerArtowrk) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerArtowrk) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRegisterArtowrk(conn *clientConn) node.RegisterArtowrk {
	return &registerArtowrk{
		conn:   conn,
		client: pb.NewRegisterArtowrkClient(conn),
	}
}
