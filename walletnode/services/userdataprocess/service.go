package userdataprocess

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb/network/walletnode/node"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix = "userdata"
)

// Service represents a service for the userdata process
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

// Task returns the task of the userdata process by the given id.
func (service *Service) Task(id string) *Task {
	return service.Worker.Task(id).(*Task)
}

// AddTask runs a new task of the userdata process and returns its taskID.
func (service *Service) AddTask(request *userdata.ProcessRequest, retrieve string) string {
	task := NewTask(service, request, retrieve)
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
