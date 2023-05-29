package cascaderegister

import (
	"context"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
)

// CascadeRegistrationTask is the task of registering new nft.
type CascadeRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler *common.MeshHandler
	RqHandler   *mixins.RQHandler

	service         *CascadeRegistrationService
	downloadService *download.NftDownloadingService
	Request         *common.ActionRegistrationRequest

	// data to create ticket
	creatorBlockHeight      int
	creatorBlockHash        string
	dataHash                []byte
	registrationFee         int64
	originalFileSizeInBytes int
	fileType                string

	// ticket
	creatorSignature []byte
	actionTicket     *pastel.ActionTicket
	serializedTicket []byte

	regCascadeTxid string

	// only set to true for unit tests
	skipPrimaryNodeTxidVerify bool
}

// Run starts the task
func (task *CascadeRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *CascadeRegistrationTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.WithContext(ctx).Info("Setting up mesh with Top Supernodes")
	task.StatusLog[common.FieldTaskType] = "Cascade Registration"

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
	task.StatusLog[common.FieldBlockHeight] = creatorBlockHeight

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)

	log.WithContext(ctx).Info("uploading data to supernodes")

	// send registration metadata
	if err := task.sendActionMetadata(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error sending action metadata")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSendingRegMetadata)
		return errors.Errorf("send registration metadata: %w", err)
	}
	log.WithContext(ctx).Info("action metadata has been sent")

	if err := task.uploadImage(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error uploading image")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorUploadImageFailed)
		return errors.Errorf("upload image: %w", err)
	}
	log.WithContext(ctx).Info("image has been uploaded")

	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Generating RaptorQ symbols' identifiers ",
		StatusString:  "",
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	// connect to rq serivce to get rq symbols identifier
	if err := task.RqHandler.GenRQIdentifiersFiles(ctx, task.Request.Image,
		task.creatorBlockHash, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase); err != nil {
		log.WithContext(ctx).WithError(err).Error("error generating RQIDs")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorGenRaptorQSymbolsFailed)
		return errors.Errorf("gen RaptorQ symbols' identifiers: %w", err)
	}
	log.WithContext(ctx).Info("rq-ids have been generated")

	// calculate hash of data
	nftBytes, err := task.Request.Image.Bytes()
	task.originalFileSizeInBytes = len(nftBytes)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error converting image to bytes")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorConvertingImageBytes)
		return errors.Errorf("convert image to byte stream %w", err)
	}

	//Detect the file type
	task.fileType = mimetype.Detect(nftBytes).String()

	if task.dataHash, err = utils.Sha3256hash(nftBytes); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("error creating hash")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorEncodingImage)
		return errors.Errorf("hash encoded image: %w", err)
	}

	fileDataInMb := utils.GetFileSizeInMB(nftBytes)
	fee, err := task.service.pastelHandler.GetEstimatedCascadeFee(ctx, fileDataInMb)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting estimated fee")
		return errors.Errorf("getting estimated fee %w", err)
	}
	task.registrationFee = int64(fee)

	task.StatusLog[common.FieldFileSize] = fileDataInMb
	task.StatusLog[common.FieldFee] = fee

	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Creating Registration Ticket: ",
		StatusString:  "Please Wait...",
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	if err := task.createCascadeTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating cascade ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorCreatingTicket)
		return errors.Errorf("create ticket: %w", err)
	}
	log.WithContext(ctx).Info("created cascade reg ticket")

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSigningTicket)
		return errors.Errorf("sign cascade ticket: %w", err)
	}
	log.WithContext(ctx).Info("signed cascade reg ticket")

	// UPLOAD signed ticket to supernodes to validate and register action with the network
	if err := task.uploadSignedTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error uploading signed ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorUploadingTicket)
		return errors.Errorf("send signed cascade ticket: %w", err)
	}
	task.StatusLog[common.FieldRegTicketTxnID] = task.regCascadeTxid
	task.UpdateStatus(common.StatusTicketAccepted)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validating Cascade Reg TXID: ",
		StatusString:  task.regCascadeTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Info("uploaded signed cascade reg ticket to SNs")

	now := time.Now()
	if task.downloadService != nil {
		if err := common.DownloadWithRetry(ctx, task, now, now.Add(common.RetryTime*time.Minute)); err != nil {
			log.WithContext(ctx).WithField("reg_tx_id", task.regCascadeTxid).WithError(err).Error("error validating cascade ticket data")
			task.StatusLog[common.FieldErrorDetail] = err.Error()
			task.StatusLog[common.FieldMeshNodes] = task.MeshHandler.Nodes.String()
			task.UpdateStatus(common.StatusErrorDownloadFailed)
			task.UpdateStatus(&common.EphemeralStatus{
				StatusTitle:   "Error validating cascade ticket data",
				StatusString:  task.regCascadeTxid,
				IsFailureBool: false,
				IsFinalBool:   false,
			})

			task.UpdateStatus(common.StatusTaskRejected)
			task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

			return errors.Errorf("error validating cascade ticket data")
		}
	}
	log.WithContext(ctx).Info("data-validation completed for cascade-reg ticket")

	log.WithContext(ctx).Infof("Cascade Reg Ticket registered. Cascade Registration Ticket txid: %s", task.regCascadeTxid)
	log.WithContext(ctx).Info("Closing SNs connections")

	// don't need SNs anymore
	err = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error closing SNs connections")
	}

	log.WithContext(ctx).Infof("Waiting Confirmations for Cascade Reg Ticket - Ticket txid: %s", task.regCascadeTxid)

	// new context because the old context already cancelled
	newCtx := log.ContextWithPrefix(context.Background(), "cascade")
	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.regCascadeTxid, int64(task.service.config.CascadeRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		log.WithContext(ctx).WithError(err).Error("error waiting for Reg TXID confirmations")
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validated Cascade Reg TXID: ",
		StatusString:  task.regCascadeTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	log.WithContext(ctx).Info("cascade reg ticket confirmed, now activating cascade reg ticket")

	// activate cascade ticket registered at previous step by SN
	activateTxID, err := task.activateActionTicket(newCtx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error activating action ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorActivatingTicket)
		return errors.Errorf("active action ticket: %w", err)
	}
	task.StatusLog[common.FieldActivateTicketTxnID] = activateTxID
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Cascade Action Activated - Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Cascade ticket activated. Activation ticket txid: %s, waiting for confirmations now", activateTxID)
	// Wait until activateTxID is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, activateTxID, int64(task.service.config.CascadeActTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error waiting for activated Reg TXID confirmations")
		return errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activated Cascade Action Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Cascade Activation ticket is confirmed. Activation ticket txid: %s", activateTxID)

	return nil
}

