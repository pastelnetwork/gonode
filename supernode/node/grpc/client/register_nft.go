package client

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registerNft struct {
	conn   *clientConn
	client pb.RegisterNftClient
	sessID string
}

func (service *registerNft) SessID() string {
	return service.sessID
}

// Session implements node.RegisterNft.Session()
func (service *registerNft) Session(ctx context.Context, nodeID, sessID string) error {
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

// SendSignedDDAndFingerprints implements SendSignedDDAndFingerprints
func (service *registerNft) SendSignedDDAndFingerprints(ctx context.Context, sessionID string, fromNodeID string, compressedDDAndFingerprints []byte) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.SendSignedDDAndFingerprints(ctx, &pb.SendSignedDDAndFingerprintsRequest{
		SessID:                    sessionID,
		NodeID:                    fromNodeID,
		ZstdCompressedFingerprint: compressedDDAndFingerprints,
	})

	return err
}

// SendNftTicketSignature implements SendNftTicketSignature
func (service *registerNft) SendNftTicketSignature(ctx context.Context, nodeID string, signature []byte) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.SendNftTicketSignature(ctx, &pb.SendTicketSignatureRequest{
		NodeID:    nodeID,
		Signature: signature,
	})

	return err
}

func newRegisterNft(conn *clientConn) node.RegisterNftInterface {
	return &registerNft{
		conn:   conn,
		client: pb.NewRegisterNftClient(conn),
	}
}
