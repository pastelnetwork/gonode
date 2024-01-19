package storagechallenge

import (
	"context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
)

// ProcessBroadcastedStorageChallengeMetrics worker process the broadcasted metrics received from other SNs
func (task *SCTask) ProcessBroadcastedStorageChallengeMetrics(ctx context.Context, req types.ProcessBroadcastChallengeMetricsRequest) error {
	logger := log.WithContext(ctx).WithField("sender_id", req.SenderID)

	logger.Info("storage-challenge metrics broadcast msg has been received")

	return nil
}
