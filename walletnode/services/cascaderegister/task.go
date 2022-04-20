package cascaderegister

import (
	"context"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// CascadeRegistrationTask is the task of registering new nft.
type CascadeRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler *common.MeshHandler
	RqHandler   *mixins.RQHandler

	service *CascadeRegistrationService
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

	// send registration metadata
	if err := task.sendActionMetadata(ctx); err != nil {
		return errors.Errorf("send registration metadata: %w", err)
	}

	if err := task.uploadImage(ctx); err != nil {
		return errors.Errorf("upload image: %w", err)
	}

	// connect to rq serivce to get rq symbols identifier
	if err := task.RqHandler.GenRQIdentifiersFiles(ctx, task.Request.Image,
		task.creatorBlockHash, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase); err != nil {
		task.UpdateStatus(common.StatusErrorGenRaptorQSymbolsFailed)
		return errors.Errorf("gen RaptorQ symbols' identifiers: %w", err)
	}

	// calculate hash of data
	imgBytes, err := task.Request.Image.Bytes()
	if err != nil {
		return errors.Errorf("convert image to byte stream %w", err)
	}
	if task.dataHash, err = utils.Sha3256hash(imgBytes); err != nil {
		return errors.Errorf("hash encoded image: %w", err)
	}

	fileDataInMb := float64(len(imgBytes)) / (1024 * 1024)
	fee, err := task.service.pastelHandler.GetEstimatedCascadeFee(ctx, fileDataInMb)
	if err != nil {
		return errors.Errorf("getting estimated fee %w", err)
	}
	task.registrationFee = int64(fee)

	if err := task.createCascadeTicket(ctx); err != nil {
		return errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("sign cascade ticket: %w", err)
	}

	// UPLOAD signed ticket to supernodes to validate and register action with the network
	if err := task.uploadSignedTicket(ctx); err != nil {
		return errors.Errorf("send signed cascade ticket: %w", err)
	}
	task.UpdateStatus(common.StatusTicketAccepted)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validating Cascade Reg TXID: ",
		StatusString:  task.regCascadeTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.regCascadeTxid, int64(task.service.config.CascadeRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validated Cascade Reg TXID: ",
		StatusString:  task.regCascadeTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	// activate cascade ticket registered at previous step by SN
	activateTxID, err := task.activateActionTicket(newCtx)
	if err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("active action ticket: %w", err)
	}
	log.Debugf("Active action ticket txid: %s", activateTxID)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activating Cascade Action Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	// Wait until activateTxID is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, activateTxID, int64(task.service.config.CascadeActTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
	if err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	log.Debugf("Active txid is confirmed")
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activated Cascade Action Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	// Send ActionAct request to primary node
	if err := task.uploadActionAct(newCtx, task.regCascadeTxid); err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("upload action act: %w", err)
	}

	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
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

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		cascadeRegNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() (err error) {
			err = cascadeRegNode.SendRegMetadata(ctx, regMetadata)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", cascadeRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

func (task *CascadeRegistrationTask) uploadImage(ctx context.Context) error {
	group, _ := errgroup.WithContext(ctx)

	for _, someNode := range task.MeshHandler.Nodes {
		cascadeNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			err := cascadeNode.UploadAsset(ctx, task.Request.Image)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", someNode).Error("upload image with thumbnail failed")
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
		APITicketData: &pastel.APICascadeTicket{
			DataHash: task.dataHash,
			RQIc:     task.RqHandler.RQIDsIc,
			RQMax:    task.service.config.RQIDsMax,
			RQIDs:    task.RqHandler.RQIDs,
			RQOti:    task.RqHandler.RQEncodeParams.Oti,
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

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		cascadeRegNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxid, err := cascadeRegNode.SendSignedTicket(ctx, task.serializedTicket, task.creatorSignature, rqidsFile, encoderParams)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", cascadeRegNode).Error("send signed ticket failed")
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

// uploadActionAct uploads action act to primary node
func (task *CascadeRegistrationTask) uploadActionAct(ctx context.Context, activateTxID string) error {
	group, _ := errgroup.WithContext(ctx)

	for _, someNode := range task.MeshHandler.Nodes {
		cascadeRegNode, ok := someNode.SuperNodeAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CascadeRegistrationNode", someNode.String())
		}

		someNode := someNode
		if someNode.IsPrimary() {
			group.Go(func() error {
				return cascadeRegNode.SendActionAct(ctx, activateTxID)
			})
		}
	}
	return group.Wait()
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
		WalletNodeTask: task,
		service:        service,
		Request:        request,
		MeshHandler:    common.NewMeshHandler(meshHandlerOpts),
		RqHandler: mixins.NewRQHandler(service.rqClient, service.pastelHandler,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir,
			service.config.NumberRQIDSFiles, service.config.RQIDsMax),
	}
}
