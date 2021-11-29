package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/messaging"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
)

type service struct {
	actor                                messaging.Actor
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

// NewService retuns new domain service instance with actor model
func NewService(cfg *Config, secConn node.Client, p2p p2p.Client) (svc StorageChallenge, stopActor func()) {
	if cfg == nil {
		panic("domain service configuration not found")
	}
	localActor := messaging.NewActor(actor.NewActorSystem())
	pid, err := localActor.RegisterActor(newDomainActor(secConn), "domain-actor")
	if err != nil {
		panic(err)
	}
	return &service{
		actor:                                localActor,
		domainActorID:                        pid,
		storageChallengeExpiredAsNanoseconds: cfg.StorageChallengeExpiredDuration.Nanoseconds(),
		repository:                           newRepository(p2p),
		nodeID:                               cfg.PastelID,
		numberOfChallengeReplicas:            cfg.NumberOfChallengeReplicas,
	}, localActor.Stop
}
