package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	"google.golang.org/grpc/metadata"
)

// RegisterCascade represents common grpc service for registration cascade.
type RegisterCascade struct {
	*cascaderegister.Service
}

// SessID retrieves SessID from the metadata.
func (service *RegisterCascade) SessID(ctx context.Context) (string, bool) {
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
func (service *RegisterCascade) TaskFromMD(ctx context.Context) (*cascaderegister.Task, error) {
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

// NewRegisterCascade returns a new RegisterCascade instance.
func NewRegisterCascade(service *cascaderegister.Service) *RegisterCascade {
	return &RegisterCascade{
		Service: service,
	}
}
