package services

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
)

// API errors.
var (
	ErrNotFoundMethod = errors.New("not found method")
)

// Service represents a service for handling API request.
type Service interface {
	Handle(ctx context.Context, method string, params []string) (interface{}, error)
}

// Services represents muiltipe Service.
type Services []Service

// Handle looks for the first matching handler that takes input parameters that do not return an `ErrNotFoundMethod` error.
func (services Services) Handle(ctx context.Context, method string, params []string) (interface{}, error) {
	for _, service := range services {
		resp, err := service.Handle(ctx, method, params)
		if err == ErrNotFoundMethod {
			continue
		}
		return resp, err
	}
	return nil, ErrNotFoundMethod
}
