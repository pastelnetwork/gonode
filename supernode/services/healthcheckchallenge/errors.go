package healthcheckchallenge

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

var (
	//ErrUniqueConstraint is the error returned if the unique constraint violates
	ErrUniqueConstraint = errors.New("UNIQUE constraint failed")
)

func logError(ctx context.Context, method string, err error) {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
		return
	}

	log.WithContext(ctx).WithField("method", method).WithError(err).Debug(err.Error())
}
