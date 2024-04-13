package storagechallenge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"os"
	"sync"
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

// The goal of the storage challenge system is to ensure that supernodes are hosting
// the files they should be.  In order to do this, storage challenges are sent out
// with each new block.  For each storage challenge, multiple challengers are deterministically
// selected, who then deterministically select: which supernodes to check, which files to check,
// and which other supernodes should verify those checks.

// The current state of the storage challenge system:
// Storage challenge is functioning, but if a storage challenge verification fails, nothing happens.
// Callbacks in task/stateStorage can be adjusted when proper consequences are determined.

const (
	broadcastSCMetricRegularInterval           = 30 * time.Minute
	retryThreshold                             = 3
	processStorageChallengeScoreEventsInterval = 12 * time.Minute
	getChallengesForScoreAggregationInterval   = 10 * time.Minute
)

// SCService keeps track of the supernode's nodeID and passes this, the pastel client,
// and node client interfaces to the tasks it controls.  The run method contains a ticker timer
// that will check for a new block and generate storage challenges as necessary if a new block
// is detected.
type SCService struct {
	*common.SuperNodeService
	config *Config

	nodeID                        string
	nodeClient                    node.ClientInterface
	pastelHandler                 *mixins.PastelHandler
	storageChallengeExpiredBlocks int32
	numberOfChallengeReplicas     int
	numberOfVerifyingNodes        int
	historyDB                     queries.LocalStoreInterface

	currentBlockCount int32
	// currently unimplemented, default always used instead.
	challengeStatusObserver SaveChallengeState

	localKeys            sync.Map
	localKeysLastFetchAt time.Time
	eventRetryMap        map[string]int
}

// CheckNextBlockAvailable calls pasteld and checks if a new block is available
func (service *SCService) CheckNextBlockAvailable(ctx context.Context) bool {
	blockCount, err := service.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		return false
	}
	if blockCount > int32(service.currentBlockCount) {
		atomic.StoreInt32(&service.currentBlockCount, blockCount)
		return true
	}

	return false
}

const (
	defaultTimerBlockCheckDuration = 30 * time.Second
	defaultLocalKeysFetchInterval  = 20 * time.Minute
)

// RunLocalKeysFetchWorker : This worker will periodically fetch the queries keys from the queries storage
func (service *SCService) RunLocalKeysFetchWorker(ctx context.Context) {
	service.localKeysLastFetchAt = time.Now().UTC()
	if err := service.populateLocalKeys(ctx, nil); err != nil {
		log.WithContext(ctx).WithError(err).Error("RunLocalKeysFetchWorker:populateLocalKeys")
	}

	for {
		select {
		case <-time.After(defaultLocalKeysFetchInterval):
			if err := service.populateLocalKeys(ctx, &service.localKeysLastFetchAt); err != nil {
				log.WithContext(ctx).WithError(err).Error("RunLocalKeysFetchWorker:populateLocalKeys")
			}
		case <-ctx.Done():
			log.Println("Context done being called in queries keys fetch worker in service.go")
			return
		}
	}
}

func (service *SCService) populateLocalKeys(ctx context.Context, from *time.Time) error {
	fromStr := ""
	if from != nil {
		fromStr = from.String()
	}
	log.WithContext(ctx).WithField("fromStr", fromStr).WithField("to", service.localKeysLastFetchAt).Info("Populating queries keys")
	keys, err := service.P2PClient.GetLocalKeys(ctx, from, service.localKeysLastFetchAt)
	if err != nil {
		return fmt.Errorf("GetLocalKeys: %w", err)
	}

	for i := 0; i < len(keys); i++ {
		service.localKeys.Store(keys[i], true)
	}

	return nil
}

