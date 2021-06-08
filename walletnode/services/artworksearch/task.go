package artworksearch

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
)

// Task is the task of searching for artwork.
type Task struct {
	task.Task
	*Service

	resultChan chan *pastel.Ticket
	Request    *artworks.ArtSearchPayload
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")

	if err := task.run(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Warnf("Task failed")
		return nil
	}

	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err := task.actTickets(ctx)
	if err != nil {
		task.UpdateStatus(StatusTaskFailure)
		return err
	}

	// TODO: iterate through tickets & get reg tickets
	// In parallel; filter reg tickets based on search params
	// & fetch thumbnail as soon as one matches
	// send the thumbnail along with data through resultChan

	return nil
}

// SubscribeSearchResult returns a new search resultof the state.
func (task *Task) SubscribeSearchResult() chan *pastel.Ticket {
	return task.resultChan
}

func (task *Task) actTickets(ctx context.Context) (tickets pastel.ActivationTickets, err error) {
	actTickets, err := task.pastelClient.ActTickets(ctx, pastel.ActTicketAll, task.Request.MinBlock)
	if err != nil {
		return tickets, err
	}

	for i := 0; i < len(actTickets); i++ {
		if task.Request.MaxBlock != nil && actTickets[i].ArtistHeight > *task.Request.MaxBlock {
			continue
		}
		if task.Request.Artist != nil && *task.Request.Artist != actTickets[i].PastelID {
			continue
		}

		tickets = append(tickets, actTickets[i])
	}

	return tickets, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, request *artworks.ArtSearchPayload) *Task {
	return &Task{
		Task:       task.New(StatusTaskStarted),
		Service:    service,
		Request:    request,
		resultChan: make(chan *pastel.Ticket),
	}
}
