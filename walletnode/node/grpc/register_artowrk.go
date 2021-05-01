package grpc

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type RegisterArtowrk struct {
	conn *Connection
	pb.WalletNode_RegisterArtowrkClient
}

func (stream *RegisterArtowrk) Handshake(ctx context.Context, connID string, IsPrimary bool) error {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_Handshake{
			Handshake: &pb.RegisterArtworkRequest_HandshakeRequest{
				ConnID:    connID,
				IsPrimary: IsPrimary,
			},
		},
	}

	log.WithContext(ctx).Debugf("Request Handshake")
	if err := stream.Send(req); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Request Handshake")
		return errors.New(err)
	}

	res, err := stream.Recv()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Response Handshake")
		return errors.New(err)
	}

	resp := res.GetHandshake()
	if resp == nil {
		return errors.New("Wrong response type, Handshake")
	}

	log.WithContext(ctx).Debugf("Response Handshake")

	if err := resp.Error; err.Status == pb.RegisterArtworkReply_Error_ERR {
		return errors.New(err.ErrMsg)
	}
	return nil
}

func (stream *RegisterArtowrk) PrimaryAcceptSecondary(ctx context.Context) (node.SuperNodes, error) {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_PrimaryAcceptSecondary{
			PrimaryAcceptSecondary: &pb.RegisterArtworkRequest_PrimaryAcceptSecondaryRequest{},
		},
	}

	log.WithContext(ctx).Debugf("Request PrimaryAcceptSecondary")
	if err := stream.Send(req); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Request PrimaryAcceptSecondary")
		return nil, errors.New(err)
	}

	res, err := stream.Recv()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Response PrimaryAcceptSecondary")
		return nil, errors.New(err)
	}
	resp := res.GetPrimayAcceptSecondary()
	if resp == nil {
		return nil, errors.New("Wrong response type, PrimaryAcceptSecondary")
	}

	log.WithContext(ctx).Debugf("Response PrimaryAcceptSecondary")

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

func (stream *RegisterArtowrk) SecondaryConnectToPrimary(ctx context.Context, nodeKey string) error {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_SecondaryConnectToPrimary{
			SecondaryConnectToPrimary: &pb.RegisterArtworkRequest_SecondaryConnectToPrimaryRequest{
				NodeKey: nodeKey,
			},
		},
	}

	log.WithContext(ctx).Debugf("Request SecondaryConnectToPrimary")
	if err := stream.Send(req); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Request SecondaryConnectToPrimary")
		return errors.New(err)
	}

	res, err := stream.Recv()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Response SecondaryConnectToPrimary")
		return errors.New(err)
	}
	resp := res.GetSecondaryConnectToPrimary()
	if resp == nil {
		return errors.New("Wrong response type, SecondaryConnectToPrimary")
	}
	log.WithContext(ctx).Debugf("Response SecondaryConnectToPrimary")

	if err := resp.Error; err.Status == pb.RegisterArtworkReply_Error_ERR {
		return errors.New(err.ErrMsg)
	}
	return nil
}

func NewRegisterArtowrk(conn *Connection, stream pb.WalletNode_RegisterArtowrkClient) *RegisterArtowrk {
	return &RegisterArtowrk{
		conn:                             conn,
		WalletNode_RegisterArtowrkClient: stream,
	}
}
