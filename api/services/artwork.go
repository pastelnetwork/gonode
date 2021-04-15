package services

import (
	"context"
	"time"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/walletnode/services/artwork"
)

const (
	imageTTL = time.Second * 30
)

type serviceArtwork struct {
	store   *artwork.Store
	artwork *artwork.Service
}

// Add a new ticket and return its ID.
func (service *serviceArtwork) Register(ctx context.Context, p *artworks.RegisterPayload) (res string, err error) {
	image := service.store.Get(p.ImageID)
	if image == nil {
		return "", artworks.MakeBadRequest(errors.Errorf("invalid image_id %q", p.ImageID))
	}

	artwork := &artwork.Artwork{
		Name:             p.Name,
		Description:      p.Description,
		Keywords:         p.Keywords,
		SeriesName:       p.SeriesName,
		IssuedCopies:     p.IssuedCopies,
		Image:            image,
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
	image := artwork.NewImage(p.Bytes)
	id := service.store.Add(image)
	expiresIn := time.Now().Add(imageTTL)

	go func() {
		time.AfterFunc(time.Until(expiresIn), func() {
			service.store.Remove(id)
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
		store:   artwork.NewStore(),
		artwork: artwork.NewService(),
	}
}
