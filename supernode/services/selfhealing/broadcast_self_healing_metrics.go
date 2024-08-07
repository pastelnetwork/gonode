package selfhealing

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"sync"
	"time"
)

// BroadcastSelfHealingMetrics worker broadcasts the metrics to the entire network
func (task *SHTask) BroadcastSelfHealingMetrics(ctx context.Context) error {
	log.WithContext(ctx).Debug("BroadcastSelfHealingMetrics has been invoked")

	pingInfos, err := task.historyDB.GetAllPingInfoForOnlineNodes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving ping info for online nodes")
		return errors.Errorf("error retrieving ping info")
	}
	log.WithContext(ctx).WithField("total_node_infos", len(pingInfos)).
		Debug("online node info has been retrieved for metrics broadcast")

	sem := make(chan struct{}, 5) // Limit set to 10
	var wg sync.WaitGroup
	for _, nodeInfo := range pingInfos {
		nodeInfo := nodeInfo

		if nodeInfo.SupernodeID == task.nodeID {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }() // Release the token back into the channel

			logger := log.WithContext(ctx).WithField("node_address", nodeInfo.IPAddress)

			executionMetrics, generationMetrics, err := task.getExecutionAndGenerationMetrics(nodeInfo)
			if err != nil {
				logger.WithError(err).Error("error retrieving execution & generation metrics for node")
				return
			}

			generationMetricBatches := splitSelfHealingGenerationMetrics(generationMetrics, 50)
			executionMetricBatches := splitSelfHealingExecutionMetrics(executionMetrics, 50)

			// Now you can iterate over these batches and process/send them
			for _, batch := range generationMetricBatches {
				dataBytes, err := task.compressGenerationMetricsData(batch)
				if err != nil {
					logger.WithError(err).Debug("error compressing generation metrics data")
					return
				}

				msg := types.ProcessBroadcastMetricsRequest{
					Type:     types.GenerationSelfHealingMetricType,
					Data:     dataBytes,
					SenderID: task.nodeID,
				}

				if err := task.SendBroadcastMessage(ctx, msg, nodeInfo.IPAddress); err != nil {
					log.WithError(err).Debug("error broadcasting generation metrics")
					return
				}

				if err := task.historyDB.UpdateGenerationMetricsBroadcastTimestamp(nodeInfo.SupernodeID); err != nil {
					log.WithContext(ctx).WithField("node_id", nodeInfo.SupernodeID).
						Error("error updating generation metrics broadcastAt timestamp")
				}
			}

			for _, batch := range executionMetricBatches {
				dataBytes, err := task.compressExecutionMetricsData(batch)
				if err != nil {
					logger.WithError(err).Debug("error compressing execution metrics data")
					return
				}

				msg := types.ProcessBroadcastMetricsRequest{
					Type:     types.ExecutionSelfHealingMetricType,
					Data:     dataBytes,
					SenderID: task.nodeID,
				}

				if err := task.SendBroadcastMessage(ctx, msg, nodeInfo.IPAddress); err != nil {
					log.WithError(err).Debug("error broadcasting execution metrics")
					return
				}
			}

			if err := task.historyDB.UpdateExecutionMetricsBroadcastTimestamp(nodeInfo.SupernodeID); err != nil {
				log.WithContext(ctx).WithField("node_id", nodeInfo.SupernodeID).
					Error("error updating execution metrics broadcast at timestamp")
			}
		}()
	}

	wg.Wait()

	log.WithContext(ctx).Debug("self-healing metrics have been broadcast")

	return nil
}

