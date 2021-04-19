package state

// Subscription represents the state of the task consumed by subscribers each time the state changes.
type Subscription struct {
	msgCh  chan *Message
	doneCh chan struct{}
}

// Msg returns message channel.
func (sub *Subscription) Msg() chan *Message {
	return sub.msgCh
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
func (sub *Subscription) Pub(msgs ...*Message) {
	for _, msg := range msgs {
		select {
		case <-sub.Done():
			return
		case sub.msgCh <- msg:
			if msg.Latest {
				sub.Close()
				return
			}
		}
	}
}

// NewSubscription returns a new Subscription instance.
func NewSubscription() *Subscription {
	return &Subscription{
		msgCh:  make(chan *Message),
		doneCh: make(chan struct{}),
	}
}
