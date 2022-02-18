package storagechallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "storage-challenge-task"
)

// Storage challenge task will manage response to storage challenge requests
type StorageChallengeTask struct {
	*common.SuperNodeTask
	*StorageChallengeService

	storage      *common.StorageHandler
	stateStorage SaveChallengeState
}

//	RunHelper's cleanup function is currently nil as WIP will determine what needs to be cleaned.
func (task *StorageChallengeTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, nil)
}

// Task returns the task of the Storage Challenge by the id
func (service *StorageChallengeService) Task(id string) *StorageChallengeTask {
	return service.Worker.Task(id).(*StorageChallengeTask)
}

// NewStorageChallengeTask returns a new Task instance.
func NewStorageChallengeTask(service *StorageChallengeService) *StorageChallengeTask {
	task := &StorageChallengeTask{
		SuperNodeTask:           common.NewSuperNodeTask(logPrefix),
		StorageChallengeService: service,
	}
	return task
}

//utils below

type SaveChallengeState interface {
	OnSent(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnResponded(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnSucceeded(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnFailed(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnTimeout(ctx context.Context, challengeID, nodeID string, sentBlock int32)
}

func (task *StorageChallengeTask) SaveChallengMessageState(ctx context.Context, status string, challengeID, nodeID string, sentBlock int32) {
	var caller func(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	switch status {
	case "sent":
		caller = task.stateStorage.OnSent
	case "respond":
		caller = task.stateStorage.OnResponded
	case "succeeded":
		caller = task.stateStorage.OnSucceeded
	case "failed":
		caller = task.stateStorage.OnFailed
	case "timeout":
		caller = task.stateStorage.OnTimeout
	}

	caller(ctx, challengeID, nodeID, sentBlock)
}
