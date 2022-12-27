package storage

import (
	"context"

	"github.com/pastelnetwork/gonode/common/types"
)

// LocalStoreInterface is interface for local sqlite store
type LocalStoreInterface interface {
	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)
	InsertFailedStorageChallenge(challenge types.FailedStorageChallenge) (hID int, err error)
	UpdateFailedStorageChallenge(challenge types.FailedStorageChallenge) error
	QueryFailedStorageChallenges() (challenges []types.FailedStorageChallenge, err error)
	CloseHistoryDB(ctx context.Context)
}
