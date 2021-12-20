package client

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client connection.
type clientConn struct {
	*commongrpc.ClientConn

	actorPID *actor.PID
	id       string
}

// RegisterArtwork implements node.Connection.RegisterArtwork()
func (conn *clientConn) RegisterArtwork(withActor ...bool) node.RegisterArtwork {
	return newRegisterArtwork(conn, withActor...)
}

// ProcessUserdata implements node.Connection.ProcessUserdata()
func (conn *clientConn) ProcessUserdata(withActor ...bool) node.ProcessUserdata {
	return newProcessUserdata(conn, withActor...)
}

// StorageChallenge implements node.Connection.StorageChallenge()
func (conn *clientConn) StorageChallenge(withActor ...bool) node.StorageChallenge {
	return newStorageChallenge(conn, withActor...)
}

func newClientConn(id string, conn *grpc.ClientConn) node.Connection {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
