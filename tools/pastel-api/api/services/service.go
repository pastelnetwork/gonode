package services

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
)

var (
	ErrNotFoundMethod = errors.New("not found method")
)

type Service interface {
	Handle(ctx context.Context, method string, params []string) (interface{}, error)
}

type Services []Service

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
