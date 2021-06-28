package artworkdownload

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "artwork"
)

// Service represents a service for the registration artwork.
type Service struct {
	*task.Worker

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.Client
}

// Run starts worker.
func (service *Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *Service) Tasks() []*Task {
	var tasks []*Task
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*Task))
	}
	return tasks
}

// Task returns the task of the artwork downloading by the given id.
func (service *Service) Task(id string) *Task {
	if t := service.Worker.Task(id); t != nil {
		return t.(*Task)
	}
	return nil
}

// AddTask adds a new task of the artwork downloading and returns its taskID.
func (service *Service) AddTask(ticket *Ticket) string {
	task := NewTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.Client) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
	}
}
