package selfhealing

import (
	"context"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
)

// ProcessBroadcastedSelfHealingMetrics worker process the broadcasted metrics received from other SNs
func (task *SHTask) ProcessBroadcastedSelfHealingMetrics(ctx context.Context, req types.ProcessBroadcastMetricsRequest) error {
	logger := log.WithContext(ctx).WithField("sender_id", req.SenderID)
	switch req.Type {
	case types.GenerationSelfHealingMetricType:
		genMetrics, err := task.decompressGenerationMetricsData(req.Data)
		if err != nil {
			logger.WithError(err).Error("error decompressing generation metrics data")
			return err
		}

		for _, metric := range genMetrics {
			err := task.StoreSelfHealingGenerationMetrics(ctx, metric)
			if err != nil {
				logger.WithField("trigger_id", metric.TriggerID).
					WithField("message_type", metric.MessageType).
					Error("error storing generation metric")
				return err
			}
		}
		logger.Info("generation metrics have been stored")
	case types.ExecutionSelfHealingMetricType:
		execMetrics, err := task.decompressExecutionMetricsData(req.Data)
		if err != nil {
			logger.WithError(err).Error("error decompressing execution metrics data")
			return err
		}
		log.WithContext(ctx).Info("execution metric data has been decompressed")

		for _, metric := range execMetrics {
			err := task.StoreSelfHealingExecutionMetrics(ctx, metric)
			if err != nil {
				logger.WithField("trigger_id", metric.TriggerID).
					WithField("message_type", metric.MessageType).
					WithField("challenge_id", metric.ChallengeID).
					Error("error storing execution metric")
				return err
			}
		}
		logger.Info("execution metrics have been stored")
	}

	return nil
}

// decompressGenerationMetricsData decompresses the received generation metrics data using gzip.
func (task *SHTask) decompressGenerationMetricsData(compressedData []byte) ([]types.SelfHealingGenerationMetric, error) {
	var generationMetrics []types.SelfHealingGenerationMetric

	decompressedData, err := utils.Decompress(compressedData)
	if err != nil {
		return nil, err
	}

	// Deserialize JSON back to data structure
	if err = json.Unmarshal(decompressedData, &generationMetrics); err != nil {
		return nil, err
	}

	return generationMetrics, nil
}

// decompressExecutionMetricsData decompresses the received execution metric data using gzip.
func (task *SHTask) decompressExecutionMetricsData(compressedData []byte) ([]types.SelfHealingExecutionMetric, error) {
	var executionMetrics []types.SelfHealingExecutionMetric

	decompressedData, err := utils.Decompress(compressedData)
	if err != nil {
		return nil, err
	}

	// Deserialize JSON back to data structure
	if err = json.Unmarshal(decompressedData, &executionMetrics); err != nil {
		return nil, err
	}

	return executionMetrics, nil
}
