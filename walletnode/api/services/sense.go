package services

import (
	"context"
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	goahttp "goa.design/goa/v3/http"
)

type Sense struct {
	*Common
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
		OpenApiUploadImageDecoderFunc(ctx, service),
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

	// TODO: Register task_id
	task_id, _ := random.String(64, random.Base62Chars)

	res = &sense.ImageUploadResult{
		TaskID: task_id,
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
func NewOpenAPI() *Sense {
	return &Sense{
		Common:   NewCommon(),
		imageTTL: defaultImageTTL,
	}
}

// OpenApiUploadImageDecoderFunc implements the multipart decoder for service "artworks" endpoint "UploadImage".
// The decoder must populate the argument p after encoding.
func OpenApiUploadImageDecoderFunc(ctx context.Context, service *Sense) server.SenseUploadImageDecoderFunc {
	return func(reader *multipart.Reader, p **sense.UploadImagePayload) error {
		var res sense.UploadImagePayload

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				return sense.MakeInternalServerError(errors.Errorf("could not read next part: %w", err))
			}

			if part.FormName() != imagePartName {
				continue
			}

			contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
			if err != nil {
				return sense.MakeBadRequest(errors.Errorf("could not parse Content-Type: %w", err))
			}

			if !strings.HasPrefix(contentType, contentTypePrefix) {
				return sense.MakeBadRequest(errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix))
			}

			// TODO: doing something with the uploaded image
			filename := part.FileName()
			log.WithContext(ctx).Debugf("Upload image %q", filename)

			res.Filename = &filename
		}

		*p = &res
		return nil
	}
}
