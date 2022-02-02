package cascaderegister

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// Task is the task of registering new Sense.
type CascadeRegistrationTask struct {
	*common.SuperNodeTask
	*common.RegTaskHelper
	*CascadeRegistrationService

	storage *common.StorageHandler

	cascadeRegMetadata *types.ActionRegMetadata
	Ticket             *pastel.ActionTicket
	Asset              *files.File

	Oti []byte

	// signature of ticket data signed by this node's pastelID
	ownSignature []byte

	creatorSignature []byte
	registrationFee  int64

	rawRqFile []byte
	rqIDFiles [][]byte
}

func (task *CascadeRegistrationTask) removeArtifacts() {
	task.RemoveFile(task.Asset)
}

// NewSenseRegistrationTask returns a new Task instance.
func NewCascadeRegistrationTask(service *CascadeRegistrationService) *CascadeRegistrationTask {
	task := &CascadeRegistrationTask{
		SuperNodeTask:              common.NewSuperNodeTask(logPrefix),
		CascadeRegistrationService: service,
		storage: common.NewStorageHandler(service.P2PClient, service.RQClient,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir),
	}

	task.RegTaskHelper = common.NewRegTaskHelper(task.SuperNodeTask,
		task.config.PastelID, task.config.PassPhrase,
		common.NewNetworkHandler(task.SuperNodeTask, service.nodeClient,
			RegisterCascadeNodeMaker{}, service.PastelClient,
			task.config.PastelID,
			service.config.NumberConnectedNodes),
		service.PastelClient,
	)

	return task
}
