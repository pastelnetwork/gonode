package services

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/cascade/server"
	"github.com/pastelnetwork/gonode/walletnode/services/cascaderegister"
	goahttp "goa.design/goa/v3/http"
)

// cascade - cascade service
type Cascade struct {
	*Common
	register *cascaderegister.Service
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *Cascade) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
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
		CascadeUploadImageDecoderFunc,
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// UploadImage - Uploads an image and return unique image id
func (service *Cascade) UploadImage(ctx context.Context, p *cascade.UploadImagePayload) (res *cascade.Image, err error) {
	if p.Filename == nil {
		return nil, cascade.MakeBadRequest(errors.New("file not specified"))
	}

	res, err = service.register.StoreImage(ctx, p)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	return res, nil
}

// ActionDetails - Starts a action data task
func (service *Cascade) ActionDetails(ctx context.Context, p *cascade.ActionDetailsPayload) (res *cascade.ActionDetailResult, err error) {
	// get imageData from storage based on imageID
	imgData, err := service.register.GetImgData(p.ImageID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(errors.Errorf("get image data: %w", err))
	}

	ImgSizeInMb := int64(len(imgData)) / (1024 * 1024)

	// Validate image signature
	err = service.register.VerifyImageSignature(ctx, imgData, p.ActionDataSignature, p.PastelID)
	if err != nil {
		return nil, cascade.MakeInternalServerError(err)
	}

	fee, err := service.register.GetEstimatedFee(ctx, ImgSizeInMb)
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
func (service *Cascade) StartProcessing(_ context.Context, p *cascade.StartProcessingPayload) (res *cascade.StartProcessingResult, err error) {

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
func (service *Cascade) RegisterTaskState(ctx context.Context, p *cascade.RegisterTaskStatePayload, stream cascade.RegisterTaskStateServerStream) (err error) {
	defer stream.Close()

	task := service.register.Task(p.TaskID)
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
				return artworks.MakeInternalServerError(err)
			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// NewCascade returns the swagger OpenAPI implementation.
func NewCascade(register *cascaderegister.Service) *Cascade {
	return &Cascade{
		Common:   NewCommon(),
		register: register,
	}
}

// CascadeUploadImageDecoderFunc implements the multipart decoder function for the UploadImage service
func CascadeUploadImageDecoderFunc(reader *multipart.Reader, p **cascade.UploadImagePayload) error {
	var res cascade.UploadImagePayload

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filename := part.FileName()
		if filename != "" {
			slurp, err := ioutil.ReadAll(part)
			if err != nil {
				return errors.Errorf("read file: %w", err)
			}

			res.Filename = &filename
			res.Bytes = slurp
			break
		}
	}

	*p = &res
	return nil
}
