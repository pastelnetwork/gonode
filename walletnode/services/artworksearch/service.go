package artworksearch

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
)

const (
	logPrefix = "artwork"
)

// Service represents a service for the artwork search.
type Service struct {
	*task.Worker
	p2pClient    p2p.Client
	pastelClient pastel.Client
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

// Task returns the task of the artwork search by the given id.
func (service *Service) Task(id string) *Task {
	return service.Worker.Task(id).(*Task)
}

// AddTask runs a new task of the artwork search and returns its taskID.
func (service *Service) AddTask(request *artworks.ArtSearchPayload) string {
	task := NewTask(service, request)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(pastelClient pastel.Client, p2pClient p2p.Client) *Service {
	return &Service{
		pastelClient: pastelClient,
		p2pClient:    p2pClient,
		Worker:       task.NewWorker(),
	}
}
