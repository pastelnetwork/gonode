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
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	goahttp "goa.design/goa/v3/http"
)

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
		SenseUploadImageDecoderFunc,
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

	res, err = service.register.StoreImage(ctx, p)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	return res, nil
}

// ActionDetails - Starts a action data task
func (service *Sense) ActionDetails(ctx context.Context, p *sense.ActionDetailsPayload) (res *sense.ActionDetailResult, err error) {
	// get imageData from storage based on imageID
	imgData, err := service.register.GetImgData(p.ImageID)
	if err != nil {
		return nil, sense.MakeInternalServerError(errors.Errorf("get image data: %w", err))
	}

	ImgSizeInMb := int64(len(imgData)) / (1024 * 1024)

	// Validate image signature
	err = service.register.VerifyImageSignature(ctx, imgData, p.ActionDataSignature, p.PastelID)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	fee, err := service.register.GetEstimatedFee(ctx, ImgSizeInMb)
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
func (service *Sense) StartProcessing(ctx context.Context, p *sense.StartProcessingPayload) (res *sense.StartProcessingResult, err error) {

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

	task := service.register.Task(p.TaskID)
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
				return artworks.MakeInternalServerError(err)
			}

			if status.IsFinal() {
				return nil
			}
		}
	}
}

// NewOpenAPI returns the swagger OpenAPI implementation.
func NewSense(register *senseregister.Service) *Sense {
	return &Sense{
		Common:   NewCommon(),
		register: register,
	}
}

// SenseUploadImageDecoderFunc implements the multipart decoder function for the UploadImage service
func SenseUploadImageDecoderFunc(reader *multipart.Reader, p **sense.UploadImagePayload) error {
	var res sense.UploadImagePayload

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
