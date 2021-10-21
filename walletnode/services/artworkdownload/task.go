package artworkdownload

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload/node"
)

// Task is the task of downloading artwork.
type Task struct {
	task.Task
	*Service

	Ticket *Ticket
	File   []byte
	err    error
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debug("States updated")
	})

	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")

	if err := task.run(ctx); err != nil {
		task.err = err
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Warn("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

// Error returns the error after running task
func (task *Task) Error() error {
	return task.err
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ttxid, err := task.pastelClient.TicketOwnership(ctx, task.Ticket.Txid, task.Ticket.PastelID, task.Ticket.PastelIDPassphrase)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("txid", task.Ticket.Txid).WithField("pastelid", task.Ticket.PastelID).Error("Could not get ticket ownership")
		return errors.Errorf("get ticket ownership: %w", err)
	}

	// Sign current-timestamp with PsstelID passed in request
	timestamp := time.Now().Format(time.RFC3339)
	signature, err := task.pastelClient.Sign(ctx, []byte(timestamp), task.Ticket.PastelID, task.Ticket.PastelIDPassphrase, "ed448")
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("timestamp", timestamp).WithField("pastelid", task.Ticket.PastelID).Error("Could not sign timestamp")
		return errors.Errorf("sign timestamp: %w", err)
	}

	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return errors.Errorf("get top rank nodes: %w", err)
	}
	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorNotEnoughSuperNode)
		return errors.New("unable to find enough supernodes")
	}

	var connectedNodes node.List
	var errs error
	secInfo := &alts.SecInfo{
		PastelID:   task.Ticket.PastelID,
		PassPhrase: task.Ticket.PastelIDPassphrase,
		Algorithm:  "ed448",
	}
	// Connect to top supernodes
	for _, node := range topNodes {
		if err := node.Connect(ctx, task.config.connectToNodeTimeout, secInfo); err == nil {
			log.WithContext(ctx).WithError(err).WithField("address", node.String()).WithField("pastelid", node.PastelID()).Debug("Connected to supernode")
			connectedNodes.Add(node)
		} else {
			log.WithContext(ctx).WithError(err).WithField("address", node.String()).WithField("pastelid", node.PastelID()).Debug("Could not connect to supernode")
			errs = errors.Append(errs, err)
		}
	}

	if len(connectedNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorNotEnoughSuperNode)
		return errors.Errorf("could not connect enough %d supernodes: %w", task.config.NumberSuperNodes, errs)
	}
	task.UpdateStatus(StatusConnected)

	// Send download request to all top supernodes.
	var nodes node.List
	var downloadErrs error
	err = connectedNodes.Download(ctx, task.Ticket.Txid, timestamp, string(signature), ttxid, downloadErrs)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("txid", task.Ticket.Txid).Error("Could not download files")
		task.UpdateStatus(StatusErrorDownloadFailed)
		return errors.Errorf("download files from supernodes: %w", err)
	}

	nodes = connectedNodes.Active()

	if len(nodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorNotEnoughFiles)
		return errors.Errorf("could not download enough files from %d supernodes: %w", task.config.NumberSuperNodes, downloadErrs)
	}
	task.UpdateStatus(StatusDownloaded)

	// Disconnect supernodes that did not return file.
	topNodes.DisconnectInactive()

	// Check files are the same
	err = nodes.MatchFiles()
	if err != nil {
		task.UpdateStatus(StatusErrorFilesNotMatch)
		return errors.Errorf("files are different between supernodes: %w", err)
	}

	// Store file to send to the caller
	task.File = nodes.File()

	// Disconnect all noded after finished downloading.
	nodes.Disconnect()

	// Wait for all connections to disconnect.
	return nil
}

func (task *Task) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		nodes.Add(node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, ticket *Ticket) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Ticket:  ticket,
	}
}
