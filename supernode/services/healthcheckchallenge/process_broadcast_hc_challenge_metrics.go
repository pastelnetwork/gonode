package healthcheckchallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
)

// ProcessBroadcastHealthCheckChallengeMetrics worker process the broadcast metrics received from other SNs
func (task *HCTask) ProcessBroadcastHealthCheckChallengeMetrics(ctx context.Context, req types.ProcessBroadcastHealthCheckChallengeMetricsRequest) error {
	logger := log.WithContext(ctx).WithField("sender_id", req.SenderID)

	logger.Debug("healthcheck-challenge metrics broadcast msg has been received")

	metrics, err := task.decompressData(req.Data)
	if err != nil {
		logger.WithError(err).Error("error decompressing metrics data")
		return err
	}

	if err := task.BatchStoreHCMetrics(ctx, metrics); err != nil {
		logger.WithError(err).Error("error storing hc metrics batch")
		return err
	}

	logger.Debug("hc metrics have been stored")
	return nil
}
