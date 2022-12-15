package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
	"google.golang.org/grpc/metadata"
)

// StorageChallenge represents common grpc service for registration NFTs.
type StorageChallenge struct {
	*storagechallenge.SCService
}

// SessID retrieves SessID from the metadata.
func (service *StorageChallenge) SessID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	mdVals := md.Get(proto.MetadataKeySessID)
	if len(mdVals) == 0 {
		return "", false
	}
	return mdVals[0], true
}

// TaskFromMD returns task by SessID from the metadata.
func (service *StorageChallenge) TaskFromMD(ctx context.Context) (*storagechallenge.SCTask, error) {
	sessID, ok := service.SessID(ctx)
	if !ok {
		log.WithContext(ctx).Error("Error not a valid session id")

		return nil, errors.New("not found sessID in metadata")
	}

	task := service.Task(sessID)
	if task == nil {
		log.WithContext(ctx).Error("task not found")
		return nil, errors.Errorf("not found %q task", sessID)
	}

	return task, nil
}

// NewStorageChallenge returns a new StorageChallenge instance.
func NewStorageChallenge(service *storagechallenge.SCService) *StorageChallenge {
	return &StorageChallenge{
		SCService: service,
	}
}
