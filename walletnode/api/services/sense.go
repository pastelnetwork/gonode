package services

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	goahttp "goa.design/goa/v3/http"
)

type Sense struct {
	*Common
	register *senseregister.Service
	db       storage.KeyValue
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

// UploadImage - Uploads an image and return unique image id
func (service *Sense) UploadImage(_ context.Context, p *sense.UploadImagePayload) (res *sense.Image, err error) {
	if p.Filename == nil {
		return nil, sense.MakeBadRequest(errors.New("file not specified"))
	}

	// Generate unique image_id and write to db
	id, _ := random.String(8, random.Base62Chars)
	if err := service.db.Set(id, []byte(*p.Filename)); err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	// Set storage expire time
	file, err := service.register.Storage.File(*p.Filename)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}
	file.RemoveAfter(service.imageTTL)

	res = &sense.Image{
		ImageID:   id,
		ExpiresIn: time.Now().Add(service.imageTTL).Format(time.RFC3339),
	}

	return res, nil
}

// StartTask - Starts a action data task
func (service *Sense) ActionDetails(ctx context.Context, p *sense.ActionDetailsPayload) (res *sense.ActionDetailResult, err error) {
	// get image filename from storage based on image_id
	filename, err := service.db.Get(p.ImageID)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	// get image data from storage
	file, err := service.register.Storage.File(string(filename))
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	// Validate image signature
	err = service.register.VerifyImageSignature(ctx, file, p.ActionDataSignature, p.PastelID)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	// get estimated fee
	ticket := senseregister.GetEstimatedFeeRequest{
		Image: file,
	}

	fee, err := service.register.GetEstimatedFee(ctx, &ticket)
	if err != nil {
		return nil, sense.MakeInternalServerError(err)
	}

	// Return data
	res = &sense.ActionDetailResult{
		EstimatedFee: fee,
	}

	return res, nil
}

// NewOpenAPI returns the swagger OpenAPI implementation.
func NewSense(register *senseregister.Service) *Sense {
	return &Sense{
		Common:   NewCommon(),
		register: register,
		db:       memory.NewKeyValue(),
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
