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

// RegisterArtwork implements node.ConnectionInterface.RegisterArtwork()
func (conn *clientConn) RegisterArtwork() node.RegisterNftInterface {
	return newRegisterArtwork(conn)
}

// DownloadArtwork implements node.ConnectionInterface.DownloadArtwork()
func (conn *clientConn) DownloadArtwork() node.DownloadNftInterface {
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
