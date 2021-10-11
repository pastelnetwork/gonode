package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/externaldupedetection"
	"google.golang.org/grpc/metadata"
)

// ExternalDupeDetection represents common grpc service for registration artwork.
type ExternalDupeDetection struct {
	*externaldupedetection.Service
}

// SessID retrieves SessID from the metadata.
func (service *ExternalDupeDetection) SessID(ctx context.Context) (string, bool) {
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
func (service *ExternalDupeDetection) TaskFromMD(ctx context.Context) (*externaldupedetection.Task, error) {
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

// NewExternalDupeDetection returns a new ExternalDupeDetection instance.
func NewExternalDupeDetection(service *externaldupedetection.Service) *ExternalDupeDetection {
	return &ExternalDupeDetection{
		Service: service,
	}
}
