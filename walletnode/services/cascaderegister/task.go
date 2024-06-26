package cascaderegister

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
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

	regCascadeFinishedAt time.Time
	regCascadeStartedAt  time.Time

	actAttemptID int64
}

// Run starts the task
func (task *CascadeRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *CascadeRegistrationTask) run(ctx context.Context) error {
	regTxid, actTxid, err := task.runTicketRegActTask(ctx)
	if err != nil {
		attemptErr := task.service.HandleTaskRegistrationErrorAttempts(ctx, task.ID(), regTxid, actTxid, task.Request.RegAttemptID, task.actAttemptID, err)
		if attemptErr != nil {
			return attemptErr
		}

		return err
	}

	taskID := task.ID()
	if regTxid != "" && actTxid != "" {
		doneBlock, err := task.service.pastelHandler.PastelClient.GetBlockCount(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error retrieving block-count")
			return err
		}

		file, err := task.service.GetFileByTaskID(taskID)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error retrieving file")
			return nil
		}

		file.DoneBlock = int(doneBlock)
		file.RegTxid = regTxid
		file.ActivationTxid = actTxid
		file.IsConcluded = true
		err = task.service.ticketDB.UpsertFile(*file)
		if err != nil {
			log.Errorf("Error in file upsert: %v", err.Error())
			return nil
		}

		// Upsert registration attempts

		ra, err := task.service.GetRegistrationAttemptsByID(int(task.Request.RegAttemptID))
		if err != nil {
			log.Errorf("Error retrieving file reg attempt: %v", err.Error())
			return nil
		}

		ra.FinishedAt = time.Now().UTC()
		ra.IsSuccessful = true
		_, err = task.service.UpdateRegistrationAttempts(*ra)
		if err != nil {
			log.Errorf("Error in registration attempts upsert: %v", err.Error())
			return nil
		}

		actAttempt, err := task.service.GetActivationAttemptByID(int(task.actAttemptID))
		if err != nil {
			log.Errorf("Error retrieving file act attempt: %v", err.Error())
			return nil
		}

		actAttempt.IsSuccessful = true
		_, err = task.service.UpdateActivationAttempts(*actAttempt)
		if err != nil {
			log.Errorf("Error in activation attempts upsert: %v", err.Error())
			return nil
		}

		return nil
	}

	return nil
}

