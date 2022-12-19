package selfhealing

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

// SHService keeps track of the supernode's nodeID and passes this, the pastel client,
// and node client interfaces to the tasks it controls.  The run method contains a ticker timer
// that will check for a new block and generate storage challenges as necessary if a new block
// is detected.
type SHService struct {
	*common.SuperNodeService
	config *Config

	nodeID            string
	nodeClient        node.ClientInterface
	currentBlockCount int32
}

// CheckNextBlockAvailable calls pasteld and checks if a new block is available
func (service *SHService) CheckNextBlockAvailable(ctx context.Context) bool {
	blockCount, err := service.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("method", "checkNextBlockAvailable.GetBlockCount").Warn("could not get block count")
		return false
	}
	if blockCount > service.currentBlockCount {
		atomic.StoreInt32(&service.currentBlockCount, blockCount)
		return true
	}

	return false
}

const defaultTimerBlockCheckDuration = 10 * time.Second

// Run : self-healing service will run continuously to generate self-healings.
func (service *SHService) Run(ctx context.Context) error {
	log.WithContext(ctx).Info("self-healing service run func has been invoked")
	//does this need to be in its own goroutine?
	go func() {
		if err := service.RunHelper(ctx, service.config.PastelID, logPrefix); err != nil {
			log.WithContext(ctx).WithError(err).Error("SelfHealingChallengeService:RunHelper")
		}
	}()

	for {
		select {
		case <-time.After(defaultTimerBlockCheckDuration):
			log.WithContext(ctx).Info("Ticker has ticked")

			if service.CheckNextBlockAvailable(ctx) {
				//Generate Self-Healing Challenge for symbol files

				log.WithContext(ctx).Info("Would normally generate a self-healing challenge")
			} else {
				log.WithContext(ctx).Info("Block not available")
			}
		case <-ctx.Done():
			log.Println("Context done being called in generate self-healing challenge")
			return nil
		}
	}
}

// NewSCTask : self-healing task handles the duties of generating, processing, and verifying self-healing challenges
func (service *SHService) NewSCTask() *SHTask {
	task := NewSHTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the self-healing by the id
func (service *SHService) Task(id string) *SHTask {
	scTask, ok := service.Worker.Task(id).(*SHTask)
	if !ok {
		log.Error("Error typecasting task to self-healing task")
		return nil
	}

	log.Info("type casted successfully")
	return scTask
}

// NewService : Create a new self-healing service
//
//	Inheriting from SuperNodeService allows us to use common methods for pastel-client, p2p, and rqClient.
func NewService(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2p p2p.Client, rqClient rqnode.ClientInterface) *SHService {
	return &SHService{
		config:           config,
		SuperNodeService: common.NewSuperNodeService(fileStorage, pastelClient, p2p, rqClient),
		nodeClient:       nodeClient,
		// repository:                    newRepository(p2p, pastelClient, challengeStatusObserver),
		nodeID: config.PastelID,
	}
}

// ListSymbolFileKeysFromNFTTicket : Get an NFT Ticket's associated raptor q ticket file id's.  These can then be accessed through p2p.
func (service *SHService) ListSymbolFileKeysFromNFTTicket(ctx context.Context) ([]string, error) {
	var keys = make([]string, 0)
	regTickets, err := service.SuperNodeService.PastelClient.RegTickets(ctx)
	if err != nil {
		return keys, err
	}
	if len(regTickets) == 0 {
		log.WithContext(ctx).WithField("count", len(regTickets)).Info("no reg tickets retrieved")
		return keys, nil
	}

	log.WithContext(ctx).WithField("count", len(regTickets)).Info("Reg tickets retrieved")
	for _, regTicket := range regTickets {

		decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}

		regTicket.RegTicketData.NFTTicketData = *decTicket
		keys = append(keys, regTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs...)
	}

	return keys, nil
}

// GetSymbolFileByKey : Wrapper for p2p file storage service - retrieves a file from kademlia based on its key. Here, they should be raptorq symbol files.
func (service *SHService) GetSymbolFileByKey(ctx context.Context, key string, getFromLocalOnly bool) ([]byte, error) {
	return service.P2PClient.Retrieve(ctx, key, getFromLocalOnly)
}

// StoreSymbolFile : Wrapper for p2p file storage service - stores a file in kademlia based on its key
func (service *SHService) StoreSymbolFile(ctx context.Context, data []byte) (key string, err error) {
	return service.P2PClient.Store(ctx, data)
}

// RemoveSymbolFileByKey : Wrapper for p2p file storage service - removes a file from kademlia based on its key
func (service *SHService) RemoveSymbolFileByKey(ctx context.Context, key string) error {
	return service.P2PClient.Delete(ctx, key)
}

// GetListOfSupernode : Access the supernode service to get a list of all supernodes, including their id's and addresses.
// This is used to enumerate supernodes both for calculation and connection
func (service *SHService) GetListOfSupernode(ctx context.Context) ([]string, error) {
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

// FilterOutSupernodes : FilterOutSupernodes gets the full list of supernodes and removes the nodesToBeIgnored
func (service *SHService) FilterOutSupernodes(listOfSupernodes []string, nodesToBeIgnored []string) []string {
	mapOfNodesToBeIgnored := make(map[string]bool)
	for _, node := range nodesToBeIgnored {
		mapOfNodesToBeIgnored[node] = true
	}

	var sliceOfNodesWithoutIgnoredNodes []string
	for _, node := range listOfSupernodes {
		if !mapOfNodesToBeIgnored[node] {
			sliceOfNodesWithoutIgnoredNodes = append(sliceOfNodesWithoutIgnoredNodes, node)
		}
	}

	return sliceOfNodesWithoutIgnoredNodes
}

// GetNClosestSupernodeIDsToComparisonString : Wrapper for a utility function that does xor string comparison to a list of strings and returns the smallest distance.
func (service *SHService) GetNClosestSupernodeIDsToComparisonString(_ context.Context, n int, comparisonString string, listSupernodes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listSupernodes, ignores...)
}

// GetNClosestFileHashesToAGivenComparisonString : Wrapper for a utility function that does xor string comparison to a list of strings and returns the smallest distance.
func (service *SHService) GetNClosestFileHashesToAGivenComparisonString(_ context.Context, n int, comparisonString string, listFileHashes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listFileHashes, ignores...)
}
