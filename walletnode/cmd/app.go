package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"

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
	"github.com/pastelnetwork/gonode/walletnode/services/nftdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/nftregister"
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
		cli.NewFlag("log-level", &config.LogConfig.Level).SetUsage("Set the log `level`.").SetValue(config.LogConfig.Level),
		cli.NewFlag("log-file", &config.LogConfig.File).SetUsage("The log `file` to write to."),
		cli.NewFlag("quiet", &config.Quiet).SetUsage("Disallows log output to stdout.").SetAliases("q"),
		// API
		cli.NewFlag("swagger", &config.API.Swagger).SetUsage("Enable Swagger UI."),
	)

	app.SetActionFunc(func(ctx context.Context, args []string) error {
		ctx = log.ContextWithPrefix(ctx, "app")

		if configFile != "" {
			if err := configurer.ParseFile(configFile, config); err != nil {
				return fmt.Errorf("error parsing walletnode config file: %v", err)
			}
		}

		if pastelConfigFile != "" {
			if err := configurer.ParseFile(pastelConfigFile, config.Pastel); err != nil {
				return fmt.Errorf("error parsing pastel config: %v", err)
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

		if err := log.SetLevelName(config.LogConfig.Level); err != nil {
			return errors.Errorf("--log-level %q, %w", config.LogConfig.Level, err)
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(func() {
		cancel()
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
	})

	// entities
	pastelClient := pastel.NewClient(config.Pastel, config.Pastel.BurnAddress())

	// Business logic services
	// ----NftApiHandler Services----
	nodeClient := grpc.NewClient(pastelClient)

	db := memory.NewKeyValue()
	fileStorage := fs.NewFileStorage(config.TempDir)

	// raptorq client
	config.ArtworkRegister.RaptorQServiceAddress = fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	config.ArtworkRegister.RqFilesDir = config.RqFilesDir
	rqClient := rqgrpc.NewClient()

	if config.RegTxMinConfirmations > 0 {
		config.ArtworkRegister.NFTRegTxMinConfirmations = config.RegTxMinConfirmations
		config.SenseRegister.SenseRegTxMinConfirmations = config.RegTxMinConfirmations
	}

	if config.ActTxMinConfirmations > 0 {
		config.ArtworkRegister.NFTActTxMinConfirmations = config.ActTxMinConfirmations
		config.SenseRegister.SenseActTxMinConfirmations = config.ActTxMinConfirmations
	}

	// business logic services
	nftRegister := nftregister.NewService(&config.ArtworkRegister, pastelClient, nodeClient, fileStorage, db, rqClient)
	nftSearch := nftsearch.NewNftSearchService(&config.ArtworkSearch, pastelClient, nodeClient)
	nftDownload := nftdownload.NewNftFownloadService(&config.ArtworkDownload, pastelClient, nodeClient)

	// ----UserdataApiHandler Services----
	userdataNodeClient := grpc.NewClient(pastelClient)
	userdataProcess := userdataprocess.NewService(&config.UserdataProcess, pastelClient, userdataNodeClient)

	// SenseApiHandler services
	senseRegister := senseregister.NewService(&config.SenseRegister, pastelClient, nodeClient, fileStorage, db)

	// api service
	server := api.NewAPIServer(config.API,
		services.NewNftApiHandler(nftRegister, nftSearch, nftDownload),
		services.NewUserdataApiHandler(userdataProcess),
		services.NewSenseApiHandler(senseRegister),
		services.NewSwagger(),
	)

	log.WithContext(ctx).Infof("Config: %s", config)

	return runServices(ctx, server, nftRegister, nftSearch, nftDownload, userdataProcess, senseRegister)
}
