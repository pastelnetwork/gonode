package selfhealing

import (
	"context"

	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "self-healing-task"
)

// SHTask : Self healing task will manage response to self healing challenges requests
type SHTask struct {
	*common.SuperNodeTask
	*SHService
}

// Run : RunHelper's cleanup function is currently nil as WIP will determine what needs to be cleaned.
func (task *SHTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.RemoveArtifacts)
}

// RemoveArtifacts : Cleanup function defined here, can be filled in later
func (task *SHTask) RemoveArtifacts() {
}

// NewSHTask returns a new Task instance.
func NewSHTask(service *SHService) *SHTask {
	task := &SHTask{
		SuperNodeTask: common.NewSuperNodeTask(logPrefix),
		SHService:     service,
	}

	return task
}
