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

// The goal of the storage challenge system is to ensure that supernodes are hosting
// the files they should be.  In order to do this, storage challenges are sent out
// with each new block.  For each storage challenge, multiple challengers are deterministically
// selected, who then deterministically select: which supernodes to check, which files to check,
// and which other supernodes should verify those checks.

// The current state of the storage challenge system:
// Storage challenge is functioning, but if a storage challenge verification fails, nothing happens.
// Callbacks in task/stateStorage can be adjusted when proper consequences are determined.

// SCService keeps track of the supernode's nodeID and passes this, the pastel client,
// and node client interfaces to the tasks it controls.  The run method contains a ticker timer
// that will check for a new block and generate storage challenges as necessary if a new block
// is detected.
type SCService struct {
	*common.SuperNodeService
	config *Config

	nodeID                        string
	nodeClient                    node.ClientInterface
	storageChallengeExpiredBlocks int32
	numberOfChallengeReplicas     int
	numberOfVerifyingNodes        int
	// repository                    Repository
	currentBlockCount int32
	// currently unimplemented, default always used instead.
	challengeStatusObserver SaveChallengeState
}

//CheckNextBlockAvailable calls pasteld and checks if a new block is available
func (service *SCService) CheckNextBlockAvailable(ctx context.Context) bool {
	blockCount, err := service.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithField("method", "checkNextBlockAvailable.GetBlockCount").Warn("could not get block count")
		return false
	}
	if blockCount > int32(service.currentBlockCount) {
		atomic.StoreInt32(&service.currentBlockCount, blockCount)
		return true
	}

	return false
}

const defaultTimerBlockCheckDuration = 10 * time.Second

// Run : storage challenge service will run continuously to generate storage challenges.
func (service *SCService) Run(ctx context.Context) error {
	ticker := time.NewTicker(defaultTimerBlockCheckDuration)
	//does this need to be in its own goroutine?
	go service.RunHelper(ctx, service.config.PastelID, logPrefix)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//log.Println("Ticker has ticked")

			// if service.CheckNextBlockAvailable(ctx) {
			// newCtx := context.Background()
			// task := service.NewSCTask()
			//task.GenerateStorageChallenges(newCtx)
			//log.WithContext(ctx).Println("Would normally generate a storage challenge")
			// } else {
			//log.WithContext(ctx).Println("Block not available")
			// }
		case <-ctx.Done():
			log.Println("Context done being called in generatestoragechallenge loop in service.go")
			return nil
		}
	}
}

// NewSCTask : Storage challenge task handles the duties of generating, processing, and verifying storage challenges
func (service *SCService) NewSCTask() *SCTask {
	task := NewSCTask(service)
	service.Worker.AddTask(task)
	return task
}

// NewService : Create a new storage challenge service
//  Inheriting from SuperNodeService allows us to use common methods for pastelclient, p2p, and rqClient.
func NewService(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2p p2p.Client, rqClient rqnode.ClientInterface, challengeStatusObserver SaveChallengeState) *SCService {
	return &SCService{
		config:                        config,
		SuperNodeService:              common.NewSuperNodeService(fileStorage, pastelClient, p2p, rqClient),
		nodeClient:                    nodeClient,
		storageChallengeExpiredBlocks: config.StorageChallengeExpiredBlocks,
		// repository:                    newRepository(p2p, pastelClient, challengeStatusObserver),
		nodeID:                    config.PastelID,
		numberOfChallengeReplicas: config.NumberOfChallengeReplicas,
		numberOfVerifyingNodes:    config.NumberOfVerifyingNodes,
		challengeStatusObserver:   challengeStatusObserver,
	}
}

//utils below that call pasteld or p2p - mostly just wrapping other functions in better names

//ListSymbolFileKeysFromNFTTicket : Get an NFT Ticket's associated raptor q ticket file id's.  These can then be accessed through p2p.
func (service *SCService) ListSymbolFileKeysFromNFTTicket(ctx context.Context) ([]string, error) {
	var keys = make([]string, 0)
	regTickets, err := service.SuperNodeService.PastelClient.RegTickets(ctx)
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

// GetSymbolFileByKey : Wrapper for p2p file storage service - retrieves a file from kademlia based on its key. Here, they should be raptorq symbol files.
func (service *SCService) GetSymbolFileByKey(ctx context.Context, key string, getFromLocalOnly bool) ([]byte, error) {
	return service.P2PClient.Retrieve(ctx, key, getFromLocalOnly)
}

// StoreSymbolFile : Wrapper for p2p file storage service - stores a file in kademlia based on its key
func (service *SCService) StoreSymbolFile(ctx context.Context, data []byte) (key string, err error) {
	return service.P2PClient.Store(ctx, data)
}

//RemoveSymbolFileByKey : Wrapper for p2p file storage service - removes a file from kademlia based on its key
func (service *SCService) RemoveSymbolFileByKey(ctx context.Context, key string) error {
	return service.P2PClient.Delete(ctx, key)
}

//GetListOfSupernode : Access the supernode service to get a list of all supernodes, including their id's and addresses.
// This is used to enumerate supernodes both for calculation and connection
func (service *SCService) GetListOfSupernode(ctx context.Context) ([]string, error) {
	var ret = make([]string, 0)
	listMN, err := service.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		return ret, err
	}

	for _, node := range listMN {
		ret = append(ret, node.ExtKey)
	}

	return ret, nil
}

// GetNClosestSupernodeIDsToComparisonString : Wrapper for a utility function that does xor string comparison to a list of strings and returns the smallest distance.
func (service *SCService) GetNClosestSupernodeIDsToComparisonString(_ context.Context, n int, comparisonString string, listSupernodes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listSupernodes, ignores...)
}

// GetNClosestSupernodesToAGivenFileUsingKademlia : Wrapper for a utility function that accesses kademlia's distributed hash table to determine which nodes should be closest to a given string (hence hosting it)
func (service *SCService) GetNClosestSupernodesToAGivenFileUsingKademlia(ctx context.Context, n int, comparisonString string, ignores ...string) []string {
	return service.P2PClient.NClosestNodes(ctx, n, comparisonString, ignores...)
}

// GetNClosestFileHashesToAGivenComparisonString : Wrapper for a utility function that does xor string comparison to a list of strings and returns the smallest distance.
func (service *SCService) GetNClosestFileHashesToAGivenComparisonString(_ context.Context, n int, comparisonString string, listFileHashes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listFileHashes, ignores...)
}
