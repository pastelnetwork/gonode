package services

import (
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/pastelnetwork/go-commons/errors"
	artworks "github.com/pastelnetwork/walletnode/api/gen/artworks"
)

const contentTypePrefix = "image/"

// UploadImageDecoderFunc implements the multipart decoder for service "artworks"
// endpoint "UploadImage". The decoder must populate the argument p after encoding.
func UploadImageDecoderFunc(reader *multipart.Reader, p **artworks.UploadImagePayload) error {
	var res artworks.UploadImagePayload

	part, err := reader.NextPart()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return artworks.MakeInternalServerError(errors.Errorf("could not read next part, %s", err))
	}

	_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
	if err != nil {
		return artworks.MakeBadRequest(errors.Errorf("could not parse Content-Disposition, %s", err))
	}

	contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err != nil {
		return artworks.MakeBadRequest(errors.Errorf("could not parse Content-Type, %s", err))
	}

	if !strings.HasPrefix(contentType, contentTypePrefix) {
		return artworks.MakeBadRequest(errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix))
	}

	if params["name"] == "file" {
		bytes, err := ioutil.ReadAll(part)
		if err != nil {
			return artworks.MakeInternalServerError(errors.Errorf("could not read multipart file, %s", err))
		}
		res.Bytes = bytes
	}

	*p = &res
	return nil
}
