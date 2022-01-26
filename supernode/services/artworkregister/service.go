package artworkregister

import (
	"context"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
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
	*files.Storage

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.Client
	p2pClient    p2p.Client
	rqClient     rqnode.ClientInterface
	ddClient     ddclient.DDServerClient
}

// run starts task
func (service *Service) run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if service.config.PastelID == "" {
		return errors.New("PastelID is not specified in the config file")
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	group.Go(func() error {
		return service.Storage.Run(ctx)
	})
	return group.Wait()
}

// Run starts task
func (service *Service) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):
			if err := service.run(ctx); err != nil {
				if utils.IsContextErr(err) {
					return err
				}

				service.Worker = task.NewWorker()
				log.WithContext(ctx).WithError(err).Error("registration failed, retrying")
			} else {
				return nil
			}
		}
	}
}

// NewTask runs a new task of the registration artwork and returns its taskID.
func (service *Service) NewTask() *Task {
	task := NewTask(service)
	service.Worker.AddTask(task)

	return task
}

// Task returns the task of the registration artwork by the given id.
func (service *Service) Task(id string) *Task {
	return service.Worker.Task(id).(*Task)
}

// NewService returns a new Service instance.
func NewService(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.Client, p2pClient p2p.Client, rqClient rqnode.ClientInterface, ddClient ddclient.DDServerClient) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		p2pClient:    p2pClient,
		rqClient:     rqClient,
		ddClient:     ddClient,
		Worker:       task.NewWorker(),
		Storage:      files.NewStorage(fileStorage),
	}
}
