package artworkregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel-client"

	"github.com/pastelnetwork/gonode/walletnode/node"
	"golang.org/x/sync/errgroup"
)

const (
	logPrefix = "artwork"
)

// Service represent artwork service.
type Service struct {
	config       *Config
	db           storage.KeyValue
	pastelClient pastel.Client
	worker       *Worker
	tasks        []*Task
	nodeClient   node.Client
}

// Run starts worker
func (service *Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		defer errors.Recover(func(recErr error) { err = recErr })
		return service.worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *Service) Tasks() []*Task {
	return service.tasks
}

// Task returns the task of the registration artwork.
func (service *Service) Task(taskID string) *Task {
	for _, task := range service.tasks {
		if task.ID == taskID {
			return task
		}
	}
	return nil
}

// AddTask runs a new task of the registration artwork and returns its taskID.
func (service *Service) AddTask(ctx context.Context, ticket *Ticket) (string, error) {
	// NOTE: for testing
	task := NewTask(service, ticket)
	service.tasks = append(service.tasks, task)
	service.worker.AddTask(ctx, task)

	return task.ID, nil
}

// NewService returns a new Service instance.
func NewService(config *Config, db storage.KeyValue, pastelClient pastel.Client, nodeClient node.Client) *Service {
	return &Service{
		config:       config,
		db:           db,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		worker:       NewWorker(),
	}
}
