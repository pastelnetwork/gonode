package selfhealing

import (
	"context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
)

// ProcessBroadcastedSelfHealingMetrics worker process the broadcasted metrics received from other SNs
func (task *SHTask) ProcessBroadcastedSelfHealingMetrics(ctx context.Context, req types.ProcessBroadcastMetricsRequest) error {
	logger := log.WithContext(ctx).WithField("sender_id", req.SenderID)

	logger.Info("self-healing metrics broadcast msg has been received")

	return nil
}
