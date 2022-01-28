package nftdownload

import (
	"bytes"
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

type downFile struct {
	file     []byte
	paslteID string
}

// Task is the task of downloading nft.
type NftDownloadTask struct {
	*common.WalletNodeTask

	MeshHandler *common.MeshHandler

	service *NftDownloadService
	Request *NftDownloadRequest

	files []downFile
	File  []byte

	err error
}

// Run starts the task
func (task *NftDownloadTask) Run(ctx context.Context) error {
	task.err = task.RunHelper(ctx, task.run, task.removeArtifacts)
	return nil
}

func (task *NftDownloadTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ttxid, err := task.service.pastelHandler.PastelClient.TicketOwnership(ctx, task.Request.Txid, task.Request.PastelID, task.Request.PastelIDPassphrase)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).WithField("pastelid", task.Request.PastelID).Error("Could not get ticket ownership")
		return errors.Errorf("get ticket ownership: %w", err)
	}

	// Sign current-timestamp with PsstelID passed in request
	timestamp := time.Now().Format(time.RFC3339)
	signature, err := task.service.pastelHandler.PastelClient.Sign(ctx, []byte(timestamp), task.Request.PastelID, task.Request.PastelIDPassphrase, "ed448")
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("timestamp", timestamp).WithField("pastelid", task.Request.PastelID).Error("Could not sign timestamp")
		return errors.Errorf("sign timestamp: %w", err)
	}

	if err = task.MeshHandler.ConnectToNSuperNodes(ctx, task.service.config.NumberSuperNodes); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)

	downloadErrs, err := task.Download(ctx, task.Request.Txid, timestamp, string(signature), ttxid)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).Error("Could not download files")
		task.UpdateStatus(common.StatusErrorDownloadFailed)
		return errors.Errorf("download files from supernodes: %w: %v", err, downloadErrs)
	}

	if len(task.files) < task.service.config.NumberSuperNodes {
		log.WithContext(ctx).WithField("DownloadedNodes", len(task.files)).Info("Not enough number of downloaded files")
		task.UpdateStatus(common.StatusErrorNotEnoughFiles)
		return errors.Errorf("could not download enough files from %d supernodes: %v", task.service.config.NumberSuperNodes, downloadErrs)
	}

	task.UpdateStatus(common.StatusDownloaded)

	// Disconnect all nodes after finished downloading.
	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	// Check files are the same
	err = task.MatchFiles()
	if err != nil {
		task.UpdateStatus(common.StatusErrorFilesNotMatch)
		return errors.Errorf("files are different between supernodes: %w", err)
	}

	// Store file to send to the caller
	task.File = task.files[0].file

	// Wait for all connections to disconnect.
	return nil
}

// Download downloads image from supernodes.
func (task *NftDownloadTask) Download(ctx context.Context, txid, timestamp, signature, ttxid string) ([]error, error) {
	group, _ := errgroup.WithContext(ctx)
	errChan := make(chan error, len(task.MeshHandler.Nodes))

	for _, someNode := range task.MeshHandler.Nodes {
		nftDownNode, ok := someNode.SuperNodeAPIInterface.(*NftDownloadNode)
		if !ok {
			//TODO: use assert here
			return nil, errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() error {
			file, subErr := nftDownNode.Download(ctx, txid, timestamp, signature, ttxid)
			if subErr != nil {
				log.WithContext(ctx).WithField("address", someNode.String()).WithError(subErr).Error("Could not download from supernode")
				errChan <- subErr
			} else {
				log.WithContext(ctx).WithField("address", someNode.String()).Info("Downloaded from supernode")
			}
			task.files = append(task.files, downFile{file: file, paslteID: someNode.PastelID()})
			return nil
		})
	}
	err := group.Wait()

	close(errChan)

	downloadErrors := []error{}
	for subErr := range errChan {
		downloadErrors = append(downloadErrors, subErr)
	}

	return downloadErrors, err
}

// MatchFiles matches files.
func (task *NftDownloadTask) MatchFiles() error {
	for _, someFile := range task.files[1:] {
		if !bytes.Equal(task.files[0].file, someFile.file) {
			return errors.Errorf("file of nodes %q and %q didn't match", task.files[0].paslteID, someFile.paslteID)
		}
	}
	return nil
}

// Error returns task err
func (task *NftDownloadTask) Error() error {
	return task.err
}

func (task *NftDownloadTask) removeArtifacts() {
}

// NewNftDownloadTask returns a new Task instance.
func NewNftDownloadTask(service *NftDownloadService, request *NftDownloadRequest) *NftDownloadTask {
	task := &NftDownloadTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		Request:        request,
	}

	task.MeshHandler = common.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &NftDownloadNodeMaker{},
		service.pastelHandler,
		request.PastelID, request.PastelIDPassphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.AcceptNodesTimeout, service.config.ConnectToNextNodeDelay,
	)

	return task
}
