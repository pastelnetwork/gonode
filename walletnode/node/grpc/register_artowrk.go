package grpc

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc/metadata"
)

const (
	uploadImageBufferSize = 32 * 1024
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

// Handshake implements node.RegisterArtowrk.Handshake()
func (service *registerArtowrk) Handshake(ctx context.Context, IsPrimary bool) error {
	ctx = service.contextWithLogPrefix(ctx)

	req := &pb.HandshakeRequest{
		IsPrimary: IsPrimary,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Handshake request")

	resp, err := service.client.Handshake(ctx, req)
	if err != nil {
		return errors.Errorf("failed to reqeust Handshake: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Handshake response")

	service.sessID = resp.SessID
	return service.healthCheck(ctx)
}

// AcceptedNodes implements node.RegisterArtowrk.AcceptedNodes()
func (service *registerArtowrk) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.AcceptedNodesRequest{}
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")

	resp, err := service.client.AcceptedNodes(ctx, req)
	if err != nil {
		return nil, errors.Errorf("failed to request to accepted secondary nodes: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")

	var ids []string
	for _, peer := range resp.Peers {
		ids = append(ids, peer.NodeID)
	}
	return ids, nil
}

// ConnectTo implements node.RegisterArtowrk.ConnectTo()
func (service *registerArtowrk) ConnectTo(ctx context.Context, nodeID, sessID string) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.ConnectToRequest{
		NodeID: nodeID,
		SessID: sessID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return errors.Errorf("failed to request to connect to primary node: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")

	return nil
}

// SendImage implements node.RegisterArtowrk.SendImage()
func (service *registerArtowrk) SendImage(ctx context.Context, filename string) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.SendImage(ctx)
	if err != nil {
		return errors.Errorf("failed to open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := os.Open(filename)
	if err != nil {
		return errors.Errorf("failed to open file %q: %w", filename, err)
	}
	defer file.Close()

	log.WithContext(ctx).WithField("filename", filename).Debugf("SendImage uploading")
	buffer := make([]byte, uploadImageBufferSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		req := &pb.SendImageRequest{
			Payload: buffer[:n],
		}
		if err := stream.Send(req); err != nil {
			return errors.Errorf("failed to send image data: %w", err).WithField("reqID", service.conn.id)
		}
	}
	log.WithContext(ctx).Debugf("SendImage uploaded")

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return errors.Errorf("failed to receive send image response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("SendImage response")

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
