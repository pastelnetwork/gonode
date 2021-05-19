package state

import (
	"sync"
)

// State represents a state of the registering process.
type State struct {
	*Status

	actionFn func(status *Status)
	sync.RWMutex
	prevs  []*Status
	subsCh []chan *Status
}

// SetActionFunc sets a function that runs after a state update.
func (state *State) SetActionFunc(fn func(status *Status)) {
	state.actionFn = fn
}

// Events returns all prevs from the very beginning.
func (state *State) Events() []*Status {
	state.RLock()
	defer state.RUnlock()

	return append(state.prevs, state.Status)
}

// Update updates the status of the state by creating a new status with the given `status`.
func (state *State) Update(subStatus SubStatus) {
	state.Lock()
	defer state.Unlock()

	status := NewStatus(subStatus)
	state.prevs = append(state.prevs, state.Status)
	state.Status = status

	if state.actionFn != nil {
		state.actionFn(status)
	}

	for _, subCh := range state.subsCh {
		subCh := subCh
		go func() {
			subCh <- status
		}()
	}
}

// Subscribe returns a new subscription of the state.
func (state *State) Subscribe() func() <-chan *Status {
	state.RLock()
	defer state.RUnlock()

	subCh := make(chan *Status)
	state.subsCh = append(state.subsCh, subCh)

	for _, status := range append(state.prevs, state.Status) {
		status := status
		go func() {
			subCh <- status
		}()
	}

	sub := func() <-chan *Status {
		return subCh
	}
	return sub
}

// New returns a new State instance.
func New(subStatus SubStatus) *State {
	return &State{
		Status: NewStatus(subStatus),
	}
}
