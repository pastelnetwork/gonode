package services

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/cascade/server"
	"github.com/pastelnetwork/gonode/walletnode/services/cascaderegister"
	goahttp "goa.design/goa/v3/http"
)

// CascadeAPIHandler - CascadeAPIHandler service
type CascadeAPIHandler struct {
	*Common
	register *cascaderegister.CascadeRegistrationService
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *CascadeAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := cascade.NewEndpoints(service)

	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
		&websocket.Upgrader{},
		nil,
		CascadeUploadImageDecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// UploadImage - Uploads an image and return unique image id
func (service *CascadeAPIHandler) UploadImage(ctx context.Context, p *cascade.UploadImagePayload) (res *cascade.Image, err error) {
	if p.Filename == nil {
		return nil, cascade.MakeBadRequest(errors.New("file not specified"))
	}

	id, expiry, err := service.register.StoreFile(ctx, p.Filename)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	res = &cascade.Image{
		ImageID:   id,
		ExpiresIn: expiry,
	}

	return res, nil
}

// ActionDetails - Starts an action data task
func (service *CascadeAPIHandler) ActionDetails(ctx context.Context, p *cascade.ActionDetailsPayload) (res *cascade.ActionDetailResult, err error) {

	fee, err := service.register.ValidateDetailsAndCalculateFee(ctx, p.ImageID, p.ActionDataSignature, p.PastelID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	// Return data
	res = &cascade.ActionDetailResult{
		EstimatedFee: fee,
	}

	return res, nil
}

// StartProcessing - Starts a processing image task
func (service *CascadeAPIHandler) StartProcessing(_ context.Context, p *cascade.StartProcessingPayload) (res *cascade.StartProcessingResult, err error) {
	taskID, err := service.register.AddTask(p)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	res = &cascade.StartProcessingResult{
		TaskID: taskID,
	}

	return res, nil
}

// RegisterTaskState - Registers a task state
func (service *CascadeAPIHandler) RegisterTaskState(ctx context.Context, p *cascade.RegisterTaskStatePayload, stream cascade.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.GetTask(p.TaskID)
	if task == nil {
		return cascade.MakeNotFound(errors.Errorf("invalid taskId: %s", p.TaskID))
	}

	sub := task.SubscribeStatus()

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			res := &cascade.TaskState{
				Date:   status.CreatedAt.Format(time.RFC3339),
				Status: status.String(),
			}
			if err := stream.Send(res); err != nil {
				return cascade.MakeInternalServerError(err)
			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// NewCascadeAPIHandler returns the swagger OpenAPI implementation.
func NewCascadeAPIHandler(register *cascaderegister.CascadeRegistrationService) *CascadeAPIHandler {
	return &CascadeAPIHandler{
		Common:   NewCommon(),
		register: register,
	}
}
