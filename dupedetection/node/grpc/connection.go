package grpc

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/dupedetection/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client conneciton.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

func (conn *clientConn) Dupedetection(config *node.DDServerConfig) node.Dupedetection {
	return newDupedetectionImpl(conn, config)
}

func newClientConn(id string, conn *grpc.ClientConn) node.Connection {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
