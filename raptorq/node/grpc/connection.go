package grpc

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/raptorq/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client conneciton.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

func (conn *clientConn) RaptorQ() node.RaptorQ {
	return newRaptorQ(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.Connection {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
