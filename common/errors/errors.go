package errors

import (
	"fmt"
	"strings"
)

// Errors is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Errors []*Error

func (errs Errors) Error() string {
	if len(errs) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", errs[0])
	}

	points := make([]string, len(errs))
	for i, err := range errs {
		points[i] = fmt.Sprintf("* %s", err)

		if len(err.fields) > 0 {
			points[i] += fmt.Sprintf(" [ %s ]", err.fields)
		}
	}

	return fmt.Sprintf(
		"%d errors occurred:\n%s\n",
		len(errs), strings.Join(points, "\n"))
}
