package senseregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/ddstore"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "sense"
)

// SenseRegistrationService represent sense service.
type SenseRegistrationService struct {
	*common.SuperNodeService
	config *Config

	nodeClient node.ClientInterface
	ddClient   ddclient.DDServerClient
}

// Run starts task
func (service *SenseRegistrationService) Run(ctx context.Context) error {
	return service.RunHelper(ctx, service.config.PastelID, logPrefix)
}

// NewSenseRegistrationTask runs a new task of the registration Sense and returns its taskID.
func (service *SenseRegistrationService) NewSenseRegistrationTask() *SenseRegistrationTask {
	task := NewSenseRegistrationTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the Sense registration by the given id.
func (service *SenseRegistrationService) Task(id string) *SenseRegistrationTask {
	return service.Worker.Task(id).(*SenseRegistrationTask)
}

// GetDupeDetectionDatabaseHash returns the hash of the dupe detection database.
func (service *SenseRegistrationService) GetDupeDetectionDatabaseHash(ctx context.Context) (string, error) {
	db, err := ddstore.NewSQLiteDDStore(service.config.DDDatabase)
	if err != nil {
		return "", err
	}

	hash, err := db.GetDDDataHash(ctx)
	if err != nil {
		return "", err
	}

	if err := db.Close(); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to close dd database")
	}

	return hash, nil
}

// NewService returns a new Service instance.
func NewService(config *Config,
	fileStorage storage.FileStorageInterface,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	p2pClient p2p.Client,
	rqClient rqnode.ClientInterface,
	ddClient ddclient.DDServerClient,
) *SenseRegistrationService {
	return &SenseRegistrationService{
		SuperNodeService: common.NewSuperNodeService(fileStorage, pastelClient, p2pClient, rqClient),
		config:           config,
		nodeClient:       nodeClient,
		ddClient:         ddClient,
	}
}
