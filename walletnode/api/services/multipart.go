package services

import (
	"context"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/server"
)

const contentTypePrefix = "image/"

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

			switch part.FormName() {
			case "filepath":
				// image sent as file path

				data, err := ioutil.ReadAll(part)
				if err != nil {
					return artworks.MakeInternalServerError(errors.Errorf("could not read multipart \"filepath\": %w", err))
				}

				filepath := string(data)
				if _, err := os.Stat(filepath); os.IsNotExist(err) {
					continue
				}
				res.Filepath = &filepath

			case "file":
				// image sent as binary data

				if res.Filepath != nil || part.Header.Get("Content-Type") == "" {
					continue
				}

				contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
				if err != nil {
					return artworks.MakeBadRequest(errors.Errorf("could not parse Content-Type: %w", err))
				}

				if !strings.HasPrefix(contentType, contentTypePrefix) {
					return artworks.MakeBadRequest(errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix))
				}

				fileID, _ := random.String(16, random.Base62Chars)
				filename := filepath.Join(service.workDir, fileID+filepath.Ext(part.FileName()))

				file, err := os.Create(filename)
				if err != nil {
					return artworks.MakeInternalServerError(errors.Errorf("failed to create temp file %q: %w", filename, err))
				}
				defer file.Close()

				if _, err := io.Copy(file, part); err != nil {
					return artworks.MakeInternalServerError(errors.Errorf("failed to write data to %q: %w", filename, err))
				}
				log.WithContext(ctx).Debugf("Uploaded image to %q", filename)
				res.Filepath = &filename

				service.wg.Add(1)
				go func() {
					defer service.wg.Done()

					select {
					case <-ctx.Done():
					case <-time.After(service.imageTTL):
					}

					os.Remove(filename)
					log.WithContext(ctx).Debugf("Removed image to %q", filename)
				}()

			}
		}

		if res.Filepath == nil {
			return artworks.MakeBadRequest(errors.New("neither specified \"file\" nor \"filepath\""))
		}

		*p = &res
		return nil
	}
}
