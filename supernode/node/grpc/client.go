package grpc

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/middleware"
	"google.golang.org/grpc"
)

type Client struct {
}

func (client *Client) Connect(ctx context.Context, address string) (node.Connection, error) {
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(middleware.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(middleware.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, errors.Errorf("fail to dial: %v", err).WithField("address", address)
	}

	return NewConnection(conn), nil
}

func NewClient() node.Client {
	return &Client{}
}

type Connection struct {
	*grpc.ClientConn
	pb.SuperNodeClient
}

func (conn *Connection) RegisterArtowrk(ctx context.Context) (node.RegisterArtowrk, error) {
	stream, err := conn.SuperNodeClient.RegisterArtowrk(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	return NewRegisterArtowrk(stream), nil
}

func NewConnection(conn *grpc.ClientConn) *Connection {
	return &Connection{
		ClientConn:      conn,
		SuperNodeClient: pb.NewSuperNodeClient(conn),
	}
}
