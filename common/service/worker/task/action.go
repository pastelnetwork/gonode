package task

import "context"

type ActionFn func(ctx context.Context) error

// Action represents the action of the task.
type Action struct {
	fn     ActionFn
	doneCh chan struct{}
}

// NewAction returns a new Action instance.
func NewAction(fn ActionFn) *Action {
	return &Action{
		fn:     fn,
		doneCh: make(chan struct{}),
	}
}
