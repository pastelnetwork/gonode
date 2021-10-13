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

type externalDupeDetection struct {
	conn   *clientConn
	client pb.ExternalDupeDetectionClient
	sessID string
}

func (service *externalDupeDetection) SessID() string {
	return service.sessID
}

// Session implements node.ExternalDupeDetection.Session()
func (service *externalDupeDetection) Session(ctx context.Context, nodeID, sessID string) error {
	service.sessID = sessID

	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("failed to open Health stream: %w", err)
	}

	req := &pb.SessionRequest{
		NodeID: nodeID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if err := stream.Send(req); err != nil {
		return errors.Errorf("failed to send Session request: %w", err)
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
		return errors.Errorf("failed to receive Session response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Session response")

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

func (service *externalDupeDetection) SendEDDTicketSignature(ctx context.Context, nodeID string, signature []byte) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	_, err := service.client.SendEDDTicketSignature(ctx, &pb.SendTicketSignatureRequest{
		NodeID:    nodeID,
		Signature: signature,
	})

	return err
}

func (service *externalDupeDetection) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *externalDupeDetection) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newExternalDupeDetection(conn *clientConn) node.ExternalDupeDetection {
	return &externalDupeDetection{
		conn:   conn,
		client: pb.NewExternalDupeDetectionClient(conn),
	}
}
