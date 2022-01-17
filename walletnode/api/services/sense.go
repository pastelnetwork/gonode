package services

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	goahttp "goa.design/goa/v3/http"
)

// Sense - Sense service
type Sense struct {
	*Common
	register *senseregister.Service
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *Sense) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := sense.NewEndpoints(service)

	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
		&websocket.Upgrader{},
		nil,
		SenseUploadImageDecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// UploadImage - Uploads an image and return unique image id
func (service *Sense) UploadImage(ctx context.Context, p *sense.UploadImagePayload) (res *sense.Image, err error) {
	if p.Filename == nil {
		return nil, sense.MakeBadRequest(errors.New("file not specified"))
	}

	id, expiry, err := service.register.ImageHandler.StoreFileNameIntoStorage(ctx, p.Filename)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	res = &sense.Image{
		ImageID:   id,
		ExpiresIn: expiry,
	}

	return res, nil
}

// ActionDetails - Starts an action data task
func (service *Sense) ActionDetails(ctx context.Context, p *sense.ActionDetailsPayload) (res *sense.ActionDetailResult, err error) {
	// get imageData from storage based on imageID
	imgData, err := service.register.ImageHandler.GetImgData(p.ImageID)
	if err != nil {
		return nil, sense.MakeInternalServerError(errors.Errorf("get image data: %w", err))
	}

	ImgSizeInMb := int64(len(imgData)) / (1024 * 1024)

	// Validate image signature
	err = service.register.PastelHandler.VerifySignature(ctx, imgData, p.ActionDataSignature, p.PastelID)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	fee, err := service.register.PastelHandler.GetEstimatedActionFee(ctx, ImgSizeInMb)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	// Return data
	res = &sense.ActionDetailResult{
		EstimatedFee: fee,
	}

	return res, nil
}

// StartProcessing - Starts a processing image task
func (service *Sense) StartProcessing(_ context.Context, p *sense.StartProcessingPayload) (res *sense.StartProcessingResult, err error) {

	taskID, err := service.register.AddTask(p)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	res = &sense.StartProcessingResult{
		TaskID: taskID,
	}

	return res, nil
}

// RegisterTaskState - Registers a task state
func (service *Sense) RegisterTaskState(ctx context.Context, p *sense.RegisterTaskStatePayload, stream sense.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.GetTask(p.TaskID)
	if task == nil {
		return sense.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			res := &sense.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return sense.MakeInternalServerError(err)
			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// NewSense returns the swagger OpenAPI implementation.
func NewSense(register *senseregister.Service) *Sense {
	return &Sense{
		Common:   NewCommon(),
		register: register,
	}
}
