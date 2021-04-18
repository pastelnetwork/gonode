package state

type Subscription struct {
	msgCh  chan *Message
	doneCh chan struct{}
}

func (sub *Subscription) Msg() chan *Message {
	return sub.msgCh
}

func (sub *Subscription) Close() {
	select {
	case <-sub.Done():
		return
	default:
		close(sub.doneCh)
	}
}

func (sub *Subscription) Done() <-chan struct{} {
	return sub.doneCh
}

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

func NewSubscription() *Subscription {
	return &Subscription{
		msgCh:  make(chan *Message),
		doneCh: make(chan struct{}),
	}
}
