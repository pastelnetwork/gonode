package artworkregister

import (
	"context"
	"fmt"
	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister/state"
)

// Task is the task of registering new artwork.
type Task struct {
	*Service

	ID    string
	State *state.State

	acceptMu    sync.Mutex
	accpetNodes Nodes

	connectNode *Node

	actionCh chan func(ctx context.Context) error

	doneMu sync.Mutex
	doneCh chan struct{}
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	defer task.Cancel()

	for {
		select {
		case <-ctx.Done():
			return nil
		case action := <-task.actionCh:
			if err := action(ctx); err != nil {
				return err
			}
		}
	}
}

// Cancel stops the task, which causes all connections associated with that task to be closed.
func (task *Task) Cancel() {
	task.doneMu.Lock()
	defer task.doneMu.Unlock()

	select {
	case <-task.Done():
		return
	default:
		log.WithPrefix(fmt.Sprintf("%s-%s", logPrefix, task.ID)).Debug("Task canceled")
		close(task.doneCh)
	}
}

// Done returns a channel when the task is canceled.
func (task *Task) Done() <-chan struct{} {
	return task.doneCh
}

func (task *Task) context(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID))
}

// Handshake is handshake wallet to supernode
func (task *Task) Handshake(ctx context.Context, isPrimary bool) error {
	ctx = task.context(ctx)

	if err := task.requiredStatus(state.StatusTaskStarted); err != nil {
		return err
	}

	if isPrimary {
		log.WithContext(ctx).Debugf("Acts as primary node")
		task.State.Update(ctx, state.NewStatus(state.StatusHandshakePrimaryNode))
		return nil
	}

	log.WithContext(ctx).Debugf("Acts as secondary node")
	task.State.Update(ctx, state.NewStatus(state.StatusHandshakeSecondaryNode))
	return nil
}

// AcceptedNodes waits for connection supernodes, as soon as there is the required amount returns them.
func (task *Task) AcceptedNodes(ctx context.Context) (Nodes, error) {
	ctx = task.context(ctx)

	if err := task.requiredStatus(state.StatusHandshakePrimaryNode); err != nil {
		return nil, err
	}
	log.WithContext(ctx).Debugf("Waiting for supernodes to connect")

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-task.State.Updated():
		if err := task.requiredStatus(state.StatusAcceptedNodes); err != nil {
			return nil, err
		}
		return task.accpetNodes, nil
	}
}

// HandshakeNode accepts secondary node
func (task *Task) HandshakeNode(ctx context.Context, nodeID string) error {
	ctx = task.context(ctx)

	task.acceptMu.Lock()
	defer task.acceptMu.Unlock()

	if err := task.requiredStatus(state.StatusHandshakePrimaryNode); err != nil {
		return err
	}

	if node := task.accpetNodes.ByID(nodeID); node != nil {
		return errors.Errorf("node %q is already registered", nodeID)
	}

	node, err := task.pastelNodeByExtKey(ctx, nodeID)
	if err != nil {
		return err
	}
	task.accpetNodes.Add(node)

	log.WithContext(ctx).WithField("nodeID", nodeID).Debugf("Accept secondary node")

	if len(task.accpetNodes) >= task.config.NumberConnectedNodes {
		task.State.Update(ctx, state.NewStatus(state.StatusAcceptedNodes))
	}
	return nil
}

// ConnectTo connects to primary node
func (task *Task) ConnectTo(_ context.Context, nodeID, sessID string) error {
	if err := task.requiredStatus(state.StatusHandshakeSecondaryNode); err != nil {
		return err
	}

	task.actionCh <- func(ctx context.Context) error {
		return task.connectTo(ctx, nodeID, sessID)
	}
	return nil
}

func (task *Task) connectTo(ctx context.Context, nodeID, sessID string) error {
	ctx = task.context(ctx)

	node, err := task.pastelNodeByExtKey(ctx, nodeID)
	if err != nil {
		return err
	}

	if err := node.connect(ctx); err != nil {
		return err
	}

	go func() {
		select {
		case <-task.Done():
			node.conn.Close()
		case <-node.conn.Done():
			task.Cancel()
		}
	}()

	if err := node.Handshake(ctx, task.config.PastelID, sessID); err != nil {
		return err
	}

	task.connectNode = node
	task.State.Update(ctx, state.NewStatus(state.StatusConnectedToNode))
	return nil
}

func (task *Task) pastelNodeByExtKey(ctx context.Context, nodeID string) (*Node, error) {
	masterNodes, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}

	for _, masterNode := range masterNodes {
		if masterNode.ExtKey != nodeID {
			continue
		}
		node := &Node{
			client:  task.Service.nodeClient,
			ID:      masterNode.ExtKey,
			Address: masterNode.ExtAddress,
		}
		return node, nil
	}

	return nil, errors.Errorf("node %q not found", nodeID)
}

func (task *Task) requiredStatus(statusType state.StatusType) error {
	latest := task.State.Latest()
	if latest == nil {
		return errors.New("not found latest status")
	}

	if latest.Type != statusType {
		return errors.Errorf("wrong order, current task status %q, ", latest.Type)
	}
	return nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service) *Task {
	taskID, _ := random.String(8, random.Base62Chars)

	return &Task{
		Service:  service,
		ID:       taskID,
		State:    state.New(state.NewStatus(state.StatusTaskStarted)),
		doneCh:   make(chan struct{}),
		actionCh: make(chan func(ctx context.Context) error),
	}
}
