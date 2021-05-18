package state

import (
	"sync"
)

// State represents a state of the registering process.
type State struct {
	*Event

	actionFn func(event *Event)
	sync.RWMutex
	prevs  []*Event
	subsCh []chan *Event
}

// SetActionFunc sets a function that runs after a state update.
func (state *State) SetActionFunc(fn func(event *Event)) {
	state.actionFn = fn
}

// Events returns all prevs from the very beginning.
func (state *State) Events() []*Event {
	state.RLock()
	defer state.RUnlock()

	return append(state.prevs, state.Event)
}

// Update updates the event of the state by creating a new event with the given `status`.
func (state *State) Update(status Status) {
	state.Lock()
	defer state.Unlock()

	event := NewEvent(status)
	state.prevs = append(state.prevs, state.Event)
	state.Event = event

	if state.actionFn != nil {
		state.actionFn(event)
	}

	for _, subCh := range state.subsCh {
		go func() {
			subCh <- event
		}()
	}
}

// Subscribe returns a new subscription of the state.
func (state *State) Subscribe() func() <-chan *Event {
	state.RLock()
	defer state.RUnlock()

	subCh := make(chan *Event)
	state.subsCh = append(state.subsCh, subCh)

	for _, event := range append(state.prevs, state.Event) {
		event := event
		go func() {
			subCh <- event
		}()
	}

	sub := func() <-chan *Event {
		return subCh
	}
	return sub
}

// New returns a new State instance.
func New(status Status) *State {
	return &State{
		Event: NewEvent(status),
	}
}
