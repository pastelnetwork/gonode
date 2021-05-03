package artworkregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/sync/errgroup"
)

// Worker represents a task handler of registering artworks.
type Worker struct {
	taskCh chan *Task
}

// AddTask adds the new task.
func (worker *Worker) AddTask(ctx context.Context, task *Task) {
	select {
	case <-ctx.Done():
	case worker.taskCh <- task:
	}
}

// Run waits for new tasks, starts handling eche of them in a new goroutine.
func (worker *Worker) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return group.Wait()

		case task := <-worker.taskCh:
			group.Go(func() (err error) {
				defer errors.Recover(func(recErr error) { err = recErr })
				return task.Run(ctx)
			})
		}
	}
}

// NewWorker returns a new Worker instance.
func NewWorker() *Worker {
	return &Worker{
		taskCh: make(chan *Task),
	}
}
