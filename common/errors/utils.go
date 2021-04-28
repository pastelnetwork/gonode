package errors

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pastelnetwork/gonode/common/log"
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

// CheckErrorAndExit checks if there is an error, display it in the console and exit with a non-zero exit code. Otherwise, exit 0.
// Note that if the debugMode is true, this will print out the stack trace.
func CheckErrorAndExit(err error) {
	defer os.Exit(ExitCode(err))

	if err == nil || IsContextCanceled(err) {
		return
	}

	errorFields := ExtractFields(err)

	if log.DebugMode() {
		log.WithError(err).WithFields(map[string]interface{}(errorFields)).Error(ErrorStack(err))
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: %s %s\n", errorFields.String(), err)
	}
}

// Log logs a statement of the given error. If log.DebugMode is true, the log record contains a stack of errors.
func Log(err error) {
	errorFields := ExtractFields(err)

	if log.DebugMode() {
		log.WithError(err).WithFields(map[string]interface{}(errorFields)).Error(ErrorStack(err))
	} else {
		log.WithFields(map[string]interface{}(errorFields)).Error(err)
	}
}
