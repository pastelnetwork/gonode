package externalstorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/externaldupedetection/node"
	"golang.org/x/crypto/sha3"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Request *Request

	// information of nodes
	nodes node.List

	// task data to create RegArt ticket
	creatorBlockHeight int
	creatorBlockHash   string
	datahash           []byte
	rqids              []string

	// TODO: need to update rqservice code to return the following info
	rqSymbolIDFiles map[string][]byte
	rqEncodeParams  rqnode.EncoderParameters

	// TODO: call cNodeAPI to get the following info
	// blockTxID    string
	burnTxID               string
	regExternalStorageTxID string
	storageFee             int64

	// ticket
	creatorSignature []byte
	ticket           *pastel.ExternalStorageTicket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debugf("States updated")
	})

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")

	defer task.removeArtifacts()
	if err := task.run(ctx); err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Warnf("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.UpdateStatus(ErrorInsufficientFee)
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Request.MaximumFee)
	}

	// Retrieve supernodes with highest ranks.
	if err := task.connectToTopRankNodes(ctx); err != nil {
		return errors.Errorf("failed to connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := make(chan struct{})
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return task.nodes.WaitConnClose(ctx, nodesDone)
	})

	// get block height + hash
	if err := task.getBlock(ctx); err != nil {
		return errors.Errorf("failed to get current block heigth %w", err)
	}

	// upload image and thumbnail coordinates to supernode(s)
	if err := task.uploadImage(ctx); err != nil {
		return errors.Errorf("failed to upload image: %w", err)
	}

	// connect to rq serivce to get rq symbols identifier
	if err := task.genRQIdentifiersFiles(ctx); err != nil {
		return errors.Errorf("gen RaptorQ symbols' identifiers failed %w", err)
	}

	if err := task.createEDDTicket(ctx); err != nil {
		return errors.Errorf("failed to create ticket %w", err)
	}

	// sign ticket with artist signature
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("failed to sign NFT ticket %w", err)
	}

	// send signed ticket to supernodes to calculate registration fee
	if err := task.sendSignedTicket(ctx); err != nil {
		return errors.Errorf("failed to send signed NFT ticket: %w", err)
	}

	// validate if address has enough psl
	if err := task.checkCurrentBalance(ctx, float64(task.storageFee)); err != nil {
		return errors.Errorf("failed to check current balance: %w", err)
	}

	// send preburn-txid to master node(s)
	// master node will create reg-art ticket and returns transaction id
	if err := task.preburnedDetectionFee(ctx); err != nil {
		return errors.Errorf("faild to pre-burnt ten percent of registration fee: %w", err)
	}

	log.WithContext(ctx).Debugf("close connection")
	close(nodesDone)
	for i := range task.nodes {
		if err := task.nodes[i].Connection.Close(); err != nil {
			log.WithContext(ctx).WithError(err).Debugf("failed to close connection to node %s", task.nodes[i].PastelID())
		}
	}

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.waitTxidValid(newCtx, task.regExternalStorageTxID, int64(task.config.RegArtTxMinConfirmations), task.config.RegArtTxTimeout, 15*time.Second); err != nil {
		return errors.Errorf("failed to wait reg-nft ticket valid %w", err)
	}

	// activate reg-art ticket at previous step
	actTxid, err := task.registerActTicket(newCtx)
	if err != nil {
		return errors.Errorf("failed to register act ticket %w", err)
	}
	log.Debugf("reg-act-txid: %s", actTxid)

	// Wait until actTxid is valid
	err = task.waitTxidValid(newCtx, actTxid, int64(task.config.RegActTxMinConfirmations), task.config.RegActTxTimeout, 15*time.Second)
	if err != nil {
		return errors.Errorf("failed to reg-act ticket valid %w", err)
	}
	log.Debugf("reg-act-tixd is confirmed")

	return nil
}

func (task *Task) checkCurrentBalance(ctx context.Context, storageFee float64) error {
	if storageFee == 0 {
		return errors.Errorf("detection fee cannot be 0.0")
	}

	if storageFee > task.Request.MaximumFee {
		return errors.Errorf("detection fee is to expensive - maximum-fee (%f) < registration-fee(%f)", task.Request.MaximumFee, storageFee)
	}

	balance, err := task.pastelClient.GetBalance(ctx, task.Request.SendingAddress)
	if err != nil {
		return errors.Errorf("failed to get current balance of address(%s): %w", task.Request.SendingAddress, err)
	}

	if balance < storageFee {
		return errors.Errorf("not enough PSL - balance(%f) < registration-fee(%f)", balance, storageFee)
	}

	return nil
}

