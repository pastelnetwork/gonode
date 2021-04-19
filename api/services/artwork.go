package services

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/walletnode/services/artwork"
	"github.com/pastelnetwork/walletnode/services/artwork/register"
	"github.com/pastelnetwork/walletnode/storage"
	"github.com/pastelnetwork/walletnode/storage/memory"
)

const (
	imageTTL = time.Second * 3600 // 1 hour
)

var (
	imageID uint32
)

type serviceArtwork struct {
	artwork *artwork.Service
	storage storage.KeyValue
}

// RegisterTaskState streams the state of the registration process.
func (service *serviceArtwork) RegisterTaskState(ctx context.Context, p *artworks.RegisterTaskStatePayload, stream artworks.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.artwork.Task(p.TaskID)
	if task == nil {
		return artworks.MakeNotFound(errors.Errorf("invalid taskId: %d", p.TaskID))
	}

	sub, err := task.State.Subscribe()
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
			res := &artworks.TaskState{
				Date:   msg.CreatedAt.Format(time.RFC3339),
				Status: msg.Status.String(),
			}
			if err := stream.Send(res); err != nil {
				return artworks.MakeInternalServerError(err)
			}

		}
	}
}

// RegisterTask returns a single task.
func (service *serviceArtwork) RegisterTask(ctx context.Context, p *artworks.RegisterTaskPayload) (res *artworks.Task, err error) {
	task := service.artwork.Task(p.TaskID)
	if task == nil {
		return nil, artworks.MakeNotFound(errors.Errorf("invalid taskId: %d", p.TaskID))
	}

	res = &artworks.Task{
		ID:     p.TaskID,
		Status: task.State.Latest().Status.String(),
	}

	for _, msg := range task.State.All() {
		res.States = append(res.States, &artworks.TaskState{
			Date:   msg.CreatedAt.Format(time.RFC3339),
			Status: msg.Status.String(),
		})
	}

	return res, nil
}

// RegisterTask returns list of all tasks.
func (service *serviceArtwork) RegisterTasks(ctx context.Context) (res artworks.TaskCollection, err error) {
	tasks := service.artwork.Tasks()
	for _, task := range tasks {
		res = append(res, &artworks.Task{
			ID:     task.ID(),
			Status: task.State.Latest().Status.String(),
		})
	}
	return res, nil
}

// Register runs registers process for the new artwork.
func (service *serviceArtwork) Register(ctx context.Context, p *artworks.RegisterPayload) (res *artworks.RegisterResult, err error) {
	key := strconv.Itoa(p.ImageID)

	image, err := service.storage.Get(key)
	if err == storage.ErrKeyNotFound {
		return nil, artworks.MakeBadRequest(errors.Errorf("invalid image_id: %q", p.ImageID))
	}
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}

	ticket := &register.Ticket{
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

	taskID, err := service.artwork.Register(ctx, ticket)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	res = &artworks.RegisterResult{
		TaskID: taskID,
	}
	return res, nil
}

// UploadImage uploads an image and return unique image id.
func (service *serviceArtwork) UploadImage(ctx context.Context, p *artworks.UploadImagePayload) (res *artworks.Image, err error) {
	id := int(atomic.AddUint32(&imageID, 1))
	key := strconv.Itoa(id)

	if err := service.storage.Set(key, p.Bytes); err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	expiresIn := time.Now().Add(imageTTL)

	go func() {
		time.AfterFunc(time.Until(expiresIn), func() {
			service.storage.Delete(key)
		})
	}()

	res = &artworks.Image{
		ImageID:   id,
		ExpiresIn: expiresIn.Format(time.RFC3339),
	}
	return res, nil
}

// NewArtwork returns the artworks serviceArtwork implementation.
func NewArtwork(artworks *artwork.Service) artworks.Service {
	return &serviceArtwork{
		artwork: artworks,
		storage: memory.NewKeyValue(),
	}
}
