package client

import (
	commongrpc "github.com/pastelnetwork/gonode/common/net/grpc"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc"
)

// clientConn represents grpc client connection.
type clientConn struct {
	*commongrpc.ClientConn

	id string
}

// RegisterNft implements node.ConnectionInterface.RegisterNft()
func (conn *clientConn) RegisterNft() node.RegisterNftInterface {
	return newRegisterNft(conn)
}

// RegisterSense implements node.ConnectionInterface.RegisterSense()
func (conn *clientConn) RegisterSense() node.RegisterSenseInterface {
	return newRegisterSense(conn)
}

// RegisterCascade implements node.ConnectionInterface.RegisterSense()
func (conn *clientConn) RegisterCascade() node.RegisterCascadeInterface {
	return newRegisterCascade(conn)
}

// StorageChallenge implements node.Connection.StorageChallenge()
func (conn *clientConn) StorageChallenge() node.StorageChallengeInterface {
	return newStorageChallengeGRPCClient(conn)
}

// SelfHealingChallenge implements node.Connection.SelfHealingChallenge()
func (conn *clientConn) SelfHealingChallenge() node.SelfHealingChallengeInterface {
	return newSelfHealingGRPCClient(conn)
}

func newClientConn(id string, conn *grpc.ClientConn) node.ConnectionInterface {
	return &clientConn{
		ClientConn: commongrpc.NewClientConn(conn),
		id:         id,
	}
}
