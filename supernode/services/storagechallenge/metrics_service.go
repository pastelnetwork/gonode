package storagechallenge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const broadcastMetricRegularInterval = 60 * time.Minute

// MetricService broadcast storage-challenge execution and triggering metrics to the entire network
type MetricService struct {
	*common.SuperNodeService
	config        *Config
	pastelHandler *mixins.PastelHandler
	nodeID        string
	nodeClient    node.ClientInterface
	historyDB     storage.LocalStoreInterface
}

// Run : storage-challenge metric service will run continuously to generate storage-challenge.
func (service *MetricService) Run(ctx context.Context) error {
	log.WithContext(ctx).Info("storage-challenge-metric service run func has been invoked")

	startTime := calculateStartTime(service.nodeID)

	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	// Wait until the start time
	time.Sleep(time.Until(startTime))

	// Run the first task immediately
	service.executeTask(ctx)

	for {
		select {
		case <-time.After(broadcastMetricRegularInterval):
			service.executeTask(context.Background())
		case <-ctx.Done():
			log.Println("Context done being called in worker")
			return nil
		}
	}
}

// NewStorageChallengeMetricTask : storage-challenge-metric task handles the broadcasting of metrics
func (service *MetricService) NewStorageChallengeMetricTask() *MetricTask {
	task := NewStorageChallengeMetricTask(service)
	service.Worker.AddTask(task)
	return task
}

// executeTask executes the self-healing metric task.
func (service *MetricService) executeTask(ctx context.Context) {
	newCtx := context.Background()
	task := service.NewStorageChallengeMetricTask()
	task.BroadcastStorageChallengeMetrics(newCtx)

	log.WithContext(ctx).Debug("storage-challenge metric broadcasted")
}

// NewStorageChallengeMetricTask returns a new Storage-challenge Metric Task instance.
func NewStorageChallengeMetricTask(service *MetricService) *MetricTask {
	task := &MetricTask{
		MetricService: service,
	}

	return task
}

// calculateStartTime calculates the start time for the periodic task based on the hash of the PastelID.
func calculateStartTime(pastelID string) time.Time {
	hash := sha256.Sum256([]byte(pastelID))
	hashString := hex.EncodeToString(hash[:])
	delayMinutes := int(hashString[0]) % 60 // simplistic hash-based delay calculation

	return time.Now().Add(time.Duration(delayMinutes) * time.Minute)
}

// NewMetricsService : Create a new metrics-storage-challenge service
func NewMetricsService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface,
	historyDB storage.LocalStoreInterface) *SCService {
	return &SCService{
		config:        config,
		nodeClient:    nodeClient,
		nodeID:        config.PastelID,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		historyDB:     historyDB,
	}
}
