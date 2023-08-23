package storagechallenge

import "github.com/pastelnetwork/gonode/common/errors"

var (
	//UniqueConstraintErr is the error returned if the unique constraint violates
	UniqueConstraintErr = errors.New("UNIQUE constraint failed")
)
