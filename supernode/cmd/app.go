package cmd

import (
	"context"
	"io/ioutil"

	"github.com/pastelnetwork/gonode/pastel-client"
	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/supernode/configs"
	"github.com/pastelnetwork/gonode/supernode/server/grpc"
	"github.com/pastelnetwork/gonode/supernode/server/grpc/services"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"github.com/pastelnetwork/gonode/supernode/storage/memory"
)

const (
	appName  = "supernode"
	appUsage = "SuperNode" // TODO: Write a clear description.

	defaultConfigFile = ""
)

// NewApp inits a new command line interface.
func NewApp() *cli.App {
	configFile := defaultConfigFile
	config := configs.New()

	app := cli.NewApp(appName)
	app.SetUsage(appUsage)
	app.SetVersion(version.Version())

	app.AddFlags(
		// Main
		cli.NewFlag("config-file", &configFile).SetUsage("Set `path` to the config file.").SetValue(configFile).SetAliases("c"),
		cli.NewFlag("log-level", &config.LogLevel).SetUsage("Set the log `level`.").SetValue(config.LogLevel),
		cli.NewFlag("log-file", &config.LogFile).SetUsage("The log `file` to write to."),
		cli.NewFlag("quiet", &config.Quiet).SetUsage("Disallows log output to stdout.").SetAliases("q"),
	)

	app.SetActionFunc(func(ctx context.Context, args []string) error {
		if configFile != "" {
			if err := configurer.ParseFile(configFile, config); err != nil {
				return err
			}
		}

		if config.Quiet {
			log.SetOutput(ioutil.Discard)
		} else {
			log.SetOutput(app.Writer)
		}

		if config.LogFile != "" {
			fileHook := hooks.NewFileHook(config.LogFile)
			log.AddHook(fileHook)
		}

		if err := log.SetLevelName(config.LogLevel); err != nil {
			return errors.Errorf("--log-level %q, %s", config.LogLevel, err)
		}

		return runApp(ctx, config)
	})

	return app
}

func runApp(ctx context.Context, config *configs.Config) error {
	log.Debug("[app] start")
	defer log.Debug("[app] end")

	log.Debugf("[app] config: %s", config)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.Info("[app] Interrupt signal received. Gracefully shutting down...")
	})

	// entities
	pastel := pastel.NewClient(config.Pastel)
	db := memory.NewKeyValue()

	// business logic services
	artworkRegister := artworkregister.NewService(config.ArtworkRegister, db, pastel)

	// server
	grpc := grpc.NewServer(config.Server,
		services.NewArtwork(artworkRegister),
	)

	return runServices(ctx, artworkRegister, grpc)
}
