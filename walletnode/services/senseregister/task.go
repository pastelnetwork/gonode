package senseregister

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	"math/rand"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
)

// SenseRegisterTask is the task of registering new artwork.
type SenseRegisterTask struct {
	*common.WalletNodeTask

	senseHandler *SenseRegisterHandler

	service *SenseRegisterService
	Request *SenseRegisterRequest

	// task data to create RegArt ticket
	creatorBlockHeight int
	creatorBlockHash   string
	fingerprint        []byte
	datahash           []byte

	// signatures from SN1, SN2 & SN3 over dd_and_fingerprints data from dd-server
	signatures [][]byte

	// ticket
	creatorSignature []byte
	ticket           *pastel.ActionTicket

	// TODO: call cNodeAPI to get the following info
	//regSenseTxid string
	//registrationFee int64
}

// Run starts the task
func (task *SenseRegisterTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *SenseRegisterTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	/* Step 3,4: Find tops supernodes and validate top 3 SNs and create mesh network of 3 SNs */
	creatorBlockHeight, creatorBlockHash, err := task.senseHandler.ConnectToTopRankNodes(ctx)
	if err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	task.creatorBlockHeight = creatorBlockHeight
	task.creatorBlockHash = creatorBlockHash

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := make(chan struct{})
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return task.senseHandler.Nodes.WaitConnClose(ctx, nodesDone)
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

	if err := task.createSenseTicket(ctx); err != nil {
		return errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("sign sense ticket: %w", err)
	}

	// send signed ticket to supernodes to validate and register action with the network
	if err := task.sendSignedTicket(ctx); err != nil {
		return errors.Errorf("send signed sense ticket: %w", err)
	}
	task.UpdateStatus(common.StatusTicketAccepted)

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.PastelHandler.WaitTxidValid(newCtx, task.regSenseTxid, int64(task.service.config.RegArtTxMinConfirmations), 15*time.Second); err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}

	task.UpdateStatus(common.StatusTicketRegistered)

	// activate reg-art ticket at previous step
	activateTxID, err := task.activateActionTicket(newCtx)
	if err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("active action ticket: %w", err)
	}
	log.Debugf("Active action ticket txid: %s", activateTxID)

	// Wait until activateTxID is valid
	err = task.PastelHandler.WaitTxidValid(newCtx, activateTxID, int64(task.service.config.RegActTxMinConfirmations), 15*time.Second)
	if err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	log.Debugf("Active txid is confirmed")

	// Send ActionAct request to primary node
	if err := task.nodes.UploadActionAct(newCtx); err != nil {
		task.closeSNsConnections(ctx, nodesDone)
		return errors.Errorf("upload action act: %w", err)
	}

	err = task.closeSNsConnections(ctx, nodesDone)
	return err
}

func (task *SenseRegisterTask) sendActionMetadata(ctx context.Context) error {
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

	return task.senseHandler.SendRegMetadata(ctx, regMetadata)
}

func (task *SenseRegisterTask) probeImage(ctx context.Context) error {
	log.WithContext(ctx).WithField("filename", task.Request.Image.Name()).Debug("probe image")

	// Send image to supernodes for probing.
	if err := task.senseHandler.ProbeImage(ctx, task.Request.Image); err != nil {
		return errors.Errorf("send image: %w", err)
	}

	var signatures [][]byte
	// Match signatures received from supernodes.
	for i := 0; i < len(task.senseHandler.Nodes); i++ {
		// Validate burn_txid transaction with SNs
		if !task.senseHandler.ValidBurnTxID() {
			task.UpdateStatus(common.StatusErrorInvalidBurnTxID)
			return errors.New("invalid burn txid")
		}

		// Validate signatures received from supernodes.
		verified, err := task.PastelHandler.VerifySignature(ctx,
			task.nodes[i].FingerprintAndScoresBytes,
			string(task.nodes[i].Signature),
			task.nodes[i].PastelID(),
			pastel.SignAlgorithmED448)
		if err != nil {
			return errors.Errorf("probeImage: pastelClient.Verify %w", err)
		}

		if !verified {
			task.UpdateStatus(common.StatusErrorSignaturesNotMatch)
			return errors.Errorf("node[%s] signature doesn't match", task.nodes[i].PastelID())
		}

		signatures = append(signatures, task.nodes[i].Signature)
	}
	task.signatures = signatures

	// Match fingerprints received from supernodes.
	if err := task.nodes.MatchFingerprintAndScores(); err != nil {
		task.UpdateStatus(common.StatusErrorFingerprintsNotMatch)
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}

	task.fingerprintAndScores = task.nodes.FingerAndScores()
	task.UpdateStatus(common.StatusImageProbed)

	return nil
}

