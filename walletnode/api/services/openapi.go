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
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/openapi/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/openapi"
	goahttp "goa.design/goa/v3/http"
)

type OpenAPI struct {
	*Common
	imageTTL time.Duration
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *OpenAPI) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := openapi.NewEndpoints(service)

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
func (service *OpenAPI) UploadImage(_ context.Context, p *openapi.UploadImagePayload) (res *openapi.Image, err error) {
	if p.Filename == nil {
		return nil, openapi.MakeBadRequest(errors.New("file not specified"))
	}

	// TODO: Process the image
	txid, _ := random.String(64, random.Base62Chars)

	res = &openapi.Image{
		Txid: txid,
	}
	return res, nil
}

// NewOpenAPI returns the swagger OpenAPI implementation.
func NewOpenAPI() *OpenAPI {
	return &OpenAPI{
		Common:   NewCommon(),
		imageTTL: defaultImageTTL,
	}
}

// OpenApiUploadImageDecoderFunc implements the multipart decoder for service "artworks" endpoint "UploadImage".
// The decoder must populate the argument p after encoding.
func OpenApiUploadImageDecoderFunc(ctx context.Context, service *OpenAPI) server.OpenapiUploadImageDecoderFunc {
	return func(reader *multipart.Reader, p **openapi.UploadImagePayload) error {
		var res openapi.UploadImagePayload

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				return openapi.MakeInternalServerError(errors.Errorf("could not read next part: %w", err))
			}

			if part.FormName() != imagePartName {
				continue
			}

			contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
			if err != nil {
				return openapi.MakeBadRequest(errors.Errorf("could not parse Content-Type: %w", err))
			}

			if !strings.HasPrefix(contentType, contentTypePrefix) {
				return openapi.MakeBadRequest(errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix))
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
