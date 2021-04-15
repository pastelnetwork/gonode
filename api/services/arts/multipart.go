package arts

import (
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/pastelnetwork/go-commons/errors"
	arts "github.com/pastelnetwork/walletnode/api/gen/arts"
)

const contentTypePrefix = "image/"

// UploadImageDecoderFunc implements the multipart decoder for service "arts"
// endpoint "UploadImage". The decoder must populate the argument p after encoding.
func UploadImageDecoderFunc(reader *multipart.Reader, p **arts.ImageUploadPayload) error {
	var res arts.ImageUploadPayload

	part, err := reader.NextPart()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return errors.New(err)
	}

	_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
	if err != nil {
		return errors.Errorf("could not parse multipart header, %s", err)
	}

	contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err != nil {
		return errors.Errorf("could not parse mediatype, %s", err)
	}

	if !strings.HasPrefix(contentType, contentTypePrefix) {
		return errors.Errorf("wrong mediatype '%s', only '%s' types are allowed", contentType, contentTypePrefix)
	}

	if params["name"] == "file" {
		bytes, err := ioutil.ReadAll(part)
		if err != nil {
			return errors.Errorf("could not read multipart file, %s", err)
		}
		res.File = bytes
	}

	*p = &res
	return nil
}
