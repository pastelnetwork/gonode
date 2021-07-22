package userdataprocess

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
)

const (
	logPrefix = "userdata"
)

// Service represent userdata service.
type Service struct {
	*task.Worker
	config       *Config
	pastelClient pastel.Client
	nodeClient   node.Client
}

// Run starts task
func (service *Service) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if service.config.PastelID == "" {
		return errors.New("PastelID is not specified in the config file")
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Task returns the task of the registration userdata by the given id.
func (service *Service) Task(id string) *Task {
	return service.Worker.Task(id).(*Task)
}

// NewTask runs a new task of the registration userdata and returns its taskID.
func (service *Service) NewTask() *Task {
	task := NewTask(service)
	service.Worker.AddTask(task)

	return task
}

// NewService returns a new Service instance.
func NewService(config *Config, fileStorage storage.FileStorage, probeTensor probe.Tensor, pastelClient pastel.Client, nodeClient node.Client, p2pClient p2p.Client, rqClient rqnode.Client) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
	}
}
