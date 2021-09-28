package ddclient

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"google.golang.org/grpc"
)

// clientConn represents grpc client conneciton.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

func newClientConn(id string, conn *grpc.ClientConn) *clientConn {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
