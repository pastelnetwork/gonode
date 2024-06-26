package selfhealing

import (
	"context"
	"sync"

	"github.com/pastelnetwork/gonode/supernode/services/download"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "self-healing-task"
)

type signatures struct {
	SNSignature1 []byte
	SNSignature2 []byte
	SNSignature3 []byte
}

// SHTask : Self healing task will manage response to self healing challenges requests
type SHTask struct {
	FingerprintsHandler *mixins.FingerprintsHandler
	*common.SuperNodeTask
	*common.StorageHandler
	*SHService
	//response message mutex to avoid race conditions
	responseMessageMu sync.Mutex
	downloadTask      *download.NftDownloadingTask
}

// Run : RunHelper's cleanup function is currently nil as WIP will determine what needs to be cleaned.
func (task *SHTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.RemoveArtifacts)
}

// RemoveArtifacts : Cleanup function defined here, can be filled in later
func (task *SHTask) RemoveArtifacts() {
}

// NewSHTask returns a new Task instance.
func NewSHTask(service *SHService) *SHTask {
	task := &SHTask{
		SuperNodeTask: common.NewSuperNodeTask(logPrefix, service.historyDB),
		StorageHandler: common.NewStorageHandler(service.P2PClient, rqgrpc.NewClient(),
			service.config.RaptorQServiceAddress, service.config.RqFilesDir, service.rqstore),
		SHService:           service,
		FingerprintsHandler: mixins.NewFingerprintsHandler(service.pastelHandler),
		downloadTask:        download.NewNftDownloadingTask(service.downloadService),
	}

	return task
}

func (task *SHTask) getRQSymbolIDs(ctx context.Context, id string, rqIDsData []byte) (rqIDs []string, err error) {
	fileContent, err := utils.Decompress(rqIDsData)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("SymbolIDsFileId", id).Warn("Decompress compressed symbol IDs file failed")
	}

	log.WithContext(ctx).WithField("Content", string(fileContent)).Debugf("symbol IDs file")

	var rqData []byte
	for i := 0; i < len(fileContent); i++ {
		if fileContent[i] == pastel.SeparatorByte {
			rqData = fileContent[:i]
			if i+1 >= len(fileContent) {
				return rqIDs, errors.New("invalid rqIDs data")
			}
			break
		}
	}

	rqDataJSON, err := utils.B64Decode(rqData)
	if err != nil {
		return rqIDs, errors.Errorf("decode %s failed: %w", string(rqData), err)
	}

	file := rqnode.RawSymbolIDFile{}
	if err := json.Unmarshal(rqDataJSON, &file); err != nil {
		log.WithContext(ctx).WithError(err).WithField("Content", string(fileContent)).
			WithField("file", string(rqIDsData)).Error("rq: parsing symbolID file failure")

		return rqIDs, errors.Errorf("parsing file: %s - file content: %s - err: %w", string(rqIDsData), string(fileContent), err)
	}

	return file.SymbolIdentifiers, nil
}
