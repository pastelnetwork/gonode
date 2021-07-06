package artworksearch

import (
	"context"
	"fmt"

	b58 "github.com/jbenet/go-base58"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
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
func (service *Service) AddTask(request *ArtSearchRequest) string {
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

// FetchThumbnail gets artwork thumbnail
func (service *Service) FetchThumbnail(ctx context.Context, res *pastel.RegTicket) (data []byte, err error) {
	hash := res.RegTicketData.ArtTicketData.AppTicketData.Thumbnail1Hash
	key := b58.Encode(hash)

	return service.p2pClient.Retrieve(ctx, key)
}

// RegTicket pull art registration ticket from cNode & decodes base64 encoded fields
func (service *Service) RegTicket(ctx context.Context, RegTXID string) (*pastel.RegTicket, error) {
	regTicket, err := service.pastelClient.RegTicket(ctx, RegTXID)
	if err != nil {
		return nil, fmt.Errorf("fetch: %s", err)
	}

	if err := fromBase64(string(regTicket.RegTicketData.ArtTicket),
		&regTicket.RegTicketData.ArtTicketData); err != nil {
		return &regTicket, fmt.Errorf("convert art ticket: %s", err)
	}

	if err := fromBase64(string(regTicket.RegTicketData.ArtTicketData.AppTicket),
		&regTicket.RegTicketData.ArtTicketData.AppTicketData); err != nil {
		return &regTicket, fmt.Errorf("convert app ticket: %s", err)
	}

	return &regTicket, nil
}
