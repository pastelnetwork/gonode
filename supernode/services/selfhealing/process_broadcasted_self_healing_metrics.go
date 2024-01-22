package selfhealing

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
)

// ProcessBroadcastedSelfHealingMetrics worker process the broadcasted metrics received from other SNs
func (task *SHTask) ProcessBroadcastedSelfHealingMetrics(ctx context.Context, req types.ProcessBroadcastMetricsRequest) error {
	logger := log.WithContext(ctx).WithField("sender_id", req.SenderID)

	logger.Info("self-healing broadcast metrics msg has been received")

	execMetrics, genMetrics, err := task.decompressData(req.Data)
	if err != nil {
		logger.WithError(err).Error("error decompressing metrics data")
		return err
	}
	log.WithContext(ctx).Info("data has been decompressed")

	for _, metric := range genMetrics {
		err := task.StoreSelfHealingGenerationMetrics(ctx, metric)
		if err != nil {
			logger.WithField("trigger_id", metric.TriggerID).WithField("message_type", metric.MessageType).
				Error("error storing generation metric")
			return err
		}
	}
	logger.Info("generation metrics have been stored")

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

	return nil
}

// VerifySignature verifies the sender's signature on message
func (task *SHTask) VerifySignature(ctx context.Context, data, senderSignature []byte, senderID string) (bool, error) {
	isVerified, err := task.PastelClient.Verify(ctx, data, string(senderSignature), senderID, pastel.SignAlgorithmED448)
	if err != nil {
		return false, errors.Errorf("error verifying self-healing message: %w", err)
	}

	return isVerified, nil
}
