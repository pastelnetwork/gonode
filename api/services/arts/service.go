package arts

import (
	"context"

	"github.com/pastelnetwork/go-commons/log"
	arts "github.com/pastelnetwork/walletnode/api/gen/arts"
)

type service struct{}

// Add a new ticket and return its ID.
func (service *service) Register(ctx context.Context, p *arts.RegisterPayload) (res string, err error) {
	log.WithContext(ctx).Debug("arts.register")
	return
}

func (service *service) UploadImage(ctx context.Context, p *arts.ImageUploadPayload) (res string, err error) {
	log.WithContext(ctx).Debug("arts.upload-image")

	log.WithContext(ctx).Debugf("bytes: %d", len(p.File))
	return "", errBadRequest()
}

// New returns the arts service implementation.
func New() arts.Service {
	return &service{}
}
