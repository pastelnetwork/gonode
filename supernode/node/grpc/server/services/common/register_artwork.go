package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc/metadata"
)

// RegisterArtowrk represents common grpc service for registration artowrk.
type RegisterArtowrk struct {
	*artworkregister.Service
}

// ConnID retrieves ConnID from the metadata.
func (service *RegisterArtowrk) ConnID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("not found metadata")
	}

	mdVals := md.Get(proto.MetadataKeyConnID)
	if len(mdVals) == 0 {
		return "", errors.Errorf("not found %q in metadata", proto.MetadataKeyConnID)
	}
	return mdVals[0], nil
}

// TaskFromMD returns task by ConnID from the metadata.
func (service *RegisterArtowrk) TaskFromMD(ctx context.Context) (*artworkregister.Task, error) {
	connID, err := service.ConnID(ctx)
	if err != nil {
		return nil, err
	}

	task := service.Task(connID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", connID)
	}
	return task, nil
}

// NewRegisterArtowrk returns a new RegisterArtowrk instance.
func NewRegisterArtowrk(service *artworkregister.Service) *RegisterArtowrk {
	return &RegisterArtowrk{
		Service: service,
	}
}
