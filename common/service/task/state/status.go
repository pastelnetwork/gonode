//go:generate mockery --name=SubStatus

package state

import (
	"fmt"
	"time"
)

// SubStatus represents a sub-status that contains a description of the status.
type SubStatus interface {
	fmt.Stringer
	IsFinal() bool
	IsFailure() bool
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
		CreatedAt: time.Now().UTC(),
		SubStatus: subStatus,
	}
}
