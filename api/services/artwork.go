package services

import (
	"context"
	"time"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/random"
	"github.com/pastelnetwork/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/walletnode/dao"
	"github.com/pastelnetwork/walletnode/dao/memory"
	"github.com/pastelnetwork/walletnode/services/artwork"
)

const (
	imageTTL = time.Second * 30
)

type serviceArtwork struct {
	storage dao.KeyValue
	artwork *artwork.Service
}

// Add a new ticket and return its ID.
func (service *serviceArtwork) Register(ctx context.Context, p *artworks.RegisterPayload) (res string, err error) {
	image, err := service.storage.Get([]byte(p.ImageID))
	if err == dao.ErrKeyNotFound {
		return "", artworks.MakeBadRequest(errors.Errorf("invalid image_id %q", p.ImageID))
	}
	if err != nil {
		return "", artworks.MakeInternalServerError(err)
	}

	artwork := &artwork.Artwork{
		Image:            image,
		Name:             p.Name,
		Description:      p.Description,
		Keywords:         p.Keywords,
		SeriesName:       p.SeriesName,
		IssuedCopies:     p.IssuedCopies,
		YoutubeURL:       p.YoutubeURL,
		ArtistPastelID:   p.ArtistPastelID,
		ArtistName:       p.ArtistName,
		ArtistWebsiteURL: p.ArtistWebsiteURL,
		SpendableAddress: p.SpendableAddress,
		NetworkFee:       p.NetworkFee,
	}

	err = service.artwork.Register(artwork)
	if err != nil {
		return "", artworks.MakeInternalServerError(err)
	}
	return
}

func (service *serviceArtwork) UploadImage(ctx context.Context, p *artworks.ImageUploadPayload) (res *artworks.WalletnodeImage, err error) {
	id, err := random.String(8, random.Base62Chars)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}

	if err := service.storage.Set([]byte(id), p.Bytes); err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	expiresIn := time.Now().Add(imageTTL)

	go func() {
		time.AfterFunc(time.Until(expiresIn), func() {
			service.storage.Delete([]byte(id))
		})
	}()

	res = &artworks.WalletnodeImage{
		ImageID:   id,
		ExpiresIn: expiresIn.Format(time.RFC3339),
	}
	return res, nil
}

// NewArtwork returns the artworks serviceArtwork implementation.
func NewArtwork() artworks.Service {
	return &serviceArtwork{
		storage: memory.NewKeyValue(),
		artwork: artwork.NewService(),
	}
}
