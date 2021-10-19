package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/externalstorage"
	"google.golang.org/grpc/metadata"
)

// ExternalStorage represents common grpc service for registration artwork.
type ExternalStorage struct {
	*externalstorage.Service
}

// SessID retrieves SessID from the metadata.
func (service *ExternalStorage) SessID(ctx context.Context) (string, bool) {
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
func (service *ExternalStorage) TaskFromMD(ctx context.Context) (*externalstorage.Task, error) {
	sessID, ok := service.SessID(ctx)
	if !ok {
		return nil, errors.New("not found sessID in metadata")
	}

	task := service.Task(sessID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", sessID)
	}
	return task, nil
}

// NewExternalStorage returns a new ExternalStorage instance.
func NewExternalStorage(service *externalstorage.Service) *ExternalStorage {
	return &ExternalStorage{
		Service: service,
	}
}
