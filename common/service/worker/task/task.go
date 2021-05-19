package task

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/service/state"
	"golang.org/x/sync/errgroup"
)

// Task is the task of registering new artwork.
type Task struct {
	id    string
	State *state.State

	actionCh chan *Action

	doneMu sync.Mutex
	doneCh chan struct{}
}

// ID returns id of the task.
func (task *Task) ID() string {
	return task.id
}

// Cancel tells a task to abandon its work.
// Cancel may be called by multiple goroutines simultaneously.
// After the first call, subsequent calls to a Cancel do nothing.
func (task *Task) Cancel() {
	task.doneMu.Lock()
	defer task.doneMu.Unlock()

	select {
	case <-task.Done():
		return
	default:
		close(task.doneCh)
	}
}

// Done returns a channel when the task is canceled.
func (task *Task) Done() <-chan struct{} {
	return task.doneCh
}

// RunAction waits for new actions, starts handling eche of them in a new goroutine.
func (task *Task) RunAction(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	for {
		select {
		case <-ctx.Done():
		case <-task.Done():
			cancel()
		case action := <-task.actionCh:
			group.Go(func() (err error) {
				defer errors.Recover(func(recErr error) { err = recErr })
				defer close(action.doneCh)

				return action.fn(ctx)
			})
			continue
		}
		break
	}

	return group.Wait()
}

// NewAction creates a new action and passes for the execution.
func (task *Task) NewAction(fn ActionFn) <-chan struct{} {
	act := NewAction(fn)
	task.actionCh <- act
	return act.doneCh
}

// RequiredStatus returns an error if the current status doen't match the given one.
func (task *Task) RequiredStatus(status state.SubStatus) error {
	if task.State.Is(status) {
		return nil
	}
	return errors.Errorf("required status %q, current %q", status, task.State.Status)
}

// New returns a new Task instance.
func New(status state.SubStatus) *Task {
	taskID, _ := random.String(8, random.Base62Chars)

	return &Task{
		id:       taskID,
		State:    state.New(status),
		doneCh:   make(chan struct{}),
		actionCh: make(chan *Action),
	}
}
