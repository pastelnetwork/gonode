package senseregister

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
)

const (
	logPrefix = "sense"
)

// Task represents a task for registering an artwork.
type Task struct {
	task.Task
	*Service

	Request *Request
}

func NewTask(service *Service, Ticket *Request) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Request: Ticket,
	}
}

// Run starts the task.
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debug("States updated")
	})

	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")

	// TODO: implement task logic

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}
