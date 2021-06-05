package errors

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ErrorStack converts the given error to a string, including the stack trace if available
func ErrorStack(err error) string {
	var errStacks []string

	for {
		switch err := err.(type) {
		case interface{ ErrorStack() string }:
			errStacks = append(errStacks, err.ErrorStack())
		default:
			errStacks = append(errStacks, err.Error())
		}

		if err = errors.Unwrap(err); err == nil {
			break
		}
	}

	return strings.Join(errStacks, "\n")
}

// IsContextCanceled returns `true` if error has occurred by event `context.Canceled` which is not really an error
func IsContextCanceled(err error) bool {
	return errors.Is(err, context.Canceled)
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Recover is a method that tries to recover from panics, and if it succeeds, calls the given onPanic function with an error that
// explains the cause of the panic. This function should only be called from a defer statement.
func Recover(onPanic func(cause error)) {
	if rec := recover(); rec != nil {
		err, isError := rec.(error)
		if !isError {
			err = fmt.Errorf("%v", rec)
		}
		onPanic(newWithSkip(err, 3))
	}
}
