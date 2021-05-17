package client

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type registerArtowrk struct {
	conn   *clientConn
	client pb.RegisterArtowrkClient
	sessID string
}

func (service *registerArtowrk) SessID() string {
	return service.sessID
}

// Handshake implements node.RegisterArtowrk.Handshake()
func (service *registerArtowrk) Handshake(ctx context.Context, nodeID, sessID string) error {
	service.sessID = sessID

	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.Handshake(ctx)
	if err != nil {
		return errors.Errorf("failed to open Health stream: %w", err)
	}

	req := &pb.HandshakeRequest{
		NodeID: nodeID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

	if err := stream.Send(req); err != nil {
		return errors.Errorf("failed to send Handshake request: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		switch status.Code(err) {
		case codes.Canceled, codes.Unavailable:
			return nil
		}
		return errors.Errorf("failed to receive Handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")

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
