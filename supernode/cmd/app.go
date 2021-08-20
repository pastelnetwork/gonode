package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/common/version"
	"github.com/pastelnetwork/gonode/dupedetection"
	"github.com/pastelnetwork/gonode/metadb"
	"github.com/pastelnetwork/gonode/metadb/database"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"
	"github.com/pastelnetwork/gonode/supernode/configs"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/client"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
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

	rqliteDefaultPort = 4446
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
		cli.NewFlag("dd-service-dir", &config.DdWorkDir).SetUsage("Set `path` to the work directory of dupe detection service.").SetValue(defaultDdWorkDir),
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

func getDatabaseNodes(ctx context.Context, pastelClient pastel.Client) ([]string, error) {
	var nodeIPList []string
	// Note: Lock the db clustering feature, run it in single node first
	if test := sys.GetStringEnv("DB_CLUSTER", ""); test != "" {
		nodeList, err := pastelClient.MasterNodesList(ctx)
		if err != nil {
			return nil, err
		}
		for _, nodeInfo := range nodeList {
			address := nodeInfo.ExtAddress
			segments := strings.Split(address, ":")
			if len(segments) != 2 {
				return nil, errors.Errorf("malformed db node address: %s", address)
			}
			nodeAddress := fmt.Sprintf("%s:%d", segments[0], rqliteDefaultPort)
			nodeIPList = append(nodeIPList, nodeAddress)
		}
	}

	// Important notice !!!
	// Enable this line below to run rqlite with 1 leader and multiple replicas, if want to test the real rqlite cluster behavior
	// Otherwise, this will run rqlite cluster with every node as a rqlite leader, and data to highest rank SN can write directly to its local rqlite
	// return []string{"0.0.0.0:4041"}, nil

	return nodeIPList, nil
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
	nodeClient := client.New()
	fileStorage := fs.NewFileStorage(config.TempDir)

	// analysis tools
	config.DupeDetection.SetWorkDir(config.DdWorkDir)
	ddClient := dupedetection.NewClient(config.DupeDetection)

	// p2p service (currently using kademlia)
	config.P2P.SetWorkDir(config.WorkDir)
	p2p := p2p.New(config.P2P)

	nodeIPList, err := getDatabaseNodes(ctx, pastelClient)
	if err != nil {
		return err
	}

	// new metadb service
	config.MetaDB.SetWorkDir(config.WorkDir)
	metadb := metadb.New(config.MetaDB, config.Node.PastelID, nodeIPList)
	database := database.NewDatabaseOps(metadb, config.UserDB)

	// raptorq client
	config.ArtworkRegister.RaptorQServiceAddress = fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	config.ArtworkRegister.RqFilesDir = config.RqFilesDir
	config.ArtworkRegister.PreburntTxMinConfirmations = config.PreburntTxMinConfirmations
	config.ArtworkRegister.PreburntTxConfirmationTimeout = time.Duration(config.PreburntTxConfirmationTimeout * int(time.Minute))
	rqClient := rqgrpc.NewClient()

	// business logic services
	config.ArtworkDownload.RaptorQServiceAddress = fmt.Sprint(config.RaptorQ.Host, ":", config.RaptorQ.Port)
	config.ArtworkDownload.RqFilesDir = config.RqFilesDir

	// business logic services
	artworkRegister := artworkregister.NewService(&config.ArtworkRegister, fileStorage, pastelClient, nodeClient, p2p, rqClient, ddClient)
	artworkDownload := artworkdownload.NewService(&config.ArtworkDownload, pastelClient, p2p, rqClient)
	dupeDetection, err := dupedetection.NewService(config.DupeDetection, pastelClient, p2p)
	if err != nil {
		return errors.Errorf("can not start dupe detection service, err: %w", err)
	}

	// ----Userdata Services----
	userdataNodeClient := client.New()
	userdataProcess := userdataprocess.NewService(&config.UserdataProcess, pastelClient, userdataNodeClient, database)

	// server
	grpc := server.New(config.Server,
		walletnode.NewRegisterArtwork(artworkRegister),
		supernode.NewRegisterArtwork(artworkRegister),
		walletnode.NewDownloadArtwork(artworkDownload),
		walletnode.NewProcessUserdata(userdataProcess, database),
		supernode.NewProcessUserdata(userdataProcess, database),
	)

	return runServices(ctx, metadb, grpc, p2p, artworkRegister, artworkDownload, dupeDetection, database, userdataProcess)
}
