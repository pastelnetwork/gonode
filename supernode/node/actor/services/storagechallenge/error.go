package storagechallenge

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	reasonInvalidEmptyValue = "invalid empty value"
)

type keys []string

func (ks keys) String() string {
	var ii []string = make([]string, len(ks))
	for i, k := range ks {
		ii[i] = strings.ReplaceAll(k, " ", "_")
	}
	return strings.Join(ii, " ")
}

type validationError struct {
	keys   keys
	reason string
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%s: %s", e.keys.String(), e.reason)
}

type validationErrorStack []*validationError

func joinKeyPart(keys ...string) string {
	return strings.Join(keys, ".")
}

func validationErrorStackWrap(ves validationErrorStack) error {
	if len(ves) == 0 {
		return nil
	}

	fieldViolations := []*errdetails.BadRequest_FieldViolation{}

	for _, ve := range ves {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{Field: ve.keys.String(), Description: ve.reason})
	}

	br := errdetails.BadRequest{
		FieldViolations: fieldViolations,
	}
	st := status.New(codes.InvalidArgument, "ValidationError")
	st1, err := st.WithDetails(&br)
	if err != nil {
		log.Printf("cannot compose status detail: %v", err)
		return st.Err()
	}
	return st1.Err()
}
