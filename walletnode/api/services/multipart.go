package services

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"

	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	cassrv "github.com/pastelnetwork/gonode/walletnode/api/gen/http/cascade/server"
	nftsrv "github.com/pastelnetwork/gonode/walletnode/api/gen/http/nft/server"
	sensrv "github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	nftreg "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	goa "goa.design/goa/v3/pkg"
)

const (
	contentTypePrefix = "image/"
	imagePartName     = "file"
)

// NftRegUploadImageDecoderFunc implements the multipart decoder for service "nftreg" endpoint "UploadImage".
// The decoder must populate the argument p after encoding.
func NftRegUploadImageDecoderFunc(ctx context.Context, service *NftAPIHandler) nftsrv.NftUploadImageDecoderFunc {
	return func(reader *multipart.Reader, p **nftreg.UploadImagePayload) error {
		var res nftreg.UploadImagePayload

		filename, errType, err := handleUploadImage(ctx, reader, service.register.ImageHandler.FileStorage, true)
		if err != nil {
			return &goa.ServiceError{
				Name:    errType,
				ID:      goa.NewErrorID(),
				Message: err.Error(),
			}
		}

		res.Filename = &filename
		*p = &res
		return nil
	}
}

// SenseUploadImageDecoderFunc decodes image uploaded from request
func SenseUploadImageDecoderFunc(ctx context.Context, service *SenseAPIHandler) sensrv.SenseUploadImageDecoderFunc {
	return func(reader *multipart.Reader, p **sense.UploadImagePayload) error {
		var res sense.UploadImagePayload

		filename, errType, err := handleUploadImage(ctx, reader, service.register.ImageHandler.FileStorage, true)
		if err != nil {
			return &goa.ServiceError{
				Name:    errType,
				ID:      goa.NewErrorID(),
				Message: err.Error(),
			}
		}

		res.Filename = &filename
		*p = &res
		return nil
	}
}

// CascadeUploadAssetDecoderFunc decodes image uploaded from request
func CascadeUploadAssetDecoderFunc(ctx context.Context, service *CascadeAPIHandler) cassrv.CascadeUploadAssetDecoderFunc {
	return func(reader *multipart.Reader, p **cascade.UploadAssetPayload) error {
		var res cascade.UploadAssetPayload

		filename, errType, err := handleUploadImage(ctx, reader, service.register.ImageHandler.FileStorage, false)
		if err != nil {
			return &goa.ServiceError{
				Name:    errType,
				ID:      goa.NewErrorID(),
				Message: err.Error(),
			}
		}

		res.Filename = &filename
		*p = &res
		return nil
	}
}

// handleUploadImage -- save image to service storage
func handleUploadImage(ctx context.Context, reader *multipart.Reader, storage *files.Storage, onlyImage bool) (string, string, error) {
	var filename string

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "InternalServerError", errors.Errorf("could not read next part: %w", err)
		}

		if part.FormName() != imagePartName {
			continue
		}

		var data []byte
		data, err = ioutil.ReadAll(part)
		if err != nil {
			return "", "BadRequest", errors.Errorf("could not read part: %w", err)
		}

		contentType := mimetype.Detect(data)

		filename = part.FileName()
		log.WithContext(ctx).Debugf("Upload image %q", filename)

		image := storage.NewFile()
		if onlyImage {
			if !strings.HasPrefix(contentType.String(), contentTypePrefix) {
				return "", "BadRequest", errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix)
			}

			if err := image.SetFormatFromExtension(contentType.Extension()); err != nil {
				return "", "BadRequest", errors.Errorf("could not set format from extension: %w", err)
			}
		}

		if err := storage.Update(image.Name(), filename, image); err != nil {
			return "", "BadRequest", errors.Errorf("could not update file name: %w", err)
		}

		fl, err := image.Create()
		if err != nil {
			return "", "InternalServerError", errors.Errorf("failed to create temp file %q: %w", filename, err)
		}
		defer fl.Close()

		if _, err := io.Copy(fl, bytes.NewBuffer(data)); err != nil {
			return "", "InternalServerError", errors.Errorf("failed to write data to %q: %w", filename, err)
		}

		log.WithContext(ctx).Debugf("Uploaded image to %q", filename)
	}

	return filename, "", nil
}
