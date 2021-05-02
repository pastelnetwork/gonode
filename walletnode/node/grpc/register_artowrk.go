package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registerArtowrk struct {
	conn *clientConn
	pb.WalletNode_RegisterArtowrkClient

	isClosed bool

	recvCh chan *pb.RegisterArtworkReply
	errCh  chan error
}

// Handshake implements node.RegisterArtowrk.Handshake()
func (stream *registerArtowrk) Handshake(ctx context.Context, connID string, IsPrimary bool) error {
	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_Handshake{
			Handshake: &pb.RegisterArtworkRequest_HandshakeRequest{
				ConnID:    connID,
				IsPrimary: IsPrimary,
			},
		},
	}

	res, err := stream.sendRecv(ctx, req)
	if err != nil {
		return err
	}

	resp := res.GetHandshake()
	if resp == nil {
		return errors.Errorf("wrong response, %q", res.String())
	}
	if err := resp.Error; err.Status == pb.RegisterArtworkReply_Error_ERR {
		return errors.New(err.ErrMsg)
	}
	return nil
}

// SecondaryNodes implements node.RegisterArtowrk.SecondaryNodes()
func (stream *registerArtowrk) SecondaryNodes(ctx context.Context) (node.SuperNodes, error) {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_SecondaryNodes{
			SecondaryNodes: &pb.RegisterArtworkRequest_SecondaryNodesRequest{},
		},
	}

	res, err := stream.sendRecv(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := res.GetSecondaryNodes()
	if resp == nil {
		return nil, errors.Errorf("wrong response, %q", res.String())
	}
	if err := resp.Error; err.Status == pb.RegisterArtworkReply_Error_ERR {
		return nil, errors.New(err.ErrMsg)
	}

	var nodes node.SuperNodes
	for _, peer := range resp.Peers {
		nodes = append(nodes, &node.SuperNode{
			Key: peer.NodeKey,
		})
	}
	return nodes, nil
}

// ConnectToPrimary implements node.RegisterArtowrk.ConnectToPrimary()
func (stream *registerArtowrk) ConnectToPrimary(ctx context.Context, nodeKey string) error {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_ConnectToPrimary{
			ConnectToPrimary: &pb.RegisterArtworkRequest_ConnectToPrimaryRequest{
				NodeKey: nodeKey,
			},
		},
	}

	res, err := stream.sendRecv(ctx, req)
	if err != nil {
		return err
	}

	resp := res.GetConnectToPrimary()
	if resp == nil {
		return errors.Errorf("wrong response, %q", res.String())
	}
	if err := resp.Error; err.Status == pb.RegisterArtworkReply_Error_ERR {
		return errors.New(err.ErrMsg)
	}
	return nil
}

func (stream *registerArtowrk) sendRecv(ctx context.Context, req *pb.RegisterArtworkRequest) (*pb.RegisterArtworkReply, error) {
	if err := stream.send(ctx, req); err != nil {
		return nil, err
	}

	resp, err := stream.recv(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (stream *registerArtowrk) send(ctx context.Context, req *pb.RegisterArtworkRequest) error {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	if stream.isClosed {
		return errors.New("stream closed")
	}

	log.WithContext(ctx).WithField("req", req.String()).Debugf("Sending")
	if err := stream.SendMsg(req); err != nil {
		switch status.Code(err) {
		case codes.Canceled:
			log.WithContext(ctx).WithError(err).Debugf("Sending canceled")
		default:
			log.WithContext(ctx).WithError(err).Errorf("Sending")
		}
		return err
	}

	return nil
}

func (stream *registerArtowrk) recv(ctx context.Context) (*pb.RegisterArtworkReply, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	case resp := <-stream.recvCh:
		return resp, nil
	case err := <-stream.errCh:
		return nil, err
	}
}

func (stream *registerArtowrk) start(ctx context.Context) error {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	go func() {
		defer func() {
			stream.isClosed = true
			stream.conn.Close()
		}()

		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.WithContext(ctx).Debugf("Stream closed by peer")
					return
				}
				switch status.Code(err) {
				case codes.Canceled, codes.Unavailable:
					log.WithContext(ctx).WithError(err).Debugf("Stream closed")
				default:
					log.WithContext(ctx).WithError(err).Errorf("Stream")
				}

				stream.errCh <- errors.New(err)
				break
			}
			log.WithContext(ctx).WithField("resp", resp.String()).Debugf("Receiving")

			stream.recvCh <- resp
		}
	}()

	return nil
}

func newRegisterArtowrk(conn *clientConn, client pb.WalletNode_RegisterArtowrkClient) node.RegisterArtowrk {
	return &registerArtowrk{
		conn:                             conn,
		WalletNode_RegisterArtowrkClient: client,

		recvCh: make(chan *pb.RegisterArtworkReply),
		errCh:  make(chan error),
	}
}
