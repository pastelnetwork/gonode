package state

import (
	"time"
)

// Status represents state of the artwork registration.
type Status struct {
	CreatedAt time.Time
	Type      StatusType
}

// NewStatus returns a new Status instance.
func NewStatus(statusType StatusType) *Status {
	return &Status{
		CreatedAt: time.Now(),
		Type:      statusType,
	}
}
