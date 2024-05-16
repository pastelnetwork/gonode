package storagechallenge

import (
	"context"

	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "storage-challenge-task"
)

// Storage challenge tasking is used to match other service patterns.  Largely just a wrapper for
//  the individual generation, processing, and verification functions. Does also include
//  the callbacks for what happens on different storage challenge states.  It is therefore
//  possible that further development could link the three functions and perform operations
//  based on state.

// SCTask : Storage challenge task will manage response to storage challenge requests
type SCTask struct {
	*common.SuperNodeTask
	*SCService

	storage *common.StorageHandler
}

// Run : RunHelper's cleanup function is currently nil as WIP will determine what needs to be cleaned.
func (task *SCTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.RemoveArtifacts)
}

// RemoveArtifacts : Cleanup function defined here, can be filled in later
func (task *SCTask) RemoveArtifacts() {
}

// NewSCTask returns a new Task instance.
func NewSCTask(service *SCService) *SCTask {
	task := &SCTask{
		SuperNodeTask: common.NewSuperNodeTask(logPrefix, service.historyDB),
		SCService:     service,
	}
	return task
}
