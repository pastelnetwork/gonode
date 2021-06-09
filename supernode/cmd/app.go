package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/metadb"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/probe"
	"github.com/pastelnetwork/gonode/probe/tfmodel"
	"github.com/pastelnetwork/gonode/supernode/configs"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/client"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
)

const (
	appName  = "supernode"
	appUsage = "SuperNode" // TODO: Write a clear description.

	tfmodelDir = "./tfmodels" // relatively from work-dir
	p2pDir     = "./p2p"      // relatively from work-dir
)

var (
	defaultPath = configurer.DefaultPath()

	defaultTempDir          = filepath.Join(os.TempDir(), appName)
	defaultWorkDir          = filepath.Join(defaultPath, appName)
	defaultConfigFile       = filepath.Join(defaultPath, appName+".yml")
	defaultPastelConfigFile = filepath.Join(defaultPath, "pastel.conf")
)

// NewApp inits a new command line interface.
func NewApp() *cli.App {
	var configFile string
	var pastelConfigFile string

	config := configs.New()

	app := cli.NewApp(appName)
	app.SetUsage(appUsage)
	app.SetVersion(version.Version())

	app.AddFlags(
		// Main
		cli.NewFlag("config-file", &configFile).SetUsage("Set `path` to the config file.").SetDefaultText(defaultConfigFile).SetAliases("c"),
		cli.NewFlag("pastel-config-file", &pastelConfigFile).SetUsage("Set `path` to the pastel config file.").SetDefaultText(defaultPastelConfigFile),
		cli.NewFlag("work-dir", &config.WorkDir).SetUsage("Set `path` for storing work data.").SetValue(defaultWorkDir),
		cli.NewFlag("temp-dir", &config.TempDir).SetUsage("Set `path` for storing temp data.").SetValue(defaultTempDir),
		cli.NewFlag("log-level", &config.LogLevel).SetUsage("Set the log `level`.").SetValue(config.LogLevel),
		cli.NewFlag("log-file", &config.LogFile).SetUsage("The log `file` to write to."),
		cli.NewFlag("quiet", &config.Quiet).SetUsage("Disallows log output to stdout.").SetAliases("q"),
	)

	app.SetActionFunc(func(ctx context.Context, args []string) error {
		ctx = log.ContextWithPrefix(ctx, "app")

		if configFile != "" {
			if err := configurer.ParseFile(configFile, config); err != nil {
				return err
			}
		}
		if pastelConfigFile != "" {
			if err := configurer.ParseFile(pastelConfigFile, config.Pastel.ExternalConfig); err != nil {
				log.WithContext(ctx).Debug(err)
			}
		}

		if config.Quiet {
			log.SetOutput(ioutil.Discard)
		} else {
			log.SetOutput(app.Writer)
		}

		if config.LogFile != "" {
			log.AddHook(hooks.NewFileHook(config.LogFile))
		}
		log.AddHook(hooks.NewDurationHook())

		if err := log.SetLevelName(config.LogLevel); err != nil {
			return errors.Errorf("--log-level %q, %w", config.LogLevel, err)
		}

		if err := os.MkdirAll(config.TempDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create temp-dir %q, %w", config.TempDir, err)
		}

		if err := os.MkdirAll(config.WorkDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create work-dir %q, %w", config.WorkDir, err)
		}

		return runApp(ctx, config)
	})

	return app
}

func runApp(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Info("Start")
	defer log.WithContext(ctx).Info("End")

	log.WithContext(ctx).Infof("Config: %s", config)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
	})

	// entities
	pastelClient := pastel.NewClient(config.Pastel)
	nodeClient := client.New()
	fileStorage := fs.NewFileStorage(config.TempDir)

	// analysis tools
	probeTensor := probe.NewTensor(filepath.Join(config.WorkDir, tfmodelDir), tfmodel.AllConfigs)

	// p2p service (currently using kademlia)
	config.P2P.SetDefaultDataDir(filepath.Join(config.WorkDir, p2pDir))
	p2p := p2p.New(config.P2P)

	// business logic services
	artworkRegister := artworkregister.NewService(&config.ArtworkRegister, fileStorage, probeTensor, pastelClient, nodeClient, p2p)

	// server
	grpc := server.New(config.Server,
		walletnode.NewRegisterArtwork(artworkRegister),
		supernode.NewRegisterArtwork(artworkRegister),
	)

	// new metadb service
	metadb := metadb.New(config.MetaDB, config.Node.PastelID)

	return runServices(ctx, metadb, grpc, p2p, artworkRegister)
}
