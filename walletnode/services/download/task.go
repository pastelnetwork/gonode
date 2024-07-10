package download

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/common"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"golang.org/x/sync/errgroup"
)

const (
	requiredSuccessfulCount = 1
	concurrentDownloads     = 1
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

	Filename       string
	TicketDataHash []byte
	matchHash      bool

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
	task.Filename = tInfo.Filename
	if task.Filename == "" {
		task.Filename = task.Request.Txid
	}

	hash := tInfo.DataHash
	if task.Request.HashOnly {
		hash, err = utils.Sha3256hash(tInfo.DataHash)
		if err != nil {
			return errors.Errorf("hash data hash: %w", err)
		}
	}
	task.TicketDataHash = hash

	// inetentionally verbose
	if task.Request.Type != pastel.ActionTypeSense {
		task.matchHash = true
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

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic in defer: %v", r)
		}

		task.close(ctx)

		// Disconnect all nodes after finished downloading.
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
	}()

	tries := 0
	for {

		addSkipNodes := func(badNodes []string) {
			for _, node := range task.MeshHandler.Nodes {
				for _, badNode := range badNodes {
					if node.PastelID() == badNode {
						skipNodes = append(skipNodes, node.PastelID())
						break
					}
				}
			}
		}

		tries++
		if tries > 3 {
			task.UpdateStatus(common.StatusErrorDownloadFailed)
			return errors.Errorf("task failed after too many tries")
		}

		log.WithContext(ctx).WithField("txid", task.Request.Txid).WithField("skip nodes", skipNodes).WithField("tries", tries).Info("Connecting to supernodes")
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
		// Sign current-timestamp with PastelID passed in request

		downloadErrs, badNodes, err := task.Download(ctx, task.Request.Txid, ttxid, task.Request.Type, tInfo.EstimatedDownloadTime)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).WithField("download errors", downloadErrs).Error("Could not download files")
			addSkipNodes(badNodes)
			continue
		}

		task.UpdateStatus(common.StatusDownloaded)
		// Check files are the same
		n, badNodes, err := task.MatchFiles()
		if err != nil {
			addSkipNodes(badNodes)
			continue
		}

		// Store file to send to the caller
		task.File = task.files[n].file
		if task.matchHash {
			gotHash, err := utils.Sha3256hash(task.File)
			if err != nil {
				// add all 3 nodes in skip nodes because all 3 agreed on a hash that doesn't match with original hash
				for _, file := range task.files {
					badNodes = append(badNodes, file.pastelID)
				}
				addSkipNodes(badNodes)

				log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).Error("Could not hash file")
				continue
			}

			if !bytes.Equal(gotHash, task.TicketDataHash) {
				for _, file := range task.files {
					badNodes = append(badNodes, file.pastelID)
				}
				addSkipNodes(badNodes)

				log.WithContext(ctx).WithField("txid", task.Request.Txid).Error("Hashes do not match")
				continue
			}
		}

		break
	}

	// Wait for all connections to disconnect.
	return nil
}

