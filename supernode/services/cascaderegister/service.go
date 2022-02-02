package cascaderegister

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
	logPrefix = "cascade"
)

// Service represent sense service.
type CascadeRegistrationService struct {
	*common.SuperNodeService
	config *Config

	nodeClient node.ClientInterface
}

// Run starts task
func (service *CascadeRegistrationService) Run(ctx context.Context) error {
	return service.RunHelper(ctx, service.config.PastelID, logPrefix)
}

// NewSenseRegistrationTask runs a new task of the registration Sense and returns its taskID.
func (service *CascadeRegistrationService) NewTask() *CascadeRegistrationTask {
	task := NewCascadeRegistrationTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the Sense registration by the given id.
func (service *CascadeRegistrationService) Task(id string) *CascadeRegistrationTask {
	return service.Worker.Task(id).(*CascadeRegistrationTask)
}

// NewService returns a new Service instance.
func NewService(config *Config,
	fileStorage storage.FileStorageInterface,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	p2pClient p2p.Client,
	rqClient rqnode.ClientInterface,
) *CascadeRegistrationService {
	return &CascadeRegistrationService{
		SuperNodeService: common.NewSuperNodeService(fileStorage, pastelClient, p2pClient, rqClient),
		config:           config,
		nodeClient:       nodeClient,
	}
}
