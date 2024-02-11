package storage

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils/metrics"
)

// LocalStoreInterface is interface for local sqlite store
type LocalStoreInterface interface {
	CloseHistoryDB(ctx context.Context)

	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)

	InsertStorageChallengeMessage(challenge types.StorageChallengeLogMessage) error
	InsertBroadcastMessage(challenge types.BroadcastLogMessage) error
	QueryStorageChallengeMessage(challengeID string, messageType int) (challenge types.StorageChallengeLogMessage, err error)
	CleanupStorageChallenges() (err error)

	UpsertPingHistory(pingInfo types.PingInfo) error
	GetPingInfoBySupernodeID(supernodeID string) (*types.PingInfo, error)
	GetAllPingInfos() (types.PingInfos, error)
	GetWatchlistPingInfo() ([]types.PingInfo, error)
	GetAllPingInfoForOnlineNodes() (types.PingInfos, error)
	UpdatePingInfo(supernodeID string, isOnWatchlist, isAdjusted bool) error
	UpdateMetricsBroadcastTimestamp(nodeID string) error
	UpdateGenerationMetricsBroadcastTimestamp(nodeID string) error
	UpdateExecutionMetricsBroadcastTimestamp(nodeID string) error

	QueryMetrics(ctx context.Context, from time.Time, to *time.Time) (m metrics.Metrics, err error)
	BatchInsertSelfHealingChallengeEvents(ctx context.Context, event []types.SelfHealingChallengeEvent) error
	UpdateSHChallengeEventProcessed(challengeID string, isProcessed bool) error
	GetSelfHealingChallengeEvents() ([]types.SelfHealingChallengeEvent, error)
	CleanupSelfHealingChallenges() (err error)
	InsertSelfHealingGenerationMetrics(metrics types.SelfHealingGenerationMetric) error
	InsertSelfHealingExecutionMetrics(metrics types.SelfHealingExecutionMetric) error
	QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error)
	GetSelfHealingGenerationMetrics(timestamp time.Time) ([]types.SelfHealingGenerationMetric, error)
	GetSelfHealingExecutionMetrics(timestamp time.Time) ([]types.SelfHealingExecutionMetric, error)

	GetLastNSHChallenges(ctx context.Context, n int) (types.SelfHealingChallengeReports, error)
	GetSHChallengeReport(ctx context.Context, challengeID string) (types.SelfHealingChallengeReports, error)
}
