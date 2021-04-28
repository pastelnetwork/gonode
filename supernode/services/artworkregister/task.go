package artworkregister

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister/state"
)

var taskID uint32

// Task is the task of registering new artwork.
type Task struct {
	*Service

	ID     int
	ConnID string
	State  *state.State

	nodes  map[string]*SuperNode
	doneCh chan struct{}
	sync.Mutex
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	return nil
}

// Cancel closes the task.
func (task *Task) Cancel() {
	select {
	case <-task.Done():
		return
	default:
		close(task.doneCh)
	}
}

// Done returns a channel when the task is canceled.
func (task *Task) Done() <-chan struct{} {
	return task.doneCh
}

// RegisterSecondaryNode registers secondary node
func (task *Task) RegisterSecondaryNode(ctx context.Context, nodeKey string) error {
	task.Lock()
	defer task.Unlock()

	if err := task.requiredStatus(state.StatusWaitForSecondaryNodes); err != nil {
		return err
	}

	if _, ok := task.nodes[nodeKey]; ok {
		return errors.Errorf("node %q is already registered", nodeKey)
	}

	node, err := task.findNode(ctx, nodeKey)
	if err != nil {
		return err
	}
	task.nodes[nodeKey] = node

	fmt.Println("primary", nodeKey)
	return nil
}

// AcceptSecondaryNodes accepts secondary nodes
func (task *Task) AcceptSecondaryNodes(ctx context.Context) error {
	if err := task.requiredStatus(state.StatusTaskStarted); err != nil {
		return err
	}

	task.State.Update(state.NewStatus(state.StatusWaitForSecondaryNodes))

	// TODO: wait for two connections

	task.State.Update(state.NewStatus(state.StatusAcceptedSecondaryNodes))

	return nil
}

// ConnectPrimaryNode connects to primary node
func (task *Task) ConnectPrimaryNode(ctx context.Context, nodeKey string) error {
	if err := task.requiredStatus(state.StatusTaskStarted); err != nil {
		return err
	}

	node, err := task.findNode(ctx, nodeKey)
	if err != nil {
		return err
	}

	fmt.Println(node.Address)
	// TODO: connect to primary node

	task.State.Update(state.NewStatus(state.StatusConnectedToPrimaryNode))

	return nil
}

// RegisterConnection registers a new connetion
func (task *Task) RegisterConnection(ctx context.Context, connID string) error {
	task.ConnID = connID

	return nil
}

func (task *Task) findNode(ctx context.Context, nodeKey string) (*SuperNode, error) {
	masterNodes, err := task.pastel.TopMasterNodes(ctx)
	if err != nil {
		return nil, err
	}

	for _, masterNode := range masterNodes {
		if masterNode.ExtKey != nodeKey {
			continue
		}
		node := &SuperNode{
			Address: masterNode.ExtAddress,
			Key:     masterNode.ExtKey,
			Fee:     masterNode.Fee,
		}
		return node, nil
	}

	return nil, errors.Errorf("node %q not found", nodeKey)
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
	return &Task{
		Service: service,
		ID:      int(atomic.AddUint32(&taskID, 1)),
		State:   state.New(state.NewStatus(state.StatusTaskStarted)),
		nodes:   make(map[string]*SuperNode),
	}
}
