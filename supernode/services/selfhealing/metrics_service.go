package selfhealing

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

// MetricService broadcast self-healing execution and triggering metrics to the entire network
type MetricService struct {
	*common.SuperNodeService
	config        *Config
	pastelHandler *mixins.PastelHandler
	nodeID        string
	nodeClient    node.ClientInterface
	historyDB     storage.LocalStoreInterface
}

// Run : self-healing metric service will run continuously to generate self-healing.
func (service *MetricService) Run(ctx context.Context) error {
	log.WithContext(ctx).Info("self-healing-metric service run func has been invoked")

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
			log.Println("Context done being called in file-healing worker")
			return nil
		}
	}
}

// NewSelfHealingMetricTask : self-healing-metric task handles the broadcasting of metrics
func (service *MetricService) NewSelfHealingMetricTask() *MetricTask {
	task := NewSelfHealingMetricTask(service)
	service.Worker.AddTask(task)
	return task
}

// executeTask executes the self-healing metric task.
func (service *MetricService) executeTask(ctx context.Context) {
	newCtx := context.Background()
	task := service.NewSelfHealingMetricTask()
	task.BroadcastSelfHealingMetrics(newCtx)

	log.WithContext(ctx).Debug("self-healing metric broadcasted")
}

// NewSelfHealingMetricTask returns a new Self-Healing Metric Task instance.
func NewSelfHealingMetricTask(service *MetricService) *MetricTask {
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

// NewMetricsService : Create a new metrics-self-healing service
func NewMetricsService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface,
	historyDB storage.LocalStoreInterface) *SHService {
	return &SHService{
		config:        config,
		nodeClient:    nodeClient,
		nodeID:        config.PastelID,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		historyDB:     historyDB,
	}
}
