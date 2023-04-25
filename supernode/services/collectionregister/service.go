package collectionregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "collection"
)

// CollectionRegistrationService represent collection service.
type CollectionRegistrationService struct {
	*common.SuperNodeService
	config *Config

	nodeClient node.ClientInterface
}

// Run starts task
func (service *CollectionRegistrationService) Run(ctx context.Context) error {
	return service.RunHelper(ctx, service.config.PastelID, logPrefix)
}

// NewCollectionRegistrationTask runs a new task of the registration Collection and returns its taskID.
func (service *CollectionRegistrationService) NewCollectionRegistrationTask() *CollectionRegistrationTask {
	task := NewCollectionRegistrationTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the Collection registration by the given id.
func (service *CollectionRegistrationService) Task(id string) *CollectionRegistrationTask {
	return service.Worker.Task(id).(*CollectionRegistrationTask)
}

// NewService returns a new Service instance.
func NewService(config *Config,
	fileStorage storage.FileStorageInterface,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	p2pClient p2p.Client,
	rqClient rqnode.ClientInterface,
) *CollectionRegistrationService {
	return &CollectionRegistrationService{
		SuperNodeService: common.NewSuperNodeService(fileStorage, pastelClient, p2pClient, rqClient),

		config:     config,
		nodeClient: nodeClient,
	}
}
