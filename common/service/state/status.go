package state

import (
	"time"
)

type subStatus interface {
	String() string
	IsFinal() bool
}

// Status represents a state of the task.
type Status struct {
	CreatedAt time.Time
	subStatus
}

// Is returns true if the current `Status` matches to the given `statuses`.
func (status *Status) Is(subStatus subStatus) bool {
	return subStatus == status.subStatus
}

// NewStatus returns a new Status instance.
func NewStatus(subStatus subStatus) *Status {
	return &Status{
		CreatedAt: time.Now(),
		subStatus: subStatus,
	}
}
