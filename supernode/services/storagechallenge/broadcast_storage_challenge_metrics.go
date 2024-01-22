package storagechallenge

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"

	"sync"
)

// BroadcastStorageChallengeMetrics worker broadcasts the metrics to the entire network
func (task *MetricTask) BroadcastStorageChallengeMetrics(ctx context.Context) error {
	log.WithContext(ctx).Infoln("BroadcastStorageChallengeMetrics has been invoked")

	pingInfos, err := task.historyDB.GetAllPingInfoForOnlineNodes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving ping info for online nodes")
		return errors.Errorf("error retrieving ping info")
	}

	var wg sync.WaitGroup
	for _, nodeInfo := range pingInfos {
		nodeInfo := nodeInfo

		wg.Add(1)
		go func() {
			defer wg.Done()
			logger := log.WithContext(ctx).WithField("node_address", nodeInfo.IPAddress)

			metrics, err := task.getBroadcastMessageMetrics(nodeInfo)
			if err != nil {
				logger.Error("error retrieving metrics for node")
				return
			}

			dataBytes, err := task.compressData(metrics)
			if err != nil {
				logger.Error("error compressing metrics for node")
				return
			}

			msg := types.ProcessBroadcastChallengeMetricsRequest{
				Data:     dataBytes,
				SenderID: task.nodeID,
			}

			if err := task.SendMessage(ctx, msg, nodeInfo.IPAddress); err != nil {
				logger.Error("error broadcasting metrics")
				return
			}
		}()
	}
	wg.Wait()

	for _, nodeInfo := range pingInfos {
		if err := task.historyDB.UpdateSCMetricsBroadcastTimestamp(nodeInfo.SupernodeID); err != nil {
			log.WithContext(ctx).WithField("node_id", nodeInfo.SupernodeID).
				Error("error updating SCMetricsLastBroadcastAt timestamp")
		}
	}

	return nil
}

func (task *MetricTask) getBroadcastMessageMetrics(nodeInfo types.PingInfo) ([]types.BroadcastMessageMetrics, error) {
	var (
		zeroTime time.Time
		metrics  []types.BroadcastMessageMetrics
		err      error
	)

	if !nodeInfo.MetricsLastBroadcastAt.Valid {
		metrics, err = task.historyDB.GetBroadcastMessageMetrics(zeroTime)
		if err != nil {
			return nil, err
		}
	} else {
		metrics, err = task.historyDB.GetBroadcastMessageMetrics(nodeInfo.SCMetricsLastBroadcastAt.Time)
		if err != nil {
			return nil, err
		}
	}

	return metrics, nil
}

// compressData compresses Storage-Challenge metrics data using gzip.
func (task *MetricTask) compressData(metrics []types.BroadcastMessageMetrics) (data []byte, err error) {
	if metrics == nil {
		return nil, errors.Errorf("storage metrics is nil")
	}

	jsonData, err := json.Marshal(metrics)
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

// SendMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *MetricTask) SendMessage(ctx context.Context, msg types.ProcessBroadcastChallengeMetricsRequest, processingSupernodeAddr string) error {
	logger := log.WithContext(ctx).WithField("node_address", processingSupernodeAddr)

	logger.Info("broadcasting storage-challenge metrics")

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

	storageChallengeIF := nodeClientConn.StorageChallenge()

	return storageChallengeIF.BroadcastStorageChallengeMetrics(ctx, msg)
}

// SignMessage signs the message using sender's pastelID and passphrase
func (task *MetricTask) SignMessage(ctx context.Context, d []byte) (sig []byte, dat []byte, err error) {
	signature, err := task.PastelClient.Sign(ctx, d, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, nil, errors.Errorf("error signing storage challenge message: %w", err)
	}

	return signature, d, nil
}

// decompressData decompresses the received data using gzip.
func (task *SCTask) decompressData(compressedData []byte) ([]types.BroadcastMessageMetrics, error) {
	var broadcatMessageMetrics []types.BroadcastMessageMetrics

	// Decompress using gzip
	gz, err := gzip.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	decompressedData, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	// Deserialize JSON back to data structure
	if err = json.Unmarshal(decompressedData, &broadcatMessageMetrics); err != nil {
		return nil, err
	}

	return broadcatMessageMetrics, nil
}
