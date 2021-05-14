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

	connID string
}

func (service *registerArtowrk) ConndID() string {
	return service.connID
}

func (service *registerArtowrk) Handshake(ctx context.Context, nodeKey string) error {
	ctx = service.context(ctx)

	stream, err := service.client.Handshake(ctx)
	if err != nil {
		return errors.New("failed to open handshake stream")
	}

	req := &pb.HandshakeRequest{
		NodeKey: nodeKey,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

	if err := stream.Send(req); err != nil {
		return errors.New("failed to send handshake request")
	}

	errCh := make(chan error)
	respCh := make(chan *pb.HandshakeReply)

	go func() {
		defer service.conn.Close()

		for {
			resp, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")
			respCh <- resp
		}
	}()

	select {
	case resp := <-respCh:
		service.connID = resp.ConnID
	case err := <-errCh:
		return errors.Errorf("failed to receive Handshake: %w", err)
	}
	return nil
}

func (service *registerArtowrk) context(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeyConnID, service.connID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
	return ctx
}

func newRegisterArtowrk(conn *clientConn, connID string) node.RegisterArtowrk {
	return &registerArtowrk{
		connID: connID,
		conn:   conn,
		client: pb.NewRegisterArtowrkClient(conn),
	}
}
