package userdataprocess

import (
	"context"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/metadb/database"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "userdata"
)

// Service represent userdata service.
type Service struct {
	*common.SuperNodeService

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.ClientInterface
	databaseOps  *database.Ops
}

// Run starts task
func (service *Service) Run(ctx context.Context) error {
	return service.RunHelper(ctx, service.config.PastelID, logPrefix)
}

// NewUserDataTask runs a new task of the registration userdata and returns its taskID.
func (service *Service) NewTask() *UserDataTask {
	task := NewUserDataTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the registration userdata by the given id.
func (service *Service) Task(id string) *UserDataTask {
	return service.Worker.Task(id).(*UserDataTask)
}

// NewService returns a new Service instance.
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface, databaseOps *database.Ops) *Service {
	return &Service{
		SuperNodeService: &common.SuperNodeService{
			Worker:  task.NewWorker(),
			Storage: nil,
		},
		config:       config,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		databaseOps:  databaseOps,
	}
}
