package artworks

import (
	"context"

	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/go-commons/random"
	"github.com/pastelnetwork/walletnode/api/gen/artworks"
)

type service struct{}

// Add a new ticket and return its ID.
func (service *service) Register(ctx context.Context, p *artworks.RegisterPayload) (res string, err error) {
	log.WithContext(ctx).Debug("artworks.register")
	return
}

func (service *service) UploadImage(ctx context.Context, p *artworks.ImageUploadPayload) (res string, err error) {
	log.WithContext(ctx).Debug("artworks.upload-image")

	random.String(8, random.Base62Chars)

	log.WithContext(ctx).Debugf("bytes: %d", len(p.File))
	return "", errBadRequest()
}

// New returns the artworks service implementation.
func New() artworks.Service {
	return &service{}
}
