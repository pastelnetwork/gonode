package errors

import (
	"errors"
	"os/exec"
	"syscall"
)

// The codes are returned by the application after completion
const (
	ExitCodeSuccess int = iota
	ExitCodeError
)

// ExitCode returns the exit code of a command. If the error contains no exit codes, a success code is returned.
func ExitCode(err error) int {
	if err, ok := err.(interface{ ExitStatus() int }); ok {
		if code := err.ExitStatus(); code != ExitCodeSuccess {
			return code
		}
	}

	if err, ok := errors.Unwrap(err).(*exec.ExitError); ok {
		status := err.Sys().(syscall.WaitStatus)
		return status.ExitStatus()
	}

	if err != nil {
		return ExitCodeError
	}
	return ExitCodeSuccess
}
