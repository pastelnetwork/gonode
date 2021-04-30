package artworkregister

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/walletnode/node"
	"github.com/pastelnetwork/walletnode/services/artworkregister/state"
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
	log.Debugf("Start task %v", task.ID)
	defer log.Debugf("End task %v", task.ID)

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

	for _, superNode := range superNodes {
		superNode.Address = "127.0.0.1:4444"
		if err := task.connect(ctx, superNode); err != nil {
			return err
		}
	}
	// if len(superNodes) < task.config.NumberSuperNodes {
	// 	task.State.Update(state.NewMessage(state.StatusErrorTooLowFee))
	// 	return NewTaskError(errors.Errorf("not found enough available SuperNodes with acceptable storage fee", task.config.NumberSuperNodes))
	// }

	return nil
}

func (task *Task) connect(ctx context.Context, superNode *node.SuperNode) error {
	log.Debugf("connect to %s", superNode.Address)

	ctxConnect, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	conn, err := task.nodeClient.Connect(ctxConnect, superNode.Address)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Debugf("connected to %s", superNode.Address)

	stream, err := conn.RegisterArtowrk(ctx)
	if err != nil {
		return err
	}

	connID, _ := random.String(8, random.Base62Chars)
	if err := stream.Handshake(connID, true); err != nil {
		return err
	}

	return nil
}

func (task *Task) findSuperNodes(ctx context.Context) (node.SuperNodes, error) {
	var superNodes node.SuperNodes

	mns, err := task.pastelClient.TopMasterNodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > task.Ticket.MaximumFee {
			continue
		}
		superNodes = append(superNodes, &node.SuperNode{
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
