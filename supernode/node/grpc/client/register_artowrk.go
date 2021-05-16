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

	nodeID string
	connID string
}

func (service *registerArtowrk) ConnID() string {
	return service.connID
}

func (service *registerArtowrk) healthCheck(ctx context.Context) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDConnID(ctx)

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

func (service *registerArtowrk) Handshake(ctx context.Context) error {
	ctx = service.contextWithLogPrefix(ctx)

	req := &pb.HandshakeRequest{
		NodeID: service.nodeID,
		ConnID: service.connID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

	resp, err := service.client.Handshake(ctx, req)
	if err != nil {
		return errors.Errorf("failed to reqeust Handshake: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")

	return service.healthCheck(ctx)
}

func (service *registerArtowrk) contextWithMDConnID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeyConnID, service.connID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerArtowrk) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRegisterArtowrk(conn *clientConn, nodeID, connID string) node.RegisterArtowrk {
	return &registerArtowrk{
		nodeID: nodeID,
		connID: connID,
		conn:   conn,
		client: pb.NewRegisterArtowrkClient(conn),
	}
}
