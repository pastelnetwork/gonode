package artworkregister

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/state"
	"golang.org/x/sync/errgroup"
)

const (
	connectToNextNodeDelay = time.Millisecond * 200
	acceptNodesTimeout     = connectToNextNodeDelay * 10 // waiting 2 seconds (10 supernodes) for secondary nodes to be accpeted by primary nodes.
	connectToNodeTimeout   = time.Second * 1
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
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID))

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")

	if err := task.run(ctx); err != nil {
		task.State.Update(ctx, state.NewMessage(state.StatusTaskRejected))
		log.WithContext(ctx).WithError(err).Warnf("Task is rejected")
		return nil
	}

	task.State.Update(ctx, state.NewMessage(state.StatusTaskCompleted))
	log.WithContext(ctx).Debugf("Task is completed")
	return nil
}

func (task *Task) run(ctx context.Context) error {
	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.State.Update(ctx, state.NewMessage(state.StatusErrorTooLowFee))
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Ticket.MaximumFee)
	}

	superNodes, err := task.findSuperNodes(ctx)
	if err != nil {
		return err
	}
	if len(superNodes) < task.config.NumberSuperNodes {
		task.State.Update(ctx, state.NewMessage(state.StatusErrorTooLowFee))
		return errors.New("not found enough available SuperNodes with acceptable storage fee: %f")
	}

	ctx, connCancel := context.WithCancel(ctx)
	defer connCancel()

	for i, primaryNode := range superNodes {
		secondaryNodes := make(node.SuperNodes, len(superNodes)-1)
		copy(secondaryNodes, superNodes[:i])
		copy(secondaryNodes[i:], superNodes[i+1:])

		if err := task.connect(ctx, primaryNode, secondaryNodes); err != nil {
			connCancel()
			continue
		}
		break
	}
	task.State.Update(ctx, state.NewMessage(state.StatusConnected))

	<-ctx.Done()

	return nil
}

// connect establishes communication between all supernodes.
func (task *Task) connect(ctx context.Context, primaryNode *node.SuperNode, secondaryNodes node.SuperNodes) error {
	connID, _ := random.String(8, random.Base62Chars)

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

	stream, err := task.connectStream(ctx, primaryNode.Address, connID, true)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() (err error) {
		defer errors.Recover(func(recErr error) { err = recErr })
		defer nextConnCancel()

		acceptCtx, acceptCancel := context.WithTimeout(ctx, acceptNodesTimeout)
		defer acceptCancel()

		nodes, err := stream.AcceptedNodes(acceptCtx)
		if err != nil {
			return err
		}

		for _, node := range nodes {
			log.WithContext(ctx).Debugf("Primary accepted %q secondary node", node.Key)
			// TODO: disconnect all not accepted connections by primary supernode.
		}
		return nil
	})

	for _, node := range secondaryNodes {
		node := node

		select {
		case <-nextConnCtx.Done():
			return nil
		case <-time.After(connectToNextNodeDelay):
			group.Go(func() (err error) {
				defer errors.Recover(func(recErr error) { err = recErr })

				stream, err := task.connectStream(ctx, node.Address, connID, false)
				if err != nil {
					return nil
				}

				if err := stream.ConnectTo(ctx, primaryNode.Key); err != nil {
					return nil
				}
				log.WithContext(ctx).Debugf("Seconary %s connected to primary", node.Address)

				return nil
			})
		}

	}
	return group.Wait()
}

func (task *Task) connectStream(ctx context.Context, address, connID string, isPrimary bool) (node.RegisterArtowrk, error) {
	connCtx, connCancel := context.WithTimeout(ctx, connectToNodeTimeout)
	defer connCancel()

	conn, err := task.nodeClient.Connect(connCtx, address)
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-conn.Done():
			// TODO: reject task
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

func (task *Task) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.pastelClient.StorageFee(ctx)
	if err != nil {
		return false, err
	}
	return fee.NetworkFee <= task.Ticket.MaximumFee, nil
}

func (task *Task) findSuperNodes(ctx context.Context) (node.SuperNodes, error) {
	var superNodes node.SuperNodes

	mns, err := task.pastelClient.MasterNodesTop(ctx)
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
	taskID, _ := random.String(8, random.Base62Chars)

	return &Task{
		Service: service,
		ID:      taskID,
		Ticket:  Ticket,
		State:   state.New(state.NewMessage(state.StatusTaskStarted)),
	}
}
