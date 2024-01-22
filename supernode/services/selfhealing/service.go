package selfhealing

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/pastelnetwork/gonode/common/types"
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
	"github.com/pastelnetwork/gonode/supernode/services/download"
)

const (
	defaultTimerBlockCheckDuration    = 5 * time.Minute
	defaultFetchNodesPingInfoInterval = 60 * time.Second
	defaultUpdateWatchlistInterval    = 70 * time.Second
	broadcastMetricRegularInterval    = 60 * time.Minute
)

// SHService keeps track of the supernode's nodeID and passes this, the pastel client,
// and node client interfaces to the tasks it controls.  The run method contains a ticker timer
// that will check for a new block and generate self-healing challenges as necessary if a new block
// is detected.
type SHService struct {
	*common.SuperNodeService
	config        *Config
	pastelHandler *mixins.PastelHandler

	nodeID                string
	nodeClient            node.ClientInterface
	historyDB             storage.LocalStoreInterface
	downloadService       *download.NftDownloaderService
	ticketsMap            map[string]bool
	SelfHealingMetricsMap map[string]types.SelfHealingMetrics

	currentBlockCount int32
	merkelroot        string
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

// PingNodes : This worker will periodically ping SNs and maintain their ping history
func (service *SHService) PingNodes(ctx context.Context) {
	for {
		select {
		case <-time.After(defaultFetchNodesPingInfoInterval):
			newCtx := context.Background()
			task := service.NewSHTask()
			task.FetchAndMaintainPingInfo(newCtx)

		case <-ctx.Done():
			log.Println("Context done being called in fetch & maintain ping info worker")
			return
		}
	}
}

// BroadcastSelfHealingMetricsWorker broadcast the self-healing metrics to the entire network
func (service *SHService) BroadcastSelfHealingMetricsWorker(ctx context.Context) {
	log.WithContext(ctx).Info("self-healing-metric worker func has been invoked")

	startTime := calculateStartTime(service.nodeID)
	log.WithContext(ctx).WithField("start_time", startTime).
		Info("self-healing-metric service will execute on the mentioned time")

	// Wait until the start time
	time.Sleep(2 * time.Minute)

	// Run the first task immediately
	service.executeMetricsBroadcastTask(ctx)

	for {
		select {
		case <-time.After(broadcastMetricRegularInterval):
			service.executeMetricsBroadcastTask(context.Background())
		case <-ctx.Done():
			log.Println("Context done being called in file-healing worker")
			return
		}
	}
}

// executeTask executes the self-healing metric task.
func (service *SHService) executeMetricsBroadcastTask(ctx context.Context) {
	newCtx := context.Background()
	task := service.NewSHTask()
	task.BroadcastSelfHealingMetrics(newCtx)

	log.WithContext(ctx).Debug("self-healing metric broadcasted")
}

// RunUpdateWatchlistWorker : This worker will periodically fetch and maintain the ping info and update watchlist field
func (service *SHService) RunUpdateWatchlistWorker(ctx context.Context) {
	for {
		select {
		case <-time.After(defaultUpdateWatchlistInterval):
			newCtx := context.Background()
			task := service.NewSHTask()
			task.UpdateWatchlist(newCtx)

		case <-ctx.Done():
			log.Println("Context done being called in update watchlist worker")
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

	go service.PingNodes(ctx)

	go service.RunUpdateWatchlistWorker(ctx)

	go service.BroadcastSelfHealingMetricsWorker(ctx)

	for {
		select {
		case <-time.After(defaultTimerBlockCheckDuration):
			newCtx := context.Background()
			task := service.NewSHTask()
			task.GenerateSelfHealingChallenge(newCtx)

			log.WithContext(ctx).Debug("Would normally invoke a self-healing worker")
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
	historyDB storage.LocalStoreInterface, downloadService *download.NftDownloaderService) *SHService {
	return &SHService{
		config:                config,
		SuperNodeService:      common.NewSuperNodeService(fileStorage, pastelClient, p2p),
		nodeClient:            nodeClient,
		nodeID:                config.PastelID,
		pastelHandler:         mixins.NewPastelHandler(pastelClient),
		historyDB:             historyDB,
		SelfHealingMetricsMap: make(map[string]types.SelfHealingMetrics),
		ticketsMap:            make(map[string]bool),
		downloadService:       downloadService,
	}
}

// GetNClosestSupernodesToAGivenFileUsingKademlia : Wrapper for a utility function that accesses kademlia's distributed hash table to determine which nodes should be closest to a given string (hence hosting it)
func (service *SHService) GetNClosestSupernodesToAGivenFileUsingKademlia(ctx context.Context, n int, comparisonString string, ignores, includingNodes []string) []string {
	return service.P2PClient.NClosestNodesWithIncludingNodeList(ctx, n, comparisonString, ignores, includingNodes)
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

// calculateStartTime calculates the start time for the periodic task based on the hash of the PastelID.
func calculateStartTime(pastelID string) time.Time {
	hash := sha256.Sum256([]byte(pastelID))
	hashString := hex.EncodeToString(hash[:])
	delayMinutes := int(hashString[0]) % 60 // simplistic hash-based delay calculation

	return time.Now().Add(time.Duration(delayMinutes) * time.Minute)
}