func (task *CascadeRegistrationTask) runTicketRegActTask(ctx context.Context) (regTxid, actTxid string, err error) {
	if r := recover(); r != nil {
		log.Errorf("Recovered from panic in cascade run: %v", r)
	}

	go func() {
		<-ctx.Done()
		log.Println("root context 'ctx' in run func was cancelled:", ctx.Err())
	}()

	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-connCtx.Done()
	}()

	task.regCascadeStartedAt = time.Now()

	nftBytes, err := task.Request.Image.Bytes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error converting image to bytes")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorConvertingImageBytes)
		return "", "", errors.Errorf("convert image to byte stream %w", err)
	}
	fileDataInMb := utils.GetFileSizeInMB(nftBytes)

	fee, err := task.service.pastelHandler.GetEstimatedCascadeFee(ctx, fileDataInMb)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting estimated fee")
		return "", "", errors.Errorf("getting estimated fee %w", err)
	}
	task.registrationFee = int64(fee)

	task.StatusLog[common.FieldFileSize] = fileDataInMb
	task.StatusLog[common.FieldFee] = fee

	log.WithContext(ctx).Info("Validating Pre-Burn Transaction")
	task.UpdateStatus(common.StatusValidateBurnTxn)
	if err := task.service.pastelHandler.ValidateBurnTxID(ctx, task.Request.BurnTxID, fee); err != nil {

		log.WithContext(ctx).WithError(err).Error("error getting confirmations on burn txn")
		return "", "", errors.Errorf("waiting on burn txn confirmations failed: %w", err)
	}
	task.UpdateStatus(common.StatusBurnTxnValidated)
	log.WithContext(ctx).Info("Pre-Burn Transaction validated successfully, hashing data now..")

	task.originalFileSizeInBytes = len(nftBytes)
	//Detect the file type
	task.fileType = mimetype.Detect(nftBytes).String()

	if task.dataHash, err = utils.Sha3256hash(nftBytes); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("error creating hash")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorEncodingImage)
		return "", "", errors.Errorf("hash encoded image: %w", err)
	}

	key := append(nftBytes, []byte(task.WalletNodeTask.ID())...)
	sortKey, _ := utils.Sha3256hash(key)

	log.WithContext(ctx).Info("Setting up mesh with Top Supernodes")
	task.StatusLog[common.FieldTaskType] = "Cascade Registration"

	/* Step 3,4: Find tops supernodes and validate top 3 SNs and create mesh network of 3 SNs */
	creatorBlockHeight, creatorBlockHash, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(connCtx, sortKey)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error setting up mesh of supernodes")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorMeshSetupFailed)

		if err := task.MeshHandler.CloseSNsConnections(ctx, nil); err != nil { // disconnect in case any of the nodes were connected
			log.WithContext(ctx).WithError(err).Error("error closing sn-connections")
		}

		return "", "", errors.Errorf("connect to top rank nodes: %w", err)
	}

	task.creatorBlockHeight = creatorBlockHeight
	task.creatorBlockHash = creatorBlockHash
	task.StatusLog[common.FieldBlockHeight] = creatorBlockHeight
	task.StatusLog[common.FieldMeshNodes] = task.MeshHandler.Nodes.String()

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	var closeErr error
	nodesDone := task.MeshHandler.ConnectionsSupervisor(connCtx, cancel)
	defer func() {
		if closeErr != nil {
			if err := task.MeshHandler.CloseSNsConnections(ctx, nodesDone); err != nil {
				log.WithContext(ctx).WithError(err).Error("error closing sn-connections")
			}
		}
	}()

	log.WithContext(ctx).Info("Sending Metdata to Supernodes")

	// send registration metadata
	if err := task.sendActionMetadata(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error sending action metadata")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		closeErr = err
		task.UpdateStatus(common.StatusErrorSendingRegMetadata)
		return "", "", errors.Errorf("send registration metadata: %w", err)
	}
	log.WithContext(ctx).Info("Uploading data to Supernodes")

	if err := task.uploadImage(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error uploading data")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorUploadImageFailed)
		closeErr = err
		return "", "", errors.Errorf("upload data: %w", err)
	}
	log.WithContext(ctx).Info("Data uploaded onto SNs, generating RaptorQ symbols' identifiers")

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
		closeErr = err
		return "", "", errors.Errorf("gen RaptorQ symbols' identifiers: %w", err)
	}
	log.WithContext(ctx).Info("RaptorQ symbols' identifiers generated")

	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Creating Registration Ticket: ",
		StatusString:  "Please Wait...",
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	if err := task.createCascadeTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating cascade ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		closeErr = err
		task.UpdateStatus(common.StatusErrorCreatingTicket)
		return "", "", errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSigningTicket)
		closeErr = err
		return "", "", errors.Errorf("sign cascade ticket: %w", err)
	}

	log.WithContext(ctx).Info("Cascade Ticket sent to Supernodes for validation & registration - waiting for Supernodes to respond..")
	if err := task.uploadSignedTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error uploading signed ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorUploadingTicket)
		closeErr = err
		return "", "", errors.Errorf("send signed cascade ticket: %w", err)
	}
	task.StatusLog[common.FieldRegTicketTxnID] = task.regCascadeTxid
	task.regCascadeFinishedAt = time.Now().UTC()
	task.UpdateStatus(common.StatusTicketAccepted)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validating Cascade Reg TXID: ",
		StatusString:  task.regCascadeTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Info("Cascade Registration Ticket accepted & registrated by Supernodes, downloading data for validation..")

	// don't need SNs anymore
	err = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
	if err != nil {
		closeErr = err
		log.WithContext(ctx).WithError(err).Error("error closing SNs connections")
	}

	now := time.Now().UTC()
	if task.downloadService != nil {
		if err := common.DownloadWithRetry(ctx, task, now, now.Add(common.RetryTime*time.Minute), true); err != nil {
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
			return task.regCascadeTxid, "", errors.Errorf("error validating cascade ticket data")
		}
	}
	log.WithContext(ctx).Infof("Cascade Registration Ticket validated through Download - Its registered with txid: %s", task.regCascadeTxid)

	log.WithContext(ctx).Infof("Waiting for enough Confirmations for Cascade Registration Ticket - Ticket txid: %s", task.regCascadeTxid)
	// new context because the old context already cancelled
	if err := task.service.pastelHandler.WaitTxidValid(ctx, task.regCascadeTxid, int64(task.service.config.CascadeRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		log.WithContext(ctx).WithError(err).Error("error waiting for Reg TXID confirmations")
		return task.regCascadeTxid, "", errors.Errorf("wait reg-nft ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validated Cascade Reg TXID: ",
		StatusString:  task.regCascadeTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	log.WithContext(ctx).Info("Cascade Registrtion Ticket confirmed, now activating it..")

	id, err := task.service.InsertActivationAttempt(types.ActivationAttempt{
		FileID:              task.Request.FileID,
		ActivationAttemptAt: time.Now().UTC(),
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error inserting activation attempt")
		return task.regCascadeTxid, "", errors.Errorf("error inserting activation attempt: %w", err)
	}
	task.actAttemptID = id
	// activate cascade ticket registered at previous step by SN
	activateTxID, err := task.activateActionTicket(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error activating action ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorActivatingTicket)
		return task.regCascadeTxid, "", errors.Errorf("active action ticket: %w", err)
	}
	task.StatusLog[common.FieldActivateTicketTxnID] = activateTxID
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Cascade Action Activated - Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Cascade Ticket activated. Activation ticket txid: %s, waiting for confirmations now", activateTxID)
	// Wait until activateTxID is valid
	err = task.service.pastelHandler.WaitTxidValid(ctx, activateTxID, int64(task.service.config.CascadeActTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error waiting for activated Reg TXID confirmations")
		return task.regCascadeTxid, activateTxID, errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activated Cascade Action Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Cascade Activation ticket is confirmed. Activation ticket txid: %s", activateTxID)

	return task.regCascadeTxid, activateTxID, nil

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
		if someNode == nil {
			return fmt.Errorf("node is nil - list of nodes: %s", task.MeshHandler.Nodes.String())
		}

		cascadeNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode

		image := task.Request.Image
		if image == nil {
			return errors.New("image is found to be nil")
		}

		if cascadeNode == nil {
			return errors.New("cascadeNode is found to be nil")
		}

		group.Go(func() error {
			err := cascadeNode.UploadAsset(gctx, image)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", someNode.String()).Error("upload image with thumbnail failed")
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

	// Debug: Check if root context 'ctx' is cancelled
	go func() {
		<-ctx.Done()
		log.Println("Root context 'ctx' was cancelled:", ctx.Err())
	}()

	rqidsFile := task.RqHandler.RQIDsFile
	encoderParams := task.RqHandler.RQEncodeParams

	group, gctx := errgroup.WithContext(ctx)

	// Debug: Check if derived context 'gctx' is cancelled
	go func() {
		<-gctx.Done()
	}()

	for _, someNode := range task.MeshHandler.Nodes {
		cascadeRegNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxid, err := cascadeRegNode.SendSignedTicket(gctx, task.serializedTicket, task.creatorSignature, rqidsFile, encoderParams)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", cascadeRegNode).Error("send signed ticket failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
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
	err := group.Wait()
	if err != nil {
		log.Println("Error from goroutine group:", err)
	}
	return err
}

func (task *CascadeRegistrationTask) activateActionTicket(ctx context.Context) (string, error) {
	request := pastel.ActivateActionRequest{
		RegTxID:          task.regCascadeTxid,
		BlockNum:         task.creatorBlockHeight,
		Fee:              task.registrationFee,
		PastelID:         task.Request.AppPastelID,
		Passphrase:       task.Request.AppPastelIDPassphrase,
		SpendableAddress: task.Request.SpendableAddress,
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
func (task *CascadeRegistrationTask) Download(ctx context.Context, hashOnly bool) error {
	// add the task to the worker queue, and worker will process the task in the background
	log.WithContext(ctx).WithField("cascade_tx_id", task.regCascadeTxid).Info("Downloading has been started")
	taskID := task.downloadService.AddTask(&nft.DownloadPayload{Txid: task.regCascadeTxid, Pid: task.Request.AppPastelID, Key: task.Request.AppPastelIDPassphrase},
		pastel.ActionTypeCascade, hashOnly)
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
	task := common.NewWalletNodeTask(logPrefix, service.historyDB)

	checkMinBalance := true
	fileBytes, _ := request.Image.Bytes()
	fee, err := service.pastelHandler.GetEstimatedCascadeFee(context.TODO(), utils.GetFileSizeInMB(fileBytes))
	if err != nil {
		log.WithContext(context.Background()).WithError(err).Error("NewCascadeRegTask: error getting estimated fee")
		checkMinBalance = false
	}

	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &RegisterCascadeNodeMaker{},
		PastelHandler: service.pastelHandler,
		NodeClient:    service.nodeClient,
		LogRequestID:  task.ID(),
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay:        service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:          service.config.ConnectToNodeTimeout,
			AcceptNodesTimeout:            service.config.AcceptNodesTimeout,
			MinSNs:                        service.config.NumberSuperNodes,
			PastelID:                      request.AppPastelID,
			Passphrase:                    request.AppPastelIDPassphrase,
			RequireSNAgreementOnMNTopList: true,
			CheckMinBalance:               checkMinBalance,
			MinBalance:                    int64(fee),
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
