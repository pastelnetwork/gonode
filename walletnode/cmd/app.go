package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/pastel"
	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/services"
	"github.com/pastelnetwork/gonode/walletnode/configs"
	"github.com/pastelnetwork/gonode/walletnode/node/grpc"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"
)

const (
	appName    = "walletnode"
	appUsage   = "WalletNode" // TODO: Write a clear description.
	rqFilesDir = "rqfiles"
)

var (
	defaultPath = configurer.DefaultPath()

	defaultTempDir          = filepath.Join(os.TempDir(), appName)
	defaultConfigFile       = filepath.Join(defaultPath, appName+".yml")
	defaultPastelConfigFile = filepath.Join(defaultPath, "pastel.conf")
	defaultRqFilesDir       = filepath.Join(defaultPath, rqFilesDir)
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
		cli.NewFlag("config-file", &configFile).SetUsage("Set `path` to the config file.").SetValue(defaultConfigFile).SetAliases("c"),
		cli.NewFlag("pastel-config-file", &pastelConfigFile).SetUsage("Set `path` to the pastel config file.").SetValue(defaultPastelConfigFile),
		cli.NewFlag("temp-dir", &config.TempDir).SetUsage("Set `path` for storing temp data.").SetValue(defaultTempDir),
		cli.NewFlag("rq-files-dir", &config.RqFilesDir).SetUsage("Set `path` for storing files for rqservice.").SetValue(defaultRqFilesDir),
		cli.NewFlag("log-level", &config.LogLevel).SetUsage("Set the log `level`.").SetValue(config.LogLevel),
		cli.NewFlag("log-file", &config.LogFile).SetUsage("The log `file` to write to."),
		cli.NewFlag("quiet", &config.Quiet).SetUsage("Disallows log output to stdout.").SetAliases("q"),
		// API
		cli.NewFlag("swagger", &config.API.Swagger).SetUsage("Enable Swagger UI."),
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

		if err := os.MkdirAll(config.RqFilesDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create rq-files-dir %q, %w", config.RqFilesDir, err)
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

	sys.RegisterInterruptHandler(func() {
		cancel()
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
	})

	// entities
	pastelClient := pastel.NewClient(config.Pastel)

	// Business logic services
	// ----Artwork Services----
	nodeClient := grpc.NewClient()
	// p2p service (currently using kademlia)
	config.P2P.SetWorkDir(config.WorkDir)
	p2p := p2p.New(config.P2P)

	db := memory.NewKeyValue()
	fileStorage := fs.NewFileStorage(config.TempDir)

	// raptorq client
	config.ArtworkRegister.RaptorQServiceAddress = fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	config.ArtworkRegister.RqFilesDir = config.RqFilesDir
	rqClient := rqgrpc.NewClient()

	// burn address
	// TODO: this should be hardcoded with format like Ptxxxxxxx
	config.ArtworkRegister.BurnAddress = config.BurnAddress
	config.ArtworkRegister.RegArtTxMinConfirmations = config.RegArtTxMinConfirmations
	config.ArtworkRegister.RegArtTxTimeout = time.Duration(config.RegArtTxTimeout * int(time.Minute))
	config.ArtworkRegister.RegActTxMinConfirmations = config.RegActTxMinConfirmations
	config.ArtworkRegister.RegActTxTimeout = time.Duration(config.RegActTxTimeout * int(time.Minute))

	// business logic services
	artworkRegister := artworkregister.NewService(&config.ArtworkRegister, db, fileStorage, pastelClient, nodeClient, rqClient)
	artworkSearch := artworksearch.NewService(&config.ArtworkSearch, pastelClient, p2p, nodeClient)
	artworkDownload := artworkdownload.NewService(&config.ArtworkDownload, pastelClient, nodeClient)

	// ----Userdata Services----
	userdataNodeClient := grpc.NewClient()
	userdataProcess := userdataprocess.NewService(&config.UserdataProcess, pastelClient, userdataNodeClient)

	// api service
	server := api.NewServer(config.API,
		services.NewArtwork(artworkRegister, artworkSearch, artworkDownload),
		services.NewUserdata(userdataProcess),
		services.NewSwagger(),
	)

	return runServices(ctx, server, artworkRegister, artworkSearch, artworkDownload, userdataProcess)
}
