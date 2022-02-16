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
	*common.DupeDetectionHandler
	*StorageChallengeService

	storage *common.StorageHandler
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
		storage: common.NewStorageHandler(service.P2PClient, service.RQClient,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir),
	}
	return task
}
