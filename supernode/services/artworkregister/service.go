package artworkregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	ddclient "github.com/pastelnetwork/gonode/dupedetection/node"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node"
)

const (
	logPrefix = "artwork"
)

// Service represent artwork service.
type Service struct {
	*task.Worker
	*artwork.Storage

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.Client
	p2pClient    p2p.Client
	rqClient     rqnode.Client
	ddClient     ddclient.Client
}

// Run starts task
func (service *Service) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if service.config.PastelID == "" {
		return errors.New("PastelID is not specified in the config file")
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return service.Storage.Run(ctx)
	})
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Task returns the task of the registration artwork by the given id.
func (service *Service) Task(id string) *Task {
	return service.Worker.Task(id).(*Task)
}

// NewTask runs a new task of the registration artwork and returns its taskID.
func (service *Service) NewTask() *Task {
	task := NewTask(service)
	service.Worker.AddTask(task)

	return task
}

// NewService returns a new Service instance.
func NewService(config *Config, fileStorage storage.FileStorage, pastelClient pastel.Client, nodeClient node.Client, p2pClient p2p.Client, rqClient rqnode.Client, ddClient ddclient.Client) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		p2pClient:    p2pClient,
		rqClient:     rqClient,
		ddClient:     ddClient,
		Worker:       task.NewWorker(),
		Storage:      artwork.NewStorage(fileStorage),
	}
}
