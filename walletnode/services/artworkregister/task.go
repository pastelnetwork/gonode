package artworkregister

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/state"
	"golang.org/x/sync/errgroup"
)

// Task is the task of registering new artwork.
type Task struct {
	*Service

	ID     string
	State  *state.State
	Ticket *Ticket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, task.ID))

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")

	if err := task.run(ctx); err != nil {
		task.State.Update(state.NewMessage(state.StatusTaskRejected))
		log.WithContext(ctx).WithField("error", err).Warnf("Task is rejected")
		return nil
	}

	task.State.Update(state.NewMessage(state.StatusTaskCompleted))
	log.WithContext(ctx).Debugf("Task is completed")
	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	superNodes, err := task.findSuperNodes(ctx)
	if err != nil {
		return err
	}

	if err := task.connect(ctx, superNodes[0], superNodes[1:3]); err != nil {
		return err
	}

	<-ctx.Done()

	// if len(superNodes) < task.config.NumberSuperNodes {
	// 	task.State.Update(state.NewMessage(state.StatusErrorTooLowFee))
	// 	return NewTaskError(errors.Errorf("not found enough available SuperNodes with acceptable storage fee", task.config.NumberSuperNodes))
	// }

	return nil
}

func (task *Task) connect(ctx context.Context, primaryNode *node.SuperNode, secondaryNodes node.SuperNodes) error {
	ctx, cancel := context.WithCancel(ctx)

	connID, _ := random.String(8, random.Base62Chars)
	stream, err := task.openStream(ctx, primaryNode.Address, connID, true)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		defer errors.Recover(func(recErr error) { err = recErr })

		nodes, err := stream.SecondaryNodes(ctx)
		if err != nil {
			cancel()
			return err
		}

		for _, node := range nodes {
			log.WithContext(ctx).Debugf("Primary accepted %q secondary node", node.Key)
		}

		return nil
	})

	for _, node := range secondaryNodes {
		node := node

		group.Go(func() (err error) {
			defer errors.Recover(func(recErr error) { err = recErr })

			stream, err := task.openStream(ctx, node.Address, connID, false)
			if err != nil {
				return err
			}

			if err := stream.ConnectToPrimary(ctx, primaryNode.Key); err != nil {
				return err
			}
			log.WithContext(ctx).Debugf("Seconary %s connected to primary", node.Address)

			return nil
		})

	}

	return group.Wait()

}

func (task *Task) openStream(ctx context.Context, address, connID string, isPrimary bool) (node.RegisterArtowrk, error) {
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
	taskID, _ := random.String(8, random.Base62Chars)

	return &Task{
		Service: service,
		ID:      taskID,
		Ticket:  Ticket,
		State:   state.New(state.NewMessage(state.StatusTaskStarted)),
	}
}
