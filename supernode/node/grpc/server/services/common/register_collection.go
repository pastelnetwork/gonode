package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/collectionregister"
	"google.golang.org/grpc/metadata"
)

// RegisterCollection represents common grpc service for registration collection.
type RegisterCollection struct {
	*collectionregister.CollectionRegistrationService
}

// SessID retrieves SessID from the metadata.
func (service *RegisterCollection) SessID(ctx context.Context) (string, bool) {
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
func (service *RegisterCollection) TaskFromMD(ctx context.Context) (*collectionregister.CollectionRegistrationTask, error) {
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

// NewRegisterCollection returns a new RegisterCollection instance.
func NewRegisterCollection(service *collectionregister.CollectionRegistrationService) *RegisterCollection {
	return &RegisterCollection{
		CollectionRegistrationService: service,
	}
}
