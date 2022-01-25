package userdataprocess

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "userdata"
)

// Service represents a service for the userdata process
type UserDataService struct {
	*task.Worker

	config        *Config
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface
}

// Run starts worker.
func (service *UserDataService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})

	return group.Wait()
}

// Tasks returns all tasks.
func (service *UserDataService) Tasks() []*UserDataTask {
	var tasks []*UserDataTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*UserDataTask))
	}
	return tasks
}

// Task returns the task of the userdata process by the given id.
func (service *UserDataService) Task(id string) *UserDataTask {
	return service.Worker.Task(id).(*UserDataTask)
}

// AddTask runs a new task of the userdata process and returns its taskID.
func (service *UserDataService) AddTask(request *userdata.ProcessRequest, userpastelid string) string {
	task := NewUserDataTask(service, request, userpastelid)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface) *UserDataService {
	return &UserDataService{
		Worker:        task.NewWorker(),
		config:        config,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		nodeClient:    nodeClient,
	}
}
