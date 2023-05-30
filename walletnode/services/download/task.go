package download

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/common"

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
	tInfo, err := task.service.pastelHandler.GetTicketInfo(ctx, task.Request.Txid, task.Request.Type)
	if err != nil {
		return errors.Errorf("validate ticket: %w", err)
	}

	var ttxid string
	if !tInfo.IsTicketPublic {
		ttxid, err = task.service.pastelHandler.PastelClient.TicketOwnership(ctx, task.Request.Txid, task.Request.PastelID, task.Request.PastelIDPassphrase)
		if err != nil && task.Request.Type != pastel.ActionTypeSense {
			log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).WithField("pastelid", task.Request.PastelID).Error("Could not get ticket ownership")
			return errors.Errorf("get ticket ownership: %w", err)
		}
	}

	var skipNodes []string
	var nodesDone chan struct{}
	tries := 0
	for {

		addSkipNodes := func(goodNodes []string) {
			for _, node := range task.MeshHandler.Nodes {
				isGoodNode := false
				for _, goodNode := range goodNodes {
					if node.PastelID() == goodNode {
						isGoodNode = true
						break
					}
				}

				if !isGoodNode {
					skipNodes = append(skipNodes, node.PastelID())
				}
			}
		}

		tries++
		if tries > 6 {
			return errors.Errorf("task failed after too many tries")
		}

		// Sign current-timestamp with PastelID passed in request
		timestamp := time.Now().Format(time.RFC3339)
		signature, err := task.service.pastelHandler.PastelClient.Sign(ctx, []byte(timestamp), task.Request.PastelID, task.Request.PastelIDPassphrase, pastel.SignAlgorithmED448)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("timestamp", timestamp).WithField("pastelid", task.Request.PastelID).Error("Could not sign timestamp")
			continue
		}

		if err = task.MeshHandler.ConnectToNSuperNodes(ctx, task.service.config.NumberSuperNodes, skipNodes); err != nil {
			log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).Error("Could not connect to supernodes")
			addSkipNodes([]string{})
			continue
		}

		// supervise the connection to top rank supernodes
		// cancel any ongoing context if the connections are broken
		nodesDone = task.MeshHandler.ConnectionsSupervisor(ctx, cancel)
		//send download requests to ALL Supernodes, number defined by mesh handler's "minNumberSuperNodes" (really just set in a config file as NumberSuperNodes)
		//or max nodes that it was able to connect with defined by mesh handler config.UseMaxNodes
		downloadErrs, goodNodes, err := task.Download(ctx, task.Request.Txid, timestamp, string(signature), ttxid, task.Request.Type, tInfo.EstimatedDownloadTime)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).WithField("download errors", downloadErrs).Error("Could not download files")
			task.UpdateStatus(common.StatusErrorDownloadFailed)
			addSkipNodes(goodNodes)
			continue
		}

		task.UpdateStatus(common.StatusDownloaded)
		// Check files are the same
		n, goodNodes, err := task.MatchFiles()
		if err != nil {
			task.UpdateStatus(common.StatusErrorFilesNotMatch)
			addSkipNodes(goodNodes)
			continue
		}

		// Store file to send to the caller
		task.File = task.files[n].file
		break
	}

	// Disconnect all nodes after finished downloading.
	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	// Wait for all connections to disconnect.
	return nil
}

// Download downloads the file from supernodes.
func (task *NftDownloadingTask) Download(ctx context.Context, txid, timestamp, signature, ttxid, ttype string, timeout time.Duration) ([]error, []string, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(task.MeshHandler.Nodes))
	var goodNodes []string

	for _, someNode := range task.MeshHandler.Nodes {
		nftDownNode, ok := someNode.SuperNodeAPIInterface.(*NftDownloadingNode)
		if !ok {
			//TODO: use assert here
			return nil, goodNodes, errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}

		someNode := someNode
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Create a new context with a timeout for this goroutine
			goroutineCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			file, subErr := nftDownNode.Download(goroutineCtx, txid, timestamp, signature, ttxid, ttype)
			if subErr != nil {
				log.WithContext(ctx).WithField("address", someNode.String()).WithField("txid", txid).WithError(subErr).Error("Could not download from supernode")
				errChan <- subErr
			} else {
				log.WithContext(ctx).WithField("address", someNode.String()).WithField("txid", txid).Info("Downloaded from supernode")

				func() {
					task.mtx.Lock()
					defer task.mtx.Unlock()
					goodNodes = append(goodNodes, someNode.PastelID())
					task.files = append(task.files, downFile{file: file, pastelID: someNode.PastelID()})
				}()
			}
		}()
	}
	wg.Wait()

	close(errChan)

	downloadErrors := []error{}
	for subErr := range errChan {
		downloadErrors = append(downloadErrors, subErr)
	}

	return downloadErrors, goodNodes, nil
}

// MatchFiles matches files. It loops through the files to find a file that matches any other two in the list
func (task *NftDownloadingTask) MatchFiles() (int, []string, error) {
	var goodNodes []string

	if len(task.files) < 3 {
		for _, file := range task.files {
			goodNodes = append(goodNodes, file.pastelID)
		}

		return 0, goodNodes, errors.Errorf("number of files is less than 3, no. of files - %d", len(task.files))
	}

	for i := 0; i < len(task.files); i++ {
		matches := 0
		log.Debugf("file of node %s - content: %q", task.files[i].pastelID, task.files[i].file)

		for j := 0; j < len(task.files); j++ {
			if i == j {
				continue
			}

			if bytes.Equal(task.files[i].file, task.files[j].file) {
				goodNodes = append(goodNodes, task.files[i].pastelID)
				goodNodes = append(goodNodes, task.files[j].pastelID)

				matches++
			}
		}

		if matches >= 2 {
			return i, goodNodes, nil
		}
	}

	return 0, goodNodes, errors.Errorf("unable to find three matching files, no. of files - %d", len(task.files))
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
			UseMaxNodes:            false,
		},
	}

	return &NftDownloadingTask{
		WalletNodeTask: task,
		service:        service,
		Request:        request,
		MeshHandler:    common.NewMeshHandler(meshHandlerOpts),
	}
}
