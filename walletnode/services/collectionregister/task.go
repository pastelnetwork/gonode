package collectionregister

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	//CollectionFinalAllowedBlockHeightDaysConversion is the formula to convert days to no of blocks
	CollectionFinalAllowedBlockHeightDaysConversion = 24 * (60 / 2.5)

	//DateTimeFormat following the go convention for request timestamp
	DateTimeFormat = "2006:01:02 15:04:05"

	collectionRegFee = 1000
)

// CollectionRegistrationTask is the task of registering new nft.
type CollectionRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler *common.MeshHandler

	service *CollectionRegistrationService
	Request *common.CollectionRegistrationRequest

	// data to create ticket
	creatorBlockHeight uint
	creatorBlockHash   string
	creationTimestamp  string
	creatorSignature   []byte

	collectionTicket *pastel.CollectionTicket
	serializedTicket []byte
	collectionTXID   string
	burnTxID         string

	// only set to true for unit tests
	skipPrimaryNodeTxidVerify bool
}

// Run starts the task
func (task *CollectionRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *CollectionRegistrationTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.WithContext(ctx).Info("Setting up mesh with Top Supernodes")
	task.StatusLog[common.FieldTaskType] = "Collection Registration"

	creatorBlockHeight, creatorBlockHash, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error setting up mesh of supernodes")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorMeshSetupFailed)
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	task.creatorBlockHeight = uint(creatorBlockHeight)
	task.creatorBlockHash = creatorBlockHash
	task.StatusLog[common.FieldBlockHeight] = creatorBlockHeight
	task.creationTimestamp = time.Now().Format(DateTimeFormat)

	log.WithContext(ctx).Info("Mesh of supernodes have been established")

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)
	log.WithContext(ctx).Info("uploading data to supernodes")

	if err := task.createCollectionTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating collection ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorCreatingTicket)
		return errors.Errorf("create ticket: %w", err)
	}
	log.WithContext(ctx).Info("created collection reg ticket")

	// sign ticket with creator signature
	if err := task.signTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.StatusLog[common.FieldRegTicketTxnID] = task.collectionTXID
		task.UpdateStatus(common.StatusErrorSigningTicket)
		return errors.Errorf("sign collection ticket: %w", err)
	}
	log.WithContext(ctx).Info("signed collection reg ticket")

	task.UpdateStatus(common.StatusPreburntRegistrationFee)
	if err := task.preBurnRegistrationFee(ctx); err != nil {
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		log.WithContext(ctx).WithError(err).Error("error pre-burning collectionreg fee")

		return err
	}
	log.WithContext(ctx).Info("pre-burned 10% of collection reg fee")

	log.WithContext(ctx).Info("upload signed collection reg ticket to SNs")
	if err := task.uploadSignedTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error uploading signed ticket")

		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorUploadingTicket)
		return errors.Errorf("send signed collection ticket: %w", err)
	}
	task.StatusLog[common.FieldRegTicketTxnID] = task.collectionTXID

	task.UpdateStatus(common.StatusTicketAccepted)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validating Collection Reg TXID: ",
		StatusString:  task.collectionTXID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Collection Reg Ticket registered. Collection Registration Ticket txid: %s", task.collectionTXID)
	log.WithContext(ctx).Info("Closing SNs connections")

	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	log.WithContext(ctx).Infof("Waiting Confirmations for Collection Reg Ticket - Ticket txid: %s", task.collectionTXID)
	newCtx := log.ContextWithPrefix(context.Background(), "Collection")
	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.collectionTXID, int64(task.service.config.CollectionRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

		log.WithContext(ctx).WithError(err).Error("error getting confirmations")
		return errors.Errorf("wait reg-collection ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validated Collection Reg TXID: ",
		StatusString:  task.collectionTXID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})

	log.WithContext(ctx).Info("Collection Reg Ticket confirmed, Activating Collection Reg Ticket")
	// activate Collection ticket registered at previous step by SN
	activateTxID, err := task.activateCollectionTicket(newCtx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error activating collection ticket")

		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorActivatingTicket)
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("active collection ticket: %w", err)
	}
	task.StatusLog[common.FieldActivateTicketTxnID] = activateTxID
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activating Collection Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Collection ticket activated. Activation ticket txid: %s", activateTxID)
	log.WithContext(ctx).Infof("Waiting Confirmations for Collection Activation Ticket - Ticket txid: %s", activateTxID)

	// Wait until activateTxID is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, activateTxID, int64(task.service.config.CollectionActTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting confirmations for activation")
		_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)
		return errors.Errorf("wait activate txid valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Activated Collection Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("Collection Activation ticket is confirmed. Activation ticket txid: %s", activateTxID)

	return nil
}

