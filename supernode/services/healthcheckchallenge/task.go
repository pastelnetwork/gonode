package healthcheckchallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "healthcheck-challenge-task"
)

// HCTask : Healthcheck challenge task will manage response to healthcheck challenge requests
type HCTask struct {
	*common.SuperNodeTask
	*HCService

	storage        *common.StorageHandler
	challengeState SaveChallengeState
}

// Run : RunHelper's cleanup function is currently nil as WIP will determine what needs to be cleaned.
func (task *HCTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.RemoveArtifacts)
}

// RemoveArtifacts : Cleanup function defined here, can be filled in later
func (task *HCTask) RemoveArtifacts() {
}

// NewHCTask returns a new Task instance.
func NewHCTask(service *HCService) *HCTask {
	task := &HCTask{
		SuperNodeTask: common.NewSuperNodeTask(logPrefix, service.historyDB),
		HCService:     service,
	}
	return task
}

//utils below

// SaveChallengeState represents the callbacks for each step in the healthcheck challenge process.
type SaveChallengeState interface {
	OnSent(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnResponded(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnSucceeded(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnFailed(ctx context.Context, challengeID, nodeID string, sentBlock int32)
	OnTimeout(ctx context.Context, challengeID, nodeID string, sentBlock int32)
}

// SaveChallengeMessageState can be called to perform the above functions based on the state of the healthcheck challenge.
func (task *HCTask) SaveChallengeMessageState(ctx context.Context, status string, challengeID, nodeID string, sentBlock int32) {
	switch status {
	case "sent":
		task.challengeState.OnSent(ctx, challengeID, nodeID, sentBlock)
	case "respond":
		task.challengeState.OnResponded(ctx, challengeID, nodeID, sentBlock)
	case "succeeded":
		task.challengeState.OnSucceeded(ctx, challengeID, nodeID, sentBlock)
	case "failed":
		task.challengeState.OnFailed(ctx, challengeID, nodeID, sentBlock)
	case "timeout":
		task.challengeState.OnTimeout(ctx, challengeID, nodeID, sentBlock)
	}
}