func (task *SHTask) getExecutionAndGenerationMetrics(nodeInfo types.PingInfo) ([]types.SelfHealingExecutionMetric, []types.SelfHealingGenerationMetric, error) {
	var (
		zeroTime          time.Time
		err               error
		executionMetrics  []types.SelfHealingExecutionMetric
		generationMetrics []types.SelfHealingGenerationMetric
	)
	if !nodeInfo.ExecutionMetricsLastBroadcastAt.Valid {
		executionMetrics, err = task.historyDB.GetSelfHealingExecutionMetrics(zeroTime)
		if err != nil {
			return nil, nil, err
		}
	} else {
		executionMetrics, err = task.historyDB.GetSelfHealingExecutionMetrics(nodeInfo.ExecutionMetricsLastBroadcastAt.Time.UTC())
		if err != nil {
			return nil, nil, err
		}
	}

	if !nodeInfo.GenerationMetricsLastBroadcastAt.Valid {
		generationMetrics, err = task.historyDB.GetSelfHealingGenerationMetrics(zeroTime)
		if err != nil {
			return nil, nil, err
		}
	} else {
		generationMetrics, err = task.historyDB.GetSelfHealingGenerationMetrics(nodeInfo.GenerationMetricsLastBroadcastAt.Time.UTC())
		if err != nil {
			return nil, nil, err
		}
	}

	return executionMetrics, generationMetrics, nil
}

// compressExecutionMetricsData compresses Self-Healing metrics data using gzip.
func (task *SHTask) compressExecutionMetricsData(executionMetrics []types.SelfHealingExecutionMetric) (data []byte, err error) {
	if executionMetrics == nil {
		return nil, errors.Errorf("no new execution metrics")
	}

	jsonData, err := json.Marshal(executionMetrics)
	if err != nil {
		return nil, err
	}

	data, err = utils.Compress(jsonData, 1)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// compressExecutionMetricsData compresses Self-Healing metrics data using gzip.
func (task *SHTask) compressGenerationMetricsData(generationMetrics []types.SelfHealingGenerationMetric) (data []byte, err error) {
	if generationMetrics == nil {
		return nil, errors.Errorf("no new generation metrics")
	}

	jsonData, err := json.Marshal(generationMetrics)
	if err != nil {
		return nil, err
	}

	data, err = utils.Compress(jsonData, 1)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SendBroadcastMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *SHTask) SendBroadcastMessage(ctx context.Context, msg types.ProcessBroadcastMetricsRequest, processingSupernodeAddr string) error {
	//Connect over grpc
	newCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nodeClientConn, err := task.nodeClient.ConnectSN(newCtx, processingSupernodeAddr)
	if err != nil {
		logError(ctx, "BroadcastSelfHealingMetrics", err)
		return fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
	}
	defer func() {
		if nodeClientConn != nil {
			nodeClientConn.Close()
		}
	}()

	selfHealingIF := nodeClientConn.SelfHealingChallenge()

	return selfHealingIF.BroadcastSelfHealingMetrics(ctx, msg)
}

// SignBroadcastMessage signs the message using sender's pastelID and passphrase
func (task *SHTask) SignBroadcastMessage(ctx context.Context, d []byte) (sig []byte, err error) {
	signature, err := task.PastelClient.Sign(ctx, d, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, errors.Errorf("error signing self-healing metric message: %w", err)
	}

	return signature, nil
}

func splitSelfHealingExecutionMetrics(metrics []types.SelfHealingExecutionMetric, batchSize int) [][]types.SelfHealingExecutionMetric {
	var batches [][]types.SelfHealingExecutionMetric
	for batchSize < len(metrics) {
		metrics, batches = metrics[batchSize:], append(batches, metrics[0:batchSize:batchSize])
	}
	batches = append(batches, metrics)
	return batches
}

func splitSelfHealingGenerationMetrics(metrics []types.SelfHealingGenerationMetric, batchSize int) [][]types.SelfHealingGenerationMetric {
	var batches [][]types.SelfHealingGenerationMetric
	for batchSize < len(metrics) {
		metrics, batches = metrics[batchSize:], append(batches, metrics[0:batchSize:batchSize])
	}
	batches = append(batches, metrics)
	return batches
}
