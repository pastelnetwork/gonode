package storagechallenge

import (
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
	GenerateStorageChallenges(ctx context.Context, blockHash, challengingMasternodeID string, challengesPerNodePerBlock int) error
	// ProcessStorageChallenge func
	ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
	// VerifyStorageChallenge func
	VerifyStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error
}

// NewService retuns new domain service instance with actor model
func NewService(cfg *Config, secConn node.Client, p2p p2p.Client, pClient pastel.Client) (svc StorageChallenge, stopActor func()) {
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
		repository:                    newRepository(p2p, pClient),
		nodeID:                        cfg.PastelID,
		numberOfChallengeReplicas:     cfg.NumberOfChallengeReplicas,
	}, localActor.Stop
}

func (s *service) checkNextBlockAvailable(ctx context.Context) bool {
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

func (s *service) getMerkleRoot(ctx context.Context, blockHeight int32) (string, error) {
	blockVerbose, err := s.pclient.GetBlockVerbose1(ctx, blockHeight)
	if err != nil {
		return "", err
	}

	return blockVerbose.MerkleRoot, nil
}

const defaultTimerBlockCheckDuration = 10 * time.Second

func (s *service) StartGenerateStorageChallengeCheckLoop(ctx context.Context) (cncl func()) {
	ticker := time.NewTicker(defaultTimerBlockCheckDuration)
	defer ticker.Stop()
	var done = make(chan bool)
	cncl = func() {
		close(done)
	}
	for {
		select {
		case _, ok := <-ticker.C:
			if ok && s.checkNextBlockAvailable(ctx) {
				merkleroot, err := s.getMerkleRoot(ctx, s.currentBlockCount)
				if err != nil {
					log.WithField("method", "StartGenerateStorageChallengeCheckLoop.getMerkleRoot.GetBlockVerbose1").WithError(err).Warn("could not get block verbose")
					continue
				}
				nodes := s.repository.GetNClosestXORDistanceMasternodesToComparisionString(ctx, 1, merkleroot)
				if len(nodes) != 1 {
					log.WithField("method", "StartGenerateStorageChallengeCheckLoop.repository.GetNClosestXORDistanceMasternodesToComparisionString").Warn("got empty closets masternodes")
					continue
				}
				newCtx := context.Background()
				err = s.GenerateStorageChallenges(newCtx, merkleroot, string(nodes[0].ID), 1)
				if len(nodes) != 1 {
					log.WithField("method", "StartGenerateStorageChallengeCheckLoop.GenerateStorageChallenges").WithError(err).Warn("could not generate storage challenge")
					continue
				}
			}
		case <-done:
			return
		}
	}
}