// Run : storage challenge service will run continuously to generate storage challenges.
func (service *SCService) Run(ctx context.Context) error {
	log.WithContext(ctx).Info("Storage challenge service run has been invoked")
	//does this need to be in its own goroutine?
	go func() {
		if err := service.RunHelper(ctx, service.config.PastelID, logPrefix); err != nil {
			log.WithContext(ctx).WithError(err).Error("StorageChallengeService:RunHelper")
		}
	}()

	go service.RunLocalKeysFetchWorker(ctx)

	go service.BroadcastStorageChallengeMetricsWorker(ctx)

	go service.GetChallengesForScoreAggregation(ctx)

	go service.ProcessAggregationChallenges(ctx)

	if !service.config.IsTestConfig {
		time.Sleep(15 * time.Minute)
	}

	for {
		select {
		case <-time.After(defaultTimerBlockCheckDuration):

			if service.CheckNextBlockAvailable(ctx) && os.Getenv("INTEGRATION_TEST_ENV") != "true" {
				newCtx := log.ContextWithPrefix(context.Background(), "storage-challenge")
				task := service.NewSCTask()

				task.GenerateStorageChallenges(newCtx)
				log.WithContext(newCtx).Debug("Would normally generate a storage challenge")
			}
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

// Task returns the task of the Storage Challenge by the id
func (service *SCService) Task(id string) *SCTask {
	scTask, ok := service.Worker.Task(id).(*SCTask)
	if !ok {
		log.Error("Error typecasting task to storage challenge task")
		return nil
	}

	log.Info("type casted successfully")
	return scTask
}

// executeTask executes the self-healing metric task.
func (service *SCService) executeMetricsBroadcastTask(ctx context.Context) {
	newCtx := context.Background()
	task := service.NewSCTask()
	task.BroadcastStorageChallengeMetrics(newCtx)

	log.WithContext(ctx).Debug("storage challenge metrics broadcasted")
}

// executeTask executes the self-healing metric task.
func (service *SCService) executeGetChallengeIDsForScoreWorker(ctx context.Context) {
	newCtx := context.Background()
	task := service.NewSCTask()
	task.GetScoreAggregationChallenges(newCtx)

	log.WithContext(ctx).Debug("storage challenge score events have been fetched")
}

// BroadcastStorageChallengeMetricsWorker broadcast the storage challenge metrics to the entire network
func (service *SCService) BroadcastStorageChallengeMetricsWorker(ctx context.Context) {
	log.WithContext(ctx).Info("storage challenge metrics worker func has been invoked")

	startTime := calculateStartTime(service.nodeID)
	log.WithContext(ctx).WithField("start_time", startTime).
		Info("storage challenge metric worker will execute on the mentioned time")

	// Wait until the start time
	time.Sleep(time.Until(startTime))

	// Run the first task immediately
	service.executeMetricsBroadcastTask(ctx)

	for {
		select {
		case <-time.After(broadcastSCMetricRegularInterval):
			service.executeMetricsBroadcastTask(context.Background())
		case <-ctx.Done():
			log.Println("Context done being called in file-healing worker")
			return
		}
	}
}

// GetChallengesForScoreAggregation get the challenges for score aggregation
func (service *SCService) GetChallengesForScoreAggregation(ctx context.Context) {
	log.WithContext(ctx).Info("storage challenge score-aggregation worker func has been invoked")

	for {
		select {
		case <-time.After(getChallengesForScoreAggregationInterval):
			service.executeGetChallengeIDsForScoreWorker(context.Background())
		case <-ctx.Done():
			log.Println("Context done being called in file-healing worker")
			return
		}
	}
}

// calculateStartTime calculates the start time for the periodic task based on the hash of the PastelID.
func calculateStartTime(pastelID string) time.Time {
	hash := sha256.Sum256([]byte(pastelID))
	hashString := hex.EncodeToString(hash[:])
	delayMinutes := int(hashString[0]) % 60 // simplistic hash-based delay calculation

	return time.Now().Add(time.Duration(delayMinutes) * time.Minute)
}

// ProcessAggregationChallenges process the challenges stored in DB for aggregation
func (service *SCService) ProcessAggregationChallenges(ctx context.Context) {
	log.WithContext(ctx).Debug("process-sc-aggregation-challenges worker func has been invoked")

	for {
		select {
		case <-time.After(processStorageChallengeScoreEventsInterval):
			service.processStorageChallengeScoreEvents(ctx)
		case <-ctx.Done():
			log.Println("Context done being called in file-healing worker")
			return
		}
	}
}

// processStorageChallengeScoreEvents executes the self-healing events.
func (service *SCService) processStorageChallengeScoreEvents(ctx context.Context) {
	newCtx := context.Background()
	task := service.NewSCTask()

	store, err := queries.OpenHistoryDB()
	if err != nil {
		return
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	events, err := store.GetStorageChallengeScoreEvents()
	if err != nil {
		return
	}

	for i := 0; i < len(events); i++ {
		event := events[i]
		service.eventRetryMap[event.ChallengeID]++

		if service.eventRetryMap[event.ChallengeID] >= retryThreshold {
			err = store.UpdateScoreChallengeEvent(event.ChallengeID)
			if err != nil {
				log.WithContext(ctx).
					WithField("score_sc_challenge_event_id", event.ChallengeID).
					WithError(err).Error("error updating score-storage-challenge-event")
				continue

			}
			delete(service.eventRetryMap, event.ChallengeID)

			continue
		}

		err = task.AccumulateStorageChallengeScoreData(newCtx, event.ChallengeID)
		if err != nil {
			log.WithContext(ctx).
				WithField("score_sc_challenge_event_id", event.ChallengeID).
				WithError(err).Error("error processing score-storage-challenge-event")

			continue
		}

		err = store.UpdateScoreChallengeEvent(event.ChallengeID)
		if err != nil {
			continue
		}

		delete(service.eventRetryMap, event.ChallengeID)
	}

	log.WithContext(ctx).Debug("self-healing events have been processed")
}

// NewService : Create a new storage challenge service
//
//	Inheriting from SuperNodeService allows us to use common methods for pastelclient, p2p, and rqClient.
func NewService(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface,
	p2p p2p.Client, challengeStatusObserver SaveChallengeState, historyDB queries.LocalStoreInterface) *SCService {
	return &SCService{
		config:                        config,
		SuperNodeService:              common.NewSuperNodeService(fileStorage, pastelClient, p2p),
		nodeClient:                    nodeClient,
		storageChallengeExpiredBlocks: config.StorageChallengeExpiredBlocks,
		// repository:                    newRepository(p2p, pastelClient, challengeStatusObserver),
		nodeID:                    config.PastelID,
		numberOfChallengeReplicas: config.NumberOfChallengeReplicas,
		numberOfVerifyingNodes:    config.NumberOfVerifyingNodes,
		challengeStatusObserver:   challengeStatusObserver,
		localKeys:                 sync.Map{},
		historyDB:                 historyDB,
		eventRetryMap:             make(map[string]int),
	}
}

//utils below that call pasteld or p2p - mostly just wrapping other functions in better names

// ListSymbolFileKeysFromNFTAndActionTickets : Get an NFT and Action Ticket's associated raptor q ticket file id's.
// These can then be accessed through p2p.
func (service *SCService) ListSymbolFileKeysFromNFTAndActionTickets(ctx context.Context) ([]string, error) {
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

	for i := 0; i < len(regTickets); i++ {
		decTicket, err := pastel.DecodeNFTTicket(regTickets[i].RegTicketData.NFTTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}

		regTickets[i].RegTicketData.NFTTicketData = *decTicket
		keys = append(keys, regTickets[i].RegTicketData.NFTTicketData.AppTicketData.RQIDs...)
	}

	actionTickets, err := service.SuperNodeService.PastelClient.ActionTickets(ctx)
	if err != nil {
		return keys, err
	}
	if len(actionTickets) == 0 {
		log.WithContext(ctx).WithField("count", len(actionTickets)).Info("no action tickets retrieved")
		return keys, nil
	}
	log.WithContext(ctx).WithField("count", len(actionTickets)).Info("Action tickets retrieved")

	for i := 0; i < len(actionTickets); i++ {
		decTicket, err := pastel.DecodeActionTicket(actionTickets[i].ActionTicketData.ActionTicket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to decode reg ticket")
			continue
		}
		actionTickets[i].ActionTicketData.ActionTicketData = *decTicket

		switch actionTickets[i].ActionTicketData.ActionType {
		case pastel.ActionTypeCascade:
			cascadeTicket, err := actionTickets[i].ActionTicketData.ActionTicketData.APICascadeTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTickets[i]).
					Warnf("Could not get cascade ticket for action ticket data")
				continue
			}

			keys = append(keys, cascadeTicket.RQIDs...)
		case pastel.ActionTypeSense:
			senseTicket, err := actionTickets[i].ActionTicketData.ActionTicketData.APISenseTicket()
			if err != nil {
				log.WithContext(ctx).WithField("actionRegTickets.ActionTicketData", actionTickets[i]).
					Warnf("Could not get sense ticket for action ticket data")
				continue
			}

			keys = append(keys, senseTicket.DDAndFingerprintsIDs...)
		}
	}

	return keys, nil
}

// GetSymbolFileByKey : Wrapper for p2p file storage service - retrieves a file from kademlia based on its key. Here, they should be raptorq symbol files.
func (service *SCService) GetSymbolFileByKey(ctx context.Context, key string, getFromLocalOnly bool) ([]byte, error) {
	if err := service.P2PClient.EnableKey(ctx, key); err != nil {
		log.WithContext(ctx).WithError(err).Error("error enabling the symbol file")
	}

	return service.P2PClient.Retrieve(ctx, key, getFromLocalOnly)
}

// StoreSymbolFile : Wrapper for p2p file storage service - stores a file in kademlia based on its key
func (service *SCService) StoreSymbolFile(ctx context.Context, data []byte) (key string, err error) {
	return service.P2PClient.Store(ctx, data, common.P2PDataRaptorQSymbol)
}

// RemoveSymbolFileByKey : Wrapper for p2p file storage service - removes a file from kademlia based on its key
func (service *SCService) RemoveSymbolFileByKey(ctx context.Context, key string) error {
	return service.P2PClient.Delete(ctx, key)
}

// GetListOfSupernode : Access the supernode service to get a list of all supernodes, including their id's and addresses.
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

// FilterOutSupernodes : FilterOutSupernodes gets the full list of supernodes and removes the nodesToBeIgnored
func (service *SCService) FilterOutSupernodes(listOfSupernodes []string, nodesToBeIgnored []string) []string {
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
func (service *SCService) GetNClosestSupernodeIDsToComparisonString(_ context.Context, n int, comparisonString string, listSupernodes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listSupernodes, ignores...)
}

// GetNClosestSupernodesToAGivenFileUsingKademlia : Wrapper for a utility function that accesses kademlia's distributed hash table to determine which nodes should be closest to a given string (hence hosting it)
func (service *SCService) GetNClosestSupernodesToAGivenFileUsingKademlia(ctx context.Context, n int, comparisonString string, ignores ...string) []string {
	log.WithContext(ctx).WithField("file_hash", comparisonString).Debug("file_hash against which closest sns required")
	return service.P2PClient.NClosestNodes(ctx, n, comparisonString, ignores...)
}

// GetNClosestFileHashesToAGivenComparisonString : Wrapper for a utility function that does xor string comparison to a list of strings and returns the smallest distance.
func (service *SCService) GetNClosestFileHashesToAGivenComparisonString(_ context.Context, n int, comparisonString string, listFileHashes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listFileHashes, ignores...)
}
