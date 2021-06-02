//go:generate mockery --name=State

package test

import "github.com/pastelnetwork/gonode/common/service/task/state"

// State mock state.State. redefine state.State interface
// TODO : go:generate does not work with accros module on current CircleCI build
// not sure how to fix it
type State interface {
	// Status returns the current status.
	Status() *state.Status

	// SetStatusNotifyFunc sets a function to be called after the state is updated.
	SetStatusNotifyFunc(fn func(status *state.Status))

	// RequiredStatus returns an error if the current status doen't match the given one.
	RequiredStatus(subStatus state.SubStatus) error

	// StatusHistory returns all history from the very beginning.
	StatusHistory() []*state.Status

	// UpdateStatus updates the status of the state by creating a new status with the given `status`.
	UpdateStatus(subStatus state.SubStatus)

	// SubscribeStatus returns a new subscription of the state.
	SubscribeStatus() func() <-chan *state.Status
}
