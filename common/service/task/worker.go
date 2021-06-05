//go:generate mockery --name=Worker

package task

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/common/errgroup"
)

// Worker represents a pool of the task.
type Worker interface {

	// Tasks returns all tasks.
	Tasks() []Task

	// Task returns the task  by the given id.
	Task(taskID string) Task

	// AddTask adds the new task.
	AddTask(task Task)

	// RemoveTask removes the task.
	RemoveTask(subTask Task)

	// Run waits for new tasks, starts handling eche of them in a new goroutine.
	Run(ctx context.Context) error
}

// implement Worker interface
type worker struct {
	sync.Mutex

	tasks  []Task
	taskCh chan Task
}

// Tasks implements Worker.Tasks
func (worker *worker) Tasks() []Task {
	return worker.tasks
}

// Task implements Worker.Task
func (worker *worker) Task(taskID string) Task {
	worker.Lock()
	defer worker.Unlock()

	for _, task := range worker.tasks {
		if task.ID() == taskID {
			return task
		}
	}
	return nil
}

// AddTask implements Worker.AddTask
func (worker *worker) AddTask(task Task) {
	worker.Lock()
	defer worker.Unlock()

	worker.tasks = append(worker.tasks, task)
	worker.taskCh <- task
}

// RemoveTask implements Worker.RemoveTask
func (worker *worker) RemoveTask(subTask Task) {
	worker.Lock()
	defer worker.Unlock()

	for i, task := range worker.tasks {
		if task == subTask {
			worker.tasks = append(worker.tasks[:i], worker.tasks[i+1:]...)
			return
		}
	}
}

// Run implements Worker.Run
func (worker *worker) Run(ctx context.Context) error {
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
func NewWorker() Worker {
	return &worker{
		taskCh: make(chan Task),
	}
}
