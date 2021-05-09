package services

import (
	"context"
	"net/http"
)

// Service represents a service for handling API request.
type Service interface {
	Handle(ctx context.Context, r *http.Request, method string, params []string) (interface{}, error)
}

// Services represents muiltipe Service.
type Services []Service

// Handle looks for the first matching handler that takes input parameters that do not return an `ErrNotFoundMethod` error.
func (services Services) Handle(ctx context.Context, r *http.Request, method string, params []string) (interface{}, error) {
	for _, service := range services {
		resp, err := service.Handle(ctx, r, method, params)
		if err == ErrNotFoundMethod {
			continue
		}
		return resp, err
	}
	return nil, ErrNotFoundMethod
}
