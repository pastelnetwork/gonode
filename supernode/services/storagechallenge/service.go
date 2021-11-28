package storagechallenge

import (
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/messaging"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
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

func NewService(remoter *messaging.Remoter, secConn node.Client, p2p p2p.Client) StorageChallenge {
	if remoter == nil {
		remoter = messaging.NewRemoter(actor.NewActorSystem(), nil)
	}
	pid, err := remoter.RegisterActor(newDomainActor(secConn), "domain-actor")
	if err != nil {
		log.Panic(err)
	}
	return &service{
		domainActorID:                        pid,
		storageChallengeExpiredAsNanoseconds: int64(time.Minute.Nanoseconds()),
		repository:                           newRepository(p2p),
	}
}
