package cmd

import (
	"context"
	"fmt"
	"io"
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
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"
	"github.com/pastelnetwork/gonode/supernode/configs"
	"github.com/pastelnetwork/gonode/supernode/debug"
	healthcheck_lib "github.com/pastelnetwork/gonode/supernode/healthcheck"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/client"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/bridge"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/healthcheck"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/hermes"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	"github.com/pastelnetwork/gonode/supernode/services/download"
	"github.com/pastelnetwork/gonode/supernode/services/nftregister"
	"github.com/pastelnetwork/gonode/supernode/services/selfhealing"
	"github.com/pastelnetwork/gonode/supernode/services/senseregister"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
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
		cli.NewFlag("log-level", &config.LogConfig.Levels.CommonLogLevel).SetUsage("Set the log `level` of application services.").SetValue(config.LogConfig.Levels.CommonLogLevel),
		// cli.NewFlag("metadb-log-level", &config.LogConfig.Levels.MetaDBLogLevel).SetUsage("Set the log `level` for metadb.").SetValue(config.LogConfig.Levels.MetaDBLogLevel),
		cli.NewFlag("p2p-log-level", &config.LogConfig.Levels.P2PLogLevel).SetUsage("Set the log `level` for p2p.").SetValue(config.LogConfig.Levels.P2PLogLevel),
		cli.NewFlag("dd-log-level", &config.LogConfig.Levels.P2PLogLevel).SetUsage("Set the log `level` for dupedetection.").SetValue(config.LogConfig.Levels.DupeDetectionLogLevel),
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
			log.SetOutput(io.Discard)
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

		if err := config.P2P.Validate(); err != nil {
			return err
		}

		// if err := config.MetaDB.Validate(); err != nil {
		// 	return err
		// }

		if err := config.RaptorQ.Validate(); err != nil {
			return err
		}

		if err := log.SetLevelName(config.LogConfig.Levels.CommonLogLevel); err != nil {
			return errors.Errorf("--log-level %q, %w", config.LogConfig.Levels.CommonLogLevel, err)
		}

		if err := log.SetP2PLogLevelName(config.LogConfig.Levels.P2PLogLevel); err != nil {
			return errors.Errorf("--p2p-log-level %q, %w", config.LogConfig.Levels.P2PLogLevel, err)
		}

		if err := log.SetMetaDBLogLevelName(config.LogConfig.Levels.MetaDBLogLevel); err != nil {
			return errors.Errorf("--metadb-log-level %q, %w", config.LogConfig.Levels.MetaDBLogLevel, err)
		}

		if err := log.SetDDLogLevelName(config.LogConfig.Levels.DupeDetectionLogLevel); err != nil {
			return errors.Errorf("--dd-log-level %q, %w", config.LogConfig.Levels.DupeDetectionLogLevel, err)
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
	pastelClient := pastel.NewClient(config.Pastel, config.Pastel.BurnAddress())
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

	// Because of rqlite failures, we're going to disable metadb for now, this consequently disables user data processing until we
	//  either fix rqlite or develop a workaround.
	//	NB: Removed protobuf and grpc comms files as well to prevent malicious behavior.

	// new metadb service
	// config.MetaDB.SetWorkDir(config.WorkDir)
	// metadb := metadb.New(config.MetaDB, config.Node.PastelID, pastelClient)
	// database := database.NewDatabaseOps(metadb, config.UserDB)

	rqAddr := fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	// raptorq client
	config.NftRegister.RaptorQServiceAddress = rqAddr
	config.NftRegister.DDDatabase = filepath.Join(config.DdWorkDir, "support_files", "registered_image_fingerprints_db.sqlite")
	config.SenseRegister.DDDatabase = filepath.Join(config.DdWorkDir, "support_files", "registered_image_fingerprints_db.sqlite")
	config.NftRegister.RqFilesDir = config.RqFilesDir
	config.CascadeRegister.RaptorQServiceAddress = rqAddr
	config.CascadeRegister.RqFilesDir = config.RqFilesDir

	if config.NumberConnectedNodes > 0 {
		config.NftRegister.NumberConnectedNodes = config.NumberConnectedNodes
		config.CascadeRegister.NumberConnectedNodes = config.NumberConnectedNodes
	}

	if config.PreburntTxMinConfirmations > 0 {
		config.NftRegister.PreburntTxMinConfirmations = config.PreburntTxMinConfirmations
		config.CascadeRegister.PreburntTxMinConfirmations = config.PreburntTxMinConfirmations
	}

	rqClient := rqgrpc.NewClient()

	// prepare for dd client
	config.DDServer.SetWorkDir(config.WorkDir)
	ddClient := ddclient.NewDDServerClient(config.DDServer)
	if err := os.MkdirAll(config.DDServer.DDFilesDir, os.ModePerm); err != nil {
		return errors.Errorf("could not create dd-temp-file-dir %q, %w", config.DDServer.DDFilesDir, err)
	}

	// business logic services
	config.NftDownload.RaptorQServiceAddress = rqAddr
	config.NftDownload.RqFilesDir = config.RqFilesDir
	config.SelfHealingChallenge.RaptorQServiceAddress = rqAddr
	config.SelfHealingChallenge.RqFilesDir = config.RqFilesDir

	// business logic services
	nftRegister := nftregister.NewService(&config.NftRegister, fileStorage, pastelClient, nodeClient, p2p, rqClient, ddClient)
	nftDownload := download.NewService(&config.NftDownload, pastelClient, p2p, rqClient)
	senseRegister := senseregister.NewService(&config.SenseRegister, fileStorage, pastelClient, nodeClient, p2p, rqClient, ddClient)
	cascadeRegister := cascaderegister.NewService(&config.CascadeRegister, fileStorage, pastelClient, nodeClient, p2p, rqClient)
	storageChallenger := storagechallenge.NewService(&config.StorageChallenge, fileStorage, pastelClient, nodeClient, p2p, rqClient, nil)
	selfHealing := selfhealing.NewService(&config.SelfHealingChallenge, fileStorage, pastelClient, nodeClient, p2p, rqClient)
	// // ----Userdata Services----
	// userdataNodeClient := client.New(pastelClient, secInfo)
	// userdataProcess := userdataprocess.NewService(&config.UserdataProcess, pastelClient, userdataNodeClient, database)

	// create stats manager
	statsMngr := healthcheck_lib.NewStatsMngr(config.HealthCheck)
	statsMngr.Add("p2p", p2p)
	// statsMngr.Add("mdl", metadb)
	statsMngr.Add("pasteld", healthcheck_lib.NewPastelStatsClient(pastelClient))

	// Debug service
	debugSerivce := debug.NewService(config.DebugService, p2p, storageChallenger)

	// server
	grpc := server.New(config.Server,
		"service",
		pastelClient,
		secInfo,
		walletnode.NewRegisterNft(nftRegister),
		supernode.NewRegisterNft(nftRegister),
		walletnode.NewRegisterSense(senseRegister),
		supernode.NewRegisterSense(senseRegister),
		walletnode.NewRegisterCascade(cascadeRegister),
		supernode.NewRegisterCascade(cascadeRegister),
		walletnode.NewDownloadNft(nftDownload),
		bridge.NewDownloadData(nftDownload),
		supernode.NewStorageChallengeGRPC(storageChallenger),
		supernode.NewSelfHealingChallengeGRPC(selfHealing),
		healthcheck.NewHealthCheck(statsMngr),
		hermes.NewHermes(p2p),
	)

	log.WithContext(ctx).Infof("Config: %s", config)

	return runServices(ctx, grpc, p2p, nftRegister, nftDownload, senseRegister, cascadeRegister, statsMngr, debugSerivce, storageChallenger, selfHealing)
}
