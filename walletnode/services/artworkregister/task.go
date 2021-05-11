package artworkregister

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
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
		task.State.Update(ctx, state.NewStatus(state.StatusTaskRejected))
		log.WithContext(ctx).WithError(err).Warnf("Task is rejected")
		return nil
	}

	task.State.Update(ctx, state.NewStatus(state.StatusTaskCompleted))
	log.WithContext(ctx).Debugf("Task is completed")
	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.State.Update(ctx, state.NewStatus(state.StatusErrorTooLowFee))
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Ticket.MaximumFee)
	}

	topNodes, err := task.findTopNodes(ctx)
	if err != nil {
		return err
	}
	if len(topNodes) < task.config.NumberSuperNodes {
		task.State.Update(ctx, state.NewStatus(state.StatusErrorTooLowFee))
		return errors.New("not found enough available SuperNodes with acceptable storage fee: %f")
	}

	var nodes Nodes
	for primaryRank := range topNodes {
		nodes, err = task.meshNodes(ctx, topNodes, primaryRank)
		if err == nil {
			break
		}
	}
	nodes.activate()
	topNodes.disconnectInactive()

	group, _ := errgroup.WithContext(ctx)
	for _, node := range nodes {
		group.Go(func() (err error) {
			defer cancel()

			select {
			case <-ctx.Done():
				return errors.Errorf("task was canceled")
			case <-node.conn.Done():
				return errors.Errorf("%q closed the connection", node.address)
			}
		})
	}

	task.State.Update(ctx, state.NewStatus(state.StatusConnected))

	<-ctx.Done()

	return group.Wait()
}

// connect establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes Nodes, primaryIndex int) (Nodes, error) {
	var meshNodes Nodes

	connID, _ := random.String(8, random.Base62Chars)

	primary := nodes[primaryIndex]
	if err := primary.connect(ctx, connID, true); err != nil {
		return nil, err
	}

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

	var secondaries Nodes
	go func() {
		for i, node := range nodes {
			node := node

			if i == primaryIndex {
				continue
			}

			select {
			case <-nextConnCtx.Done():
				return
			case <-time.After(connectToNextNodeDelay):
				go func() {
					defer errors.Recover(errors.CheckErrorAndExit)

					if err := node.connect(ctx, connID, false); err != nil {
						return
					}
					secondaries.add(node)

					if err := node.ConnectTo(ctx, primary.pastelID); err != nil {
						return
					}
					log.WithContext(ctx).Debugf("Seconary %s connected to primary", node.address)
				}()
			}
		}
	}()

	acceptCtx, acceptCancel := context.WithTimeout(ctx, acceptNodesTimeout)
	defer acceptCancel()

	accepted, err := primary.AcceptedNodes(acceptCtx)
	if err != nil {
		return nil, err
	}

	meshNodes.add(primary)
	for _, pastelID := range accepted {
		log.WithContext(ctx).Debugf("Primary accepted %q secondary node", pastelID)

		node := secondaries.findByPastelID(pastelID)
		if node == nil {
			return nil, errors.New("not found accepted node")
		}
		meshNodes.add(node)
	}
	return meshNodes, nil
}

func (task *Task) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.pastelClient.StorageFee(ctx)
	if err != nil {
		return false, err
	}
	return fee.NetworkFee <= task.Ticket.MaximumFee, nil
}

func (task *Task) findTopNodes(ctx context.Context) (Nodes, error) {
	var nodes Nodes

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > task.Ticket.MaximumFee {
			continue
		}
		nodes = append(nodes, &Node{
			client:   task.Service.nodeClient,
			address:  mn.ExtAddress,
			pastelID: mn.ExtKey,
		})
	}

	return nodes, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Ticket) *Task {
	taskID, _ := random.String(8, random.Base62Chars)

	return &Task{
		Service: service,
		ID:      taskID,
		Ticket:  Ticket,
		State:   state.New(state.NewStatus(state.StatusTaskStarted)),
	}
}
