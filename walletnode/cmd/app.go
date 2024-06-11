package cmd

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/ticketstore"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/services"
	"github.com/pastelnetwork/gonode/walletnode/configs"
	"github.com/pastelnetwork/gonode/walletnode/node/grpc"
	"github.com/pastelnetwork/gonode/walletnode/services/cascaderegister"
	"github.com/pastelnetwork/gonode/walletnode/services/collectionregister"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
	"github.com/pastelnetwork/gonode/walletnode/services/healthcheckchallenge"
	"github.com/pastelnetwork/gonode/walletnode/services/metrics"
	"github.com/pastelnetwork/gonode/walletnode/services/nftregister"
	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"
	"github.com/pastelnetwork/gonode/walletnode/services/score"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	"github.com/pastelnetwork/gonode/walletnode/services/storagechallenge"
)

const (
	appName        = "walletnode"
	appUsage       = "WalletNode" // TODO: Write a clear description.
	rqFilesDir     = "rqfiles"
	staticFilesDir = "files"
)

var (
	defaultPath = configurer.DefaultPath()

	defaultTempDir          = filepath.Join(os.TempDir(), appName)
	defaultConfigFile       = filepath.Join(defaultPath, appName+".yml")
	defaultPastelConfigFile = filepath.Join(defaultPath, "pastel.conf")
	defaultRqFilesDir       = filepath.Join(defaultPath, rqFilesDir)
	defaultStaticFilesDir   = filepath.Join(defaultPath, staticFilesDir)
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
		cli.NewFlag("temp-dir", &config.TempDir).SetUsage("Set `path` for storing temp data.").SetValue(defaultTempDir),
		cli.NewFlag("rq-files-dir", &config.RqFilesDir).SetUsage("Set `path` for storing files for rqservice.").SetValue(defaultRqFilesDir),
		cli.NewFlag("log-level", &config.LogConfig.Level).SetUsage("Set the log `level`.").SetValue(config.LogConfig.Level),
		cli.NewFlag("log-file", &config.LogConfig.File).SetUsage("The log `file` to write to."),
		cli.NewFlag("quiet", &config.Quiet).SetUsage("Disallows log output to stdout.").SetAliases("q"),
		// API
		cli.NewFlag("swagger", &config.API.Swagger).SetUsage("Enable Swagger UI."),
	)

	//SetActionFunc is the default, and in our case only, executed action.
	//Sets up configs and logging, and also returns the "app" function to main.go to be called (Run) there.
	app.SetActionFunc(func(ctx context.Context, args []string) error {
		//Sets logging prefix to pastel-app
		ctx = log.ContextWithPrefix(ctx, "pastel-app")

		//Parse config files supplied in arguments
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

		if err := os.MkdirAll(config.TempDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create temp-dir %q, %w", config.TempDir, err)
		}

		if err := os.MkdirAll(config.RqFilesDir, os.ModePerm); err != nil {
			return errors.Errorf("could not create rq-files-dir %q, %w", config.RqFilesDir, err)
		}
		//This function, defined below, is what will be run when "Run" is called on our app
		return runApp(ctx, config)
	})

	//while "app" is returned, the default action func we just set as a parameter to it is going to be called when run
	return app
}