// sendActionMetadata sends Action Ticket metadata to supernodes
func (task *CascadeRegistrationTask) sendActionMetadata(ctx context.Context) error {
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
	task.StatusLog[common.FieldBurnTxnID] = task.Request.BurnTxID

	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		cascadeRegNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() (err error) {
			err = cascadeRegNode.SendRegMetadata(gctx, regMetadata)
			if err != nil {
				task.StatusLog[common.FieldErrorDetail] = err.Error()
				log.WithContext(gctx).WithError(err).WithField("node", cascadeRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

func (task *CascadeRegistrationTask) uploadImage(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)

	for _, someNode := range task.MeshHandler.Nodes {
		cascadeNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		image := task.Request.Image
		group.Go(func() error {
			err := cascadeNode.UploadAsset(gctx, image)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", someNode).Error("upload image with thumbnail failed")
				return err
			}
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Errorf("upload image: %w", err)
	}

	task.UpdateStatus(common.StatusImageAndThumbnailUploaded)

	return nil
}

func (task *CascadeRegistrationTask) createCascadeTicket(_ context.Context) error {
	if task.dataHash == nil {
		return common.ErrEmptyDatahash
	}
	if task.RqHandler.IsEmpty() {
		return common.ErrEmptyRaptorQSymbols
	}

	ticket := &pastel.ActionTicket{
		Version:    1,
		Caller:     task.Request.AppPastelID,
		BlockNum:   task.creatorBlockHeight,
		BlockHash:  task.creatorBlockHash,
		ActionType: pastel.ActionTypeCascade,
		APITicketData: pastel.APICascadeTicket{
			FileName:                task.Request.FileName,
			DataHash:                task.dataHash,
			RQIc:                    task.RqHandler.RQIDsIc,
			RQMax:                   task.service.config.RQIDsMax,
			RQIDs:                   task.RqHandler.RQIDs,
			RQOti:                   task.RqHandler.RQEncodeParams.Oti,
			OriginalFileSizeInBytes: task.originalFileSizeInBytes,
			FileType:                task.fileType,
			MakePubliclyAccessible:  task.Request.MakePubliclyAccessible,
		},
	}

	task.actionTicket = ticket
	return nil
}

func (task *CascadeRegistrationTask) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeActionTicket(task.actionTicket)
	if err != nil {
		return errors.Errorf("encode cascade ticket %w", err)
	}

	task.creatorSignature, err = task.service.pastelHandler.PastelClient.Sign(ctx, data, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign cascade ticket %w", err)
	}
	task.serializedTicket = data
	return nil
}

// uploadSignedTicket uploads cascade ticket  and its signature to super nodes
func (task *CascadeRegistrationTask) uploadSignedTicket(ctx context.Context) error {
	if task.serializedTicket == nil {
		return errors.Errorf("uploading ticket: serializedTicket is empty")
	}
	if task.creatorSignature == nil {
		return errors.Errorf("uploading ticket: creatorSignature is empty")
	}

	rqidsFile := task.RqHandler.RQIDsFile
	encoderParams := task.RqHandler.RQEncodeParams

	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		cascadeRegNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxid, err := cascadeRegNode.SendSignedTicket(gctx, task.serializedTicket, task.creatorSignature, rqidsFile, encoderParams)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", cascadeRegNode).Error("send signed ticket failed")
				return err
			}
			if !someNode.IsPrimary() && ticketTxid != "" && !task.skipPrimaryNodeTxidCheck() {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, someNode.PastelID())
			}
			if someNode.IsPrimary() {
				if ticketTxid == "" {
					return errors.Errorf("primary node - %s, returned empty txid", someNode.PastelID())
				}
				task.regCascadeTxid = ticketTxid
			}
			return nil
		})
	}
	return group.Wait()
}

