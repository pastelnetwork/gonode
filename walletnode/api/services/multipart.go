package services

import (
	"io/ioutil"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
	"bytes"
	"github.com/mitchellh/mapstructure"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/server"

	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	mdlserver "github.com/pastelnetwork/gonode/walletnode/api/gen/http/userdatas/server"
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


// UserdatasProcessUserdataDecoderFunc implements the multipart decoder for service "userdatas" endpoint "/update".
// The decoder must populate the argument p after encoding.
func UserdatasProcessUserdataDecoderFunc(ctx context.Context, service *Userdata) mdlserver.UserdatasProcessUserdataDecoderFunc {
	return func(reader *multipart.Reader, p **userdatas.ProcessUserdataPayload) error {

		var response *userdatas.ProcessUserdataPayload
		formValues := make(map[string]interface{})

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				return userdatas.MakeInternalServerError(errors.Errorf("could not read next part: %w", err))
			}

			if part.FormName() != imagePartName {
				// Process for other field that's not a file

				buffer, err := ioutil.ReadAll(part)
				if err != nil {
					return userdatas.MakeInternalServerError(errors.Errorf("could not process fields: %w", err))
				}
				formValues[part.FormName()] = string(buffer)
				log.WithContext(ctx).Debugf("Multipart process field: %q", formValues[part.FormName()])

			} else {
				// Process for the field that have a files
				contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
				if err != nil {
					return userdatas.MakeBadRequest(errors.Errorf("could not parse Content-Type: %w", err))
				}
	
				if !strings.HasPrefix(contentType, contentTypePrefix) {
					return userdatas.MakeBadRequest(errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix))
				}
	
				filename := part.FileName()
				filePart := new(bytes.Buffer)
				if _, err := io.Copy(filePart, part); err != nil {
					return userdatas.MakeInternalServerError(errors.Errorf("failed to write data to %q: %w", filename, err))
				}
				
				formValues[part.FormName()] = userdatas.UserImageUploadPayload{filePart.Bytes(), &filename}
				log.WithContext(ctx).Debugf("Multipart process image: %q", filename)
			}

			res, err := CreateFromMap (formValues) // Convert entire map object contain all form fields and its value, into a go struct
			if err != nil {
				return userdatas.MakeBadRequest(errors.Errorf("Could not convert formValues to object: %w", err))
			}

			response = &res
		}

		*p = response
		return nil
	}
}

// Convert map object contain fields to ProcessUserdataPayload object
func CreateFromMap(m map[string]interface{}) (userdatas.ProcessUserdataPayload, error) {
    var result userdatas.ProcessUserdataPayload
    err := mapstructure.Decode(m, &result)
    return result, err
}