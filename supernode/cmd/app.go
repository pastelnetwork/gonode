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
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	"github.com/pastelnetwork/gonode/dupedetection/ddscan"
	"github.com/pastelnetwork/gonode/metadb"
	"github.com/pastelnetwork/gonode/metadb/database"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"
	"github.com/pastelnetwork/gonode/supernode/configs"
	"github.com/pastelnetwork/gonode/supernode/debug"
	healthcheck_lib "github.com/pastelnetwork/gonode/supernode/healthcheck"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/client"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/healthcheck"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/supernode"
	grpcstoragechallenge "github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/supernode/storagechallenge"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
)

const (
	appName  = "supernode"
	appUsage = "SuperNode" // TODO: Write a clear description.

	tfmodelDir = "./tfmodels" // relatively from work-dir
	rqFilesDir = "rqfiles"
	ddWorkDir  = "pastel_dupe_detection_service"
)

var (
	defaultPath = configurer.DefaultPath()
	homePath, _ = os.UserHomeDir()

	defaultTempDir          = filepath.Join(os.TempDir(), appName)
	defaultWorkDir          = filepath.Join(defaultPath, appName)
	defaultConfigFile       = filepath.Join(defaultPath, appName+".yml")
	defaultPastelConfigFile = filepath.Join(defaultPath, "pastel.conf")

	defaultRqFilesDir = filepath.Join(defaultPath, rqFilesDir)
	defaultDdWorkDir  = filepath.Join(homePath, ddWorkDir)
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
		cli.NewFlag("work-dir", &config.WorkDir).SetUsage("Set `path` for storing work data.").SetValue(defaultWorkDir),
		cli.NewFlag("temp-dir", &config.TempDir).SetUsage("Set `path` for storing temp data.").SetValue(defaultTempDir),
		cli.NewFlag("rq-files-dir", &config.RqFilesDir).SetUsage("Set `path` for storing files for rqservice.").SetValue(defaultRqFilesDir),
		cli.NewFlag("dd-service-dir", &config.DdWorkDir).SetUsage("Set `path` to the directory of dupe detection service - contains database file").SetValue(defaultDdWorkDir),
		cli.NewFlag("log-level", &config.LogLevel).SetUsage("Set the log `level`.").SetValue(config.LogLevel),
		cli.NewFlag("log-file", &config.LogFile).SetUsage("The log `file` to write to."),
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

		if config.LogFile != "" {
			log.AddHook(hooks.NewFileHook(config.LogFile))
		}
		log.AddHook(hooks.NewDurationHook())

		if err := config.P2P.Validate(); err != nil {
			return err
		}

		if err := config.MetaDB.Validate(); err != nil {
			return err
		}

		if err := config.RaptorQ.Validate(); err != nil {
			return err
		}

		if err := log.SetLevelName(config.LogLevel); err != nil {
			return errors.Errorf("--log-level %q, %w", config.LogLevel, err)
		}

		if err := os.MkdirAll(config.TempDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create temp-dir %q, %w", config.TempDir, err)
		}

		if err := os.MkdirAll(config.WorkDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create work-dir %q, %w", config.WorkDir, err)
		}

		if err := os.MkdirAll(config.RqFilesDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create rq-files-dir %q, %w", config.RqFilesDir, err)
		}

		if err := os.MkdirAll(config.DdWorkDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create dd-service-dir %q, %w", config.DdWorkDir, err)
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
	pastelClient := pastel.NewClient(config.Pastel)
	secInfo := &alts.SecInfo{
		PastelID:   config.PastelID,
		PassPhrase: config.PassPhrase,
		Algorithm:  "ed448",
	}
	nodeClient := client.New(pastelClient, secInfo)
	fileStorage := fs.NewFileStorage(config.TempDir)

	// p2p service (currently using kademlia)
	config.P2P.SetWorkDir(config.WorkDir)
	config.P2P.ID = config.PastelID
	p2p := p2p.New(config.P2P, pastelClient, secInfo)

	// new metadb service
	config.MetaDB.SetWorkDir(config.WorkDir)
	metadb := metadb.New(config.MetaDB, config.Node.PastelID, pastelClient)
	database := database.NewDatabaseOps(metadb, config.UserDB)

	// raptorq client
	config.ArtworkRegister.RaptorQServiceAddress = fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	config.ArtworkRegister.RqFilesDir = config.RqFilesDir

	if config.NumberConnectedNodes > 0 {
		config.ArtworkRegister.NumberConnectedNodes = config.NumberConnectedNodes
	}

	if config.PreburntTxMinConfirmations > 0 {
		config.ArtworkRegister.PreburntTxMinConfirmations = config.PreburntTxMinConfirmations
	}

	rqClient := rqgrpc.NewClient()

	// prepare for dd client
	config.DDServer.SetWorkDir(config.WorkDir)
	ddClient := ddclient.NewDDServerClient(config.DDServer)
	if err := os.MkdirAll(config.DDServer.DDFilesDir, os.ModePerm); err != nil {
		return errors.Errorf("could not create dd-temp-file-dir %q, %w", config.DDServer.DDFilesDir, err)
	}

	// business logic services
	config.ArtworkDownload.RaptorQServiceAddress = fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	config.ArtworkDownload.RqFilesDir = config.RqFilesDir

	// business logic services
	artworkRegister := artworkregister.NewService(&config.ArtworkRegister, fileStorage, pastelClient, nodeClient, p2p, rqClient, ddClient)
	artworkDownload := artworkdownload.NewService(&config.ArtworkDownload, pastelClient, p2p, rqClient)

	ddScanConfig := ddscan.NewConfig()
	ddScanConfig.SetWorkDir(config.DdWorkDir)

	ddScan, err := ddscan.NewService(ddScanConfig, pastelClient, p2p)
	if err != nil {
		return errors.Errorf("can not start dupe detection service, err: %w", err)
	}

	// ----Userdata Services----
	secClient := client.New(pastelClient, secInfo)
	userdataProcess := userdataprocess.NewService(&config.UserdataProcess, pastelClient, secClient, database)

	// ----Storage Challenges----
	stService := storagechallenge.NewService(nil, secClient, p2p)

	// create stats manager
	statsMngr := healthcheck_lib.NewStatsMngr(config.HealthCheck)
	statsMngr.Add("p2p", p2p)
	statsMngr.Add("mdl", metadb)
	statsMngr.Add("pasteld", healthcheck_lib.NewPastelStatsClient(pastelClient))
	statsMngr.Add("ddscan", ddScan)

	// Debug service
	debugSerivce := debug.NewService(config.DebugService, p2p)

	// server
	grpc := server.New(config.Server,
		"service",
		pastelClient,
		secInfo,
		walletnode.NewRegisterArtwork(artworkRegister),
		supernode.NewRegisterArtwork(artworkRegister),
		walletnode.NewDownloadArtwork(artworkDownload),
		walletnode.NewProcessUserdata(userdataProcess, database),
		supernode.NewProcessUserdata(userdataProcess, database),
		healthcheck.NewHealthCheck(statsMngr),
		grpcstoragechallenge.NewStorageChallenge(stService),
	)

	log.WithContext(ctx).Infof("Config: %s", config)

	return runServices(ctx, metadb, grpc, p2p, artworkRegister, artworkDownload, ddScan, database, userdataProcess, statsMngr, debugSerivce)
}
