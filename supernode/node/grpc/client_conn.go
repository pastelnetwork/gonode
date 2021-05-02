package grpc

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc"
)

type streamService interface {
	start(ctx context.Context) error
}

// clientConn represents grpc client conneciton.
type clientConn struct {
	*commongrpc.ClientConn
	pb.SuperNodeClient

	id string
}

// RegisterArtowrk implements node.Connection.RegisterArtowrk()
func (conn *clientConn) RegisterArtowrk(ctx context.Context) (node.RegisterArtowrk, error) {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logClientPrefix, conn.id))

	client, err := conn.SuperNodeClient.RegisterArtowrk(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("address", conn.Target()).Errorf("Failed to get client")
		return nil, errors.New(err)
	}

	stream := newRegisterArtowrk(conn, client)
	if stream, ok := stream.(streamService); ok {
		if err := stream.start(ctx); err != nil {
			return nil, err
		}
	}

	return stream, nil
}

func newClientConn(id string, conn *grpc.ClientConn) node.Connection {
	return &clientConn{
		ClientConn:      commongrpc.NewClientConn(conn),
		SuperNodeClient: pb.NewSuperNodeClient(conn),
		id:              id,
	}
}
