package collectionregister

import (
	"context"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// CollectionRegistrationTask is the task of registering new nft.
type CollectionRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler *common.MeshHandler

	service *CollectionRegistrationService
	Request *common.CollectionRegistrationRequest
}

// Run starts the task
func (task *CollectionRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *CollectionRegistrationTask) run(_ context.Context) error {
	//registration will go here
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
