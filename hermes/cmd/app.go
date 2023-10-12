package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/hermes/config"
	"github.com/pastelnetwork/gonode/hermes/debug"
	"github.com/pastelnetwork/gonode/hermes/service/hermes"
	"github.com/pastelnetwork/gonode/hermes/service/node/grpc"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	appName   = "hermes"
	workDir   = "supernode"
	ddWorkDir = "pastel_dupe_detection_service"
	appUsage  = "SuperNode" // TODO: Write a clear description.
)

var (
	defaultPath = configurer.DefaultPath()
	homePath, _ = os.UserHomeDir()

	defaultWorkDir          = filepath.Join(defaultPath, workDir)
	defaultConfigFile       = filepath.Join(defaultPath, appName+".yml")
	defaultPastelConfigFile = filepath.Join(defaultPath, "pastel.conf")
	defaultDdWorkDir        = filepath.Join(homePath, ddWorkDir)
)

// NewApp inits a new command line interface.
func NewApp() *cli.App {
	var configFile string
	var pastelConfigFile string

	config := config.NewConfig()

	app := cli.NewApp(appName)
	app.SetUsage(appUsage)
	app.SetVersion(version.Version())

	app.AddFlags(
		// Main
		cli.NewFlag("config-file", &configFile).SetUsage("Set `path` to the config file.").SetValue(defaultConfigFile).SetAliases("c"),
		cli.NewFlag("pastel-config-file", &pastelConfigFile).SetUsage("Set `path` to the pastel config file.").SetValue(defaultPastelConfigFile),
		cli.NewFlag("work-dir", &config.WorkDir).SetUsage("Set `path` for storing work data.").SetValue(defaultWorkDir),
		cli.NewFlag("dd-service-dir", &config.DdWorkDir).SetUsage("Set `path` to the directory of dupe detection service - contains database file").SetValue(defaultDdWorkDir),
		cli.NewFlag("log-level", &config.LogConfig.Levels.CommonLogLevel).SetUsage("Set the log `level` of application services.").SetValue(config.LogConfig.Levels.CommonLogLevel),
		cli.NewFlag("log-file", &config.LogConfig.File).SetUsage("The log `file` to write to."),
		cli.NewFlag("quiet", &config.Quiet).SetUsage("Disallows log output to stdout.").SetAliases("q"),
	)

	app.SetActionFunc(func(ctx context.Context, args []string) error {
		ctx = log.ContextWithPrefix(ctx, "app")

		if configFile != "" {
			if err := configurer.ParseFile(configFile, config); err != nil {
				return fmt.Errorf("error parsing supernode config file: %v", err)
			}
		}
		if pastelConfigFile != "" {
			if err := configurer.ParseFile(pastelConfigFile, config.Pastel); err != nil {
				return fmt.Errorf("error parsing pastel config file: %v", err)
			}
		}

		if config.Quiet {
			log.SetOutput(ioutil.Discard)
		} else {
			log.SetOutput(app.Writer)
		}

		if config.LogConfig.File != "" {
			rotateHook := hooks.NewFileHook(config.LogConfig.File)
			rotateHook.SetMaxSizeInMB(config.LogConfig.MaxSizeInMB)
			rotateHook.SetMaxAgeInDays(config.LogConfig.MaxAgeInDays)
			rotateHook.SetMaxBackups(config.LogConfig.MaxBackups)
			rotateHook.SetCompress(config.LogConfig.Compress)
			log.AddHook(rotateHook)
		}
		log.AddHook(hooks.NewDurationHook())

		if err := log.SetLevelName(config.LogConfig.Levels.CommonLogLevel); err != nil {
			return errors.Errorf("--log-level %q, %w", config.LogConfig.Levels.CommonLogLevel, err)
		}

		if err := os.MkdirAll(config.WorkDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create work-dir %q, %w", config.WorkDir, err)
		}

		if err := os.MkdirAll(config.DdWorkDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create dd-service-dir %q, %w", config.DdWorkDir, err)
		}

		return runApp(ctx, config)
	})

	return app
}

func runApp(ctx context.Context, conf *config.Config) error {
	log.WithContext(ctx).Info("Start")
	defer log.WithContext(ctx).Info("End")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(func() {
		cancel()
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
	})

	// entities
	pastelClient := pastel.NewClient(conf.Pastel, conf.Pastel.BurnAddress())

	if conf.PassPhrase == "" {
		return errors.New("passphrase is empty, please provide passphrase in config file")
	}

	// Try to get PastelID from cnode API MasterNodeStatus
	extKey, err := mixins.GetPastelIDfromMNConfig(ctx, pastelClient, conf.PastelID)
	if err != nil {
		return fmt.Errorf("get pastelID from mn config: %w", err)
	}
	conf.PastelID = extKey

	// Validate PastelID and passphrase
	if !mixins.ValidateUser(ctx, pastelClient, conf.PastelID, conf.PassPhrase) {
		return errors.New("invalid pastelID or passphrase")
	}

	secInfo := &alts.SecInfo{
		PastelID:   conf.PastelID,
		PassPhrase: conf.PassPhrase,
		Algorithm:  "ed448",
	}

	nodeClient := grpc.NewSNClient(pastelClient, secInfo)

	hermesConfig := hermes.NewConfig()
	hermesConfig.SNHost = conf.SNHost
	hermesConfig.SNPort = conf.SNPort
	hermesConfig.CreatorPastelID = conf.PastelID
	hermesConfig.CreatorPastelIDPassphrase = conf.PassPhrase
	hermesConfig.SetWorkDir(conf.DdWorkDir)

	service, err := hermes.NewService(hermesConfig, pastelClient, nodeClient)
	if err != nil {
		return fmt.Errorf("start hermes: %w", err)
	}

	// Debug service
	debugSerivce := debug.NewService(conf.DebugService, service)

	log.WithContext(ctx).Infof("Config: %s", conf)

	return runServices(ctx, service, debugSerivce)
}
