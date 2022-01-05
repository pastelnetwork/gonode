package userdataprocess

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/metadb/database"
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
	databaseOps  *database.Ops
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
				log.WithContext(ctx).WithError(err).Error("failed to run userdata process, retrying.")
			} else {
				return nil
			}
		}
	}
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
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.Client, databaseOps *database.Ops) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
		databaseOps:  databaseOps,
	}
}