func (task *Task) waitTxidValid(ctx context.Context, txID string, expectedConfirms int64, timeout time.Duration, interval time.Duration) error {
	log.WithContext(ctx).Debugf("Need %d confirmation for txid %s, timeout %v", expectedConfirms, txID, timeout)

	checkTimeout := func(checked chan<- struct{}) {
		time.Sleep(timeout)
		close(checked)
	}

	timeoutCh := make(chan struct{})
	go checkTimeout(timeoutCh)

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done %w", ctx.Err())
		case <-time.After(interval):
			checkConfirms := func() error {
				subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				result, err := task.pastelClient.GetRawTransactionVerbose1(subCtx, txID)
				log.WithContext(ctx).Debugf("getrawtransaction result: %s %v", txID, result)
				if err != nil {
					return errors.Errorf("failed to get transaction %w", err)
				}

				if result.Confirmations >= expectedConfirms {
					return nil
				}

				return errors.Errorf("confirmations not meet: expected %d, got %d", expectedConfirms, result.Confirmations)
			}

			err := checkConfirms()
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("check confirmations failed")
			} else {
				return nil
			}
		case <-timeoutCh:
			return errors.Errorf("timeout")
		}
	}
}

// meshNodes establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
	var meshNodes node.List
	secInfo := &alts.SecInfo{
		PastelID:   task.Request.PastelID,
		PassPhrase: task.Request.PastelIDPassphrase,
		Algorithm:  "ed448",
	}

	primary := nodes[primaryIndex]
	if err := primary.Connect(ctx, task.config.ConnectToNodeTimeout, secInfo); err != nil {
		return nil, err
	}
	if err := primary.Session(ctx, true); err != nil {
		return nil, err
	}

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

	// FIXME: ugly hack here. Need to make the Node and List to be safer
	secondariesMtx := &sync.Mutex{}
	var secondaries node.List
	go func() {
		for i, node := range nodes {
			node := node

			if i == primaryIndex {
				continue
			}

			select {
			case <-nextConnCtx.Done():
				return
			case <-time.After(task.config.connectToNextNodeDelay):
				go func() {
					defer errors.Recover(log.Fatal)

					if err := node.Connect(ctx, task.config.ConnectToNodeTimeout, secInfo); err != nil {
						return
					}
					if err := node.Session(ctx, false); err != nil {
						return
					}
					go func() {
						secondariesMtx.Lock()
						defer secondariesMtx.Unlock()
						secondaries.Add(node)
					}()

					if err := node.ConnectTo(ctx, primary.PastelID(), primary.SessID()); err != nil {
						return
					}
					log.WithContext(ctx).Debugf("Seconary %q connected to primary", node)
				}()
			}
		}
	}()

	acceptCtx, acceptCancel := context.WithTimeout(ctx, task.config.acceptNodesTimeout)
	defer acceptCancel()

	accepted, err := primary.AcceptedNodes(acceptCtx)
	if err != nil {
		return nil, err
	}

	primary.SetPrimary(true)
	meshNodes.Add(primary)

	secondariesMtx.Lock()
	defer secondariesMtx.Unlock()
	for _, pastelID := range accepted {
		log.WithContext(ctx).Debugf("Primary accepted %q secondary node", pastelID)

		node := secondaries.FindByPastelID(pastelID)
		if node == nil {
			return nil, errors.New("not found accepted node")
		}
		meshNodes.Add(node)
	}
	return meshNodes, nil
}

func (task *Task) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.pastelClient.StorageNetworkFee(ctx)
	if err != nil {
		return false, err
	}
	return fee <= task.Request.MaximumFee, nil
}

func (task *Task) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > task.Request.MaximumFee {
			continue
		}
		nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// determine current block height & hash of it
func (task *Task) getBlock(ctx context.Context) error {
	// Get block num
	blockNum, err := task.pastelClient.GetBlockCount(ctx)
	task.creatorBlockHeight = int(blockNum)
	if err != nil {
		return errors.Errorf("failed to get block num: %w", err)
	}

	// Get block hash string
	blockInfo, err := task.pastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		return errors.Errorf("failed to get block info with given block num %d: %w", blockNum, err)
	}

	// Decode hash string to byte
	task.creatorBlockHash = blockInfo.Hash
	if err != nil {
		return errors.Errorf("failed to convert hash string %s to bytes: %w", blockInfo.Hash, err)
	}

	return nil
}

