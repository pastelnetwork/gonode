package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pastelnetwork/gonode/bridge/node/grpc"
	"github.com/pastelnetwork/gonode/bridge/services/bridge"
	"github.com/pastelnetwork/gonode/bridge/services/download"
	"github.com/pastelnetwork/gonode/bridge/services/server"
	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/bridge/configs"
	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
)

const (
	appName  = "bridge"
	appUsage = "create & maintain connections with SNs"
)

var (
	defaultPath = configurer.DefaultPath()

	//defaultTempDir    = filepath.Join(os.TempDir(), appName)
	defaultConfigFile       = filepath.Join(defaultPath, appName+".yml")
	defaultPastelConfigFile = filepath.Join(defaultPath, "pastel.conf")
)

// NewApp configures our app by parsing command line flags, config files, and setting up logging and temporary directories
func NewApp() *cli.App {
	var configFile string
	var pastelConfigFile string

	config := configs.New()

	app := cli.NewApp(appName)
	app.SetUsage(appUsage)
	app.SetVersion(version.Version())

	app.AddFlags(
		// Main
		cli.NewFlag("config-file", &configFile).SetUsage("Set `path` to the config file.").SetValue(defaultConfigFile).SetAliases("c"),
		cli.NewFlag("pastel-config-file", &pastelConfigFile).SetUsage("Set `path` to the pastel config file.").SetValue(defaultPastelConfigFile),
	)

	//SetActionFunc is the default, and in our case only, executed action.
	//Sets up configs and logging, and also returns the "app" function to main.go to be called (Run) there.
	app.SetActionFunc(func(ctx context.Context, args []string) error {
		//Sets logging prefix to pastel-app
		ctx = log.ContextWithPrefix(ctx, "pastel-app")

		//Parse config files supplied in arguments
		if configFile != "" {
			if err := configurer.ParseFile(configFile, config); err != nil {
				return fmt.Errorf("error parsing bridge config file: %v", err)
			}
		}

		if pastelConfigFile != "" {
			if err := configurer.ParseFile(pastelConfigFile, config.Pastel); err != nil {
				return fmt.Errorf("error parsing pastel config: %v", err)
			}
		}
		//Write to stdout by default, no output if quiet is set in config
		if config.Quiet {
			log.SetOutput(ioutil.Discard)
		} else {
			log.SetOutput(app.Writer)
		}

		//Configure log rotation IFF a file is specified for output
		if config.LogConfig.File != "" {
			rotateHook := hooks.NewFileHook(config.LogConfig.File)
			rotateHook.SetMaxSizeInMB(config.LogConfig.MaxSizeInMB)
			rotateHook.SetMaxAgeInDays(config.LogConfig.MaxAgeInDays)
			rotateHook.SetMaxBackups(config.LogConfig.MaxBackups)
			rotateHook.SetCompress(config.LogConfig.Compress)
			log.AddHook(rotateHook)
		}
		//By default, log time since start of app
		log.AddHook(hooks.NewDurationHook())

		//Parse other flags like log level, temp directory, and raptorq files directory
		if err := log.SetLevelName(config.LogConfig.Level); err != nil {
			return errors.Errorf("--log-level %q, %w", config.LogConfig.Level, err)
		}

		//This function, defined below, is what will be run when "Run" is called on our app
		return runApp(ctx, config)
	})

	//while "app" is returned, the default action func we just set as a parameter to it is going to be called when run
	return app
}

func runApp(ctx context.Context, config *configs.Config) error {
	//Create App Context
	log.WithContext(ctx).Info("Start")
	defer log.WithContext(ctx).Info("End")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//Register system interrupts for shutdown
	sys.RegisterInterruptHandler(func() {
		cancel()
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
	})

	// pastelClient reads in the json-formatted hostname, port, username, and password and
	//  connects over gRPC to cNode for access to Blockchain, Masternodes, Tickets, and PastelID databases
	pastelClient := pastel.NewClient(config.Pastel, config.Pastel.BurnAddress())
	// wrap pastelClient in security functionality then allow for calling of nfts, userdata, and sense functions
	nodeClient := grpc.NewSNClient(pastelClient)

	log.WithContext(ctx).Infof("Config: %s", config)

	downloadService := download.NewService(config.Download, pastelClient, nodeClient)
	grpc := server.New(server.NewConfig(),
		"service",
		bridge.NewService(downloadService),
	)

	return runServices(ctx, downloadService, grpc)
}
