package worker

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/sync/errgroup"
)

// Task represent a worker task.
type Task interface {
	ID() string
	Run(ctx context.Context) error
}

// Worker represents a task handler.
type Worker struct {
	sync.Mutex

	tasks  []Task
	taskCh chan Task
}

// Tasks returns all tasks.
func (worker *Worker) Tasks() []Task {
	return worker.tasks
}

// Task returns the task  by the given id.
func (worker *Worker) Task(taskID string) Task {
	worker.Lock()
	defer worker.Unlock()

	for _, task := range worker.tasks {
		if task.ID() == taskID {
			return task
		}
	}
	return nil
}

// AddTask adds the new task.
func (worker *Worker) AddTask(task Task) {
	worker.Lock()
	defer worker.Unlock()

	worker.tasks = append(worker.tasks, task)
	worker.taskCh <- task
}

// RemoveTask removes the task.
func (worker *Worker) RemoveTask(subTask Task) {
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

// New returns a new Worker instance.
func New() *Worker {
	return &Worker{
		taskCh: make(chan Task),
	}
}
