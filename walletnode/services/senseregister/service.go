package senseregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// Service represents a service for the registration artwork.
type Service struct {
	*task.Worker
	*artwork.Storage

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.Client
	db           storage.KeyValue
}

// Run starts worker.
func (service Service) Run(ctx context.Context) error {
	// TODO: implement service run logic

	return nil
}

// AddTask adds a task to the worker.
func (service *Service) AddTask(ticket *Request) string {
	task := NewTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance
func NewService(
	config *Config,
	fileStorage storage.FileStorage,
	pastelClient pastel.Client,
	nodeClient node.Client,
	db storage.KeyValue,
) *Service {
	return &Service{
		config:       config,
		db:           db,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
		Storage:      artwork.NewStorage(fileStorage),
	}
}
