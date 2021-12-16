package services

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	goahttp "goa.design/goa/v3/http"
)

type Sense struct {
	*Common
	register *senseregister.Service
	imageTTL time.Duration
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
		SenseUploadImageDecoderFunc(ctx, service),
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// UploadImage - Uploads an image to the blockchain.
func (service *Sense) UploadImage(_ context.Context, p *sense.UploadImagePayload) (res *sense.ImageUploadResult, err error) {
	if p.Filename == nil {
		return nil, sense.MakeBadRequest(errors.New("file not specified"))
	}

	file, err := service.register.Storage.File(string(*p.Filename))
	if err != nil {
		return nil, sense.MakeBadRequest(errors.Errorf("could not get file: %w", err))
	}

	ticket := &senseregister.Request{
		Image: file,
	}

	taskID := service.register.AddTask(ticket)
	res = &sense.ImageUploadResult{
		TaskID: taskID,
	}

	return res, nil
}

// StartTask - Starts a action data task
func (service *Sense) StartTask(_ context.Context, p *sense.StartTaskPayload) (res *sense.StartActionDataResult, err error) {
	// Validate payload

	// Process data

	// Return data
	res = &sense.StartActionDataResult{
		EstimatedFee: 100,
	}

	return res, nil
}

// NewOpenAPI returns the swagger OpenAPI implementation.
func NewSense(register *senseregister.Service) *Sense {
	return &Sense{
		Common:   NewCommon(),
		register: register,
		imageTTL: defaultImageTTL,
	}
}

// OpenApiUploadImageDecoderFunc implements the multipart decoder for service "artworks" endpoint "UploadImage".
// The decoder must populate the argument p after encoding.
func SenseUploadImageDecoderFunc(ctx context.Context, service *Sense) server.SenseUploadImageDecoderFunc {
	return func(reader *multipart.Reader, p **sense.UploadImagePayload) error {
		var res sense.UploadImagePayload

		filename, err := handleUploadImage(ctx, reader, service.register.Storage)
		if err != nil {
			return sense.MakeBadRequest(err)
		}

		res.Filename = &filename
		*p = &res
		return nil
	}
}
