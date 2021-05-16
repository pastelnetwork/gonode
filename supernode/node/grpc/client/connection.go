package client

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client conneciton.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

// RegisterArtowrk implements node.Connection.RegisterArtowrk()
func (conn *clientConn) RegisterArtowrk() node.RegisterArtowrk {
	return newRegisterArtowrk(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.Connection {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
