package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc/metadata"
)

// RegisterArtowrk represents grpc service for registration artowrk.
type RegisterArtowrk struct {
	*artworkregister.Service
}

func (service *RegisterArtowrk) Task(ctx context.Context) (*artworkregister.Task, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("not found metadata")
	}

	mdVals := md.Get(proto.MetadataKeyTaskID)
	if len(mdVals) == 0 {
		return nil, errors.Errorf("not found %q in metadata", proto.MetadataKeyTaskID)
	}
	taskID := mdVals[0]
	log.WithContext(ctx).WithField("taskID", taskID).Debugf("Metadata")

	task := service.Service.Task(taskID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", taskID)
	}

	return task, nil
}

// NewRegisterArtowrk returns a new RegisterArtowrk instance.
func NewRegisterArtowrk(service *artworkregister.Service) *RegisterArtowrk {
	return &RegisterArtowrk{
		Service: service,
	}
}