func (task *SenseRegisterTask) closeSNsConnections(ctx context.Context, nodesDone chan struct{}) error {
	var err error

	close(nodesDone)

	log.WithContext(ctx).Debug("close connections to supernodes")

	for i := range task.nodes {
		if err := task.nodes[i].ConnectionInterface.Close(); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"pastelId": task.nodes[i].PastelID(),
				"addr":     task.nodes[i].String(),
			}).WithError(err).Errorf("close supernode connection failed")
		}
	}

	return err
}

// generateDDAndFingerprintsIDs generates redundant IDs and assigns to task.redundantIDs
func (task *SenseRegisterTask) generateDDAndFingerprintsIDs() error {
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
	task.ddAndFingerprintsIDs, _, err = pastel.GetIDFiles(ddFpFile, task.ddAndFingerprintsIc, task.service.config.DDAndFingerprintsMax)
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

func (task *SenseRegisterTask) createSenseTicket(_ context.Context) error {
	if task.datahash == nil {
		return errEmptyDatahash
	}

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	ticket := &pastel.ActionTicket{
		Version:    1,
		Caller:     task.Request.AppPastelID,
		BlockNum:   task.creatorBlockHeight,
		BlockHash:  task.creatorBlockHash,
		ActionType: pastel.ActionTypeSense,
		APITicketData: &pastel.APISenseTicket{
			DataHash:             task.datahash,
			DDAndFingerprintsIc:  task.ddAndFingerprintsIc,
			DDAndFingerprintsMax: task.service.config.DDAndFingerprintsMax,
			DDAndFingerprintsIDs: task.ddAndFingerprintsIDs,
		},
	}

	task.ticket = ticket
	return nil
}

func (task *SenseRegisterTask) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeActionTicket(task.ticket)
	if err != nil {
		return errors.Errorf("encode sense ticket %w", err)
	}

	task.creatorSignature, err = task.PastelClient.Sign(ctx, data, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign sense ticket %w", err)
	}
	return nil
}

func (task *SenseRegisterTask) activateActionTicket(ctx context.Context) (string, error) {
	request := pastel.ActivateActionRequest{
		RegTxID:    task.regSenseTxid,
		BlockNum:   task.creatorBlockHeight,
		Fee:        task.registrationFee,
		PastelID:   task.Request.AppPastelID,
		Passphrase: task.Request.AppPastelIDPassphrase,
	}

	return task.PastelClient.ActivateActionTicket(ctx, request)
}

func (task *SenseRegisterTask) sendSignedTicket(ctx context.Context) error {
	buf, err := pastel.EncodeActionTicket(task.ticket)
	if err != nil {
		return errors.Errorf("marshal ticket: %w", err)
	}
	log.Debug(string(buf))

	if err := task.nodes.UploadSignedTicket(ctx, buf, task.creatorSignature, task.ddAndFpFile); err != nil {
		return errors.Errorf("upload signed ticket: %w", err)
	}

	task.regSenseTxid = task.nodes.RegActionTicketID()
	if task.regSenseTxid == "" {
		return errors.Errorf("empty regSenseTxid")
	}

	return nil
}

func (task *SenseRegisterTask) removeArtifacts() {
	if task.Request != nil {
		task.RemoveFile(task.Request.Image)
	}
}

// NewSenseRegisterTask returns a new SenseRegisterTask instance.
// TODO: make config interface and pass it instead of individual items
func NewSenseRegisterTask(service *SenseRegisterService, request *SenseRegisterRequest) *SenseRegisterTask {
	task := SenseRegisterTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		Request:        request,
	}
	meshHandler := mixins.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &SenseRegisterNodeMaker{},
		service.pastelHandler,
		request.AppPastelID, request.AppPastelIDPassphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.acceptNodesTimeout, service.config.connectToNextNodeDelay,
	)

	senseHandler := &SenseRegisterHandler{MeshHandler: meshHandler}

	task.senseHandler = senseHandler
	return &task
}
