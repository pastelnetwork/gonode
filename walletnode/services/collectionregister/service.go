package collectionregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/collection"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	logPrefix = "collection"
)

// CollectionRegistrationService represents a service for Cascade Open API
type CollectionRegistrationService struct {
	config *Config

	*task.Worker

	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface
	historyDB     storage.LocalStoreInterface
}

// Run starts worker.
func (service *CollectionRegistrationService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// Run worker service
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *CollectionRegistrationService) Tasks() []*CollectionRegistrationTask {
	var tasks []*CollectionRegistrationTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*CollectionRegistrationTask))
	}
	return tasks
}

// GetTask returns the task of the Collection OpenAPI by the given id.
func (service *CollectionRegistrationService) GetTask(id string) *CollectionRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*CollectionRegistrationTask)
	}
	return nil
}

// AddTask create ticket request and start a new task with the given payload
func (service *CollectionRegistrationService) AddTask(p *collection.RegisterCollectionPayload) (string, error) {
	request := FromCollectionRegistrationPayload(p)

	task := NewCollectionRegistrationTask(service, request)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// ValidateUser validates user
func (service *CollectionRegistrationService) ValidateUser(ctx context.Context, id string, pass string) bool {
	return common.ValidateUser(ctx, service.pastelHandler.PastelClient, id, pass) &&
		common.IsPastelIDTicketRegistered(ctx, service.pastelHandler.PastelClient, id)
}

// NewService returns a new Service instance
func NewService(
	config *Config,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	historyDB storage.LocalStoreInterface,
) *CollectionRegistrationService {
	return &CollectionRegistrationService{
		config:        config,
		Worker:        task.NewWorker(),
		nodeClient:    nodeClient,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		historyDB:     historyDB,
	}
}
