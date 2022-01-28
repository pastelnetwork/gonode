package grpc

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client connection.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

// RegisterNft implements node.ConnectionInterface.RegisterNft()
func (conn *clientConn) RegisterNft() node.RegisterNftInterface {
	return newRegisterArtwork(conn)
}

// DownloadNft implements node.ConnectionInterface.DownloadNft()
func (conn *clientConn) DownloadNft() node.DownloadNftInterface {
	return newDownloadArtwork(conn)
}

// ProcessUserdata implements node.ConnectionInterface.ProcessUserdata()
func (conn *clientConn) ProcessUserdata() node.ProcessUserdataInterface {
	return newProcessUserdata(conn)
}

// RegisterSense implements node.ConnectionInterface.RegisterSense()
func (conn *clientConn) RegisterSense() node.RegisterSenseInterface {
	return newRegisterSense(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.ConnectionInterface {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
