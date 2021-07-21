package client

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/metadb/network/supernode/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client conneciton.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

// ProcessUserdata implements node.Connection.ProcessUserdata()
func (conn *clientConn) ProcessUserdata() node.ProcessUserdata {
	return newRegisterArtwork(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.Connection {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
