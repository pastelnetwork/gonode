package cascaderegister

import (
	"context"
	"github.com/pastelnetwork/gonode/common/storage/queries"

	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "cascade"
)

// CascadeRegistrationService represent sense service.
type CascadeRegistrationService struct {
	*common.SuperNodeService
	config *Config

	nodeClient node.ClientInterface
	historyDB  queries.LocalStoreInterface
}

// Run starts task
func (service *CascadeRegistrationService) Run(ctx context.Context) error {
	return service.RunHelper(ctx, service.config.PastelID, logPrefix)
}

// NewCascadeRegistrationTask runs a new task of the registration Sense and returns its taskID.
func (service *CascadeRegistrationService) NewCascadeRegistrationTask() *CascadeRegistrationTask {
	task := NewCascadeRegistrationTask(service)
	service.Worker.AddTask(task)

	return task
}

// Task returns the task of the Sense registration by the given id.
func (service *CascadeRegistrationService) Task(id string) *CascadeRegistrationTask {
	if service.Worker.Task(id) == nil {
		return nil
	}

	return service.Worker.Task(id).(*CascadeRegistrationTask)
}

// NewService returns a new Service instance.
func NewService(config *Config, fileStorage storage.FileStorageInterface,
	pastelClient pastel.Client, nodeClient node.ClientInterface, p2pClient p2p.Client, historyDB queries.LocalStoreInterface) *CascadeRegistrationService {

	return &CascadeRegistrationService{
		SuperNodeService: common.NewSuperNodeService(fileStorage, pastelClient, p2pClient),
		config:           config,
		nodeClient:       nodeClient,
		historyDB:        historyDB,
	}
}
