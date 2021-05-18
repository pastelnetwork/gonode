package state

import (
	"time"
)

type SubStatus interface {
	String() string
	IsFinal() bool
}

// Status represents a state of the task.
type Status struct {
	CreatedAt time.Time
	SubStatus
}

// Is returns true if the current `Status` matches to the given `statuses`.
func (status *Status) Is(subStatus SubStatus) bool {
	return status.SubStatus == subStatus
}

// NewStatus returns a new Status instance.
func NewStatus(subStatus SubStatus) *Status {
	return &Status{
		CreatedAt: time.Now(),
		SubStatus: subStatus,
	}
}
