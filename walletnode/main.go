package main

import (
	"os"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/walletnode/cmd"
)

const (
	debugModeEnvName = "WALLETNODE_DEBUG"
)

var (
	debugMode = sys.GetBoolEnv(debugModeEnvName, false)
)

func main() {
	defer errors.Recover(log.CheckErrorAndExit)

	app := cmd.NewApp()
	err := app.Run(os.Args)

	log.CheckErrorAndExit(err)
}

func init() {
	log.SetDebugMode(debugMode)
}
