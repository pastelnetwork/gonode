package storage

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/types"
)

// Metrics is the struct for metrics
type Metrics struct {
	SCMetrics          []byte
	SHTriggerMetrics   []byte
	SHExecutionMetrics []byte
}

// LocalStoreInterface is interface for local sqlite store
type LocalStoreInterface interface {
	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)
	InsertStorageChallengeMessage(challenge types.StorageChallengeLogMessage) error
	InsertBroadcastMessage(challenge types.BroadcastLogMessage) error
	QueryStorageChallengeMessage(challengeID string, messageType int) (challenge types.StorageChallengeLogMessage, err error)
	CleanupStorageChallenges() (err error)
	CleanupSelfHealingChallenges() (err error)
	InsertSelfHealingChallenge(challenge types.SelfHealingLogMessage) error
	InsertSelfHealingMetrics(metrics types.SelfHealingMetrics) error
	QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error)
	UpsertPingHistory(pingInfo types.PingInfo) error
	GetPingInfoBySupernodeID(supernodeID string) (*types.PingInfo, error)
	GetAllPingInfos() (types.PingInfos, error)
	GetWatchlistPingInfo() ([]types.PingInfo, error)
	UpdatePingInfo(supernodeID string, isOnWatchlist, isAdjusted bool) error
	CloseHistoryDB(ctx context.Context)
	QueryMetrics(from time.Time, to *time.Time) (m Metrics, err error)
}
