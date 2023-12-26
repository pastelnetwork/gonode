package selfhealing

import "github.com/pastelnetwork/gonode/common/errors"

var (
	//ErrUniqueConstraint is the error returned if the unique constraint violates
	ErrUniqueConstraint = errors.New("UNIQUE constraint failed")
)
