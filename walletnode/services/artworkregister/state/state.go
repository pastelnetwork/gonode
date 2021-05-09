package state

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
)

// State represents a states of the registering process.
type State struct {
	statuses []*Status
	subs     []*Subscription
}

// All returns all states from the very beginning.
func (states *State) All() []*Status {
	return states.statuses
}

// Latest returns the message of the latest states.
func (states *State) Latest() *Status {
	if last := len(states.statuses); last > 0 {
		return states.statuses[last-1]
	}
	return nil
}

// Update updates the last states of the states by adding the message that contains properties of the current states.
func (states *State) Update(ctx context.Context, status *Status) {
	for _, sub := range states.subs {
		go sub.Pub(status)
	}
	log.WithContext(ctx).WithField("statuses", status.Type.String()).Debugf("State updated")
	states.statuses = append(states.statuses, status)
}

// Subscribe returns a new subscription of the states.
func (states *State) Subscribe() (*Subscription, error) {
	sub := NewSubscription()
	states.subs = append(states.subs, sub)

	go func() {
		sub.Pub(states.statuses...)
		<-sub.Done()
		states.unsubscribe(sub)
	}()

	return sub, nil
}

func (states *State) unsubscribe(elem *Subscription) {
	for i, sub := range states.subs {
		if sub == elem {
			states.subs = append(states.subs[:i], states.subs[i+1:]...)
			return
		}
	}
}

// New returns a new State instance.
func New(status *Status) *State {
	return &State{
		statuses: []*Status{status},
	}
}
