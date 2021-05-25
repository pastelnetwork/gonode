package artworkregister

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"golang.org/x/sync/errgroup"
)

const (
	connectToNextNodeDelay = time.Millisecond * 200
	acceptNodesTimeout     = connectToNextNodeDelay * 10 // waiting 2 seconds (10 supernodes) for secondary nodes to be accpeted by primary nodes.
	connectToNodeTimeout   = time.Second * 1

	thumbnailWidth, thumbnailHeight = 224, 224
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Ticket *Ticket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status.String()).Debugf("States updated")
	})

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")

	if err := task.run(ctx); err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithError(err).Warnf("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	log.WithContext(ctx).Debugf("Task is completed")
	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.UpdateStatus(ErrorInsufficientFee)
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Ticket.MaximumFee)
	}

	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return err
	}
	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(ErrorInsufficientFee)
		return errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

	var nodes Nodes
	for primaryRank := range topNodes {
		nodes, err = task.meshNodes(ctx, topNodes, primaryRank)
		if err == nil {
			break
		}
		log.WithContext(ctx).WithError(err).Warnf("Could not get mesh of nodes")
	}
	nodes.activate()
	topNodes.disconnectInactive()

	group, _ := errgroup.WithContext(ctx)
	for _, node := range nodes {
		node := node
		group.Go(func() (err error) {
			defer errors.Recover(func(recErr error) { err = recErr })
			defer cancel()

			select {
			case <-ctx.Done():
				return errors.Errorf("task was canceled")
			case <-node.conn.Done():
				return errors.Errorf("%q unexpectedly closed the connection", node.Address)
			}
		})
	}
	task.UpdateStatus(StatusConnected)

	log.WithContext(ctx).WithField("filename", task.Ticket.Image.Name()).Debugf("Copy image")
	thumbnail, err := task.Ticket.Image.Copy()
	if err != nil {
		return err
	}

	log.WithContext(ctx).WithField("filename", thumbnail.Name()).Debugf("Resize image to thumbnail %dx%d", thumbnailWidth, thumbnailHeight)
	if err := thumbnail.ResizeImage(thumbnailWidth, thumbnailHeight); err != nil {
		return err
	}

	if err := nodes.sendImage(ctx, thumbnail); err != nil {
		return err
	}

	task.UpdateStatus(StatusImageUploaded)

	<-ctx.Done()

	return group.Wait()
}

// meshNodes establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes Nodes, primaryIndex int) (Nodes, error) {
	var meshNodes Nodes

	primary := nodes[primaryIndex]
	if err := primary.connect(ctx); err != nil {
		return nil, err
	}
	if err := primary.Session(ctx, true); err != nil {
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
					defer errors.Recover(log.Fatal)

					if err := node.connect(ctx); err != nil {
						return
					}
					if err := node.Session(ctx, false); err != nil {
						return
					}
					secondaries.add(node)

					if err := node.ConnectTo(ctx, primary.PastelID, primary.SessID()); err != nil {
						return
					}
					log.WithContext(ctx).Debugf("Seconary %s connected to primary", node.Address)
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

func (task *Task) pastelTopNodes(ctx context.Context) (Nodes, error) {
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
			Address:  mn.ExtAddress,
			PastelID: mn.ExtKey,
		})
	}

	return nodes, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Ticket) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Ticket:  Ticket,
	}
}
