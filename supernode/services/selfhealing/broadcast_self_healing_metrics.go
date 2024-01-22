package selfhealing

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"io"
	"sync"
	"time"
)

// BroadcastSelfHealingMetrics worker broadcasts the metrics to the entire network
func (task *SHTask) BroadcastSelfHealingMetrics(ctx context.Context) error {
	log.WithContext(ctx).Infoln("BroadcastSelfHealingMetrics has been invoked")

	pingInfos, err := task.historyDB.GetAllPingInfoForOnlineNodes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving ping info for online nodes")
		return errors.Errorf("error retrieving ping info")
	}
	log.WithContext(ctx).WithField("total_node_infos", len(pingInfos)).
		Info("online node info has been retrieved for metrics broadcast")

	var wg sync.WaitGroup
	for _, nodeInfo := range pingInfos {
		nodeInfo := nodeInfo

		if nodeInfo.SupernodeID == task.nodeID {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			logger := log.WithContext(ctx).WithField("node_address", nodeInfo.IPAddress)

			executionMetrics, generationMetrics, err := task.getExecutionAndGenerationMetrics(nodeInfo)
			if err != nil {
				logger.Error("error retrieving execution & generation metrics for node")
				return
			}

			dataBytes, err := task.compressData(executionMetrics, generationMetrics)
			if err != nil {
				logger.WithError(err).Error("error compressing data")
				return
			}
			logger.Info("execution & generation metrics have been compressed")

			msg := types.ProcessBroadcastMetricsRequest{
				Data:            dataBytes,
				SenderID:        task.nodeID,
				SenderSignature: nil,
			}

			if err := task.SendBroadcastMessage(ctx, msg, nodeInfo.IPAddress); err != nil {
				logger.Error("error broadcasting metrics")
				return
			}

			if err := task.historyDB.UpdateMetricsBroadcastTimestamp(nodeInfo.SupernodeID); err != nil {
				log.WithContext(ctx).WithField("node_id", nodeInfo.SupernodeID).
					Error("error updating broadcastAt timestamp")
			}
		}()
	}

	wg.Wait()

	log.WithContext(ctx).Info("self-healing metrics have been broadcast")

	return nil
}

func (task *SHTask) getExecutionAndGenerationMetrics(nodeInfo types.PingInfo) ([]types.SelfHealingExecutionMetric, []types.SelfHealingGenerationMetric, error) {
	var (
		zeroTime          time.Time
		executionMetrics  []types.SelfHealingExecutionMetric
		generationMetrics []types.SelfHealingGenerationMetric
		err               error
	)
	if !nodeInfo.MetricsLastBroadcastAt.Valid {
		executionMetrics, err = task.historyDB.GetSelfHealingExecutionMetrics(zeroTime)
		if err != nil {
			return nil, nil, err
		}

		generationMetrics, err = task.historyDB.GetSelfHealingGenerationMetrics(zeroTime)
		if err != nil {
			return nil, nil, err
		}
	} else {
		executionMetrics, err = task.historyDB.GetSelfHealingExecutionMetrics(nodeInfo.MetricsLastBroadcastAt.Time)
		if err != nil {
			return nil, nil, err
		}

		generationMetrics, err = task.historyDB.GetSelfHealingGenerationMetrics(nodeInfo.MetricsLastBroadcastAt.Time)
		if err != nil {
			return nil, nil, err
		}
	}

	return executionMetrics, generationMetrics, nil
}

// compressData compresses Self-Healing metrics data using gzip.
func (task *SHTask) compressData(executionMetrics []types.SelfHealingExecutionMetric, generationMetrics []types.SelfHealingGenerationMetric) (data []byte, err error) {
	if executionMetrics == nil && generationMetrics == nil {
		return nil, errors.Errorf("no new execution and generation metrics")
	}

	combinedMetrics := types.SelfHealingCombinedMetrics{
		ExecutionMetrics:  executionMetrics,
		GenerationMetrics: generationMetrics,
	}

	jsonData, err := json.Marshal(combinedMetrics)
	if err != nil {
		return nil, err
	}

	// Compress using gzip
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err = gz.Write(jsonData); err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decompressData decompresses the received data using gzip.
func (task *SHTask) decompressData(compressedData []byte) ([]types.SelfHealingExecutionMetric, []types.SelfHealingGenerationMetric, error) {
	var combinedMetrics struct {
		ExecutionMetrics  []types.SelfHealingExecutionMetric
		GenerationMetrics []types.SelfHealingGenerationMetric
	}

	// Decompress using gzip
	gz, err := gzip.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		return nil, nil, err
	}
	defer gz.Close()

	decompressedData, err := io.ReadAll(gz)
	if err != nil {
		return nil, nil, err
	}

	// Deserialize JSON back to data structure
	if err = json.Unmarshal(decompressedData, &combinedMetrics); err != nil {
		return nil, nil, err
	}

	return combinedMetrics.ExecutionMetrics, combinedMetrics.GenerationMetrics, nil
}

// SendBroadcastMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *SHTask) SendBroadcastMessage(ctx context.Context, msg types.ProcessBroadcastMetricsRequest, processingSupernodeAddr string) error {
	logger := log.WithContext(ctx).WithField("node_address", processingSupernodeAddr)

	logger.Info("broadcasting self-healing metrics")

	//Connect over grpc
	newCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nodeClientConn, err := task.nodeClient.Connect(newCtx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
		logger.WithField("method", "SendMessage").Warn(err.Error())
		return err
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