func (task *CollectionRegistrationTask) createCollectionTicket(_ context.Context) error {
	ticket := &pastel.CollectionTicket{
		CollectionTicketVersion:                 1,
		CollectionName:                          task.Request.CollectionName,
		ItemType:                                task.Request.ItemType,
		BlockNum:                                task.creatorBlockHeight,
		BlockHash:                               task.creatorBlockHash,
		Creator:                                 task.Request.AppPastelID,
		ListOfPastelIDsOfAuthorizedContributors: task.Request.ListOfPastelIDsOfAuthorizedContributors,
		MaxCollectionEntries:                    uint(task.Request.MaxCollectionEntries),
		CollectionItemCopyCount:                 uint(task.Request.CollectionItemCopyCount),
		CollectionFinalAllowedBlockHeight:       task.creatorBlockHeight + uint(task.Request.NoOfDaysToFinalizeCollection*CollectionFinalAllowedBlockHeightDaysConversion),
		Royalty:                                 task.Request.Royalty / 100,
		Green:                                   task.Request.Green,
		AppTicketData: pastel.AppTicket{
			MaxPermittedOpenNSFWScore:                      task.Request.MaxPermittedOpenNSFWScore,
			MinimumSimilarityScoreToFirstEntryInCollection: task.Request.MinimumSimilarityScoreToFirstEntryInCollection,
		},
	}

	task.collectionTicket = ticket
	return nil
}

func (task *CollectionRegistrationTask) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeCollectionTicket(task.collectionTicket)
	if err != nil {
		return errors.Errorf("encode collection ticket %w", err)
	}

	task.creatorSignature, err = task.service.pastelHandler.PastelClient.SignCollectionTicket(ctx, data, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign collection ticket %w", err)
	}
	task.serializedTicket = data

	return nil
}

func (task *CollectionRegistrationTask) preBurnRegistrationFee(ctx context.Context) (err error) {
	address := task.Request.SpendableAddress
	balance, err := task.service.pastelHandler.PastelClient.GetBalance(ctx, address)
	if err != nil {
		task.UpdateStatus(common.StatusErrorCheckBalance)
		return errors.Errorf("get balance of address(%s): %w", address, err)
	}

	if balance < collectionRegFee {
		task.UpdateStatus(common.StatusErrorInsufficientBalance)
		return errors.Errorf("not enough PSL - balance(%f) < registration-fee(%d)", balance, collectionRegFee)
	}

	task.burnTxID, err = task.service.pastelHandler.BurnSomeCoins(ctx, address,
		collectionRegFee, 10)
	if err != nil {
		task.UpdateStatus(common.StatusErrorPreburnRegFeeGetTicketID)
		return fmt.Errorf("burn some coins: %w", err)
	}

	task.StatusLog[common.FieldBurnTxnID] = task.burnTxID

	return nil
}

// uploadSignedTicket uploads Collection ticket  and its signature to super nodes
func (task *CollectionRegistrationTask) uploadSignedTicket(ctx context.Context) error {
	if task.serializedTicket == nil {
		return errors.Errorf("uploading ticket: serializedTicket is empty")
	}
	if task.creatorSignature == nil {
		return errors.Errorf("uploading ticket: creatorSignature is empty")
	}

	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		CollectionRegNode, ok := someNode.SuperNodeAPIInterface.(*CollectionRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not CollectionRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxID, err := CollectionRegNode.SendTicketForSignature(gctx, task.serializedTicket, task.creatorSignature, task.burnTxID)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", CollectionRegNode).Error("send signed ticket failed")
				return err
			}
			if !someNode.IsPrimary() && ticketTxID != "" && !task.skipPrimaryNodeTxidCheck() {
				return errors.Errorf("receive response %s from secondary node %s", ticketTxID, someNode.PastelID())
			}
			if someNode.IsPrimary() {
				if ticketTxID == "" {
					return errors.Errorf("primary node - %s, returned empty txid", someNode.PastelID())
				}
				task.collectionTXID = ticketTxID
			}
			return nil
		})
	}
	return group.Wait()
}

func (task *CollectionRegistrationTask) activateCollectionTicket(ctx context.Context) (string, error) {
	request := pastel.ActivateCollectionRequest{
		RegTxID:    task.collectionTXID,
		BlockNum:   int(task.creatorBlockHeight),
		Fee:        collectionRegFee,
		PastelID:   task.Request.AppPastelID,
		Passphrase: task.Request.AppPastelIDPassphrase,
	}

	return task.service.pastelHandler.PastelClient.ActivateCollectionTicket(ctx, request)
}

func (task *CollectionRegistrationTask) removeArtifacts() {
	//remove artifacts to be implemented here
}

// Error returns task err
func (task *CollectionRegistrationTask) Error() error {
	return task.WalletNodeTask.Error()
}

func (task *CollectionRegistrationTask) skipPrimaryNodeTxidCheck() bool {
	return task.skipPrimaryNodeTxidVerify || os.Getenv("INTEGRATION_TEST_ENV") == "true"
}

// NewCollectionRegistrationTask returns a new CollectionRegistrationTask instance.
// TODO: make config interface and pass it instead of individual items
func NewCollectionRegistrationTask(service *CollectionRegistrationService, request *common.CollectionRegistrationRequest) *CollectionRegistrationTask {
	task := common.NewWalletNodeTask(logPrefix)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &RegisterCollectionNodeMaker{},
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

	return &CollectionRegistrationTask{
		WalletNodeTask: task,
		service:        service,
		Request:        request,
		MeshHandler:    common.NewMeshHandler(meshHandlerOpts),
	}
}
