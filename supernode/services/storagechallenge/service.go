package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/messaging"
	"github.com/pastelnetwork/gonode/pastel"
)

type service struct {
	remoter                              *messaging.Remoter
	domainActorID                        *actor.PID
	nodeID                               string
	pclient                              pastel.Client
	storageChallengeExpiredAsNanoseconds int64
	numberOfChallengeReplicas            int
	repository                           repository
}

// StorageChallenge interface
type StorageChallenge interface {
	// GenerateStorageChallenges func
	GenerateStorageChallenges(ctx context.Context, blockHash, challengingMasternodeID string, challengesPerNodePerBlock int) error
	// ProcessStorageChallenge func
	ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
	// VerifyStorageChallenge func
	VerifyStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
}
