package grpc

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client connection.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

// HermesP2P implements node.ConnectionInterface.HermesP2P()
func (conn *clientConn) HermesP2P() node.HermesP2PInterface {
	return newHermesP2P(conn)
}

// RegisterNft implements node.ConnectionInterface.RegisterNft()
func (conn *clientConn) RegisterNft() node.RegisterNftInterface {
	return newRegisterNft(conn)
}

// DownloadNft implements node.ConnectionInterface.DownloadNft()
func (conn *clientConn) DownloadNft() node.DownloadNftInterface {
	return newDownloadNft(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.ConnectionInterface {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
