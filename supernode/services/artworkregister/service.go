package artworkregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/probe"
	"github.com/pastelnetwork/gonode/supernode/node"
	"golang.org/x/sync/errgroup"
)

const (
	logPrefix = "artwork"
)

// Service represent artwork service.
type Service struct {
	*task.Worker
	*artwork.Storage

	config       *Config
	db           storage.KeyValue
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
	group.Go(func() (err error) {
		defer errors.Recover(func(recErr error) { err = recErr })
		return service.Storage.Run(ctx)
	})
	group.Go(func() (err error) {
		defer errors.Recover(func(recErr error) { err = recErr })
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
func NewService(config *Config, db storage.KeyValue, fileStorage storage.FileStorage, probe probe.Probe, pastelClient pastel.Client, nodeClient node.Client) *Service {
	return &Service{
		config:       config,
		db:           db,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
		Storage:      artwork.NewStorage(fileStorage),
	}
}
