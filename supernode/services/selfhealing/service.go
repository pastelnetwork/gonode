package selfhealing

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/mixins"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	defaultTimerBlockCheckDuration    = 10 * time.Second
	defaultFetchNodesPingInfoInterval = 60 * time.Second
)

// SHService keeps track of the supernode's nodeID and passes this, the pastel client,
// and node client interfaces to the tasks it controls.  The run method contains a ticker timer
// that will check for a new block and generate self-healing challenges as necessary if a new block
// is detected.
type SHService struct {
	*common.SuperNodeService
	config        *Config
	pastelHandler *mixins.PastelHandler

	nodeID            string
	nodeClient        node.ClientInterface
	currentBlockCount int32
	historyDB         storage.LocalStoreInterface
}

// CheckNextBlockAvailable calls pasteld and checks if a new block is available
func (service *SHService) CheckNextBlockAvailable(ctx context.Context) bool {
	blockCount, err := service.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		return false
	}
	if blockCount > service.currentBlockCount {
		atomic.StoreInt32(&service.currentBlockCount, blockCount)
		return true
	}

	return false
}

// RunFetchNodesPingInfoWorker : This worker will periodically fetch and maintain the ping info about other SNs in the network
func (service *SHService) RunFetchNodesPingInfoWorker(ctx context.Context) {
	for {
		select {
		case <-time.After(defaultFetchNodesPingInfoInterval):
			newCtx := context.Background()
			task := service.NewSHTask()
			task.FetchAndMaintainPingInfo(newCtx)

		case <-ctx.Done():
			log.Println("Context done being called in local keys fetch worker in service.go")
			return
		}
	}
}

// Run : self-healing service will run continuously to generate self-healing.
func (service *SHService) Run(ctx context.Context) error {
	log.WithContext(ctx).Info("self-healing service run func has been invoked")
	//does this need to be in its own goroutine?
	go func() {
		if err := service.RunHelper(ctx, service.config.PastelID, logPrefix); err != nil {
			log.WithContext(ctx).WithError(err).Error("SelfHealingChallengeService:RunHelper")
		}
	}()

	go service.RunFetchNodesPingInfoWorker(ctx)

	for {
		select {
		case <-time.After(defaultTimerBlockCheckDuration):

			if service.CheckNextBlockAvailable(ctx) {
				//newCtx := context.Background()
				//task := service.NewSCTask()
				//task.ExecuteFileHealingWorker(newCtx)

				log.WithContext(ctx).Debug("Would normally invoke a self-healing worker")
			}
		case <-ctx.Done():
			log.Println("Context done being called in file-healing worker")
			return nil
		}
	}
}

// NewSHTask : self-healing task handles the duties of generating, processing, and verifying self-healing challenges
func (service *SHService) NewSHTask() *SHTask {
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
func NewService(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2p p2p.Client,
	historyDB storage.LocalStoreInterface) *SHService {
	return &SHService{
		config:           config,
		SuperNodeService: common.NewSuperNodeService(fileStorage, pastelClient, p2p),
		nodeClient:       nodeClient,
		nodeID:           config.PastelID,
		pastelHandler:    mixins.NewPastelHandler(pastelClient),
		historyDB:        historyDB,
	}
}

// MapSymbolFileKeysFromNFTAndActionTickets : Get an NFT and Action Ticket's associated raptor q ticket file id's and creates a map of them
func (service *SHService) MapSymbolFileKeysFromNFTAndActionTickets(ctx context.Context) (regTicketKeys map[string]string, actionTicketKeys map[string]string, err error) {
	regTicketKeys = make(map[string]string)
	actionTicketKeys = make(map[string]string)
	regTickets, err := service.SuperNodeService.PastelClient.RegTickets(ctx)
	if err != nil {
		return regTicketKeys, actionTicketKeys, err
	}

	if len(regTickets) == 0 {
		log.WithContext(ctx).WithField("count", len(regTickets)).Info("no reg tickets retrieved")
		return regTicketKeys, actionTicketKeys, nil
	}

	log.WithContext(ctx).WithField("count", len(regTickets)).Info("Reg tickets retrieved")
	for _, regTicket := range regTickets {

		decTicket, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}

		regTicket.RegTicketData.NFTTicketData = *decTicket

		for _, rqID := range regTicket.RegTicketData.NFTTicketData.AppTicketData.RQIDs {
			regTicketKeys[rqID] = regTicket.TXID
		}
	}

	actionTickets, err := service.SuperNodeService.PastelClient.ActionTickets(ctx)
	if err != nil {
		return regTicketKeys, actionTicketKeys, err
	}

	if len(actionTickets) == 0 {
		log.WithContext(ctx).WithField("count", len(regTickets)).Info("no action reg tickets retrieved")
		return regTicketKeys, actionTicketKeys, nil
	}
	log.WithContext(ctx).WithField("count", len(regTickets)).Info("Action tickets retrieved")

	for _, actionTicket := range actionTickets {
		decTicket, err := pastel.DecodeActionTicket(actionTicket.ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		actionTicket.ActionTicketData.ActionTicketData = *decTicket

		cascadeTicket, err := actionTicket.ActionTicketData.ActionTicketData.APICascadeTicket()
		if err != nil {
			log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTicket.TXID).
				Warnf("Could not get cascade ticket for action ticket data self healing")
			continue
		}

		for _, rqID := range cascadeTicket.RQIDs {
			actionTicketKeys[rqID] = actionTicket.TXID
		}
	}

	return regTicketKeys, actionTicketKeys, nil
}

// GetNClosestSupernodeIDsToComparisonString : Wrapper for a utility function that does xor string comparison to a list of strings and returns the smallest distance.
func (service *SHService) GetNClosestSupernodeIDsToComparisonString(_ context.Context, n int, comparisonString string, listSupernodes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listSupernodes, ignores...)
}
