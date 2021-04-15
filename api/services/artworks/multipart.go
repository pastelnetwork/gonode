package artworks

import (
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"strings"

	artworks "github.com/pastelnetwork/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/walletnode/api/log"
)

const contentTypePrefix = "image/"

// UploadImageDecoderFunc implements the multipart decoder for service "artworks"
// endpoint "UploadImage". The decoder must populate the argument p after encoding.
func UploadImageDecoderFunc(reader *multipart.Reader, p **artworks.ImageUploadPayload) error {
	var res artworks.ImageUploadPayload

	part, err := reader.NextPart()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Errorf("could not read next part, %s", err)
		return errInternalServerError("Could not read multipart")
	}

	_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
	if err != nil {
		log.Debugf("could not parse Content-Disposition, %s", err)
		return errBadRequest("Could not parse Content-Disposition")
	}

	contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
	if err != nil {
		log.Debugf("could not parse Content-Type, %s", err)
		return errBadRequest("Could not parse Content-Type")
	}

	if !strings.HasPrefix(contentType, contentTypePrefix) {
		log.Debugf("wrong mediatype '%s', only '%s' types are allowed", contentType, contentTypePrefix)
		return errBadRequest("Content-Type is not an image type")
	}

	if params["name"] == "file" {
		bytes, err := ioutil.ReadAll(part)
		if err != nil {
			log.Debugf("could not read multipart file, %s", err)
			return errBadRequest("Could not read multipart file")
		}
		res.File = bytes
	}

	*p = &res
	return nil
}
