package artwork

import (
	"context"

	"github.com/pastelnetwork/walletnode/services/artwork/register"
	"github.com/pastelnetwork/walletnode/storage"
)

// const logPrefix = "[artwork]"

// Service represent artwork service.
type Service struct {
	db     storage.KeyValue
	worker *register.Worker
	tasks  []*register.Task
}

// Tasks returns all tasks.
func (service *Service) Tasks() []*register.Task {
	return service.tasks
}

// Run starts worker
func (service *Service) Run(ctx context.Context) error {
	return service.worker.Run(ctx)
}

// Task returns the task of the registration artwork.
func (service *Service) Task(taskID int) *register.Task {
	for _, task := range service.tasks {
		if task.ID() == taskID {
			return task
		}
	}
	return nil
}

// Register runs a new task of the registration artwork and returns its taskID.
func (service *Service) Register(ctx context.Context, ticket *register.Ticket) (int, error) {
	// NOTE: for testing
	task := register.NewTask(ticket)
	service.tasks = append(service.tasks, task)
	service.worker.AddTask(ctx, task)

	return task.ID(), nil
}

// New returns a new Service instance.
func New(db storage.KeyValue) *Service {
	return &Service{
		db:     db,
		worker: register.NewWorker(),
	}
}
