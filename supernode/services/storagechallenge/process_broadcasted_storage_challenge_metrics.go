package storagechallenge

import (
	"context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"strings"
)

// ProcessBroadcastedStorageChallengeMetrics worker process the broadcasted metrics received from other SNs
func (task *SCTask) ProcessBroadcastedStorageChallengeMetrics(ctx context.Context, req types.ProcessBroadcastChallengeMetricsRequest) error {
	logger := log.WithContext(ctx).WithField("sender_id", req.SenderID)

	logger.Info("storage-challenge metrics broadcast msg has been received")

	metrics, err := task.decompressData(req.Data)
	if err != nil {
		logger.WithError(err).Error("error decompressing metrics data")
		return err
	}

	for _, metric := range metrics {
		broadcastMetric := types.BroadcastLogMessage{
			ChallengeID: metric.ChallengeID,
			Data:        metric.Data,
			Recipient:   metric.Recipient,
			Observers:   metric.Observers,
			Challenger:  metric.Challenger,
		}
		err = task.historyDB.InsertBroadcastMessage(broadcastMetric)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				log.WithContext(ctx).
					WithField("challenge_id", metric.ChallengeID).
					WithField("data", metric.Data).
					Info("message already exists, not updating")
				return nil
			}

			log.WithContext(ctx).WithError(err).Error("Error storing challenge message to DB")
			return err
		}
	}

	logger.Info("metrics have been stored")
	return nil
}
