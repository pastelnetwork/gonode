//go:generate mockery --name=SvcInterface

package service

import (
	"context"
)

// SvcInterface ...
type SvcInterface interface {
	// Run starts task
	Run(ctx context.Context) error

	// Stats returns current status of service
	Stats(ctx context.Context) (map[string]interface{}, error)
}
