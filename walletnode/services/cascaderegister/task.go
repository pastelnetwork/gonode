package cascaderegister

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/cascaderegister/node"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Request *Request

	// information of 3 nodes
	nodes node.List

	// task data to create RegArt ticket
	creatorBlockHeight   int
	creatorBlockHash     string
	fingerprintAndScores *pastel.DDAndFingerprints
	fingerprint          []byte
	datahash             []byte

	// signatures from SN1, SN2 & SN3 over dd_and_fingerprints data from dd-server
	signatures          [][]byte
	ddAndFingerprintsIc uint32
	// ddAndFingerprintsIDs are Base58(SHA3_256(compressed(Base64URL(dd_and_fingerprints).
	// Base64URL(signatureSN1).Base64URL(signatureSN2).Base64URL(signatureSN3).dd_and_fingerprints_ic)))
	ddAndFingerprintsIDs []string
	// ddAndFpFile is Base64(compressed(Base64(dd_and_fingerprints).Base64(signatureSN1).Base64(signatureSN3)))
	ddAndFpFile []byte

	// TODO: call cNodeAPI to get the following info
	regCascadeTxid string

	// ticket
	creatorSignature []byte
	ticket           *pastel.ActionTicket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debug("States updated")
	})

	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")

	defer task.removeArtifacts()
	if err := task.run(ctx); err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Error("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	/* Step 3,4: Find tops supernodes and validate top 3 SNs and create mesh network of 3 SNs */
	if err := task.connectToTopRankNodes(ctx); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := make(chan struct{})
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return task.nodes.WaitConnClose(ctx, nodesDone)
	})

	/* Step 5: Send image, burn txid to SNs */

	// send registration metadata
	if err := task.sendActionMetadata(ctx); err != nil {
		return errors.Errorf("send registration metadata: %w", err)
	}

	// probe image for average rareness, nsfw and seen score
	if err := task.probeImage(ctx); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	// generateDDAndFingerprintsIDs generates dd & fp IDs
	if err := task.generateDDAndFingerprintsIDs(); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	if err := task.createCascadeTicket(ctx); err != nil {
		return errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("sign cascade ticket: %w", err)
	}

	// send signed ticket to supernodes to validate and register action with the network
	if err := task.sendSignedTicket(ctx); err != nil {
		return errors.Errorf("send signed cascade ticket: %w", err)
	}
	task.UpdateStatus(StatusTicketAccepted)

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.waitTxidValid(newCtx, task.regCascadeTxid, int64(task.config.RegArtTxMinConfirmations), 15*time.Second); err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}

	task.UpdateStatus(StatusTicketRegistered)

	// activate reg-art ticket at previous step
	activateTxID, err := task.activateActionTicket(newCtx)
	if err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("active action ticket: %w", err)
	}
	log.Debugf("Active action ticket txid: %s", activateTxID)

	// Wait until activateTxID is valid
	err = task.waitTxidValid(newCtx, activateTxID, int64(task.config.RegActTxMinConfirmations), 15*time.Second)
	if err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(StatusTicketActivated)
	log.Debugf("Active txid is confirmed")

	// Send ActionAct request to primary node
	if err := task.nodes.UploadActionAct(newCtx); err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("upload action act: %w", err)
	}

	err = task.closeSNsConnections(ctx, nodesDone)
	return err
}

func (task *Task) closeSNsConnections(ctx context.Context, nodesDone chan struct{}) error {
	var err error

	close(nodesDone)

	log.WithContext(ctx).Debug("close connections to supernodes")

	for i := range task.nodes {
		if err := task.nodes[i].Connection.Close(); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"pastelId": task.nodes[i].PastelID(),
				"addr":     task.nodes[i].String(),
			}).WithError(err).Errorf("close supernode connection failed")
		}
	}

	return err
}

// generateDDAndFingerprintsIDs generates redundant IDs and assigns to task.redundantIDs
func (task *Task) generateDDAndFingerprintsIDs() error {
	ddDataJSON, err := json.Marshal(task.fingerprintAndScores)
	if err != nil {
		return errors.Errorf("failed to marshal dd-data: %w", err)
	}

	ddEncoded := utils.B64Encode(ddDataJSON)

	var buffer bytes.Buffer
	buffer.Write(ddEncoded)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[0])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[1])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[2])
	ddFpFile := buffer.Bytes()

	task.ddAndFingerprintsIc = rand.Uint32()
	task.ddAndFingerprintsIDs, _, err = pastel.GetIDFiles(ddFpFile, task.ddAndFingerprintsIc, task.config.DDAndFingerprintsMax)
	if err != nil {
		return fmt.Errorf("get ID Files: %w", err)
	}

	comp, err := zstd.CompressLevel(nil, ddFpFile, 22)
	if err != nil {
		return errors.Errorf("compress: %w", err)
	}
	task.ddAndFpFile = utils.B64Encode(comp)

	return nil
}

