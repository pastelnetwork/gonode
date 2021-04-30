package grpc

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

type RegisterArtowrk struct {
	pb.SuperNode_RegisterArtowrkClient
}

func (stream *RegisterArtowrk) Handshake(ctx context.Context, connID, nodeKey string) error {
	req := &pb.RegisterArtworkRequest{
		Requests: &pb.RegisterArtworkRequest_Handshake{
			Handshake: &pb.RegisterArtworkRequest_HandshakeRequest{
				ConnID:  connID,
				NodeKey: nodeKey,
			},
		},
	}

	if err := stream.Send(req); err != nil {
		return errors.New(err)
	}

	res, err := stream.Recv()
	if err != nil {
		return errors.New(err)
	}
	resp := res.GetHandshake()

	if err := resp.Error; err.Status == pb.RegisterArtworkReply_Error_ERR {
		return errors.New(err.ErrMsg)
	}
	return nil
}

func NewRegisterArtowrk(stream pb.SuperNode_RegisterArtowrkClient) *RegisterArtowrk {
	return &RegisterArtowrk{
		SuperNode_RegisterArtowrkClient: stream,
	}
}
