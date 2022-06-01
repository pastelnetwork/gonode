package grpc

import (
	"github.com/pastelnetwork/gonode/bridge/node"
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"google.golang.org/grpc"
)

// clientConn represents grpc client connection.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

// DownloadData implements node.ConnectionInterface.DownloadData()
func (conn *clientConn) DownloadData() node.DownloadDataInterface {
	return newDownloadData(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.ConnectionInterface {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
