package healthcheckchallenge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	healthCheckBlockInterval         int32 = 3
	broadcastHCMetricRegularInterval       = 30 * time.Minute
	defaultTimerBlockCheckDuration         = 30 * time.Second
)

type HCService struct {
	*common.SuperNodeService
	config *Config

	nodeID                            string
	nodeClient                        node.ClientInterface
	healthCheckChallengeExpiredBlocks int32
	numberOfChallengeReplicas         int
	numberOfVerifyingNodes            int
	historyDB                         queries.LocalStoreInterface

	currentBlockCount int32
	// currently unimplemented, default always used instead.
	challengeStatusObserver SaveChallengeState

	localKeys sync.Map
}

// CheckNextBlockAvailable calls pasteld and checks if a new block is available
func (service *HCService) CheckNextBlockAvailable(ctx context.Context) bool {
	blockCount, err := service.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		return false
	}

	if blockCount > service.currentBlockCount+healthCheckBlockInterval {
		atomic.StoreInt32(&service.currentBlockCount, blockCount)
		return true
	}

	return false
}

// Run : health check challenge service will run continuously to generate health check challenges.
func (service *HCService) Run(ctx context.Context) error {
	log.WithContext(ctx).Info("Health check challenge service run has been invoked")
	//does this need to be in its own goroutine?
	go func() {
		if err := service.RunHelper(ctx, service.config.PastelID, logPrefix); err != nil {
			log.WithContext(ctx).WithError(err).Error("HealthCheckChallengeService:RunHelper")
		}
	}()

	go service.BroadcastHealthCheckChallengeMetricsWorker(ctx)

	for {
		select {
		case <-time.After(defaultTimerBlockCheckDuration):

			if service.CheckNextBlockAvailable(ctx) {
				service.RunGenerateHealthCheckChallenges(ctx)
				log.WithContext(ctx).Debug("Would normally generate a healthcheck challenge")
			}
		case <-ctx.Done():
			log.Println("Context done being called in generatehealthcheckchallenge loop in service.go")
			return nil
		}
	}
}

// NewHCTask : HealthCheck challenge task handles the duties of generating, processing, and verifying HealthCheck challenges
func (service *HCService) NewHCTask() *HCTask {
	task := NewHCTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the HealthCheck Challenge by the id
func (service *HCService) Task(id string) *HCTask {
	scTask, ok := service.Worker.Task(id).(*HCTask)
	if !ok {
		log.Error("Error typecasting task to healthcheck challenge task")
		return nil
	}

	log.Info("type casted successfully")
	return scTask
}

// executeTask executes the health check  metric task.
func (service *HCService) executeMetricsBroadcastTask(ctx context.Context) {
	newCtx := context.Background()
	task := service.NewHCTask()
	task.BroadcastHealthCheckChallengeMetrics(newCtx)

	log.WithContext(ctx).Debug("health-check challenge metrics broadcasted")
}

func (service *HCService) RunGenerateHealthCheckChallenges(ctx context.Context) {
	startTime := calculateStartTimeForHealthCheckTrigger(service.nodeID)
	log.WithContext(ctx).WithField("start_time", startTime).
		Info("health check challenge metric worker will execute on the mentioned time")

	// Wait until the start time
	time.Sleep(time.Until(startTime))

	newCtx := context.Background()
	task := service.NewHCTask()

	task.GenerateHealthCheckChallenges(newCtx)
}

// BroadcastHealthCheckChallengeMetricsWorker broadcast the health check challenge metrics to the entire network
func (service *HCService) BroadcastHealthCheckChallengeMetricsWorker(ctx context.Context) {
	log.WithContext(ctx).Info("health check challenge metrics worker func has been invoked")

	startTime := calculateStartTime(service.nodeID)
	log.WithContext(ctx).WithField("start_time", startTime).
		Info("health check challenge metric worker will execute on the mentioned time")

	// Wait until the start time
	time.Sleep(time.Until(startTime))

	// Run the first task immediately
	service.executeMetricsBroadcastTask(ctx)

	for {
		select {
		case <-time.After(broadcastHCMetricRegularInterval):
			service.executeMetricsBroadcastTask(context.Background())
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

// calculateStartTime calculates the start time for the periodic task based on the hash of the PastelID.
func calculateStartTimeForHealthCheckTrigger(pastelID string) time.Time {
	hash := sha256.Sum256([]byte(pastelID))
	hashString := hex.EncodeToString(hash[:])
	delaySecs := int(hashString[0]) % 20 // simplistic hash-based delay calculation

	return time.Now().Add(time.Duration(delaySecs) * time.Second)
}

// NewService : Create a new healthcheck challenge service
//
//	Inheriting from SuperNodeService allows us to use common methods for pastelclient, p2p, and rqClient.
func NewService(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface,
	p2p p2p.Client, challengeStatusObserver SaveChallengeState, historyDB queries.LocalStoreInterface) *HCService {
	return &HCService{
		config:                            config,
		SuperNodeService:                  common.NewSuperNodeService(fileStorage, pastelClient, p2p),
		nodeClient:                        nodeClient,
		healthCheckChallengeExpiredBlocks: config.HealthCheckChallengeExpiredBlocks,
		// repository:                    newRepository(p2p, pastelClient, challengeStatusObserver),
		nodeID:                    config.PastelID,
		numberOfChallengeReplicas: config.NumberOfChallengeReplicas,
		numberOfVerifyingNodes:    config.NumberOfVerifyingNodes,
		challengeStatusObserver:   challengeStatusObserver,
		localKeys:                 sync.Map{},
		historyDB:                 historyDB,
	}
}

// GetListOfSupernode : Access the supernode service to get a list of all supernodes, including their id's and addresses.
// This is used to enumerate supernodes both for calculation and connection
func (service *HCService) GetListOfSupernode(ctx context.Context) ([]string, error) {
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
func (service *HCService) FilterOutSupernodes(listOfSupernodes []string, nodesToBeIgnored []string) []string {
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
func (service *HCService) GetNClosestSupernodeIDsToComparisonString(_ context.Context, n int, comparisonString string, listSupernodes []string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, comparisonString, listSupernodes, ignores...)
}
