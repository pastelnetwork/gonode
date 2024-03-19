package storagechallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
)

// ProcessBroadcastStorageChallengeMetrics worker process the broadcast metrics received from other SNs
func (task *SCTask) ProcessBroadcastStorageChallengeMetrics(ctx context.Context, req types.ProcessBroadcastChallengeMetricsRequest) error {
	logger := log.WithContext(ctx).WithField("sender_id", req.SenderID)

	logger.Debug("storage-challenge metrics broadcast msg has been received")

	metrics, err := task.decompressData(req.Data)
	if err != nil {
		logger.WithError(err).Error("error decompressing metrics data")
		return err
	}

	if err := task.BatchStoreSCMetrics(ctx, metrics); err != nil {
		logger.WithError(err).Error("error storing sc metrics batch")
		return err
	}

	logger.Debug("sc metrics have been stored")
	return nil
}
