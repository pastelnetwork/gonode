package externaldupedetetion

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/externaldupedetection/node"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Request *Request

	// information of nodes
	nodes node.List

	// task data to create RegArt ticket
	creatorBlockHeight   int
	creatorBlockHash     string
	fingerprintAndScores *pastel.FingerAndScores
	fingerprintsHash     []byte
	fingerprint          []byte
	datahash             []byte
	rqids                []string

	// TODO: call cNodeAPI to get the reall signature instead of the fake one
	fingerprintSignature []byte

	// TODO: need to update rqservice code to return the following info
	rqSymbolIDFiles map[string][]byte
	rqEncodeParams  rqnode.EncoderParameters

	// TODO: call cNodeAPI to get the following info
	// blockTxID    string
	burnTxID     string
	regEDDTxID   string
	detectionFee int64

	// ticket
	creatorSignature []byte
	ticket           *pastel.ExDDTicket
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

	// probe image for rareness, nsfw and seen score
	if err := task.probeImage(ctx); err != nil {
		return errors.Errorf("failed to probe image: %w", err)
	}

	// upload image and thumbnail coordinates to supernode(s)
	if err := task.uploadImage(ctx); err != nil {
		return errors.Errorf("failed to upload image: %w", err)
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
	if err := task.checkCurrentBalance(ctx, float64(task.detectionFee)); err != nil {
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
	if err := task.waitTxidValid(newCtx, task.regEDDTxID, int64(task.config.RegArtTxMinConfirmations), task.config.RegArtTxTimeout, 15*time.Second); err != nil {
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

func (task *Task) checkCurrentBalance(ctx context.Context, detectionFee float64) error {
	if detectionFee == 0 {
		return errors.Errorf("detection fee cannot be 0.0")
	}

	if detectionFee > task.Request.MaximumFee {
		return errors.Errorf("detection fee is to expensive - maximum-fee (%f) < registration-fee(%f)", task.Request.MaximumFee, detectionFee)
	}

	balance, err := task.pastelClient.GetBalance(ctx, task.Request.SendingAddress)
	if err != nil {
		return errors.Errorf("failed to get current balance of address(%s): %w", task.Request.SendingAddress, err)
	}

	if balance < detectionFee {
		return errors.Errorf("not enough PSL - balance(%f) < registration-fee(%f)", balance, detectionFee)
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

func (task *Task) encodeFingerprint(ctx context.Context, fingerprint []byte, img *artwork.File) error {
	// Sign fingerprint
	ed448PubKey := []byte(task.Request.PastelID)
	ed448Signature, err := task.pastelClient.Sign(ctx, fingerprint, task.Request.PastelID, task.Request.PastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign fingerprint with pastelId and pastelPassphrase failed %w", err)
	}

	// TODO: Should be replaced with real data from the Pastel API.
	ticket, err := task.pastelClient.FindTicketByID(ctx, task.Request.PastelID)
	if err != nil {
		return errors.Errorf("failed to find register ticket of artist pastel id:%s :%w", task.Request.PastelID, err)
	}
	pqSignature, err := task.pastelClient.Sign(ctx, fingerprint, task.Request.PastelID, task.Request.PastelIDPassphrase, pastel.SignAlgorithmLegRoast)
	if err != nil {
		return errors.Errorf("failed to sign fingerprint with legroats: %w", err)
	}

	// Encode data to the image.
	encSig := qrsignature.New(
		qrsignature.Fingerprint(fingerprint),
		qrsignature.PostQuantumSignature(pqSignature),
		qrsignature.PostQuantumPubKey([]byte(ticket.Signature)),
		qrsignature.Ed448Signature(ed448Signature),
		qrsignature.Ed448PubKey(ed448PubKey),
	)
	if err := img.Encode(encSig); err != nil {
		return err
	}
	task.fingerprintSignature = ed448Signature

	// TODO: check with the change in legroast
	// Decode data from the image, to make sure their integrity.
	// decSig := qrsignature.New()
	// copyImage, _ := img.Copy()
	// if err := copyImage.Decode(decSig); err != nil {
	// 	return err
	// }

	// if !bytes.Equal(fingerprint, decSig.Fingerprint()) {
	// 	return errors.Errorf("fingerprints do not match, original len:%d, decoded len:%d\n", len(fingerprint), len(decSig.Fingerprint()))
	// }
	// if !bytes.Equal(pqSignature, decSig.PostQuantumSignature()) {
	// 	return errors.Errorf("post quantum signatures do not match, original len:%d, decoded len:%d\n", len(pqSignature), len(decSig.PostQuantumSignature()))
	// }
	// if !bytes.Equal(pqPubKey, decSig.PostQuantumPubKey()) {
	// 	return errors.Errorf("post quantum public keys do not match, original len:%d, decoded len:%d\n", len(pqPubKey), len(decSig.PostQuantumPubKey()))
	// }
	// if !bytes.Equal(ed448Signature, decSig.Ed448Signature()) {
	// 	return errors.Errorf("ed448 signatures do not match, original len:%d, decoded len:%d\n", len(ed448Signature), len(decSig.Ed448Signature()))
	// }
	// if !bytes.Equal(ed448PubKey, decSig.Ed448PubKey()) {
	// 	return errors.Errorf("ed448 public keys do not match, original len:%d, decoded len:%d\n", len(ed448PubKey), len(decSig.Ed448PubKey()))
	// }
	return nil
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

func (task *Task) createEDDTicket(_ context.Context) error {
	if task.fingerprint == nil {
		return errEmptyFingerprints
	}
	if task.fingerprintsHash == nil {
		return errEmptyFingerprintsHash
	}
	if task.fingerprintSignature == nil {
		return errEmptyFingerprintSignature
	}
	if task.datahash == nil {
		return errEmptyDatahash
	}
	if task.rqids == nil {
		return errEmptyRaptorQSymbols
	}

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	task.ticket = &pastel.ExDDTicket{
		Version:             1,
		BlockHash:           task.creatorBlockHash,
		BlockNum:            task.creatorBlockHeight,
		Identifier:          task.Request.RequestIdentifier,
		ImageHash:           task.Request.ImageHash,
		MaximumFee:          task.Request.MaximumFee,
		SendingAddress:      task.Request.SendingAddress,
		EffectiveTotalFee:   float64(task.detectionFee),
		SupernodesSignature: []map[string]string{},
	}
	return nil
}

func (task *Task) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeExDDTicket(task.ticket)
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
	return task.pastelClient.RegisterActTicket(ctx, task.regEDDTxID, task.creatorBlockHeight, task.detectionFee, task.Request.PastelID, task.Request.PastelIDPassphrase)
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

func (task *Task) probeImage(ctx context.Context) error {
	log.WithContext(ctx).WithField("filename", task.Request.Image.Name()).Debugf("probe image")

	// Send thumbnail to supernodes for probing.
	if err := task.nodes.ProbeImage(ctx, task.Request.Image); err != nil {
		return errors.Errorf("failed to probe image: %w", err)
	}

	// Match fingerprints received from supernodes.
	if err := task.nodes.MatchFingerprintAndScores(); err != nil {
		task.UpdateStatus(StatusErrorFingerprintsNotMatch)
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}
	task.fingerprintAndScores = task.nodes.FingerAndScores()
	task.UpdateStatus(StatusImageProbed)

	// As we are going to store the the compressed figerprint to kamedila
	// so we calculated the hash base on the compressed fingerprint print also
	fingerprintsHash, err := common.Sha3256hash(task.fingerprintAndScores.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("failed to hash zstd commpressed fingerprints %w", err)
	}
	task.fingerprintsHash = []byte(base58.Encode(fingerprintsHash))

	// Decompress the fingerprint as bytes for later use
	task.fingerprint, err = zstd.Decompress(nil, task.fingerprintAndScores.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("failed to decompress finger")
	}
	return nil
}

func (task *Task) uploadImage(ctx context.Context) error {
	img1, err := task.Request.Image.Copy()
	if err != nil {
		return errors.Errorf("copy image to encode failed %w", err)
	}

	log.WithContext(ctx).WithField("FileName", img1.Name()).Debugf("final image")
	if err := task.encodeFingerprint(ctx, task.fingerprint, img1); err != nil {
		return errors.Errorf("encode image with fingerprint %w", err)
	}

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
	buf, err := pastel.EncodeExDDTicket(task.ticket)
	if err != nil {
		return errors.Errorf("failed to marshal ticket %w", err)
	}
	log.Debug(string(buf))

	if err := task.nodes.UploadSignedTicket(ctx, buf, task.creatorSignature, task.rqSymbolIDFiles, task.rqEncodeParams); err != nil {
		return errors.Errorf("failed to upload signed ticket %w", err)
	}

	// check if fee is over-expection
	task.detectionFee = task.nodes.DetectionFee()

	if task.detectionFee > int64(task.Request.MaximumFee) {
		return errors.Errorf("fee too high: detection fee %d, maximum fee %d", task.detectionFee, int64(task.Request.MaximumFee))
	}
	task.detectionFee = task.nodes.DetectionFee()

	return nil
}

func (task *Task) preburnedDetectionFee(ctx context.Context) error {
	if task.detectionFee <= 0 {
		return errors.Errorf("invalid registration fee")
	}

	burnedAmount := float64(task.detectionFee) / 20
	burnTxid, err := task.pastelClient.SendFromAddress(ctx, task.Request.SendingAddress, task.config.BurnAddress, burnedAmount)
	if err != nil {
		return errors.Errorf("failed to burn 20 percent of transaction fee %w", err)
	}
	task.burnTxID = burnTxid
	log.WithContext(ctx).Debugf("preburn txid: %s", task.burnTxID)

	if err := task.nodes.SendPreBurnedFeeTxID(ctx, task.burnTxID); err != nil {
		return errors.Errorf("failed to send pre-burn-txid: %s to supernode(s): %w", task.burnTxID, err)
	}
	task.regEDDTxID = task.nodes.RegEDDTicketID()
	if task.regEDDTxID == "" {
		return errors.Errorf("regEDDTxID is empty")
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
