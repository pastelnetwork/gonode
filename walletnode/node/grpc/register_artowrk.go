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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

// Session implements node.RegisterArtowrk.Session()
func (service *registerArtowrk) Session(ctx context.Context, isPrimary bool) error {
	ctx = service.contextWithLogPrefix(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("failed to open Health stream: %w", err)
	}

	req := &pb.SessionRequest{
		IsPrimary: isPrimary,
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
	service.sessID = resp.SessID

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

// UploadImage implements node.RegisterArtowrk.UploadImage()
func (service *registerArtowrk) UploadImage(ctx context.Context, filename string) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.UploadImage(ctx)
	if err != nil {
		return errors.Errorf("failed to open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := os.Open(filename)
	if err != nil {
		return errors.Errorf("failed to open file %q: %w", filename, err)
	}
	defer file.Close()

	log.WithContext(ctx).WithField("filename", filename).Debugf("UploadImage uploading")
	buffer := make([]byte, uploadImageBufferSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		req := &pb.UploadImageRequest{
			Payload: buffer[:n],
		}
		if err := stream.Send(req); err != nil {
			return errors.Errorf("failed to send image data: %w", err).WithField("reqID", service.conn.id)
		}
	}
	log.WithContext(ctx).Debugf("UploadImage uploaded")

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return errors.Errorf("failed to receive send image response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("UploadImage response")

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
