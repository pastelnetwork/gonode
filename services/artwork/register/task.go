package register

import (
	"context"
	"sync/atomic"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

var taskID uint32

// Task is the task of registering new artwork.
type Task struct {
	*Service

	ID     int
	State  *state.State
	Ticket *Ticket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	if err := task.run(ctx); err != nil {
		if err, ok := err.(*TaskError); ok {
			task.State.Update(state.NewMessage(state.StatusTaskRejected))
			log.WithField("error", err).Debugf("Task %d is rejected", task.ID)
			return nil
		}
		return err
	}

	task.State.Update(state.NewMessage(state.StatusTaskCompleted))
	log.Debugf("Task %d is completed", task.ID)
	return nil
}

func (task *Task) run(ctx context.Context) error {
	superNodes, err := task.findSuperNodes(ctx)
	if err != nil {
		return err
	}

	if len(superNodes) < task.config.NumberSuperNodes {
		task.State.Update(state.NewMessage(state.StatusErrorTooLowFee))
		return NewTaskError(errors.Errorf("not found %d SuperNodes with acceptable storage fee", task.config.NumberSuperNodes))
	}

	return nil
}

func (task *Task) findSuperNodes(ctx context.Context) ([]SuperNode, error) {
	var superNodes []SuperNode

	mns, err := task.pastel.TopMasterNodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > task.Ticket.MaximumFee {
			continue
		}
		superNodes = append(superNodes, SuperNode{
			Address: mn.ExtAddress,
			Key:     mn.ExtKey,
			Fee:     mn.Fee,
		})
	}

	return superNodes, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Ticket) *Task {
	return &Task{
		Service: service,
		ID:      int(atomic.AddUint32(&taskID, 1)),
		Ticket:  Ticket,
		State:   state.New(state.NewMessage(state.StatusTaskStarted)),
	}
}
