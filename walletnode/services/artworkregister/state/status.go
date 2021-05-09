package state

import (
	"time"
)

// Status represents status of the registration to notify users.
type Status struct {
	CreatedAt time.Time
	Type      StatusType
	isFinal   bool
}

// NewMessage returns a new Status instance..
func NewMessage(statusType StatusType) *Status {
	return &Status{
		CreatedAt: time.Now(),
		Type:      statusType,
		isFinal:   statusType == StatusTaskCompleted || statusType == StatusTaskRejected,
	}
}
