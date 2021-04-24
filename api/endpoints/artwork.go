package endpoints

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/walletnode/api"
	"github.com/pastelnetwork/walletnode/services/artwork/register"
	"github.com/pastelnetwork/walletnode/storage"
	"github.com/pastelnetwork/walletnode/storage/memory"

	"github.com/pastelnetwork/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/walletnode/api/gen/http/artworks/server"

	goahttp "goa.design/goa/v3/http"
)

const (
	imageTTL = time.Second * 3600 // 1 hour
)

var (
	imageID uint32
)

// Artwork represents services for artworks endpoints.
type Artwork struct {
	register *register.Service
	storage  storage.KeyValue
}

// RegisterTaskState streams the state of the registration process.
func (service *Artwork) RegisterTaskState(ctx context.Context, p *artworks.RegisterTaskStatePayload, stream artworks.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.Task(p.TaskID)
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
func (service *Artwork) RegisterTask(ctx context.Context, p *artworks.RegisterTaskPayload) (res *artworks.Task, err error) {
	task := service.register.Task(p.TaskID)
	if task == nil {
		return nil, artworks.MakeNotFound(errors.Errorf("invalid taskId: %d", p.TaskID))
	}

	res = &artworks.Task{
		ID:     p.TaskID,
		Status: task.State.Latest().Status.String(),
		Ticket: toArtworkTicket(task.Ticket),
		States: toArtworkStates(task.State.All()),
	}
	return res, nil
}

// RegisterTasks returns list of all tasks.
func (service *Artwork) RegisterTasks(ctx context.Context) (res artworks.TaskCollection, err error) {
	tasks := service.register.Tasks()
	for _, task := range tasks {
		res = append(res, &artworks.Task{
			ID:     task.ID,
			Status: task.State.Latest().Status.String(),
			Ticket: toArtworkTicket(task.Ticket),
		})
	}
	return res, nil
}

// Register runs registers process for the new artwork.
func (service *Artwork) Register(ctx context.Context, p *artworks.RegisterPayload) (res *artworks.RegisterResult, err error) {
	ticket := fromRegisterPayload(p)

	key := strconv.Itoa(p.ImageID)
	ticket.Image, err = service.storage.Get(key)
	if err == storage.ErrKeyNotFound {
		return nil, artworks.MakeBadRequest(errors.Errorf("invalid image_id: %q", p.ImageID))
	}
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}

	taskID, err := service.register.Register(ctx, ticket)
	if err != nil {
		return nil, artworks.MakeInternalServerError(err)
	}
	res = &artworks.RegisterResult{
		TaskID: taskID,
	}
	return res, nil
}

// UploadImage uploads an image and return unique image id.
func (service *Artwork) UploadImage(ctx context.Context, p *artworks.UploadImagePayload) (res *artworks.Image, err error) {
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

// Mount configures the mux to serve the artworks endpoints.
func (service *Artwork) Mount(mux goahttp.Muxer) goahttp.Server {
	endpoints := artworks.NewEndpoints(service)
	srv := server.New(endpoints, nil, goahttp.RequestDecoder, goahttp.ResponseEncoder, api.ErrorHandler, nil, &websocket.Upgrader{}, nil, UploadImageDecoderFunc)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.Infof("%s HTTP %q mounted on %s %s", api.LogPrefix, m.Method, m.Verb, m.Pattern)
	}

	return srv
}

// NewArtwork returns the artworks Artwork implementation.
func NewArtwork(register *register.Service) *Artwork {
	return &Artwork{
		register: register,
		storage:  memory.NewKeyValue(),
	}
}
