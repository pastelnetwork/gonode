package client

import (
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
)

func newStorageChallenge(conn *clientConn) node.StorageChallenge {
	return pb.NewStorageChallengeClient(conn)
}