// Download downloads the file from supernodes.
func (task *NftDownloadingTask) Download(cctx context.Context, txid, ttxid, ttype string, timeout time.Duration) ([]error, []string, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(task.MeshHandler.Nodes))
	var badNodes []string
	var badMtx sync.Mutex

	// Create a new context with a cancel function
	ctx, cancel := context.WithCancel(cctx)
	defer cancel()

	// Create a buffered channel to communicate successful downloads
	successes := make(chan struct{}, requiredSuccessfulCount)

	// Semaphore channel to limit concurrent downloads to 3
	semaphore := make(chan struct{}, concurrentDownloads)

	defer close(successes)

	for _, someNode := range task.MeshHandler.Nodes {
		nftDownNode, ok := someNode.SuperNodeAPIInterface.(*NftDownloadingNode)
		if !ok {
			return nil, badNodes, errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}

		someNode := someNode
		wg.Add(1)

		semaphore <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore when goroutine finishes

			// Create a new context with a timeout for this goroutine
			goroutineCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			timestamp := time.Now().UTC().Format(time.RFC3339)
			signature, err := task.service.pastelHandler.PastelClient.Sign(ctx, []byte(timestamp), task.Request.PastelID, task.Request.PastelIDPassphrase, pastel.SignAlgorithmED448)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("txid", task.Request.Txid).WithField("timestamp", timestamp).WithField("pastelid", task.Request.PastelID).Error("Could not sign timestamp")
				errChan <- err
			}
			log.WithContext(ctx).WithField("address", someNode.String()).WithField("txid", txid).Info("Downloading from supernode")
			file, subErr := nftDownNode.Download(goroutineCtx, txid, timestamp, string(signature), ttxid, ttype, task.Request.HashOnly)
			if subErr != nil || len(file) == 0 {
				if subErr == nil {
					subErr = errors.New("empty file")
				}

				if !strings.Contains(subErr.Error(), "context canceled") {
					log.WithContext(ctx).WithField("address", someNode.String()).WithField("txid", txid).WithError(subErr).Error("Could not download from supernode")

					func() {
						badMtx.Lock()
						defer badMtx.Unlock()

						badNodes = append(badNodes, someNode.PastelID())
					}()
				}

				errChan <- subErr
			} else {
				log.WithContext(ctx).WithField("address", someNode.String()).WithField("txid", txid).WithField("len", len(file)).Info("Downloaded from supernode")

				func() {
					task.mtx.Lock()
					defer task.mtx.Unlock()
					task.files = append(task.files, downFile{file: file, pastelID: someNode.PastelID()})
				}()

				// Signal a successful download
				successes <- struct{}{}
			}
		}()
	}

	// Wait for 3 successful downloads, then cancel the context
	go func() {
		for i := 0; i < requiredSuccessfulCount; i++ {
			select {
			case <-successes:
			case <-ctx.Done():
				return
			}
		}
		cancel()
	}()
	wg.Wait()

	close(errChan)

	downloadErrors := []error{}
	for subErr := range errChan {
		downloadErrors = append(downloadErrors, subErr)
	}

	return downloadErrors, badNodes, nil
}

// MatchFiles matches files. It loops through the files to find a file that matches any other two in the list
func (task *NftDownloadingTask) MatchFiles() (int, []string, error) {
	var badNodes []string

	if len(task.files) < requiredSuccessfulCount {
		return 0, badNodes, errors.Errorf("number of files is less than 3, no. of files - %d", len(task.files))
	}

	for i := 0; i < len(task.files); i++ {
		matches := 0
		log.Debugf("file of node %s - content: %q", task.files[i].pastelID, task.files[i].file)

		for j := 0; j < len(task.files); j++ {
			if i == j {
				continue
			}

			if bytes.Equal(task.files[i].file, task.files[j].file) {
				matches++
			}
		}

		if matches >= 2 {
			return i, badNodes, nil
		}

		badNodes = append(badNodes, task.files[i].pastelID)
	}

	return 0, badNodes, errors.Errorf("unable to find three matching files, no. of files - %d", len(task.files))
}

// Error returns task err
func (task *NftDownloadingTask) Error() error {
	return task.WalletNodeTask.Error()
}

func (task *NftDownloadingTask) removeArtifacts() {
}

func (task *NftDownloadingTask) close(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, someNode := range task.MeshHandler.Nodes {
		nftDownNode, ok := someNode.SuperNodeAPIInterface.(*NftDownloadingNode)
		if !ok {
			log.WithContext(ctx).WithField("address", someNode.String()).Error("node is not NftDownloadingNode")
			continue
		}

		someNode := someNode
		g.Go(func() error {
			if err := nftDownNode.Close(); err != nil {
				log.WithContext(ctx).WithField("address", someNode.String()).WithField("message", err.Error()).Debug("Connection already closed")
			}

			return nil
		})
	}

	return g.Wait()
}

// NewNftDownloadTask returns a new Task instance.
func NewNftDownloadTask(service *NftDownloadingService, request *NftDownloadingRequest) *NftDownloadingTask {
	task := common.NewWalletNodeTask(logPrefix, service.historyDB)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &NftDownloadingNodeMaker{},
		PastelHandler: service.pastelHandler,
		NodeClient:    service.nodeClient,
		LogRequestID:  request.Txid,
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay:        service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:          service.config.DownloadConnectToNodeTimeout,
			AcceptNodesTimeout:            service.config.AcceptNodesTimeout,
			MinSNs:                        service.config.NumberSuperNodes,
			PastelID:                      request.PastelID,
			Passphrase:                    request.PastelIDPassphrase,
			UseMaxNodes:                   true,
			RequireSNAgreementOnMNTopList: false,
		},
	}

	return &NftDownloadingTask{
		WalletNodeTask: task,
		service:        service,
		Request:        request,
		MeshHandler:    common.NewMeshHandler(meshHandlerOpts),
	}
}
