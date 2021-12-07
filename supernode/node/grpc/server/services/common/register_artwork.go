package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc/metadata"
)

// RegisterArtwork represents common grpc service for registration artwork.
type RegisterArtwork struct {
	*artworkregister.Service
}

// SessID retrieves SessID from the metadata.
func (service *RegisterArtwork) SessID(ctx context.Context) (string, bool) {
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
func (service *RegisterArtwork) TaskFromMD(ctx context.Context) (*artworkregister.Task, error) {
	sessID, ok := service.SessID(ctx)
	if !ok {
		return nil, errors.New("could not find sessID in metadata")
	}

	task := service.Task(sessID)
	if task == nil {
		return nil, errors.Errorf("could not find %q task", sessID)
	}
	return task, nil
}

// NewRegisterArtwork returns a new RegisterArtwork instance.
func NewRegisterArtwork(service *artworkregister.Service) *RegisterArtwork {
	return &RegisterArtwork{
		Service: service,
	}
}
