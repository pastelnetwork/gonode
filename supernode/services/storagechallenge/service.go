package storagechallenge

import (
	baseCtx "context"
	"sync/atomic"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/messaging"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
)

type service struct {
	actor                         messaging.Actor
	domainActorID                 *actor.PID
	nodeID                        string
	pclient                       pastel.Client
	storageChallengeExpiredBlocks int32
	numberOfChallengeReplicas     int
	repository                    repository
	currentBlockCount             int32
}

// StorageChallenge interface
type StorageChallenge interface {
	// GenerateStorageChallenges func
	GenerateStorageChallenges(ctx context.Context, challengesPerNodePerBlock int) error
	// ProcessStorageChallenge func
	ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
	// VerifyStorageChallenge func
	VerifyStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
	// Run service
	Run(ctx baseCtx.Context) error
}

// NewService retuns new domain service instance with actor model
func NewService(cfg *Config, secConn node.Client, p2p p2p.Client, pClient pastel.Client, challengeStatusObserver SaveChallengeState) (svc StorageChallenge, stopActor func()) {
	if cfg == nil {
		panic("domain service configuration not found")
	}
	localActor := messaging.NewActor(actor.NewActorSystem())
	pid, err := localActor.RegisterActor(newDomainActor(secConn), "domain-actor")
	if err != nil {
		panic(err)
	}
	return &service{
		actor:                         localActor,
		domainActorID:                 pid,
		storageChallengeExpiredBlocks: cfg.StorageChallengeExpiredBlocks,
		pclient:                       pClient,
		repository:                    newRepository(p2p, pClient, challengeStatusObserver),
		nodeID:                        cfg.PastelID,
		numberOfChallengeReplicas:     cfg.NumberOfChallengeReplicas,
	}, localActor.Stop
}

func (s *service) checkNextBlockAvailable(ctx baseCtx.Context) bool {
	blockCount, err := s.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithField("method", "checkNextBlockAvailable.GetBlockCount").Warn("could not get block count")
		return false
	}
	if blockCount > int32(s.currentBlockCount) {
		atomic.StoreInt32(&s.currentBlockCount, blockCount)
		return true
	}

	return false
}

const defaultTimerBlockCheckDuration = 30 * time.Second

func (s *service) Run(ctx baseCtx.Context) error {
	ticker := time.NewTicker(defaultTimerBlockCheckDuration)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-ticker.C:
			if ok && s.checkNextBlockAvailable(ctx) {
				newCtx := context.Background()
				err := s.GenerateStorageChallenges(newCtx, 1)
				if err != nil {
					log.WithField("method", "StartGenerateStorageChallengeCheckLoop.GenerateStorageChallenges").WithError(err).Warn("could not generate storage challenge")
					continue
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}
