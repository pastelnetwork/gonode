package sys

import (
	"fmt"
	"os"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
)

// CheckErrorAndExit checks if there is an error, display it in the console and exit with a non-zero exit code. Otherwise, exit 0.
// Note that if the debugMode is true, this will print out the stack trace.
func CheckErrorAndExit(debugMode bool) func(err error) {
	return func(err error) {
		defer os.Exit(ExitCode(err))

		if err == nil || errors.IsContextCanceled(err) {
			return
		}

		errorFields := errors.ExtractFields(err)

		if debugMode {
			log.WithError(err).WithFields(map[string]interface{}(errorFields)).Error(errors.ErrorStack(err))
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %s %s\n", errorFields.String(), err)
		}
	}
}
