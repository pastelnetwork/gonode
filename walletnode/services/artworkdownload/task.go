package artworkdownload

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload/node"
)

// Task is the task of downloading artwork.
type NftDownloadTask struct {
	*common.WalletNodeTask
	*Service

	Ticket *Ticket
	File   []byte
	err    error
}

// Run starts the task
func (task *NftDownloadTask) Run(ctx context.Context) error {
	task.err = task.RunHelper(ctx, task.run, task.removeArtifacts)
	return nil
}

// Error returns the error after running task
func (task *NftDownloadTask) Error() error {
	return task.err
}

func (task *NftDownloadTask) run(ctx context.Context) error {
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
		task.UpdateStatus(common.StatusErrorNotEnoughSuperNode)
		return errors.New("unable to find enough supernodes")
	}

	//var connectedNodes node.List
	//var errs error
	secInfo := &alts.SecInfo{
		PastelID:   task.Ticket.PastelID,
		PassPhrase: task.Ticket.PastelIDPassphrase,
		Algorithm:  "ed448",
	}

	// Send download request to all top supernodes.
	var nodes node.List

	downloadErrs, err := topNodes.Download(ctx, task.Ticket.Txid, timestamp, string(signature), ttxid, task.config.connectToNodeTimeout, secInfo)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("txid", task.Ticket.Txid).Error("Could not download files")
		task.UpdateStatus(common.StatusErrorDownloadFailed)
		return errors.Errorf("download files from supernodes: %w", err)
	}

	nodes = topNodes.Active()

	if len(nodes) < task.config.NumberSuperNodes {
		log.WithContext(ctx).WithField("DownloadedNodes", len(nodes)).Info("Not enough number of downloaded node")
		task.UpdateStatus(common.StatusErrorNotEnoughFiles)
		return errors.Errorf("could not download enough files from %d supernodes: %v", task.config.NumberSuperNodes, downloadErrs)
	}

	task.UpdateStatus(common.StatusDownloaded)

	// Disconnect supernodes that did not return file.
	topNodes.DisconnectInactive()

	// Check files are the same
	err = nodes.MatchFiles()
	if err != nil {
		task.UpdateStatus(common.StatusErrorFilesNotMatch)
		return errors.Errorf("files are different between supernodes: %w", err)
	}

	// Store file to send to the caller
	task.File = nodes.File()

	// Disconnect all noded after finished downloading.
	nodes.Disconnect()

	// Wait for all connections to disconnect.
	return nil
}

func (task *NftDownloadTask) pastelTopNodes(ctx context.Context) (node.List, error) {
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

func (task *NftDownloadTask) removeArtifacts() {
}

// NewNftDownloadTask returns a new Task instance.
func NewNftDownloadTask(service *Service, ticket *Ticket) *NftDownloadTask {
	return &NftDownloadTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		Service:        service,
		Ticket:         ticket,
	}
}
