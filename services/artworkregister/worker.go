package artworkregister

import (
	"context"

	"github.com/pastelnetwork/go-commons/errors"
	"golang.org/x/sync/errgroup"
)

type task interface {
	Run(ctx context.Context) error
}

// Worker represents a task handler of registering artworks.
type Worker struct {
	taskCh chan task
}

// AddTask adds the new task.
func (worker *Worker) AddTask(ctx context.Context, task task) {
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
				defer errors.Recover(func(rec error) { err = rec })
				return task.Run(ctx)
			})
		}
	}
}

// NewWorker returns a new Worker instance.
func NewWorker() *Worker {
	return &Worker{
		taskCh: make(chan task),
	}
}
