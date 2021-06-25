package task

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/common/errgroup"
)

// Worker represents a pool of the task.
type Worker struct {
	sync.Mutex

	tasks  []Task
	taskCh chan Task
}

// Tasks returns all tasks.
func (worker *Worker) Tasks() []Task {
	worker.Lock()
	defer worker.Unlock()

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
			group.Go(func() error {
				defer worker.RemoveTask(task)

				return task.Run(ctx)
			})
		}
	}
}

// NewWorker returns a new Worker instance.
func NewWorker() *Worker {
	return &Worker{
		taskCh: make(chan Task),
	}
}
