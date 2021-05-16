package artworkregister

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/sync/errgroup"
)

// Worker represents a task handler of registering artworks.
type Worker struct {
	sync.Mutex

	tasks  []*Task
	taskCh chan *Task
}

// Task returns the task of the registration artwork by the given id.
func (worker *Worker) Task(id string) *Task {
	worker.Lock()
	defer worker.Unlock()

	for _, task := range worker.tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// AddTask adds the new task.
func (worker *Worker) AddTask(ctx context.Context, task *Task) {
	worker.Lock()
	defer worker.Unlock()

	worker.tasks = append(worker.tasks, task)
	select {
	case <-ctx.Done():
	case worker.taskCh <- task:
	}
}

// RemoveTask removes the task.
func (worker *Worker) RemoveTask(subTask *Task) {
	worker.Lock()
	defer worker.Unlock()

	for i, task := range worker.tasks {
		if task == subTask {
			worker.tasks = append(worker.tasks[:i], worker.tasks[i+1:]...)
			return
		}
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
				defer worker.RemoveTask(task)

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
