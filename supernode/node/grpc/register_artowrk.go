package grpc

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegisterArtowrk struct {
	conn *Connection
	pb.SuperNode_RegisterArtowrkClient
}

func (stream *RegisterArtowrk) Handshake(ctx context.Context, connID, nodeKey string) error {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, stream.conn.id))

	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_Handshake{
			Handshake: &pb.RegisterArtworkRequest_HandshakeRequest{
				ConnID:  connID,
				NodeKey: nodeKey,
			},
		},
	}

	log.WithContext(ctx).Debugf("Request Handshake")
	if err := stream.Send(req); err != nil {
		if status.Code(err) == codes.Canceled {
			log.WithContext(ctx).Debugf("Canceled request Handshake")
		} else {
			log.WithContext(ctx).WithError(err).Errorf("Request Handshake")
		}
		return errors.New(err)
	}

	res, err := stream.Recv()
	if err != nil {
		if status.Code(err) == codes.Canceled {
			log.WithContext(ctx).Debugf("Canceled response Handshake")
		} else {
			log.WithContext(ctx).WithError(err).Errorf("Response Handshake")
		}
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

func NewRegisterArtowrk(conn *Connection, stream pb.SuperNode_RegisterArtowrkClient) *RegisterArtowrk {
	return &RegisterArtowrk{
		conn:                            conn,
		SuperNode_RegisterArtowrkClient: stream,
	}
}
