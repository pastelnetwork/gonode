package main

import (
	"fmt"
	"os"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/go-commons/sys"
	"github.com/pastelnetwork/supernode/cli"
)

const (
	debugModeEnvName = "SUPERNODE_DEBUG"
)

var (
	debugMode = sys.GetBoolEnv(debugModeEnvName, false)
)

func main() {
	defer errors.Recover(checkErrorAndExit)

	app := cli.NewApp()
	err := app.Run(os.Args)

	checkErrorAndExit(err)
}

func checkErrorAndExit(err error) {
	defer os.Exit(sys.ExitCode(err))

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

func init() {
	log.SetDebugMode(debugMode)
}
