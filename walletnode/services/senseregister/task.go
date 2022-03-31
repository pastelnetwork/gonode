package senseregister

import (
	"context"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// SenseRegistrationTask is the task of registering new nft.
type SenseRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler         *common.MeshHandler
	FingerprintsHandler *mixins.FingerprintsHandler

	service *SenseRegistrationService
	Request *common.ActionRegistrationRequest

	// data to create ticket
	creatorBlockHeight int
	creatorBlockHash   string
	dataHash           []byte
	registrationFee    int64

	// ticket
	creatorSignature []byte
	actionTicket     *pastel.ActionTicket
	serializedTicket []byte

	regSenseTxid string

	// only set to true for unit tests
	skipPrimaryNodeTxidVerify bool
}

// Run starts the task
func (task *SenseRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

// Error returns task err
func (task *SenseRegistrationTask) Error() error {
	return task.WalletNodeTask.Error()
}

func (task *SenseRegistrationTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	/* Step 3,4: Find tops supernodes and validate top 3 SNs and create mesh network of 3 SNs */
	creatorBlockHeight, creatorBlockHash, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(ctx)
	if err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	task.creatorBlockHeight = creatorBlockHeight
	task.creatorBlockHash = creatorBlockHash

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)

	/* Step 5: Send image, burn txid to SNs */

	// send registration metadata
	if err := task.sendActionMetadata(ctx); err != nil {
		return errors.Errorf("send registration metadata: %w", err)
	}

	// probe image for average rareness, nsfw and seen score - populate FingerprintsHandler with results
	if err := task.ProbeImage(ctx, task.Request.Image, task.Request.Image.Name()); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	// generateDDAndFingerprintsIDs generates dd & fp IDs
	if err := task.FingerprintsHandler.GenerateDDAndFingerprintsIDs(ctx, task.service.config.DDAndFingerprintsMax); err != nil {
		return errors.Errorf("generate dd and fp IDs: %w", err)
	}

	// calculate hash of data
	imgBytes, err := task.Request.Image.Bytes()
	if err != nil {
		return errors.Errorf("convert image to byte stream %w", err)
	}
	if task.dataHash, err = utils.Sha3256hash(imgBytes); err != nil {
		return errors.Errorf("hash encoded image: %w", err)
	}

	fileDataInMb := int64(len(imgBytes)) / (1024 * 1024)
	fee, err := task.service.pastelHandler.GetEstimatedSenseFee(ctx, fileDataInMb)
	if err != nil {
		return errors.Errorf("getting estimated fee %w", err)
	}
	task.registrationFee = int64(fee)

	if err := task.createSenseTicket(ctx); err != nil {
		return errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("sign sense ticket: %w", err)
	}

	// UPLOAD signed ticket to supernodes to validate and register action with the network
	if err := task.uploadSignedTicket(ctx); err != nil {
		return errors.Errorf("send signed sense ticket: %w", err)
	}
	task.UpdateStatus(common.StatusTicketAccepted)

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.regSenseTxid, int64(task.service.config.SenseRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)
	log.Infof("RegSense txid is confirmed: %s\n", task.regSenseTxid)

	// activate sense ticket registered at previous step by SN
	activateTxID, err := task.activateActionTicket(newCtx)
	if err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("active action ticket: %w", err)
	}
	log.Debugf("Active action ticket txid: %s", activateTxID)

	// Wait until activateTxID is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, activateTxID, int64(task.service.config.SenseActTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
	if err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	log.Infof("Active txid is confirmed: %s\n", activateTxID)

	// Send ActionAct request to primary node
	if err := task.uploadActionAct(newCtx, task.regSenseTxid); err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("upload action act: %w", err)
	}

	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
	return nil
}

