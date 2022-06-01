package main


import (
	"os"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/bridge/cmd"
)

const (
	debugModeEnvName = "BRIDGE_DEBUG"
)

var (
	debugMode = sys.GetBoolEnv(debugModeEnvName, false)
)

func main() {
	defer errors.Recover(log.FatalAndExit)

	//configuration in here, app.Run will actually run NewApp's runApp function.
	app := cmd.NewApp()
	err := app.Run(os.Args)

	log.FatalAndExit(err)
}

func init() {
	log.SetDebugMode(debugMode)
}
