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
	imageTTL = time.Second * 3600 // 1 hour
)

type serviceArtwork struct {
	storage dao.KeyValue
	artwork *artwork.Service
}

// RegisterStatus streams the job status of the new artwork registration.
func (service *serviceArtwork) RegisterStatus(ctx context.Context, p *artworks.RegisterStatusPayload, stream artworks.RegisterStatusServerStream) (err error) {
	job := service.artwork.RegisterJob(p.JobID)
	if job == nil {
		return artworks.MakeNotFound(errors.Errorf("invalid jobId: %d", p.JobID))
	}
	subscription := job.Subscribe()

	defer stream.Close()
	for {
		select {
		case <-ctx.Done():
			return nil
		case sub, ok := <-subscription:
			if !ok {
				return nil
			}
			res := &artworks.Job{
				ID:     job.ID,
				Status: sub.Status.String(),
			}
			if err := stream.Send(res); err != nil {
				return artworks.MakeBadRequest(err)
			}

		}
	}
}

// Register runs registers process for the new artwork.
func (service *serviceArtwork) Register(ctx context.Context, p *artworks.RegisterPayload) (res *artworks.RegisterResult, err error) {
	image, err := service.storage.Get(p.ImageID)
	if err == dao.ErrKeyNotFound {
		return nil, artworks.MakeBadRequest(errors.Errorf("invalid image_id: %q", p.ImageID))
	}
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
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

	jobID, err := service.artwork.Register(artwork)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	res = &artworks.RegisterResult{
		JobID: jobID,
	}
	return res, nil
}

// UploadImage uploads an image and return unique image id.
func (service *serviceArtwork) UploadImage(ctx context.Context, p *artworks.UploadImagePayload) (res *artworks.Image, err error) {
	id, err := random.String(8, random.Base62Chars)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}

	if err := service.storage.Set(id, p.Bytes); err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	expiresIn := time.Now().Add(imageTTL)

	go func() {
		time.AfterFunc(time.Until(expiresIn), func() {
			service.storage.Delete(id)
		})
	}()

	res = &artworks.Image{
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
