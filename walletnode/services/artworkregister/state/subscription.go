package state

// Subscription represents the state of the task consumed by subscribers each time the state changes.
type Subscription struct {
	statusCh chan *Status
	doneCh   chan struct{}
}

// Status returns status channel.
func (sub *Subscription) Status() chan *Status {
	return sub.statusCh
}

// Close closes the subsription.
func (sub *Subscription) Close() {
	select {
	case <-sub.Done():
		return
	default:
		close(sub.doneCh)
	}
}

// Done returns a channel when the subscription is closed.
func (sub *Subscription) Done() <-chan struct{} {
	return sub.doneCh
}

// Pub publishes messages to subscriber.
func (sub *Subscription) Pub(msgs ...*Status) {
	for _, msg := range msgs {
		select {
		case <-sub.Done():
			return
		case sub.statusCh <- msg:
			if msg.isFinal {
				sub.Close()
				return
			}
		}
	}
}

// NewSubscription returns a new Subscription instance.
func NewSubscription() *Subscription {
	return &Subscription{
		statusCh: make(chan *Status, len(statusNames)),
		doneCh:   make(chan struct{}),
	}
}
