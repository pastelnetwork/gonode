package main

import (
	"os"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/tools/fake-pastel-api/cmd"
)

const (
	debugModeEnvName = "FAKE_PASTEL_API"
)

var (
	debugMode = sys.GetBoolEnv(debugModeEnvName, false)
)

func main() {
	defer errors.Recover(errors.CheckErrorAndExit)

	app := cmd.NewApp()
	err := app.Run(os.Args)

	errors.CheckErrorAndExit(err)
}

func init() {
	log.SetDebugMode(debugMode)
}
