package senseregister

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/duplicate"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
)

const (
	//DateTimeFormat following the go convention for request timestamp
	DateTimeFormat = "2006:01:02 15:04:05"

	taskTypeSense = "sense"
)

// SenseRegistrationTask is the task of registering new nft.
type SenseRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler         *common.MeshHandler
	FingerprintsHandler *mixins.FingerprintsHandler

	service         *SenseRegistrationService
	downloadService *download.NftDownloadingService
	Request         *common.ActionRegistrationRequest

	// data to create ticket
	creatorBlockHeight int
	creatorBlockHash   string
	creationTimestamp  string
	dataHash           []byte
	registrationFee    int64

	// ticket
	creatorSignature []byte
	actionTicket     *pastel.ActionTicket
	serializedTicket []byte

	regSenseTxid   string
	collectionTxID string
	taskType       string

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

	log.WithContext(ctx).Info("checking collection verification")
	err := task.IsValidForCollection(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("ticket is not valid to be added in collection")
		return err
	}

	log.WithContext(ctx).Info("Setting up mesh with Top Supernodes")
	task.StatusLog[common.FieldTaskType] = "Sense Registration"

	/* Step 3,4: Find tops supernodes and validate top 3 SNs and create mesh network of 3 SNs */
	creatorBlockHeight, creatorBlockHash, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error setting up mesh of supernodes")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorMeshSetupFailed)
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	task.creatorBlockHeight = creatorBlockHeight
	task.creatorBlockHash = creatorBlockHash
	task.creationTimestamp = time.Now().Format(DateTimeFormat)
	task.StatusLog[common.FieldBlockHeight] = creatorBlockHeight

	log.WithContext(ctx).Info("Mesh of supernodes have been established")

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)

	/* Step 5: Send image, burn txid to SNs */

	log.WithContext(ctx).Info("Uploading data to Supernodes")

	// send registration metadata
	if err := task.sendActionMetadata(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error sending action metadata")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSendingRegMetadata)
		return errors.Errorf("send registration metadata: %w", err)
	}

	log.WithContext(ctx).Info("action metadata has been sent")

	// calculate hash of data
	imgBytes, err := task.Request.Image.Bytes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error converting image to bytes")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorConvertingImageBytes)
		return errors.Errorf("convert image to byte stream %w", err)
	}
	if task.dataHash, err = utils.Sha3256hash(imgBytes); err != nil {
		log.WithContext(ctx).WithError(err).Error("error converting bytes to hash")

		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorEncodingImage)
		return errors.Errorf("hash encoded image: %w", err)
	}

	task.UpdateStatus(common.StatusValidateDuplicateTickets)
	dtc := duplicate.NewDupTicketsDetector(task.service.pastelHandler.PastelClient)
	if err := dtc.CheckDuplicateSenseOrNFTTickets(ctx, task.dataHash); err != nil {
		log.WithContext(ctx).WithError(err)
		return errors.Errorf("Error duplicate ticket")
	}
	log.WithContext(ctx).Info("no duplicate tickets have been found")

	task.UpdateStatus(common.StatusValidateBurnTxn)

	// probe image for average rareness, nsfw and seen score - populate FingerprintsHandler with results
	if err := task.ProbeImage(ctx, task.Request.Image, task.Request.Image.Name()); err != nil {
		log.WithContext(ctx).WithError(err).Error("error probing image")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorProbeImage)
		return errors.Errorf("probe image: %w", err)
	}

	log.WithContext(ctx).Info("image has been probed, getting DD & fingerprints from supernodes")

	// generateDDAndFingerprintsIDs generates dd & fp IDs
	if err := task.FingerprintsHandler.GenerateDDAndFingerprintsIDs(ctx, task.service.config.DDAndFingerprintsMax); err != nil {
		log.WithContext(ctx).WithError(err).Error("error generating DD & Fingerprint IDs")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorGenerateDDAndFPIds)
		return errors.Errorf("generate dd and fp IDs: %w", err)
	}

	log.WithContext(ctx).Info("fingerprints have been generated")

	fee, err := task.service.pastelHandler.GetEstimatedSenseFee(ctx, utils.GetFileSizeInMB(imgBytes))
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting estimated sense fee")
		return errors.Errorf("getting estimated fee %w", err)
	}
	task.registrationFee = int64(fee)
	task.StatusLog[common.FieldFee] = fee

	log.WithContext(ctx).Info("Create and sign Sense Reg Ticket")

	if err := task.createSenseTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating sense ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorCreatingTicket)
		return errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing sense ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSigningTicket)
		return errors.Errorf("sign sense ticket: %w", err)
	}

	log.WithContext(ctx).Info("Upload signed Sense Reg Ticket to SNs")

	// UPLOAD signed ticket to supernodes to validate and register action with the network
	if err := task.uploadSignedTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error uploading signed ticket")

		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorUploadingTicket)
		return errors.Errorf("send signed sense ticket: %w", err)
	}
	task.StatusLog[common.FieldRegTicketTxnID] = task.regSenseTxid

	task.UpdateStatus(common.StatusTicketAccepted)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validating Sense Reg TXID: ",
		StatusString:  task.regSenseTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	log.WithContext(ctx).Infof("Sense Reg Ticket registered. Sense Registration Ticket txid: %s", task.regSenseTxid)
	log.WithContext(ctx).Info("Closing SNs connections")

	now := time.Now()
	if task.downloadService != nil {
		if err := common.DownloadWithRetry(ctx, task, now, now.Add(common.RetryTime*time.Minute)); err != nil {
			log.WithContext(ctx).WithField("reg_sense_tx_id", task.regSenseTxid).WithError(err).Error("error validating sense ticket data")

			task.StatusLog[common.FieldErrorDetail] = err.Error()
			task.StatusLog[common.FieldMeshNodes] = task.MeshHandler.Nodes.String()
			task.UpdateStatus(common.StatusErrorDownloadFailed)
			task.UpdateStatus(&common.EphemeralStatus{
				StatusTitle:   "Error validating sense ticket data",
				StatusString:  task.regSenseTxid,
				IsFailureBool: false,
				IsFinalBool:   false,
			})

			task.UpdateStatus(common.StatusTaskRejected)
			task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

			return errors.Errorf("error validating sense ticket data")
		}
	}

	// don't need SNs anymore
	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	log.WithContext(ctx).Infof("Waiting Confirmations for Sense Reg Ticket - Ticket txid: %s", task.regSenseTxid)

	// new context because the old context already cancelled
	newCtx := log.ContextWithPrefix(context.Background(), "sense")
	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.regSenseTxid, int64(task.service.config.SenseRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

		log.WithContext(ctx).WithError(err).Error("error getting confirmations")
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validated Sense Reg TXID: ",
		StatusString:  task.regSenseTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	log.WithContext(ctx).Debug("Sense Reg Ticket confirmed, Activating Sense Reg Ticket")
	// activate sense ticket registered at previous step by SN
	activateTxID, err := task.activateActionTicket(newCtx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error activating sense ticket")

		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorActivatingTicket)
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("active action ticket: %w", err)
	}
	task.StatusLog[common.FieldActivateTicketTxnID] = activateTxID
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activating Action Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Sense ticket activated. Activation ticket txid: %s", activateTxID)
	log.WithContext(ctx).Infof("Waiting Confirmations for Sense Activation Ticket - Ticket txid: %s", activateTxID)

	// Wait until activateTxID is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, activateTxID, int64(task.service.config.SenseActTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting confirmations for activation")
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activated Action Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Sense Activation ticket is confirmed. Activation ticket txid: %s", activateTxID)

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
		BlockHeight:     strconv.Itoa(task.creatorBlockHeight),
		Timestamp:       task.creationTimestamp,
		CreatorPastelID: task.Request.AppPastelID,
		BurnTxID:        task.Request.BurnTxID,
		GroupID:         task.Request.GroupID,
	}

	if task.Request.CollectionTxID != "" {
		regMetadata.CollectionTxID = task.Request.CollectionTxID
	}

	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(*SenseRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegistrationNode", someNode.String())
		}
		group.Go(func() (err error) {
			err = senseRegNode.SendRegMetadata(gctx, regMetadata)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", senseRegNode).Error("send registration metadata failed")
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

	reqCtx, cancel := context.WithTimeout(ctx, 40*time.Minute)
	defer cancel()
	group, gctx := errgroup.WithContext(reqCtx)
	for _, someNode := range task.MeshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(*SenseRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() (err error) {
			compress, stateOk, errString, hashExists, err := senseRegNode.ProbeImage(gctx, file)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", senseRegNode).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			if hashExists {
				log.WithContext(gctx).WithField("node", senseRegNode).Error("image already registered")
				return errors.Errorf("remote node %s: image already registered: %s", someNode.String(), errString)
			}

			someNode.SetRemoteState(stateOk)
			if !stateOk {
				log.WithContext(gctx).WithField("node", senseRegNode).Error("probe image failed")
				return errors.Errorf("remote node %s: indicated processing error:%s", someNode.String(), errString)
			}

			fingerprintAndScores, fingerprintAndScoresBytes, signature, err := pastel.ExtractCompressSignedDDAndFingerprints(compress)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", someNode).Error("extract compressed signed DDAandFingerprints failed")
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

func (task *SenseRegistrationTask) createSenseTicket(_ context.Context) (err error) {
	if task.dataHash == nil ||
		task.FingerprintsHandler.DDAndFingerprintsIDs == nil {
		return common.ErrEmptyDatahash
	}

	ticket := &pastel.ActionTicket{
		Version:    2,
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

	if task.collectionTxID != "" {
		ticket.CollectionTxID = task.collectionTxID
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

	reqCtx, cancel := context.WithTimeout(ctx, 40*time.Minute)
	defer cancel()

	group, gctx := errgroup.WithContext(reqCtx)
	for _, someNode := range task.MeshHandler.Nodes {
		senseRegNode, ok := someNode.SuperNodeAPIInterface.(*SenseRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not SenseRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxid, err := senseRegNode.SendSignedTicket(gctx, task.serializedTicket, task.creatorSignature, ddFpFile)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", senseRegNode).Error("send signed ticket failed")
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

func (task *SenseRegistrationTask) removeArtifacts() {
	if task.Request != nil {
		task.RemoveFile(task.Request.Image)
	}
}

func (task *SenseRegistrationTask) skipPrimaryNodeTxidCheck() bool {
	return task.skipPrimaryNodeTxidVerify || os.Getenv("INTEGRATION_TEST_ENV") == "true"
}

// Download downloads the data from p2p for data validation before ticket activation
func (task *SenseRegistrationTask) Download(ctx context.Context) error {
	// add the task to the worker queue, and worker will process the task in the background
	log.WithContext(ctx).WithField("sense_tx_id", task.regSenseTxid).Info("Downloading has been started")
	taskID := task.downloadService.AddTask(&nft.DownloadPayload{Txid: task.regSenseTxid, Pid: task.Request.AppPastelID, Key: task.Request.AppPastelIDPassphrase}, pastel.ActionTypeSense)
	downloadTask := task.downloadService.GetTask(taskID)
	defer downloadTask.Cancel()

	sub := downloadTask.SubscribeStatus()
	log.WithContext(ctx).WithField("sense_tx_id", task.regSenseTxid).Info("Subscribed to status channel")

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			if status.IsFailure() {
				log.WithContext(ctx).WithField("sense_tx_id", task.regSenseTxid).WithError(task.Error())

				return errors.New("Download failed")
			}

			if status.IsFinal() {
				log.WithContext(ctx).WithField("sense_tx_id", task.regSenseTxid).Info("task has been downloaded successfully")
				return nil
			}
		case <-time.After(20 * time.Minute):
			log.WithContext(ctx).WithField("sense_tx_id", task.regSenseTxid).Info("Download request has been timed out")
			return errors.New("download request timeout, data validation failed")
		}
	}
}

// IsValidForCollection checks if the ticket is valid to be added to collection
func (task *SenseRegistrationTask) IsValidForCollection(ctx context.Context) error {
	if task.Request.CollectionTxID == "" {
		log.WithContext(ctx).Info("no collection txid found in the request, should proceed with normal registration")
		return nil
	}

	collectionActTicket, err := task.service.pastelHandler.PastelClient.CollectionActTicket(ctx, task.Request.CollectionTxID)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("collection_act_tx_id", task.Request.CollectionTxID).Error("error getting collection-act ticket")
		return errors.Errorf("error retrieving collection ticket")
	}

	collectionTicket, err := task.service.pastelHandler.PastelClient.CollectionRegTicket(ctx, collectionActTicket.CollectionActTicketData.RegTXID)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("collection_reg_tx_id", collectionActTicket.CollectionActTicketData.RegTXID).Error("error getting collection-reg ticket")
		return errors.Errorf("error retrieving collection ticket")
	}

	ticket := collectionTicket.CollectionRegTicketData.CollectionTicketData

	isValidContributor := false
	for _, pastelID := range ticket.ListOfPastelIDsOfAuthorizedContributors {
		if task.Request.AppPastelID == pastelID {
			isValidContributor = true
		}
	}

	if !isValidContributor {
		return errors.Errorf("creator's pastelid is not a found in list of authorized contributors")
	}

	if task.taskType != strings.ToLower(ticket.ItemType) {
		return errors.Errorf("ticket item type does not match with collection item type")
	}

	if collectionActTicket.CollectionActTicketData.IsFull {
		return errors.Errorf("collection does not have capacity for more items, max collection entries have been reached")
	}

	if collectionActTicket.CollectionActTicketData.IsExpiredByHeight {
		return errors.Errorf("block height exceeds than collection final allowed block height")
	}

	task.collectionTxID = task.Request.CollectionTxID

	return nil
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
		LogRequestID:  task.ID(),
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay: service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:   service.config.ConnectToNodeTimeout,
			AcceptNodesTimeout:     service.config.AcceptNodesTimeout,
			MinSNs:                 service.config.NumberSuperNodes,
			PastelID:               request.AppPastelID,
			Passphrase:             request.AppPastelIDPassphrase,
			CheckDDDatabaseHashes:  true,
			HashCheckMaxRetries:    service.config.HashCheckMaxRetries,
		},
	}

	return &SenseRegistrationTask{
		WalletNodeTask:      task,
		downloadService:     service.downloadHandler,
		service:             service,
		Request:             request,
		MeshHandler:         common.NewMeshHandler(meshHandlerOpts),
		taskType:            taskTypeSense,
		FingerprintsHandler: mixins.NewFingerprintsHandler(service.pastelHandler),
	}
}