func (task *Task) genRQIdentifiersFiles(ctx context.Context) error {
	log.Debugf("Connect to %s", task.config.RaptorQServiceAddress)
	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return errors.Errorf("failed to connect to raptorQ service %w", err)
	}
	defer conn.Close()

	content, err := task.Request.Image.Bytes()
	if err != nil {
		return errors.Errorf("failed to read image contents: %w", err)
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: task.Service.config.RqFilesDir,
	})

	log.WithContext(ctx).Debugf("Image hash %x", sha3.Sum256(content))
	// FIXME :
	// - check format of artis block hash should be base58 or not
	encodeInfo, err := rqService.EncodeInfo(ctx, content, task.config.NumberRQIDSFiles, task.creatorBlockHash, task.Request.PastelID)
	if err != nil {
		return errors.Errorf("failed to generate RaptorQ symbols' identifiers %w", err)
	}

	files := make(map[string][]byte)
	for _, rawSymbolIDFile := range encodeInfo.SymbolIDFiles {
		name, content, err := task.convertToSymbolIDFile(ctx, rawSymbolIDFile)
		if err != nil {
			task.UpdateStatus(StatusErrorGenRaptorQSymbolsFailed)
			return errors.Errorf("failed to create rqids file %w", err)
		}
		files[name] = content
	}

	if len(files) < 1 {
		return errors.Errorf("nuber of raptorq symbol identifiers files must be greater than 1")
	}
	for k := range files {
		log.WithContext(ctx).Debugf("symbol file identifier %s", k)
		task.rqids = append(task.rqids, k)
	}

	task.rqSymbolIDFiles = files
	task.rqEncodeParams = encodeInfo.EncoderParam
	return nil
}

func (task *Task) convertToSymbolIDFile(ctx context.Context, rawFile rqnode.RawSymbolIDFile) (string, []byte, error) {
	var content []byte
	appendStr := func(b []byte, s string) []byte {
		b = append(b, []byte(s)...)
		b = append(b, '\n')
		return b
	}

	content = appendStr(content, rawFile.ID)
	content = appendStr(content, rawFile.BlockHash)
	content = appendStr(content, rawFile.PastelID)
	for _, id := range rawFile.SymbolIdentifiers {
		content = appendStr(content, id)
	}

	// log.WithContext(ctx).Debugf("%s", content)
	signature, err := task.pastelClient.Sign(ctx, content, task.Request.PastelID, task.Request.PastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return "", nil, errors.Errorf("failed to sign identifiers file: %w", err)
	}
	content = append(content, signature...)

	// read all content of the data b
	// compress the contente of the file
	compressedData, err := zstd.CompressLevel(nil, content, 22)
	if err != nil {
		return "", nil, errors.Errorf("failed to compress: ")
	}
	hasher := sha3.New256()
	src := bytes.NewReader(compressedData)
	if _, err := io.Copy(hasher, src); err != nil {
		return "", nil, errors.Errorf("failed to hash identifiers file %w", err)
	}
	return base58.Encode(hasher.Sum(nil)), compressedData, nil
}

func (task *Task) createEDDTicket(_ context.Context) error {
	if task.datahash == nil {
		return errEmptyDatahash
	}
	if task.rqids == nil {
		return errEmptyRaptorQSymbols
	}

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	task.ticket = &pastel.ExternalStorageTicket{
		Version:             1,
		BlockHash:           task.creatorBlockHash,
		BlockNum:            task.creatorBlockHeight,
		Identifier:          task.Request.RequestIdentifier,
		ImageHash:           task.Request.ImageHash,
		MaximumFee:          task.Request.MaximumFee,
		SendingAddress:      task.Request.SendingAddress,
		EffectiveTotalFee:   float64(task.storageFee),
		SupernodesSignature: []map[string]string{},
	}
	return nil
}

func (task *Task) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeExternalStorageTicket(task.ticket)
	if err != nil {
		return errors.Errorf("failed to encode ticket %w", err)
	}

	task.creatorSignature, err = task.pastelClient.Sign(ctx, data, task.Request.PastelID, task.Request.PastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("failed to sign ticket %w", err)
	}
	return nil
}

func (task *Task) registerActTicket(ctx context.Context) (string, error) {
	return task.pastelClient.RegisterActTicket(ctx, task.regExternalStorageTxID, task.creatorBlockHeight, task.storageFee, task.Request.PastelID, task.Request.PastelIDPassphrase)
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Request) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Request: Ticket,
	}
}

