package client

import (
	"github.com/pastelnetwork/gonode/proto/supernode/storagechallenge"
	"github.com/pastelnetwork/gonode/supernode/node"
)

func newStorageChallenge(conn *clientConn) node.StorageChallenge {
	return storagechallenge.NewStorageChallengeClient(conn)
}
