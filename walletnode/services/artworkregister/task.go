package artworkregister

import (
	"context"
	"sync/atomic"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/state"
	"golang.org/x/sync/errgroup"
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
		task.State.Update(state.NewMessage(state.StatusTaskRejected))
		log.WithField("error", err).Warnf("Task %d is rejected", task.ID)
		return nil
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

	if err := task.connect(ctx, superNodes[0], superNodes[1:3]); err != nil {
		return err
	}

	// if len(superNodes) < task.config.NumberSuperNodes {
	// 	task.State.Update(state.NewMessage(state.StatusErrorTooLowFee))
	// 	return NewTaskError(errors.Errorf("not found enough available SuperNodes with acceptable storage fee", task.config.NumberSuperNodes))
	// }

	return nil
}

func (task *Task) connect(ctx context.Context, primaryNode *node.SuperNode, secondaryNodes node.SuperNodes) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	connID, _ := random.String(8, random.Base62Chars)

	stream, err := task.nodeStream(ctx, primaryNode.Address, connID, true)
	if err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() (err error) {
		defer errors.Recover(func(recErr error) { err = recErr })

		nodes, err := stream.PrimaryAcceptSecondary(ctx)
		if err != nil {
			return err
		}

		for _, node := range nodes {
			log.WithContext(ctx).Debugf("Primary accepted %d secondary node", node.Key)
		}

		return nil
	})

	for _, node := range secondaryNodes {
		node := node

		group.Go(func() (err error) {
			defer errors.Recover(func(recErr error) { err = recErr })

			stream, err := task.nodeStream(ctx, node.Address, connID, false)
			if err != nil {
				return err
			}

			if err := stream.SecondaryConnectToPrimary(ctx, primaryNode.Key); err != nil {
				return err
			}
			log.WithContext(ctx).Debugf("Seconary %s connected to primary", node.Address)

			return nil
		})
	}

	return group.Wait()

}

func (task *Task) nodeStream(ctx context.Context, address, connID string, isPrimary bool) (node.RegisterArtowrk, error) {
	conn, err := task.nodeClient.Connect(ctx, address)
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-conn.Done():
			// TODO Remove from the `nodes` list
		}
	}()

	stream, err := conn.RegisterArtowrk(ctx)
	if err != nil {
		return nil, err
	}

	if err := stream.Handshake(ctx, connID, isPrimary); err != nil {
		return nil, err
	}
	return stream, nil
}

func (task *Task) findSuperNodes(ctx context.Context) (node.SuperNodes, error) {
	var superNodes node.SuperNodes

	mns, err := task.pastelClient.TopMasterNodes(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("Could not get pastel top masternodes")
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