func (task *Task) connectToTopRankNodes(ctx context.Context) error {
	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return errors.Errorf("failed to call masternode top: %w", err)
	}

	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(ErrorInsufficientFee)
		return errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

	// For testing
	// // TODO: Remove this when releaset because in localnet there is no need to start 10 gonode
	// // and chase the log
	// pinned := []string{"127.0.0.1:4444", "127.0.0.1:4445", "127.0.0.1:4446"}
	// pinnedMn := make([]*node.Node, 0)
	// for i := range topNodes {
	// 	for j := range pinned {
	// 		if topNodes[i].Address() == pinned[j] {
	// 			pinnedMn = append(pinnedMn, topNodes[i])
	// 		}
	// 	}
	// }
	// log.WithContext(ctx).Debugf("%v", pinnedMn)

	// Try to create mesh of supernodes, connecting to all supernodes in a different sequences.
	var nodes node.List
	var errs error
	for primaryRank := range topNodes {
		nodes, err = task.meshNodes(ctx, topNodes, primaryRank)
		if err != nil {
			if errors.IsContextCanceled(err) {
				return err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warnf("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(nodes) < task.config.NumberSuperNodes {
		return errors.Errorf("Could not create a mesh of %d nodes: %w", task.config.NumberSuperNodes, errs)
	}

	// Activate supernodes that are in the mesh.
	nodes.Activate()
	// Disconnect supernodes that are not involved in the process.
	topNodes.DisconnectInactive()

	// Cancel context when any connection is broken.
	task.UpdateStatus(StatusConnected)
	task.nodes = nodes

	return nil
}

func (task *Task) uploadImage(ctx context.Context) error {
	img1, err := task.Request.Image.Copy()
	if err != nil {
		return errors.Errorf("copy image to encode failed %w", err)
	}

	log.WithContext(ctx).WithField("FileName", img1.Name()).Debugf("final image")

	imgBytes, err := img1.Bytes()
	if err != nil {
		return errors.Errorf("failed to convert image to byte stream %w", err)
	}

	if task.datahash, err = common.Sha3256hash(imgBytes); err != nil {
		return errors.Errorf("failed to hash encoded image: %w", err)
	}

	// Upload image with pqgsinganature and its thumb to supernodes
	if err := task.nodes.UploadImage(ctx, img1); err != nil {
		return errors.Errorf("upload encoded image failed %w", err)
	}
	task.UpdateStatus(StatusImageUploaded)

	return nil
}

func (task *Task) sendSignedTicket(ctx context.Context) error {
	buf, err := pastel.EncodeExternalStorageTicket(task.ticket)
	if err != nil {
		return errors.Errorf("failed to marshal ticket %w", err)
	}
	log.Debug(string(buf))

	if err := task.nodes.UploadSignedTicket(ctx, buf, task.creatorSignature, task.rqSymbolIDFiles, task.rqEncodeParams); err != nil {
		return errors.Errorf("failed to upload signed ticket %w", err)
	}

	// check if fee is over-expection
	task.storageFee = task.nodes.DetectionFee()

	if task.storageFee > int64(task.Request.MaximumFee) {
		return errors.Errorf("fee too high: detection fee %d, maximum fee %d", task.storageFee, int64(task.Request.MaximumFee))
	}
	task.storageFee = task.nodes.DetectionFee()

	return nil
}

func (task *Task) preburnedDetectionFee(ctx context.Context) error {
	if task.storageFee <= 0 {
		return errors.Errorf("invalid registration fee")
	}

	burnedAmount := float64(task.storageFee) / 20
	burnTxid, err := task.pastelClient.SendFromAddress(ctx, task.Request.SendingAddress, task.config.BurnAddress, burnedAmount)
	if err != nil {
		return errors.Errorf("failed to burn 20 percent of transaction fee %w", err)
	}
	task.burnTxID = burnTxid
	log.WithContext(ctx).Debugf("preburn txid: %s", task.burnTxID)

	if err := task.nodes.SendPreBurnedFeeTxID(ctx, task.burnTxID); err != nil {
		return errors.Errorf("failed to send pre-burn-txid: %s to supernode(s): %w", task.burnTxID, err)
	}
	task.regExternalStorageTxID = task.nodes.RegEDDTicketID()
	if task.regExternalStorageTxID == "" {
		return errors.Errorf("regExternalStorageTxID is empty")
	}

	return nil
}

func (task *Task) removeArtifacts() {
	removeFn := func(file *artwork.File) {
		if file != nil {
			log.Debugf("remove file: %s", file.Name())
			if err := file.Remove(); err != nil {
				log.Debugf("failed to remove filed: %s", err.Error())
			}
		}
	}
	if task.Request != nil {
		removeFn(task.Request.Image)
	}
}
