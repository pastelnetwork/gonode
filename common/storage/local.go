package storage

import (
	"context"

	"github.com/pastelnetwork/gonode/common/types"
)

// LocalStoreInterface is interface for local sqlite store
type LocalStoreInterface interface {
	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)
	InsertStorageChallenge(challenge types.StorageChallenge) (hID int, err error)
	QueryStorageChallenges() (challenges []types.StorageChallenge, err error)
	CleanupStorageChallenges() (err error)
	InsertSelfHealingChallenge(challenge types.SelfHealingChallenge) (hID int, err error)
	QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error)
	CloseHistoryDB(ctx context.Context)
}
