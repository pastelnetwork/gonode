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

// RegisterJobState streams the state of the registration process.
func (service *serviceArtwork) RegisterJobState(ctx context.Context, p *artworks.RegisterJobStatePayload, stream artworks.RegisterJobStateServerStream) (err error) {
	defer stream.Close()

	job := service.artwork.Job(p.JobID)
	if job == nil {
		return artworks.MakeNotFound(errors.Errorf("invalid jobId: %d", p.JobID))
	}

	sub, err := job.State.Subscribe()
	if err != nil {
		return artworks.MakeInternalServerError(err)
	}
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-sub.Done():
			return nil
		case msg := <-sub.Msg():
			res := &artworks.JobState{
				Date:   msg.CreatedAt.Format(time.RFC3339),
				Status: msg.Status.String(),
			}
			if err := stream.Send(res); err != nil {
				return artworks.MakeInternalServerError(err)
			}

		}
	}
}

// RegisterJob returns a single job.
func (service *serviceArtwork) RegisterJob(ctx context.Context, p *artworks.RegisterJobPayload) (res *artworks.Job, err error) {
	job := service.artwork.Job(p.JobID)
	if job == nil {
		return nil, artworks.MakeNotFound(errors.Errorf("invalid jobId: %d", p.JobID))
	}

	res = &artworks.Job{
		ID:     p.JobID,
		Status: job.State.Latest().Status.String(),
	}

	for _, msg := range job.State.All() {
		res.States = append(res.States, &artworks.JobState{
			Date:   msg.CreatedAt.Format(time.RFC3339),
			Status: msg.Status.String(),
		})
	}

	return res, nil
}

// RegisterJob returns list of all jobs.
func (service *serviceArtwork) RegisterJobs(ctx context.Context) (res artworks.JobCollection, err error) {
	jobs := service.artwork.Jobs()
	for _, job := range jobs {
		res = append(res, &artworks.Job{
			ID:     job.ID(),
			Status: job.State.Latest().Status.String(),
		})
	}
	return res, nil
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

	jobID, err := service.artwork.Register(ctx, artwork)
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
