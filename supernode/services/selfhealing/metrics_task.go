package selfhealing

import (
	"context"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// MetricTask : Self healing metric task will manage response broadcasting metrics request
type MetricTask struct {
	*common.SuperNodeTask
	*MetricService
}

// Run : RunHelper's cleanup function is currently nil as WIP will determine what needs to be cleaned.
func (task *MetricTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.RemoveArtifacts)
}

// RemoveArtifacts : Cleanup function defined here, can be filled in later
func (task *MetricTask) RemoveArtifacts() {
}
