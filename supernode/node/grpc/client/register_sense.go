package client

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
)

// this implements SN's GRPC methods that call another SN during Sense Registration
// meaning - these methods implements client side of SN to SN GRPC communication

type registerSense struct {
	conn   *clientConn
	client pb.RegisterSenseClient
	sessID string
}

func (service *registerSense) SessID() string {
	return service.sessID
}

// Session implements node.RegisterNft.Session()
func (service *registerSense) Session(ctx context.Context, nodeID, sessID string) error {
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
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.SendSignedDDAndFingerprints(ctx, &pb.SendSignedDDAndFingerprintsRequest{
		SessID:                    sessionID,
		NodeID:                    fromNodeID,
		ZstdCompressedFingerprint: compressedDDAndFingerprints,
	})

	return err
}

// SendSenseTicketSignature implements SendSenseTicketSignature
func (service *registerSense) SendSenseTicketSignature(ctx context.Context, nodeID string, signature []byte) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.SendSenseTicketSignature(ctx, &pb.SendTicketSignatureRequest{
		NodeID:    nodeID,
		Signature: signature,
	})

	return err
}

func newRegisterSense(conn *clientConn) node.RegisterSenseInterface {
	return &registerSense{
		conn:   conn,
		client: pb.NewRegisterSenseClient(conn),
	}
}