func (task *CascadeRegistrationTask) activateActionTicket(ctx context.Context) (string, error) {
	request := pastel.ActivateActionRequest{
		RegTxID:    task.regCascadeTxid,
		BlockNum:   task.creatorBlockHeight,
		Fee:        task.registrationFee,
		PastelID:   task.Request.AppPastelID,
		Passphrase: task.Request.AppPastelIDPassphrase,
	}

	return task.service.pastelHandler.PastelClient.ActivateActionTicket(ctx, request)
}

func (task *CascadeRegistrationTask) removeArtifacts() {
	if task.Request != nil {
		task.RemoveFile(task.Request.Image)
	}
}

func (task *CascadeRegistrationTask) skipPrimaryNodeTxidCheck() bool {
	return task.skipPrimaryNodeTxidVerify || os.Getenv("INTEGRATION_TEST_ENV") == "true"
}

// Error returns task err
func (task *CascadeRegistrationTask) Error() error {
	return task.WalletNodeTask.Error()
}

// Download downloads the data from p2p for data validation before ticket activation
func (task *CascadeRegistrationTask) Download(ctx context.Context) error {
	// add the task to the worker queue, and worker will process the task in the background
	log.WithContext(ctx).WithField("cascade_tx_id", task.regCascadeTxid).Info("Downloading has been started")
	taskID := task.downloadService.AddTask(&nft.DownloadPayload{Txid: task.regCascadeTxid, Pid: task.Request.AppPastelID, Key: task.Request.AppPastelIDPassphrase}, pastel.ActionTypeCascade)
	downloadTask := task.downloadService.GetTask(taskID)
	defer downloadTask.Cancel()

	sub := downloadTask.SubscribeStatus()
	log.WithContext(ctx).WithField("cascade_tx_id", task.regCascadeTxid).Info("Subscribed to status channel")

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			if status.IsFailure() {
				log.WithContext(ctx).WithField("cascade_tx_id", task.regCascadeTxid).WithError(task.Error())

				return errors.New("Download failed")
			}

			if status.IsFinal() {
				log.WithContext(ctx).WithField("cascade_tx_id", task.regCascadeTxid).Info("task has been downloaded successfully")
				return nil
			}
		case <-time.After(20 * time.Minute):
			log.WithContext(ctx).WithField("cascade_tx_id", task.regCascadeTxid).Info("Download request has been timed out")
			return errors.New("download request timeout, data validation failed")

		}
	}
}

// NewCascadeRegisterTask returns a new CascadeRegistrationTask instance.
// TODO: make config interface and pass it instead of individual items
func NewCascadeRegisterTask(service *CascadeRegistrationService, request *common.ActionRegistrationRequest) *CascadeRegistrationTask {
	task := common.NewWalletNodeTask(logPrefix)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &RegisterCascadeNodeMaker{},
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

	return &CascadeRegistrationTask{
		WalletNodeTask:  task,
		downloadService: &service.downloadHandler,
		service:         service,
		Request:         request,
		MeshHandler:     common.NewMeshHandler(meshHandlerOpts),
		RqHandler: mixins.NewRQHandler(service.rqClient, service.pastelHandler,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir,
			service.config.NumberRQIDSFiles, service.config.RQIDsMax),
	}
}