// sendActionMetadata sends Action Ticket metadata to supernodes
func (task *SenseRegistrationTask) sendActionMetadata(ctx context.Context) error {
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

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(*SenseRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegistrationNode", someNode.String())
		}
		group.Go(func() (err error) {
			err = senseRegNode.SendRegMetadata(ctx, regMetadata)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

// ProbeImage sends the image to supernodes for image analysis, such as fingerprint, rareness score, NSWF.
// Add received fingerprints into Fingerprint handler
func (task *SenseRegistrationTask) ProbeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	task.FingerprintsHandler.Clear()

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(*SenseRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() (err error) {
			compress, stateOk, err := senseRegNode.ProbeImage(ctx, file)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			someNode.SetRemoteState(stateOk)
			if !stateOk {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("probe image failed")
				return errors.Errorf("remote node %s: indicated processing error", someNode.String())
			}

			fingerprintAndScores, fingerprintAndScoresBytes, signature, err := pastel.ExtractCompressSignedDDAndFingerprints(compress)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", someNode).Error("extract compressed signed DDAandFingerprints failed")
				return errors.Errorf("node %s: extract failed: %w", someNode.String(), err)
			}
			task.FingerprintsHandler.AddNew(fingerprintAndScores, fingerprintAndScoresBytes, signature, someNode.PastelID())

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}

	if err := task.FingerprintsHandler.Match(ctx); err != nil {
		log.WithContext(ctx).WithError(err).WithField("filename", fileName).Error("probe image failed")
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}

	return nil
}

func (task *SenseRegistrationTask) createSenseTicket(_ context.Context) error {
	if task.dataHash == nil ||
		task.FingerprintsHandler.DDAndFingerprintsIDs == nil {
		return common.ErrEmptyDatahash
	}

	ticket := &pastel.ActionTicket{
		Version:    1,
		Caller:     task.Request.AppPastelID,
		BlockNum:   task.creatorBlockHeight,
		BlockHash:  task.creatorBlockHash,
		ActionType: pastel.ActionTypeSense,
		APITicketData: &pastel.APISenseTicket{
			DataHash:             task.dataHash,
			DDAndFingerprintsIc:  task.FingerprintsHandler.DDAndFingerprintsIc,
			DDAndFingerprintsMax: task.service.config.DDAndFingerprintsMax,
			DDAndFingerprintsIDs: task.FingerprintsHandler.DDAndFingerprintsIDs,
		},
	}

	task.actionTicket = ticket
	return nil
}

func (task *SenseRegistrationTask) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeActionTicket(task.actionTicket)
	if err != nil {
		return errors.Errorf("encode sense ticket %w", err)
	}

	task.creatorSignature, err = task.service.pastelHandler.PastelClient.Sign(ctx, data, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign sense ticket %w", err)
	}
	task.serializedTicket = data
	return nil
}

// uploadSignedTicket uploads sense ticket  and its signature to super nodes
func (task *SenseRegistrationTask) uploadSignedTicket(ctx context.Context) error {
	if task.serializedTicket == nil {
		return errors.Errorf("uploading ticket: serializedTicket is empty")
	}
	if task.creatorSignature == nil {
		return errors.Errorf("uploading ticket: creatorSignature is empty")
	}
	ddFpFile := task.FingerprintsHandler.DDAndFpFile

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(*SenseRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxid, err := senseRegNode.SendSignedTicket(ctx, task.serializedTicket, task.creatorSignature, ddFpFile)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", senseRegNode).Error("send signed ticket failed")
				return err
			}
			if !someNode.IsPrimary() && ticketTxid != "" && !task.skipPrimaryNodeTxidCheck() {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, someNode.PastelID())
			}
			if someNode.IsPrimary() {
				if ticketTxid == "" {
					return errors.Errorf("primary node - %s, returned empty txid", someNode.PastelID())
				}
				task.regSenseTxid = ticketTxid
			}
			return nil
		})
	}
	return group.Wait()
}

func (task *SenseRegistrationTask) activateActionTicket(ctx context.Context) (string, error) {
	request := pastel.ActivateActionRequest{
		RegTxID:    task.regSenseTxid,
		BlockNum:   task.creatorBlockHeight,
		Fee:        task.registrationFee,
		PastelID:   task.Request.AppPastelID,
		Passphrase: task.Request.AppPastelIDPassphrase,
	}

	return task.service.pastelHandler.PastelClient.ActivateActionTicket(ctx, request)
}

// uploadActionAct uploads action act to primary node
func (task *SenseRegistrationTask) uploadActionAct(ctx context.Context, activateTxID string) error {
	group, _ := errgroup.WithContext(ctx)

	for _, someNode := range task.MeshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(*SenseRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegistrationNode", someNode.String())
		}
		if someNode.IsPrimary() {
			group.Go(func() error {
				return senseRegNode.SendActionAct(ctx, activateTxID)
			})
		}
	}
	return group.Wait()
}

func (task *SenseRegistrationTask) removeArtifacts() {
	if task.Request != nil {
		task.RemoveFile(task.Request.Image)
	}
}

func (task *SenseRegistrationTask) skipPrimaryNodeTxidCheck() bool {
	return task.skipPrimaryNodeTxidVerify || os.Getenv("INTEGRATION_TEST_ENV") == "true"
}

// NewSenseRegisterTask returns a new SenseRegistrationTask instance.
// TODO: make config interface and pass it instead of individual items
func NewSenseRegisterTask(service *SenseRegistrationService, request *common.ActionRegistrationRequest) *SenseRegistrationTask {
	task := common.NewWalletNodeTask(logPrefix)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &RegisterSenseNodeMaker{},
		PastelHandler: service.pastelHandler,
		NodeClient:    service.nodeClient,
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay: service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:   service.config.ConnectToNodeTimeout,
			AcceptNodesTimeout:     service.config.AcceptNodesTimeout,
			MinSNs:                 service.config.NumberSuperNodes,
			PastelID:               request.AppPastelID,
			Passphrase:             request.AppPastelIDPassphrase,
		},
	}

	return &SenseRegistrationTask{
		WalletNodeTask:      task,
		service:             service,
		Request:             request,
		MeshHandler:         common.NewMeshHandler(meshHandlerOpts),
		FingerprintsHandler: mixins.NewFingerprintsHandler(service.pastelHandler),
	}
}
