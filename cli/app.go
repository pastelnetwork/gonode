package cli

import (
	"context"
	"io/ioutil"

	"github.com/pastelnetwork/go-commons/cli"
	"github.com/pastelnetwork/go-commons/configer"
	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/go-commons/log/hooks"
	"github.com/pastelnetwork/go-commons/sys"
	"github.com/pastelnetwork/go-commons/version"
	"github.com/pastelnetwork/walletnode/config"
	"github.com/pastelnetwork/walletnode/pastel"
	"github.com/pastelnetwork/walletnode/rest"
)

const (
	appName  = "walletnode"
	appUsage = "WalletNode" // TODO: Write a clear description.

	defaultConfigFile = ""
)

// NewApp inits a new command line interface.
func NewApp() *cli.App {
	configFile := defaultConfigFile
	config := config.New()

	app := cli.NewApp(appName)
	app.SetUsage(appUsage)
	app.SetVersion(version.Version())

	app.AddFlags(
		// Main
		cli.NewFlag("config-file", &configFile).SetUsage("Set `path` to the config file.").SetValue(configFile).SetAliases("c"),
		cli.NewFlag("log-level", &config.LogLevel).SetUsage("Set the log `level`.").SetValue(config.LogLevel),
		cli.NewFlag("log-file", &config.LogFile).SetUsage("The log `file` to write to."),
		cli.NewFlag("quiet", &config.Quiet).SetUsage("Disallows log output to stdout.").SetAliases("q"),
		// Rest
		cli.NewFlag("swagger", &config.Rest.Swagger).SetUsage("Enable Swagger UI."),
	)

	app.SetActionFunc(func(args []string) error {
		if configFile != "" {
			if err := configer.ParseFile(configFile, config); err != nil {
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

		return run(config)
	})

	return app
}

func run(config *config.Config) error {
	log.Debug("[app] start")
	defer log.Debug("[app] end")

	log.Debugf("[app] config: %s", config)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.Info("[app] Interrupt signal received. Gracefully shutting down...")
	})

	if err := pastel.Init(config.Pastel); err != nil {
		return err
	}

	err := rest.New(config.Rest).Run(ctx)
	return err
}