// What walletnode actually does when you run it!
// - Creates App Context
// - Registers system interrupts for shutdown
// - Configures RPC connections to nodeC and raptorq
// - Sets up queries file storage for NFT registration
// - Set minimum confirmation requirements for transactions
// - Create services for the functions we want to expose to users, allowing for running and tasking
// - Create API endpoints for those services
// - Run those services in their own goroutines and wait for them
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

	if _, err := os.Stat(defaultStaticFilesDir); os.IsNotExist(err) {
		// directory does not exist, create it
		errDir := os.MkdirAll(defaultStaticFilesDir, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	// pastelClient reads in the json-formatted hostname, port, username, and password and
	//  connects over gRPC to cNode for access to Blockchain, Masternodes, Tickets, and PastelID databases
	pastelClient := pastel.NewClient(config.Pastel, config.Pastel.BurnAddress())
	// wrap pastelClient in security functionality then allow for calling of nfts, userdata, and sense functions
	nodeClient := grpc.NewClient(pastelClient)

	// raptorq client
	rqAddr := fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	config.NftRegister.RaptorQServiceAddress = rqAddr
	config.NftRegister.RqFilesDir = config.RqFilesDir
	config.CascadeRegister.RaptorQServiceAddress = rqAddr
	config.CascadeRegister.RqFilesDir = config.RqFilesDir

	// NB: As part of current dev push for Sense and Cascade, we are disabling userdata handling thru rqlite.

	// create another wrapped RPC connection for use by our userdata service, basically the same as nodeClient
	//userdataNodeClient := grpc.NewClient(pastelClient)

	// start new key value storage
	db := memory.NewKeyValue()
	// Initialize temporary queries file storage
	fileStorage := fs.NewFileStorage(config.TempDir)

	//Initialize History DB
	hDB, err := queries.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error connecting history db..")
	}
	defer hDB.CloseHistoryDB(ctx)

	tDB, err := ticketstore.OpenTicketingDb()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error connecting ticket db..")
	}
	defer tDB.CloseTicketDB(ctx)

	//Set minimum confirmation requirements for transactions to ensure completion
	if config.RegTxMinConfirmations > 0 {
		config.NftRegister.NFTRegTxMinConfirmations = config.RegTxMinConfirmations
		config.SenseRegister.SenseRegTxMinConfirmations = config.RegTxMinConfirmations
		config.CascadeRegister.CascadeRegTxMinConfirmations = config.RegTxMinConfirmations
		config.CollectionRegister.CollectionRegTxMinConfirmations = config.RegTxMinConfirmations
	}

	if config.ActTxMinConfirmations > 0 {
		config.NftRegister.NFTActTxMinConfirmations = config.ActTxMinConfirmations
		config.SenseRegister.SenseActTxMinConfirmations = config.ActTxMinConfirmations
		config.CascadeRegister.CascadeActTxMinConfirmations = config.ActTxMinConfirmations
		config.CollectionRegister.CollectionActTxMinConfirmations = config.ActTxMinConfirmations
	}

	// The API Server takes our configured services and wraps them further with "Mount", creating the API endpoints.
	//  Since the API Server has access to the services, this is what finally exposes useful methods like
	//  "NftGet" and "Download".
	apiSrcvConf := &services.Config{
		StaticFilesDir: defaultStaticFilesDir,
	}
	config.API.StaticFilesDir = defaultStaticFilesDir
	config.NftDownload.StaticDir = defaultStaticFilesDir

	// These services connect the different clients and configs together to provide tasking and handling for
	//  the required functionality.  These services aren't started with these declarations, they will be run
	//	later through the API Server.
	nftSearch := nftsearch.NewNftSearchService(&config.NftSearch, pastelClient, nodeClient, hDB)
	nftDownload := download.NewNftDownloadService(&config.NftDownload, pastelClient, nodeClient, hDB)

	cascadeRegister := cascaderegister.NewService(&config.CascadeRegister, pastelClient, nodeClient, fileStorage, db, *nftDownload, hDB, tDB)
	senseRegister := senseregister.NewService(&config.SenseRegister, pastelClient, nodeClient, fileStorage, nftDownload, db, hDB)
	nftRegister := nftregister.NewService(&config.NftRegister, pastelClient, nodeClient, fileStorage, db, nftDownload, hDB)
	collectionRegister := collectionregister.NewService(&config.CollectionRegister, pastelClient, nodeClient, hDB)
	metricsService := metrics.NewMetricsService(pastelClient)
	storageChallengeService := storagechallenge.NewStorageChallengeService(pastelClient)
	scoreService := score.NewScoreService(pastelClient)
	healthCheckChallengeService := healthcheckchallenge.NewHealthCheckChallengeService(pastelClient)

	fileMappings := &sync.Map{}
	server := api.NewAPIServer(config.API, fileMappings, pastelClient,
		services.NewNftAPIHandler(apiSrcvConf, fileMappings, nftRegister, nftSearch, nftDownload),
		// services.NewUserdataAPIHandler(userdataProcess),
		services.NewSenseAPIHandler(apiSrcvConf, senseRegister, nftDownload),
		services.NewCascadeAPIHandler(apiSrcvConf, fileMappings, cascadeRegister, nftDownload),
		services.NewCollectionAPIIHandler(apiSrcvConf, collectionRegister),
		services.NewSwagger(apiSrcvConf),
		services.NewMetricsAPIHandler(metricsService),
		services.NewStorageChallengeAPIHandler(storageChallengeService),
		services.NewHealthCheckChallengeAPIHandler(healthCheckChallengeService),
		services.NewScoreAPIHandler(scoreService),
	)

	log.WithContext(ctx).Infof("Config: %s", config)

	//calls all of our services in its own goroutine with "Run" but waits for them as a group
	//NB: Tasks are defined individually for each service, and probably override the actual task struct's functions.
	//  For instance, nftRegister uses the NftRegistrationTask found in services/nftregister/task.go
	//return runServices(ctx, server, nftRegister, nftSearch, nftDownload, userdataProcess, senseRegister)
	return runServices(ctx, server, nftRegister, nftSearch, nftDownload, senseRegister, cascadeRegister, collectionRegister)
}
