package storagechallenge

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/utils"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"

	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

type StorageChallengeService struct {
	*common.SuperNodeService
	config *Config

	nodeID                        string
	pclient                       pastel.Client
	nodeClient                    node.ClientInterface
	storageChallengeExpiredBlocks int32
	numberOfChallengeReplicas     int
	// repository                    Repository
	currentBlockCount int32
}

func (s *StorageChallengeService) checkNextBlockAvailable(ctx context.Context) bool {
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

const defaultTimerBlockCheckDuration = 10 * time.Second

// Storage challenges ensure that other supernodes are properly storing the files they're supposed to.
// Storage challenge service will run continuously to generate storage challenges.
func (s *StorageChallengeService) Run(ctx context.Context) error {
	ticker := time.NewTicker(defaultTimerBlockCheckDuration)
	//does this need to be in its own goroutine?
	go s.RunHelper(ctx, s.config.PastelID, logPrefix)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-ticker.C:
			log.Println("Ticker has ticked")

			if ok && s.checkNextBlockAvailable(ctx) {
				newCtx := context.Background()
				task := s.NewStorageChallengeTask()
				task.GenerateStorageChallenges(newCtx)
			} else {
				if !ok {
					log.WithContext(ctx).Println("Ticker not okay")
				} else {
					log.WithContext(ctx).Println("Block not available")
				}
			}
		case <-ctx.Done():
			log.Println("Context done being called in generatestoragechallenge loop in service.go")
			return nil
		}
	}
}

// Storage challenge task handles the duties of responding to others' storage challenges.
func (service *StorageChallengeService) NewStorageChallengeTask() *StorageChallengeTask {
	task := NewStorageChallengeTask(service)
	service.Worker.AddTask(task)
	return task
}

func NewService(cfg *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2p p2p.Client, rqClient rqnode.ClientInterface, challengeStatusObserver SaveChallengeState) *StorageChallengeService {
	if cfg == nil {
		panic("domain service configuration not found")
	}

	return &StorageChallengeService{
		config:                        cfg,
		SuperNodeService:              common.NewSuperNodeService(fileStorage, pastelClient, p2p, rqClient),
		nodeClient:                    nodeClient,
		storageChallengeExpiredBlocks: cfg.StorageChallengeExpiredBlocks,
		pclient:                       pastelClient,
		// repository:                    newRepository(p2p, pastelClient, challengeStatusObserver),
		nodeID:                    cfg.PastelID,
		numberOfChallengeReplicas: cfg.NumberOfChallengeReplicas,
	}
}

//utils below

func (service *StorageChallengeService) ListSymbolFileKeysFromNFTTicket(ctx context.Context) ([]string, error) {
	var keys = make([]string, 0)
	regTickets, err := service.pclient.RegTickets(ctx)
	if err != nil {
		return keys, err
	}
	for _, regTicket := range regTickets {
		for _, key := range regTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs {
			keys = append(keys, string(key))
		}
	}

	return keys, nil
}

func (service *StorageChallengeService) GetSymbolFileByKey(ctx context.Context, key string, getFromLocalOnly bool) ([]byte, error) {
	return service.P2PClient.Retrieve(ctx, key)
}

func (service *StorageChallengeService) StoreSymbolFile(ctx context.Context, data []byte) (key string, err error) {
	return service.P2PClient.Store(ctx, data)
}

func (service *StorageChallengeService) RemoveSymbolFileByKey(ctx context.Context, key string) error {
	return service.P2PClient.Delete(ctx, key)
}

func (service *StorageChallengeService) GetListOfSupernode(ctx context.Context) ([]string, error) {
	var ret = make([]string, 0)
	listMN, err := service.pclient.MasterNodesExtra(ctx)
	if err != nil {
		return ret, err
	}

	for _, node := range listMN {
		ret = append(ret, node.ExtKey)
	}

	return ret, nil
}

func (service *StorageChallengeService) GetNClosestSupernodeIDsToComparisonString(_ context.Context, n int, comparisonString string, listSupernodes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listSupernodes, ignores...)
}

func (service *StorageChallengeService) GetNClosestSupernodesToAGivenFileUsingKademlia(ctx context.Context, n int, comparisonString string, ignores ...string) []string {
	return service.P2PClient.NClosestNodes(ctx, n, comparisonString, ignores...)
}

func (service *StorageChallengeService) GetNClosestFileHashesToAGivenComparisonString(_ context.Context, n int, comparisonString string, listFileHashes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listFileHashes, ignores...)
}
