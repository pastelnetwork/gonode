package grpc

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/node/grpc/middleware"
	"google.golang.org/grpc"
)

const (
	logPrefix = "grpc"
)

type Client struct{}

func (client *Client) Connect(ctx context.Context, address string) (node.Connection, error) {
	ctx = context.WithValue(ctx, log.PrefixKey, logPrefix)

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
	pb.WalletNodeClient
}

func (conn *Connection) RegisterArtowrk(ctx context.Context) (node.RegisterArtowrk, error) {
	stream, err := conn.WalletNodeClient.RegisterArtowrk(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	return NewRegisterArtowrk(stream), nil
}

func NewConnection(conn *grpc.ClientConn) *Connection {
	return &Connection{
		ClientConn:       conn,
		WalletNodeClient: pb.NewWalletNodeClient(conn),
	}
}
