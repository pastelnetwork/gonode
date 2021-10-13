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

// RegisterArtwork implements node.Connection.RegisterArtwork()
func (conn *clientConn) RegisterArtwork() node.RegisterArtwork {
	return newRegisterArtwork(conn)
}

// DownloadArtwork implements node.Connection.DownloadArtwork()
func (conn *clientConn) DownloadArtwork() node.DownloadArtwork {
	return newDownloadArtwork(conn)
}

// ProcessUserdata implements node.Connection.ProcessUserdata()
func (conn *clientConn) ProcessUserdata() node.ProcessUserdata {
	return newProcessUserdata(conn)
}

func (conn *clientConn) ExternalDupeDetection() node.ExternalDupeDetection {
	return newExternalDupeDetection(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.Connection {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
