package cmd

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/sqlite"
	"github.com/pastelnetwork/gonode/supernode/services/metamigrator"
	"io"
	"os"
	"path/filepath"

	"net/http"
	_ "net/http/pprof" //profiling

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/storage/rqstore"
	"github.com/pastelnetwork/gonode/common/storage/scorestore"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/cloud.go"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/configs"
	"github.com/pastelnetwork/gonode/supernode/debug"
	healthcheck_lib "github.com/pastelnetwork/gonode/supernode/healthcheck"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/client"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/healthcheck"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/hermes"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	"github.com/pastelnetwork/gonode/supernode/services/collectionregister"
	"github.com/pastelnetwork/gonode/supernode/services/download"
	"github.com/pastelnetwork/gonode/supernode/services/healthcheckchallenge"
	"github.com/pastelnetwork/gonode/supernode/services/nftregister"
	"github.com/pastelnetwork/gonode/supernode/services/selfhealing"
	"github.com/pastelnetwork/gonode/supernode/services/senseregister"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
)

const (
	appName       = "supernode"
	profilingPort = "8848"
	appUsage      = "SuperNode" // TODO: Write a clear description.

	tfmodelDir = "./tfmodels" // relatively from work-dir
	rqFilesDir = "rqfiles"
	ddWorkDir  = "pastel_dupe_detection_service"
	rqDB       = "rqstore.db"
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

	app.SetActionFunc(func(ctx context.Context, _ []string) error {
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

	pastelClient := pastel.NewClient(config.Pastel, config.Pastel.BurnAddress())

	if config.PassPhrase == "" {
		return errors.New("passphrase is empty, please provide passphrase in config file")
	}

	// Try to get PastelID from cnode API MasterNodeStatus
	extKey, err := mixins.GetPastelIDfromMNConfig(ctx, pastelClient, config.PastelID)
	if err != nil {
		return fmt.Errorf("get pastelID from mn config: %w", err)
	}
	config.OverridePastelIDAndPass(extKey, config.PassPhrase)

	// Validate PastelID and passphrase
	if !mixins.ValidateUser(ctx, pastelClient, config.PastelID, config.PassPhrase) {
		return errors.New("invalid pastelID or passphrase")
	}

	secInfo := &alts.SecInfo{
		PastelID:   config.PastelID,
		PassPhrase: config.PassPhrase,
		Algorithm:  "ed448",
	}
	nodeClient := client.New(pastelClient, secInfo)
	fileStorage := fs.NewFileStorage(config.TempDir)

	rqstore, err := rqstore.NewSQLiteRQStore(filepath.Join(defaultPath, rqDB))
	if err != nil {
		return errors.Errorf("could not create rqstore, %w", err)
	}
	defer rqstore.Close()

	// p2p service (currently using kademlia)
	config.P2P.SetWorkDir(config.WorkDir)
	config.P2P.ID = config.PastelID

	var cloudStorage cloud.Storage

	if config.RcloneStorageConfig.BucketName != "" && config.RcloneStorageConfig.SpecName != "" {
		cloudStorage = cloud.NewRcloneStorage(config.RcloneStorageConfig.BucketName, config.RcloneStorageConfig.SpecName)
		if err := cloudStorage.CheckCloudConnection(); err != nil {
			log.WithContext(ctx).WithError(err).Fatal("error establishing connection with the cloud")
			return fmt.Errorf("rclone connection check failed: %w", err)
		}
	}

	p2p, err := p2p.New(ctx, config.P2P, pastelClient, secInfo, rqstore, cloudStorage)
	if err != nil {
		return errors.Errorf("could not create p2p service, %w", err)
	}

	metaMigratorStore, err := sqlite.NewMigrationMetaStore(ctx, config.P2P.DataDir, cloudStorage)
	if err != nil {
		return errors.Errorf("could not create p2p service, %w", err)
	}

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
		config.CollectionRegister.NumberConnectedNodes = config.NumberConnectedNodes
	}

	if config.PreburntTxMinConfirmations > 0 {
		config.NftRegister.PreburntTxMinConfirmations = config.PreburntTxMinConfirmations
		config.CascadeRegister.PreburntTxMinConfirmations = config.PreburntTxMinConfirmations
		config.CollectionRegister.PreburntTxMinConfirmations = config.PreburntTxMinConfirmations
	}

	// prepare for dd client
	config.DDServer.SetWorkDir(config.WorkDir)
	ddClient := ddclient.NewDDServerClient(config.DDServer)
	if err := os.MkdirAll(config.DDServer.DDFilesDir, os.ModePerm); err != nil {
		return errors.Errorf("could not create dd-temp-file-dir %q, %w", config.DDServer.DDFilesDir, err)
	}

	ip, _ := utils.GetExternalIPAddress()
	// business logic services
	config.NftDownload.RaptorQServiceAddress = rqAddr
	config.NftDownload.RqFilesDir = config.RqFilesDir
	config.SelfHealingChallenge.RaptorQServiceAddress = rqAddr
	config.SelfHealingChallenge.RqFilesDir = config.RqFilesDir
	config.SelfHealingChallenge.NodeAddress = ip

	//Initialize History DB
	hDB, err := queries.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error connecting history db..")
	}
	defer hDB.CloseHistoryDB(ctx)

	//Initialize History DB
	sDB, err := scorestore.OpenScoringDb()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error connecting history db..")
	}
	defer sDB.CloseDB(ctx)

	// business logic services

	nftRegister := nftregister.NewService(&config.NftRegister, fileStorage, pastelClient, nodeClient, p2p, ddClient, hDB, rqstore)
	nftDownload := download.NewService(&config.NftDownload, pastelClient, p2p, hDB)
	senseRegister := senseregister.NewService(&config.SenseRegister, fileStorage, pastelClient, nodeClient, p2p, ddClient, hDB)
	cascadeRegister := cascaderegister.NewService(&config.CascadeRegister, fileStorage, pastelClient, nodeClient, p2p, hDB, rqstore)
	collectionRegister := collectionregister.NewService(&config.CollectionRegister, fileStorage, pastelClient, nodeClient, p2p, hDB)
	storageChallenger := storagechallenge.NewService(&config.StorageChallenge, fileStorage, pastelClient, nodeClient, p2p, hDB, sDB)
	healthCheckChallenger := healthcheckchallenge.NewService(&config.HealthCheckChallenge, fileStorage, pastelClient, nodeClient, p2p, nil, hDB, sDB)
	selfHealing := selfhealing.NewService(&config.SelfHealingChallenge, fileStorage, pastelClient, nodeClient, p2p, hDB, nftDownload, rqstore)
	metaMigratorWorker := metamigrator.NewService(metaMigratorStore)
	// // ----Userdata Services----
	// userdataNodeClient := client.New(pastelClient, secInfo)
	// userdataProcess := userdataprocess.NewService(&config.UserdataProcess, pastelClient, userdataNodeClient, database)

	// create stats manager
	statsMngr := healthcheck_lib.NewStatsMngr(config.HealthCheck)
	statsMngr.Add("p2p", p2p)
	// statsMngr.Add("mdl", metadb)
	statsMngr.Add("pasteld", healthcheck_lib.NewPastelStatsClient(pastelClient))

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
		walletnode.NewRegisterCollection(collectionRegister),
		supernode.NewRegisterCollection(collectionRegister),
		walletnode.NewDownloadNft(nftDownload),
		supernode.NewStorageChallengeGRPC(storageChallenger),
		supernode.NewHealthCheckChallengeGRPC(healthCheckChallenger),
		supernode.NewSelfHealingChallengeGRPC(selfHealing),
		healthcheck.NewHealthCheck(statsMngr),
		hermes.NewHermes(p2p),
	)

	//openConsMap := grpc.GetConnTrackerMap()

	// Debug service
	debugSerivce := debug.NewService(config.DebugService, p2p, storageChallenger, nil)

	log.WithContext(ctx).Infof("Config: %s", config)

	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%s", profilingPort), nil)
	}()

	return runServices(ctx, grpc, p2p, nftRegister, nftDownload, senseRegister, cascadeRegister, statsMngr, debugSerivce, storageChallenger, selfHealing, collectionRegister, healthCheckChallenger, metaMigratorWorker)
}
