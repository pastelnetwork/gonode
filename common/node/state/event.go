package state

import (
	"time"
)

// Event represents event of the registration to notify users.
type Event struct {
	Status
	CreatedAt time.Time
}

// Is returns true if the current `Status` matches at least one of the given `statuses`.
func (event *Event) Is(statuses ...Status) bool {
	for _, status := range statuses {
		if event.Status == status {
			return true
		}
	}
	return false
}

// NewEvent returns a new Event instance.
func NewEvent(status Status) *Event {
	return &Event{
		CreatedAt: time.Now(),
		Status:    status,
	}
}
