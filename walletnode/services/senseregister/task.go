package senseregister

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister/node"
)

// Task represents a task for registering an artwork.
type Task struct {
	task.Task
	*Service

	Request *Request

	// top 3 MNs from pastel network
	nodes node.List

	// Task data
	creatorBlockHeight int
	creatorBlockHash   string
}

func NewTask(service *Service, Ticket *Request) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Request: Ticket,
	}
}

// Run starts the task.
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debug("States updated")
	})

	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")

	if err := task.run(ctx); err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Error("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

func (task *Task) run(ctx context.Context) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get top 10MNs from pastel network
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return errors.Errorf("get top 10MNs from pastel network: %v", err)
	}

	if len(topNodes) < task.config.NumberSuperNodes {
		return errors.New("unable to find enough Supernodes from pastel network")
	}

	// Get current block height & hash
	if err := task.getBlock(ctx); err != nil {
		return errors.Errorf("get block: %v", err)
	}

	// Connect to top nodes to validate their infor
	err = task.validateMNsInfo(ctx, topNodes)
	if err != nil {
		return errors.Errorf("validate MNs info: %v", err)
	}

	return nil
}

// pastelTopNodes - get 10 top MNs from pastel network
func (task *Task) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.ExtKey == "" || mn.ExtAddress == "" {
			continue
		}

		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
		_, err = task.pastelClient.FindTicketByID(ctx, mn.ExtKey)
		if err != nil {
			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
			continue
		}
		nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// validateMNsInfo - validate MNs info, until found at least 3 valid MNs
func (task *Task) validateMNsInfo(ctx context.Context, nnondes node.List) error {
	var nodes node.List
	count := 0

	secInfo := &alts.SecInfo{
		PastelID:   task.Request.AppPastelID,
		PassPhrase: task.Request.AppPastelIDPassphrase,
		Algorithm:  "ed448",
	}

	for _, node := range nnondes {
		if err := node.Connect(ctx, task.config.ConnectToNodeTimeout, secInfo); err != nil {
			continue
		}

		count++
		nodes = append(nodes, node)
		if count == task.config.NumberSuperNodes {
			break
		}
	}

	// Close all connected connnections
	nnondes.DisconnectAll()

	if count < task.config.NumberSuperNodes {
		return errors.Errorf("validate %d Supernodes from pastel network", task.config.NumberSuperNodes)
	}

	task.nodes = nodes
	return nil
}

// determine current block height & hash of it
func (task *Task) getBlock(ctx context.Context) error {
	// Get block num
	blockNum, err := task.pastelClient.GetBlockCount(ctx)
	task.creatorBlockHeight = int(blockNum)
	if err != nil {
		return errors.Errorf("get block num: %w", err)
	}

	// Get block hash string
	blockInfo, err := task.pastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		return errors.Errorf("get block info blocknum=%d: %w", blockNum, err)
	}

	// Decode hash string to byte
	task.creatorBlockHash = blockInfo.Hash
	if err != nil {
		return errors.Errorf("convert hash string %s to bytes: %w", blockInfo.Hash, err)
	}

	return nil
}
