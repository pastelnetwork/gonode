package client

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/supernode/node"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// this implements SN's GRPC methods that call another SN during Cascade Registration
// meaning - these methods implements client side of SN to SN GRPC communication

type registerCascade struct {
	conn   *clientConn
	client pb.RegisterCascadeClient
	sessID string
}

func (service *registerCascade) SessID() string {
	return service.sessID
}

// Session implements node.RegisterNft.Session()
func (service *registerCascade) Session(ctx context.Context, nodeID, sessID string) error {
	service.sessID = sessID

	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("open Health stream: %w", err)
	}

	req := &pb.SessionRequest{
		NodeID: nodeID,
	}

	if err := stream.Send(req); err != nil {
		return errors.Errorf("send Session request: %w", err)
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
		return errors.Errorf("receive Session response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("Session response")

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

// SendCascadeTicketSignature implements SendCascadeTicketSignature
func (service *registerCascade) SendCascadeTicketSignature(ctx context.Context, nodeID string, signature []byte) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.SendCascadeTicketSignature(ctx, &pb.SendTicketSignatureRequest{
		NodeID:    nodeID,
		Signature: signature,
	})

	return err
}

func newRegisterCascade(conn *clientConn) node.RegisterCascadeInterface {
	return &registerCascade{
		conn:   conn,
		client: pb.NewRegisterCascadeClient(conn),
	}
}
