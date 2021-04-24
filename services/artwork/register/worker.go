package register

import (
	"context"
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
	for {
		select {
		case <-ctx.Done():
			return nil
		case task := <-worker.taskCh:
			go task.Run(ctx)

			// go func() {
			// 	// NOTE: for testing
			// 	time.Sleep(time.Second)
			// 	task.State.Update(state.NewMessage(state.StatusTicketAccepted))
			// 	time.Sleep(time.Second)
			// 	task.State.Update(state.NewMessage(state.StatusTicketRegistered))
			// 	time.Sleep(time.Second)
			// 	task.State.Update(state.NewMessage(state.StatusTicketActivated))
			// 	time.Sleep(time.Second)
			// 	task.State.Update(state.NewMessage(state.StatusTaskCompleted))
			// }()
		}
	}
}

// NewWorker returns a new Worker instance.
func NewWorker() *Worker {
	return &Worker{
		taskCh: make(chan *Task),
	}
}
