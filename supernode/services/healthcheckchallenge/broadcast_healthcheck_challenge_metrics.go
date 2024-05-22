package healthcheckchallenge

import (
	"context"
	"fmt"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
)

// BroadcastHealthCheckChallengeMetrics worker broadcasts the metrics to the entire network
func (task *HCTask) BroadcastHealthCheckChallengeMetrics(ctx context.Context) error {
	log.WithContext(ctx).Debug("BroadcastHealthCheckChallengeMetrics has been invoked")

	pingInfos, err := task.historyDB.GetAllPingInfoForOnlineNodes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving ping info for online nodes")
		return errors.Errorf("error retrieving ping info")
	}

	sem := make(chan struct{}, 5) // Limit set to 5
	var wg sync.WaitGroup
	for _, nodeInfo := range pingInfos {
		nodeInfo := nodeInfo

		if nodeInfo.SupernodeID == task.nodeID {
			continue
		}

		wg.Add(1)

		go func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()

			logger := log.WithContext(ctx).WithField("node_address", nodeInfo.IPAddress)

			broadcastAt := time.Now().UTC()
			metrics, err := task.getHCMetrics(nodeInfo)
			if err != nil {
				logger.Error("error retrieving hc metrics for node")
				return
			}

			if len(metrics) == 0 {
				return
			}

			scMetricBatches := splitHCMetrics(metrics, 50)

			// Now you can iterate over these batches and process/send them
			for _, batch := range scMetricBatches {
				dataBytes, err := task.compressData(batch)
				if err != nil {
					logger.Error("error compressing hc metrics for node")
					return
				}

				msg := types.ProcessBroadcastHealthCheckChallengeMetricsRequest{
					Data:     dataBytes,
					SenderID: task.nodeID,
				}

				if err := task.SendHCMetrics(ctx, msg, nodeInfo.IPAddress); err != nil {
					log.WithError(err).Debug("error broadcasting HC metrics")
					return
				}
			}

			if err := task.historyDB.UpdateHCMetricsBroadcastTimestamp(nodeInfo.SupernodeID, broadcastAt); err != nil {
				log.WithContext(ctx).WithField("node_id", nodeInfo.SupernodeID).
					Error("error updating HCMetricsLastBroadcastAt timestamp")
			}
		}()
	}
	wg.Wait()

	return nil
}

func (task *HCTask) getHCMetrics(nodeInfo types.PingInfo) ([]types.HealthCheckChallengeLogMessage, error) {
	var (
		zeroTime time.Time
		metrics  []types.HealthCheckChallengeLogMessage
		err      error
	)

	if !nodeInfo.HealthCheckMetricsLastBroadcastAt.Valid {
		metrics, err = task.historyDB.HealthCheckChallengeMetrics(zeroTime)
		if err != nil {
			return nil, err
		}
	} else {
		metrics, err = task.historyDB.HealthCheckChallengeMetrics(nodeInfo.HealthCheckMetricsLastBroadcastAt.Time)
		if err != nil {
			return nil, err
		}
	}

	return metrics, nil
}

// compressData compresses HealthCheck-Challenge metrics data using gzip.
func (task *HCTask) compressData(metrics []types.HealthCheckChallengeLogMessage) (data []byte, err error) {
	if metrics == nil {
		return nil, errors.Errorf("health check metrics is nil")
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	data, err = utils.Compress(jsonData, 1)
	if err != nil {
		return nil, errors.Errorf("error compressing hc metrics")
	}

	return data, nil
}

// SendHCMetrics sends HC metrics to the given node address
func (task *HCTask) SendHCMetrics(ctx context.Context, msg types.ProcessBroadcastHealthCheckChallengeMetricsRequest, processingSupernodeAddr string) error {
	logger := log.WithContext(ctx).WithField("node_address", processingSupernodeAddr)

	logger.Debug("sending healthcheck-challenge metrics")

	//Connect over grpc
	newCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nodeClientConn, err := task.nodeClient.ConnectSN(newCtx, processingSupernodeAddr)
	if err != nil {
		logError(ctx, "BroadcastHealthCheckChallengeMetrics", err)
		return fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
	}
	defer func() {
		if nodeClientConn != nil {
			nodeClientConn.Close()
		}
	}()

	healthCheckChallengeIF := nodeClientConn.HealthCheckChallenge()

	return healthCheckChallengeIF.BroadcastHealthCheckChallengeMetrics(ctx, msg)
}

// decompressData decompresses the received data using gzip.
func (task *HCTask) decompressData(compressedData []byte) ([]types.HealthCheckChallengeLogMessage, error) {
	var scMetrics []types.HealthCheckChallengeLogMessage

	// Decompress using gzip
	decompressedData, err := utils.Decompress(compressedData)
	if err != nil {
		return nil, errors.Errorf("error decompressing data")
	}

	// Deserialize JSON back to data structure
	if err = json.Unmarshal(decompressedData, &scMetrics); err != nil {
		return nil, err
	}

	return scMetrics, nil
}

func splitHCMetrics(metrics []types.HealthCheckChallengeLogMessage, batchSize int) [][]types.HealthCheckChallengeLogMessage {
	var batches [][]types.HealthCheckChallengeLogMessage
	for batchSize < len(metrics) {
		metrics, batches = metrics[batchSize:], append(batches, metrics[0:batchSize:batchSize])
	}
	batches = append(batches, metrics)
	return batches
}

func (task *HCTask) BatchStoreHCMetrics(ctx context.Context, metrics []types.HealthCheckChallengeLogMessage) error {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	if store != nil {
		err = store.BatchInsertHCMetrics(metrics)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error storing health check challenge metrics batch to DB")
			return err
		}
	}

	return nil
}
