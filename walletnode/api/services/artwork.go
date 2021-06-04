package services

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/server"

	goahttp "goa.design/goa/v3/http"
)

const (
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// Artwork represents services for artworks endpoints.
type Artwork struct {
	*Common
	register *artworkregister.Service
	db       storage.KeyValue
	imageTTL time.Duration
}

// RegisterTaskState streams the state of the registration process.
func (service *Artwork) RegisterTaskState(ctx context.Context, p *artworks.RegisterTaskStatePayload, stream artworks.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.Task(p.TaskID)
	if task == nil {
		return artworks.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			if status.IsFinal() {
				return nil
			}
			res := &artworks.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return artworks.MakeInternalServerError(err)
			}

		}
	}
}

// RegisterTask returns a single task.
func (service *Artwork) RegisterTask(_ context.Context, p *artworks.RegisterTaskPayload) (res *artworks.Task, err error) {
	task := service.register.Task(p.TaskID)
	if task == nil {
		return nil, artworks.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	res = &artworks.Task{
		ID:     p.TaskID,
		Status: task.Status().String(),
		Ticket: toArtworkTicket(task.Ticket),
		States: toArtworkStates(task.StatusHistory()),
	}
	return res, nil
}

// RegisterTasks returns list of all tasks.
func (service *Artwork) RegisterTasks(_ context.Context) (res artworks.TaskCollection, err error) {
	tasks := service.register.Tasks()
	for _, task := range tasks {
		res = append(res, &artworks.Task{
			ID:     task.ID(),
			Status: task.Status().String(),
			Ticket: toArtworkTicket(task.Ticket),
		})
	}
	return res, nil
}

// Register runs registers process for the new artwork.
func (service *Artwork) Register(_ context.Context, p *artworks.RegisterPayload) (res *artworks.RegisterResult, err error) {
	ticket := fromRegisterPayload(p)

	filename, err := service.db.Get(p.ImageID)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}

	file, err := service.register.Storage.File(string(filename))
	if err != nil {
		return nil, artworks.MakeBadRequest(errors.Errorf("invalid image_id: %q", p.ImageID))
	}
	ticket.Image = file

	taskID := service.register.AddTask(ticket)
	res = &artworks.RegisterResult{
		TaskID: taskID,
	}
	return res, nil
}

// UploadImage uploads an image and return unique image id.
func (service *Artwork) UploadImage(_ context.Context, p *artworks.UploadImagePayload) (res *artworks.Image, err error) {
	if p.Filename == nil {
		return nil, artworks.MakeBadRequest(errors.New("file not specified"))
	}

	id, _ := random.String(8, random.Base62Chars)
	if err := service.db.Set(id, []byte(*p.Filename)); err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}

	file, err := service.register.Storage.File(*p.Filename)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	file.RemoveAfter(service.imageTTL)

	res = &artworks.Image{
		ImageID:   id,
		ExpiresIn: time.Now().Add(service.imageTTL).Format(time.RFC3339),
	}
	return res, nil
}

// Mount configures the mux to serve the artworks endpoints.
func (service *Artwork) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := artworks.NewEndpoints(service)
	srv := server.New(endpoints, nil, goahttp.RequestDecoder, goahttp.ResponseEncoder, api.ErrorHandler, nil, &websocket.Upgrader{}, nil, UploadImageDecoderFunc(ctx, service))
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// NewArtwork returns the artworks Artwork implementation.
func NewArtwork(register *artworkregister.Service) *Artwork {
	return &Artwork{
		Common:   NewCommon(),
		register: register,
		db:       memory.NewKeyValue(),
		imageTTL: defaultImageTTL,
	}
}