func (task *Task) waitTxidValid(ctx context.Context, txID string, expectedConfirms int64, interval time.Duration) error {
	log.WithContext(ctx).Debugf("Need %d confirmation for txid %s", expectedConfirms, txID)
	blockTracker := blocktracker.New(task.pastelClient)
	baseBlkCnt, err := blockTracker.GetBlockCount()
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("failed to get block count")
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(interval):
			checkConfirms := func() error {
				subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				result, err := task.pastelClient.GetRawTransactionVerbose1(subCtx, txID)
				if err != nil {
					return errors.Errorf("get transaction: %w", err)
				}

				if result.Confirmations >= expectedConfirms {
					return nil
				}

				return errors.Errorf("not enough confirmations: expected %d, got %d", expectedConfirms, result.Confirmations)
			}

			err := checkConfirms()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("check confirmations failed")
			} else {
				return nil
			}

			currentBlkCnt, err := blockTracker.GetBlockCount()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("failed to get block count")
				continue
			}

			if currentBlkCnt-baseBlkCnt >= int32(expectedConfirms)+2 {
				return errors.Errorf("timeout when wating for confirmation of transaction %s", txID)
			}
		}
	}
}

// meshNodes establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
	var meshNodes node.List
	secInfo := &alts.SecInfo{
		PastelID:   task.Request.AppPastelID,
		PassPhrase: task.Request.AppPastelIDPassphrase,
		Algorithm:  "ed448",
	}

	primary := nodes[primaryIndex]
	log.WithContext(ctx).Debugf("Trying to connect to primary node %q", primary)
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
					// Should not run this code in go routine
					func() {
						secondariesMtx.Lock()
						defer secondariesMtx.Unlock()
						secondaries.Add(node)
					}()

					if err := node.ConnectTo(ctx, types.MeshedSuperNode{
						NodeID: primary.PastelID(),
						SessID: primary.SessID(),
					}); err != nil {
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

func (task *Task) createCascadeTicket(_ context.Context) error {
	if task.datahash == nil {
		return errEmptyDatahash
	}

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	ticket := &pastel.ActionTicket{
		Version:    1,
		Caller:     task.Request.AppPastelID,
		BlockNum:   task.creatorBlockHeight,
		BlockHash:  task.creatorBlockHash,
		ActionType: pastel.ActionTypeCascade,
		APITicketData: &pastel.APICascadeTicket{
			DataHash: task.datahash,
			// FIXME: add respective values for cascade ticket
			// DDAndFingerprintsIc:  task.ddAndFingerprintsIc,
			// DDAndFingerprintsMax: task.config.DDAndFingerprintsMax,
			// DDAndFingerprintsIDs: task.ddAndFingerprintsIDs,
		},
	}

	task.ticket = ticket
	return nil
}

func (task *Task) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeActionTicket(task.ticket)
	if err != nil {
		return errors.Errorf("encode cascade ticket %w", err)
	}

	task.creatorSignature, err = task.pastelClient.Sign(ctx, data, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign cascade ticket %w", err)
	}
	return nil
}

func (task *Task) activateActionTicket(ctx context.Context) (string, error) {
	request := pastel.ActivateActionRequest{
		RegTxID:    task.regCascadeTxid,
		BlockNum:   task.creatorBlockHeight,
		Fee:        task.Service.registrationFee,
		PastelID:   task.Request.AppPastelID,
		Passphrase: task.Request.AppPastelIDPassphrase,
	}

	return task.pastelClient.ActivateActionTicket(ctx, request)
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Request) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Request: Ticket,
	}
}

