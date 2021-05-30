package services

import (
	"context"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/server"
)

const (
	contentTypePrefix = "image/"
	imagePartName     = "file"
)

// UploadImageDecoderFunc implements the multipart decoder for service "artworks" endpoint "UploadImage".
// The decoder must populate the argument p after encoding.
func UploadImageDecoderFunc(ctx context.Context, service *Artwork) server.ArtworksUploadImageDecoderFunc {
	return func(reader *multipart.Reader, p **artworks.UploadImagePayload) error {
		var res artworks.UploadImagePayload

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				return artworks.MakeInternalServerError(errors.Errorf("could not read next part: %w", err))
			}

			if part.FormName() != imagePartName {
				continue
			}

			contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
			if err != nil {
				return artworks.MakeBadRequest(errors.Errorf("could not parse Content-Type: %w", err))
			}

			if !strings.HasPrefix(contentType, contentTypePrefix) {
				return artworks.MakeBadRequest(errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix))
			}

			filename := part.FileName()

			image := service.register.Storage.NewFile()
			if err := image.SetFormatFromExtension(filepath.Ext(filename)); err != nil {
				return artworks.MakeBadRequest(err)
			}
			filename = image.Name()

			fl, err := image.Create()
			if err != nil {
				return artworks.MakeInternalServerError(errors.Errorf("failed to create temp file %q: %w", filename, err))
			}
			defer fl.Close()

			if _, err := io.Copy(fl, part); err != nil {
				return artworks.MakeInternalServerError(errors.Errorf("failed to write data to %q: %w", filename, err))
			}
			log.WithContext(ctx).Debugf("Uploaded image to %q", filename)

			res.Filename = &filename
		}

		*p = &res
		return nil
	}
}
