package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/messaging"
	"github.com/pastelnetwork/gonode/pastel"
)

type service struct {
	remoter                               *messaging.Remoter
	domainActorID                         *actor.PID
	nodeID                                string
	pclient                               pastel.Client
	storageChallengeExpiredAsMilliSeconds int64
	numberOfChallengeReplicas             int
}

type StorageChallenge interface {
	GenerateStorageChallenges(ctx context.Context, blockHash, challengingMasternodeID string, challengesPerNodePerBlock int) error
	ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
	VerifyStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
}
