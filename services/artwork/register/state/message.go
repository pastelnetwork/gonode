package state

import (
	"time"
)

// Message represents status of the registration to notify users.
type Message struct {
	CreatedAt time.Time
	Status    Status
	isFinal   bool
}

// NewMessage returns a new Message instance..
func NewMessage(status Status) *Message {
	return &Message{
		CreatedAt: time.Now(),
		Status:    status,
		isFinal:   status == StatusTaskCompleted || status == StatusTaskRejected,
	}
}
