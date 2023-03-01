package download

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/common"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

type downFile struct {
	file     []byte
	pastelID string
}

// NftDownloadingTask is the task of downloading nft.
type NftDownloadingTask struct {
	*common.WalletNodeTask

	MeshHandler *common.MeshHandler

	service *NftDownloadingService
	Request *NftDownloadingRequest

	files []downFile
	File  []byte

	mtx sync.Mutex
}

// Run starts the task
func (task *NftDownloadingTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *NftDownloadingTask) run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//validate download requested - ensure that the pastelID belongs either to the
	//  creator of the non-sold NFT or the latest buyer of the NFT
	var ttxid string
	if !task.service.pastelHandler.IsTicketPublic(ctx, task.Request.Txid, task.Request.Type) {
		ttxid, err = task.service.pastelHandler.PastelClient.TicketOwnership(ctx, task.Request.Txid, task.Request.PastelID, task.Request.PastelIDPassphrase)
		if err != nil && task.Request.Type != pastel.ActionTypeSense {
			log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).WithField("pastelid", task.Request.PastelID).Error("Could not get ticket ownership")
			return errors.Errorf("get ticket ownership: %w", err)
		}
	}

	// Sign current-timestamp with PastelID passed in request
	timestamp := time.Now().Format(time.RFC3339)
	signature, err := task.service.pastelHandler.PastelClient.Sign(ctx, []byte(timestamp), task.Request.PastelID, task.Request.PastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("timestamp", timestamp).WithField("pastelid", task.Request.PastelID).Error("Could not sign timestamp")
		return errors.Errorf("sign timestamp: %w", err)
	}

	if err = task.MeshHandler.ConnectToNSuperNodes(ctx, task.service.config.NumberSuperNodes); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank supernodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)
	//send download requests to ALL Supernodes, number defined by mesh handler's "minNumberSuperNodes" (really just set in a config file as NumberSuperNodes)
	//or max nodes that it was able to connect with defined by mesh handler config.UseMaxNodes
	downloadErrs, err := task.Download(ctx, task.Request.Txid, timestamp, string(signature), ttxid, task.Request.Type)
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
	n, err := task.MatchFiles()
	if err != nil {
		task.UpdateStatus(common.StatusErrorFilesNotMatch)
		return errors.Errorf("files are different between supernodes: %w", err)
	}

	// Store file to send to the caller
	task.File = task.files[n].file

	// Wait for all connections to disconnect.
	return nil
}

// Download downloads image from supernodes.
func (task *NftDownloadingTask) Download(ctx context.Context, txid, timestamp, signature, ttxid, ttype string) ([]error, error) {
	group, gctx := errgroup.WithContext(ctx)
	errChan := make(chan error, len(task.MeshHandler.Nodes))

	for _, someNode := range task.MeshHandler.Nodes {
		nftDownNode, ok := someNode.SuperNodeAPIInterface.(*NftDownloadingNode)
		if !ok {
			//TODO: use assert here
			return nil, errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			file, subErr := nftDownNode.Download(gctx, txid, timestamp, signature, ttxid, ttype)
			if subErr != nil {
				log.WithContext(gctx).WithField("address", someNode.String()).WithError(subErr).Error("Could not download from supernode")
				errChan <- subErr
			} else {
				log.WithContext(gctx).WithField("address", someNode.String()).Info("Downloaded from supernode")
			}

			func() {
				task.mtx.Lock()
				defer task.mtx.Unlock()
				task.files = append(task.files, downFile{file: file, pastelID: someNode.PastelID()})
			}()

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

// MatchFiles matches files. It loops through the files to find a file that matches any other two in the list
func (task *NftDownloadingTask) MatchFiles() (int, error) {
	for a, fileA := range task.files {
		matches := 0
		log.Debugf("file of node %s - content: %q", fileA.pastelID, fileA.file)

		for b, fileB := range task.files {
			if a == b {
				continue
			}

			if bytes.Equal(fileA.file, fileB.file) {
				matches++
			}
		}

		if matches >= 2 {
			return a, nil
		}
	}

	return 0, errors.Errorf("unable to find three matching files, no. of files - %d", len(task.files))
}

// Error returns task err
func (task *NftDownloadingTask) Error() error {
	return task.WalletNodeTask.Error()
}

func (task *NftDownloadingTask) removeArtifacts() {
}

// NewNftDownloadTask returns a new Task instance.
func NewNftDownloadTask(service *NftDownloadingService, request *NftDownloadingRequest) *NftDownloadingTask {
	task := common.NewWalletNodeTask(logPrefix)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &NftDownloadingNodeMaker{},
		PastelHandler: service.pastelHandler,
		NodeClient:    service.nodeClient,
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay: service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:   service.config.ConnectToNodeTimeout,
			AcceptNodesTimeout:     service.config.AcceptNodesTimeout,
			MinSNs:                 service.config.NumberSuperNodes,
			PastelID:               request.PastelID,
			Passphrase:             request.PastelIDPassphrase,
			UseMaxNodes:            true,
		},
	}

	return &NftDownloadingTask{
		WalletNodeTask: task,
		service:        service,
		Request:        request,
		MeshHandler:    common.NewMeshHandler(meshHandlerOpts),
	}
}
