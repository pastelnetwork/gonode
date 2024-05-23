package storagechallenge

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

// BroadcastStorageChallengeMetrics worker broadcasts the metrics to the entire network
func (task *SCTask) BroadcastStorageChallengeMetrics(ctx context.Context) error {
	log.WithContext(ctx).Debug("BroadcastStorageChallengeMetrics has been invoked")

	pingInfos, err := task.historyDB.GetAllPingInfoForOnlineNodes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving ping info for online nodes")
		return errors.Errorf("error retrieving ping info")
	}

	sem := make(chan struct{}, 5) // Limit set to 10
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
			metrics, err := task.getSCMetrics(nodeInfo)
			if err != nil {
				logger.Error("error retrieving sc metrics for node")
				return
			}

			if len(metrics) == 0 {
				return
			}

			scMetricBatches := splitSCMetrics(metrics, 50)

			// Now you can iterate over these batches and process/send them
			for _, batch := range scMetricBatches {
				dataBytes, err := task.compressData(batch)
				if err != nil {
					logger.Error("error compressing sc metrics for node")
					return
				}

				msg := types.ProcessBroadcastChallengeMetricsRequest{
					Data:     dataBytes,
					SenderID: task.nodeID,
				}

				if err := task.SendSCMetrics(ctx, msg, nodeInfo.IPAddress); err != nil {
					log.WithError(err).Debug("error broadcasting SC metrics")
					return
				}
			}

			if err := task.historyDB.UpdateSCMetricsBroadcastTimestamp(nodeInfo.SupernodeID, broadcastAt); err != nil {
				log.WithContext(ctx).WithField("node_id", nodeInfo.SupernodeID).
					Error("error updating SCMetricsLastBroadcastAt timestamp")
			}
		}()
	}
	wg.Wait()

	return nil
}

func (task *SCTask) getSCMetrics(nodeInfo types.PingInfo) ([]types.StorageChallengeLogMessage, error) {
	var (
		zeroTime time.Time
		metrics  []types.StorageChallengeLogMessage
		err      error
	)

	if !nodeInfo.MetricsLastBroadcastAt.Valid {
		metrics, err = task.historyDB.StorageChallengeMetrics(zeroTime)
		if err != nil {
			return nil, err
		}
	} else {
		metrics, err = task.historyDB.StorageChallengeMetrics(nodeInfo.MetricsLastBroadcastAt.Time)
		if err != nil {
			return nil, err
		}
	}

	return metrics, nil
}

// compressData compresses Storage-Challenge metrics data using gzip.
func (task *SCTask) compressData(metrics []types.StorageChallengeLogMessage) (data []byte, err error) {
	if metrics == nil {
		return nil, errors.Errorf("storage metrics is nil")
	}

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	data, err = utils.Compress(jsonData, 1)
	if err != nil {
		return nil, errors.Errorf("error compressing sc metrics")
	}

	return data, nil
}

// SendSCMetrics sends SC metrics to the given node address
func (task *SCTask) SendSCMetrics(ctx context.Context, msg types.ProcessBroadcastChallengeMetricsRequest, processingSupernodeAddr string) error {
	logger := log.WithContext(ctx).WithField("node_address", processingSupernodeAddr)

	logger.Debug("sending storage-challenge metrics")

	//Connect over grpc
	newCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nodeClientConn, err := task.nodeClient.ConnectSN(newCtx, processingSupernodeAddr)
	if err != nil {
		logError(ctx, "BroadcastSCMetrics", err)
		return fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
	}
	defer func() {
		if nodeClientConn != nil {
			nodeClientConn.Close()
		}
	}()

	storageChallengeIF := nodeClientConn.StorageChallenge()

	return storageChallengeIF.BroadcastStorageChallengeMetrics(ctx, msg)
}

// decompressData decompresses the received data using gzip.
func (task *SCTask) decompressData(compressedData []byte) ([]types.StorageChallengeLogMessage, error) {
	var scMetrics []types.StorageChallengeLogMessage

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

func splitSCMetrics(metrics []types.StorageChallengeLogMessage, batchSize int) [][]types.StorageChallengeLogMessage {
	var batches [][]types.StorageChallengeLogMessage
	for batchSize < len(metrics) {
		metrics, batches = metrics[batchSize:], append(batches, metrics[0:batchSize:batchSize])
	}
	batches = append(batches, metrics)
	return batches
}

func (task *SCTask) BatchStoreSCMetrics(ctx context.Context, metrics []types.StorageChallengeLogMessage) error {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	if store != nil {
		err = store.BatchInsertSCMetrics(metrics)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error storing storage challenge metrics batch to DB")
			return err
		}
	}

	return nil
}
