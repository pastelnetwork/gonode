package storagechallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
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
		stateStorage:            &defaultChallengeStateLogging{},
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

type defaultChallengeStateLogging struct{}

func (cs defaultChallengeStateLogging) OnSent(ctx context.Context, challengeID, nodeID string, sentBlock int32) {
	log.WithContext(ctx).WithPrefix(logPrefix).WithField("challengeID", challengeID).WithField("nodeID", nodeID).WithField("sentBlock", sentBlock).Println("Storage Challenge Sent")
}
func (cs defaultChallengeStateLogging) OnResponded(ctx context.Context, challengeID, nodeID string, sentBlock int32) {
	log.WithContext(ctx).WithPrefix(logPrefix).WithField("challengeID", challengeID).WithField("nodeID", nodeID).WithField("sentBlock", sentBlock).Println("Storage Challenge Responded")
}
func (cs defaultChallengeStateLogging) OnSucceeded(ctx context.Context, challengeID, nodeID string, sentBlock int32) {
	log.WithContext(ctx).WithPrefix(logPrefix).WithField("challengeID", challengeID).WithField("nodeID", nodeID).WithField("sentBlock", sentBlock).Println("Storage Challenge Succeeded")
}
func (cs defaultChallengeStateLogging) OnFailed(ctx context.Context, challengeID, nodeID string, sentBlock int32) {
	log.WithContext(ctx).WithPrefix(logPrefix).WithField("challengeID", challengeID).WithField("nodeID", nodeID).WithField("sentBlock", sentBlock).Println("Storage Challenge Failed")
}
func (cs defaultChallengeStateLogging) OnTimeout(ctx context.Context, challengeID, nodeID string, sentBlock int32) {
	log.WithContext(ctx).WithPrefix(logPrefix).WithField("challengeID", challengeID).WithField("nodeID", nodeID).WithField("sentBlock", sentBlock).Println("Storage Challenge Timed Out")
}

//SaveChallengeMessageState should be a function
func (task *StorageChallengeTask) SaveChallengeMessageState(ctx context.Context, status string, challengeID, nodeID string, sentBlock int32) {
	switch status {
	case "sent":
		task.stateStorage.OnSent(ctx, challengeID, nodeID, sentBlock)
	case "respond":
		task.stateStorage.OnResponded(ctx, challengeID, nodeID, sentBlock)
	case "succeeded":
		task.stateStorage.OnSucceeded(ctx, challengeID, nodeID, sentBlock)
	case "failed":
		task.stateStorage.OnFailed(ctx, challengeID, nodeID, sentBlock)
	case "timeout":
		task.stateStorage.OnTimeout(ctx, challengeID, nodeID, sentBlock)
	}
}