// connectToTopRankNodes - find 3 supesnodes and create mesh of 3 nodes
func (task *Task) connectToTopRankNodes(ctx context.Context) error {

	/* Step 3. Select 3 SNs */
	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return errors.Errorf("call masternode top: %w", err)
	}

	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorFindTopNodes)
		return errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

	// Get current block height & hash
	if err := task.getBlock(ctx); err != nil {
		return errors.Errorf("get block: %v", err)
	}

	// Connect to top nodes to find 3SN and validate their infor
	err = task.validateMNsInfo(ctx, topNodes)
	if err != nil {
		task.UpdateStatus(StatusErrorFindResponsdingSNs)
		return errors.Errorf("validate MNs info: %v", err)
	}

	/* Step 4. Establish mesh with 3 SNs */
	var nodes node.List
	var errs error

	for primaryRank := range task.nodes {
		nodes, err = task.meshNodes(ctx, task.nodes, primaryRank)
		if err != nil {
			// close connected connections
			task.nodes.DisconnectAll()

			if errors.IsContextCanceled(err) {
				return err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warn("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(nodes) < task.config.NumberSuperNodes {
		// close connected connections
		topNodes.DisconnectAll()
		return errors.Errorf("Could not create a mesh of %d nodes: %w", task.config.NumberSuperNodes, errs)
	}

	// Activate supernodes that are in the mesh.
	nodes.Activate()

	// Cancel context when any connection is broken.
	task.UpdateStatus(StatusConnected)

	// Send all meshed supernode info to nodes - that will be used to node send info to other nodes
	meshedSNInfo := []types.MeshedSuperNode{}
	for _, node := range nodes {
		meshedSNInfo = append(meshedSNInfo, types.MeshedSuperNode{
			NodeID: node.PastelID(),
			SessID: node.SessID(),
		})
	}

	for _, node := range nodes {
		err = node.MeshNodes(ctx, meshedSNInfo)
		if err != nil {
			nodes.DisconnectAll()
			return errors.Errorf("could not send info of meshed nodes: %w", err)
		}
	}
	return nil
}

func (task *Task) sendActionMetadata(ctx context.Context) error {
	if task.creatorBlockHash == "" {
		return errors.New("empty current block hash")
	}

	if task.Request.AppPastelID == "" {
		return errors.New("empty creator pastelID")
	}

	regMetadata := &types.ActionRegMetadata{
		BlockHash:       task.creatorBlockHash,
		CreatorPastelID: task.Request.AppPastelID,
		BurnTxID:        task.Request.BurnTxID,
	}

	return task.nodes.SendRegMetadata(ctx, regMetadata)
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

func (task *Task) probeImage(ctx context.Context) error {
	log.WithContext(ctx).WithField("filename", task.Request.Image.Name()).Debug("probe image")

	// Send image to supernodes for probing.
	if err := task.nodes.ProbeImage(ctx, task.Request.Image); err != nil {
		return errors.Errorf("send image: %w", err)
	}

	signatures := [][]byte{}
	// Match signatures received from supernodes.
	for i := 0; i < len(task.nodes); i++ {
		// Validate burn_txid transaction with SNs
		if !task.nodes.ValidBurnTxID() {
			task.UpdateStatus(StatsuErrorInvalidBurnTxID)
			return errors.New("invalid burn txid")
		}

		// Validate signatures received from supernodes.
		verified, err := task.pastelClient.Verify(ctx, task.nodes[i].FingerprintAndScoresBytes, string(task.nodes[i].Signature), task.nodes[i].PastelID(), pastel.SignAlgorithmED448)
		if err != nil {
			return errors.Errorf("probeImage: pastelClient.Verify %w", err)
		}

		if !verified {
			task.UpdateStatus(StatusErrorSignaturesNotMatch)
			return errors.Errorf("node[%s] signature doesn't match", task.nodes[i].PastelID())
		}

		signatures = append(signatures, task.nodes[i].Signature)
	}
	task.signatures = signatures

	// Match fingerprints received from supernodes.
	if err := task.nodes.MatchFingerprintAndScores(); err != nil {
		task.UpdateStatus(StatusErrorFingerprintsNotMatch)
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}

	task.fingerprintAndScores = task.nodes.FingerAndScores()
	task.UpdateStatus(StatusImageProbed)

	return nil
}

func (task *Task) sendSignedTicket(ctx context.Context) error {
	buf, err := pastel.EncodeActionTicket(task.ticket)
	if err != nil {
		return errors.Errorf("marshal ticket: %w", err)
	}
	log.Debug(string(buf))

	if err := task.nodes.UploadSignedTicket(ctx, buf, task.creatorSignature, task.ddAndFpFile); err != nil {
		return errors.Errorf("upload signed ticket: %w", err)
	}

	task.regCascadeTxid = task.nodes.RegActionTicketID()
	if task.regCascadeTxid == "" {
		return errors.Errorf("empty regCascadeTxid")
	}

	return nil
}

func (task *Task) removeArtifacts() {
	removeFn := func(file *artwork.File) {
		if file != nil {
			log.Debugf("remove file: %s", file.Name())
			if err := file.Remove(); err != nil {
				log.Debugf("remove file failed: %s", err.Error())
			}
		}
	}
	if task.Request != nil {
		removeFn(task.Request.Image)
	}
}