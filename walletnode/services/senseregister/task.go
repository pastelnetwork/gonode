package senseregister

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
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

	if err := task.run(ctx); err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Error("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

func (task *Task) run(ctx context.Context) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	return nil
}
