package collectionregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	//CollectionFinalAllowedBlockHeightDaysConversion is the formula to convert days to no of blocks
	CollectionFinalAllowedBlockHeightDaysConversion = 24 * (60 / 2.5)
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
	creatorSignature   []byte
	signatures         map[string][]byte

	collectionTicket *pastel.CollectionTicket
	collectionTXID   string
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

	/* Step 3,4: Find tops supernodes and validate top 3 SNs and create mesh network of 3 SNs */
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

	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	return nil
}

func (task *CollectionRegistrationTask) createCollectionTicket(_ context.Context) error {
	ticket := &pastel.CollectionTicket{
		CollectionTicketVersion:                 1,
		CollectionName:                          task.Request.CollectionName,
		BlockNum:                                task.creatorBlockHeight,
		BlockHash:                               task.creatorBlockHash,
		Creator:                                 task.Request.AppPastelID,
		ListOfPastelIDsOfAuthorizedContributors: task.Request.ListOfPastelIDsOfAuthorizedContributors,
		MaxCollectionEntries:                    uint(task.Request.MaxCollectionEntries),
		CollectionItemCopyCount:                 uint(task.Request.CollectionItemCopyCount),
		CollectionFinalAllowedBlockHeight:       task.creatorBlockHeight + uint(task.Request.CollectionFinalAllowedBlockHeight*CollectionFinalAllowedBlockHeightDaysConversion),
		Royalty:                                 task.Request.Royalty / 100,
		Green:                                   task.Request.Green,
		AppTicketData:                           pastel.AppTicket{
			//MaxPermittedOpenNSFWScore: ,
			//MinimumSimilarityScoreToFirstEntryInCollection:
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

	task.creatorSignature, err = task.service.pastelHandler.PastelClient.Sign(ctx, data, task.Request.AppPastelID, task.Request.AppPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign collection ticket %w", err)
	}

	task.signatures = make(map[string][]byte)
	task.signatures["principal"] = task.creatorSignature

	return nil
}

func (task *CollectionRegistrationTask) removeArtifacts() {
	//remove artifacts to be implemented here
}

// Error returns task err
func (task *CollectionRegistrationTask) Error() error {
	return task.WalletNodeTask.Error()
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
