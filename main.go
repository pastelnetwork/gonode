package main

import (
	"os"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/go-commons/sys"
	"github.com/pastelnetwork/supernode/cmd"
)

const (
	debugModeEnvName = "SUPERNODE_DEBUG"
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
