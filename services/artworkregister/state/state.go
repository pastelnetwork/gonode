package state

// State represents a states of the registering process.
type State struct {
	msgs []*Message
	subs []*Subscription
}

// All returns all states from the very beginning.
func (states *State) All() []*Message {
	return states.msgs
}

// Latest returns the message of the latest states.
func (states *State) Latest() *Message {
	if last := len(states.msgs); last > 0 {
		return states.msgs[last-1]
	}
	return nil
}

// Update updates the last states of the states by adding the message that contains properties of the current states.
func (states *State) Update(msg *Message) {
	for _, sub := range states.subs {
		go sub.Pub(msg)
	}
	states.msgs = append(states.msgs, msg)
}

// Subscribe returns a new subscription of the states.
func (states *State) Subscribe() (*Subscription, error) {
	sub := NewSubscription()
	states.subs = append(states.subs, sub)

	go func() {
		sub.Pub(states.msgs...)
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
func New(msg *Message) *State {
	return &State{
		msgs: []*Message{msg},
	}
}
