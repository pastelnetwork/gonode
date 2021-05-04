package services

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
)

var errWrongParams = errors.New("wrong command")

type Service interface {
	Handle(ctx context.Context, params []string) (interface{}, error)
}
