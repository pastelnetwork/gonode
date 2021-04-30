package state

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/common/log"
)

// State represents a states of the registering process.
type State struct {
	statuses  []*Status
	updatedCh chan struct{}
	sync.Mutex
}

// All returns all states from the very beginning.
func (states *State) All() []*Status {
	states.Lock()
	defer states.Unlock()

	return states.statuses
}

// Updated waits for new state.
func (states *State) Updated() <-chan struct{} {
	return states.updatedCh
}

// Latest returns the message of the latest states.
func (states *State) Latest() *Status {
	states.Lock()
	defer states.Unlock()

	if last := len(states.statuses); last > 0 {
		return states.statuses[last-1]
	}
	return nil
}

// Update updates the last states of the states by adding the status that contains properties of the current states.
func (states *State) Update(ctx context.Context, status *Status) {
	states.Lock()
	defer states.Unlock()

	states.statuses = append(states.statuses, status)
	log.WithField("status", status.Type.String()).Debugf("[service] Status updated")

	close(states.updatedCh)
	states.updatedCh = make(chan struct{})
}

// New returns a new State instance.
func New(status *Status) *State {
	return &State{
		statuses:  []*Status{status},
		updatedCh: make(chan struct{}),
	}
}
