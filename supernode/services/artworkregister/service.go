package artworkregister

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/storage"
)

const (
	logPrefix = "artwork"
)

// Service represent artwork service.
type Service struct {
	sync.Mutex

	config       *Config
	db           storage.KeyValue
	pastelClient pastel.Client
	nodeClient   node.Client
	worker       *Worker
	tasks        []*Task
}

// Run starts worker
func (service *Service) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	nodeConfig, err := service.pastelClient.MasterNodeConfig(ctx)
	if err != nil {
		return err
	}
	if nodeConfig.ExtKey == "" {
		return errors.New("masternode configuration: \"extKey\" cannot be empty, check \"masternode list-conf\"")
	}
	service.config.node = nodeConfig

	return service.worker.Run(ctx)
}

// TaskByConnID returns the task of the registration artwork by the given connID.
func (service *Service) TaskByConnID(connID string) *Task {
	service.Lock()
	defer service.Unlock()

	for _, task := range service.tasks {
		if task.ConnID == connID {
			return task
		}
	}
	return nil
}

// NewTask runs a new task of the registration artwork and returns its taskID.
func (service *Service) NewTask(ctx context.Context) *Task {
	service.Lock()
	defer service.Unlock()

	task := NewTask(service)
	service.tasks = append(service.tasks, task)
	service.worker.AddTask(ctx, task)

	return task
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
