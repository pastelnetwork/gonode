package storage

import (
	"context"

	"github.com/pastelnetwork/gonode/common/types"
)

// LocalStoreInterface is interface for local sqlite store
type LocalStoreInterface interface {
	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)
	InsertStorageChallengeMessage(challenge types.StorageChallengeLogMessage) error
	InsertBroadcastMessage(challenge types.BroadcastLogMessage) error
	QueryStorageChallengeMessage(challengeID string, messageType int) (challenge types.StorageChallengeLogMessage, err error)
	CleanupStorageChallenges() (err error)
	CleanupSelfHealingChallenges() (err error)
	InsertSelfHealingChallenge(challenge types.SelfHealingChallenge) (hID int, err error)
	QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error)
	UpsertPingHistory(pingInfo types.PingInfo) error
	GetPingInfoBySupernodeID(supernodeID string) (*types.PingInfo, error)
	CloseHistoryDB(ctx context.Context)
}
