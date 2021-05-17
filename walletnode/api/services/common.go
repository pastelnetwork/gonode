package services

import "context"

// Common represents common service.
type Common struct{}

// Run starts serving for operations that must be performed throughout the entire operation of the service.
func (service *Common) Run(_ context.Context) error {
	return nil
}

// NewCommon returns a new Common instance.
func NewCommon() *Common {
	return &Common{}
}
