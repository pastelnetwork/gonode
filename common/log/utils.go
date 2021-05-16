package log

import (
	"fmt"
	"os"

	"github.com/pastelnetwork/gonode/common/errors"
)

// CheckErrorAndExit checks if there is an error, display it in the console and exit with a non-zero exit code. Otherwise, exit 0.
// Note that if the debugMode is true, this will print out the stack trace.
func CheckErrorAndExit(err error) {
	if err == nil || errors.IsContextCanceled(err) {
		return
	}
	defer os.Exit(errors.ExitCode(err))

	if debugMode {
		WithErrorStack(err).Fatal("Exit")
	}

	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
}
