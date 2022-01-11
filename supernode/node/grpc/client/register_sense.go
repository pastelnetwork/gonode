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

type registerSense struct {
	conn   *clientConn
	client pb.RegisterSenseClient
	sessID string
}

func (service *registerSense) SessID() string {
	return service.sessID
}

// Session implements node.RegisterArtwork.Session()
func (service *registerSense) Session(ctx context.Context, nodeID, sessID string) error {
	service.sessID = sessID

	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("open Health stream: %w", err)
	}

	req := &pb.SessionRequest{
		NodeID: nodeID,
	}
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

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

// SendSignedDDAndFingerprints implements SendSignedDDAndFingerprints
func (service *registerSense) SendSignedDDAndFingerprints(ctx context.Context, sessionID string, fromNodeID string, compressedDDAndFingerprints []byte) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	_, err := service.client.SendSignedDDAndFingerprints(ctx, &pb.SendSignedDDAndFingerprintsRequest{
		SessID:                    sessionID,
		NodeID:                    fromNodeID,
		ZstdCompressedFingerprint: compressedDDAndFingerprints,
	})

	return err
}

// SendArtTicketSignature implements SendArtTicketSignature
func (service *registerSense) SendArtTicketSignature(ctx context.Context, nodeID string, signature []byte) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	_, err := service.client.SendArtTicketSignature(ctx, &pb.SendArtTicketSignatureRequest{
		NodeID:    nodeID,
		Signature: signature,
	})

	return err
}

func (service *registerSense) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerSense) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRegisterSense(conn *clientConn) node.RegisterSense {
	return &registerSense{
		conn:   conn,
		client: pb.NewRegisterSenseClient(conn),
	}
}
