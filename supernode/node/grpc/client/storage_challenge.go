package client

import (
	"github.com/pastelnetwork/gonode/messaging"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
)

func newStorageChallenge(conn *clientConn, withActor ...bool) node.StorageChallenge {
	client := pb.NewStorageChallengeClient(conn)
	if len(withActor) > 0 && withActor[0] {
		client = newStorageChallengeClientWrapper(client, messaging.GetActorSystem(), conn.actorPID)
	}
	return client
}
