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

	connID string
}

func (service *registerArtowrk) ConndID() string {
	return service.connID
}

// Handshake implements node.RegisterArtowrk.Handshake()
func (service *registerArtowrk) Handshake(ctx context.Context, IsPrimary bool) error {
	ctx = service.context(ctx)

	stream, err := service.client.Handshake(ctx)
	if err != nil {
		return errors.New("failed to open handshake stream")
	}

	req := &pb.HandshakeRequest{
		IsPrimary: IsPrimary,
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

// AcceptedNodes implements node.RegisterArtowrk.AcceptedNodes()
func (service *registerArtowrk) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.context(ctx)

	req := &pb.AcceptedNodesRequest{}
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")

	resp, err := service.client.AcceptedNodes(ctx, req)
	if err != nil {
		return nil, errors.New("failed to request to accepted secondary nodes")
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")

	var ids []string
	for _, peer := range resp.Peers {
		ids = append(ids, peer.NodeKey)
	}
	return ids, nil
}

// ConnectTo implements node.RegisterArtowrk.ConnectTo()
func (service *registerArtowrk) ConnectTo(ctx context.Context, nodeKey, connID string) error {
	ctx = service.context(ctx)

	req := &pb.ConnectToRequest{
		NodeKey: nodeKey,
		ConnID:  connID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return errors.New("failed to request to connect to primary node")
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")

	return nil
}

// SendImage implements node.RegisterArtowrk.SendImage()
func (service *registerArtowrk) SendImage(ctx context.Context, filename string) error {
	ctx = service.context(ctx)

	stream, err := service.client.SendImage(ctx)
	if err != nil {
		return errors.New("failed to open stream")
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
			return errors.New("failed to send image data")
		}
	}
	log.WithContext(ctx).Debugf("SendImage uploaded")

	return nil
}

func (service *registerArtowrk) context(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeyConnID, service.connID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
	return ctx
}

func newRegisterArtowrk(conn *clientConn) node.RegisterArtowrk {
	return &registerArtowrk{
		conn:   conn,
		client: pb.NewRegisterArtowrkClient(conn),
	}
}
