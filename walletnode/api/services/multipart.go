package services

import (
	"bytes"
	"context"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	nftsrv "github.com/pastelnetwork/gonode/walletnode/api/gen/http/nft/server"
	sensrv "github.com/pastelnetwork/gonode/walletnode/api/gen/http/sense/server"
	nftreg "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	goa "goa.design/goa/v3/pkg"

	mdlserver "github.com/pastelnetwork/gonode/walletnode/api/gen/http/userdatas/server"
	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
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

		filename, errType, err := handleUploadImage(ctx, reader, service.register.ImageHandler.FileStorage)
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

		filename, errType, err := handleUploadImage(ctx, reader, service.register.ImageHandler.FileStorage)
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
func handleUploadImage(ctx context.Context, reader *multipart.Reader, storage *files.Storage) (string, string, error) {
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

		contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return "", "BadRequest", errors.Errorf("could not parse Content-Type: %w", err)
		}

		if !strings.HasPrefix(contentType, contentTypePrefix) {
			return "", "BadRequest", errors.Errorf("wrong mediatype %q, only %q types are allowed", contentType, contentTypePrefix)
		}

		filename = part.FileName()
		log.WithContext(ctx).Debugf("Upload image %q", filename)

		image := storage.NewFile()
		if err := image.SetFormatFromExtension(filepath.Ext(filename)); err != nil {
			return "", "BadRequest", errors.Errorf("could not set format from extension: %w", err)
		}

		//service.register.Storage.AddFile(image)
		filename = image.Name()
		log.WithContext(ctx).Debugf("Upload image new name %q", filename)

		fl, err := image.Create()
		if err != nil {
			return "", "InternalServerError", errors.Errorf("failed to create temp file %q: %w", filename, err)
		}
		defer fl.Close()

		if _, err := io.Copy(fl, part); err != nil {
			return "", "InternalServerError", errors.Errorf("failed to write data to %q: %w", filename, err)
		}
		log.WithContext(ctx).Debugf("Uploaded image to %q", filename)
	}

	return filename, "", nil
}

// UserdatasCreateUserdataDecoderFunc implements the multipart decoder for service "userdatas" endpoint "/create".
// The decoder must populate the argument p after encoding.
func UserdatasCreateUserdataDecoderFunc(ctx context.Context, _ *UserdataAPIHandler) mdlserver.UserdatasCreateUserdataDecoderFunc {
	return func(reader *multipart.Reader, p **userdatas.CreateUserdataPayload) error {

		var response *userdatas.CreateUserdataPayload
		var res userdatas.CreateUserdataPayload

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				return userdatas.MakeInternalServerError(errors.Errorf("could not read next part: %w", err))
			}

			if part.FormName() != "avatar_image" && part.FormName() != "cover_photo" {
				// Process for other field that's not a file

				buffer, err := ioutil.ReadAll(part)
				if err != nil {
					return userdatas.MakeInternalServerError(errors.Errorf("could not process fields: %w", err))
				}
				log.WithContext(ctx).Debugf("Multipart process field: %q", part.FormName())

				switch part.FormName() {
				case "artist_pastelid":
					res.UserPastelID = string(buffer)
				case "artist_pastelid_passphrase":
					res.UserPastelIDPassphrase = string(buffer)
				case "biography":
					value := string(buffer)
					res.Biography = &value
				case "categories":
					value := string(buffer)
					res.Categories = &value
				case "facebook_link":
					value := string(buffer)
					res.FacebookLink = &value
				case "location":
					value := string(buffer)
					res.Location = &value
				case "primary_language":
					value := string(buffer)
					res.PrimaryLanguage = &value
				case "native_currency":
					value := string(buffer)
					res.NativeCurrency = &value
				case "realname":
					value := string(buffer)
					res.RealName = &value
				case "twitter_link":
					value := string(buffer)
					res.TwitterLink = &value
				default:
					log.WithContext(ctx).Errorf("unknown form name: %s", part.FormName())
				}
			} else {
				// Process for the field that have a files
				filename := part.FileName()
				filePart := new(bytes.Buffer)
				if _, err := io.Copy(filePart, part); err != nil {
					return userdatas.MakeInternalServerError(errors.Errorf("failed to write data to %q: %w", filename, err))
				}

				switch part.FormName() {
				case "avatar_image":
					res.AvatarImage = &userdatas.UserImageUploadPayload{
						Content:  filePart.Bytes(),
						Filename: &filename,
					}
				case "cover_photo":
					res.CoverPhoto = &userdatas.UserImageUploadPayload{
						Content:  filePart.Bytes(),
						Filename: &filename,
					}
				default:
					log.WithContext(ctx).Errorf("unknown form name: %s", part.FormName())
				}
				log.WithContext(ctx).Debugf("Multipart process image: %q", filename)
			}
		}

		response = &res

		*p = response
		return nil
	}
}

// UserdatasUpdateUserdataDecoderFunc implements the multipart decoder for service "userdatas" endpoint "/update".
// The decoder must populate the argument p after encoding.
func UserdatasUpdateUserdataDecoderFunc(ctx context.Context, _ *UserdataAPIHandler) mdlserver.UserdatasUpdateUserdataDecoderFunc {
	return func(reader *multipart.Reader, p **userdatas.UpdateUserdataPayload) error {

		var response *userdatas.UpdateUserdataPayload
		var res userdatas.UpdateUserdataPayload

		for {
			part, err := reader.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				}
				return userdatas.MakeInternalServerError(errors.Errorf("could not read next part: %w", err))
			}

			if part.FormName() != "avatar_image" && part.FormName() != "cover_photo" {
				// Process for other field that's not a file

				buffer, err := ioutil.ReadAll(part)
				if err != nil {
					return userdatas.MakeInternalServerError(errors.Errorf("could not process fields: %w", err))
				}
				log.WithContext(ctx).Debugf("Multipart process field: %q", part.FormName())

				switch part.FormName() {
				case "artist_pastelid":
					res.UserPastelID = string(buffer)
				case "artist_pastelid_passphrase":
					res.UserPastelIDPassphrase = string(buffer)
				case "biography":
					value := string(buffer)
					res.Biography = &value
				case "categories":
					value := string(buffer)
					res.Categories = &value
				case "facebook_link":
					value := string(buffer)
					res.FacebookLink = &value
				case "location":
					value := string(buffer)
					res.Location = &value
				case "primary_language":
					value := string(buffer)
					res.PrimaryLanguage = &value
				case "native_currency":
					value := string(buffer)
					res.NativeCurrency = &value
				case "realname":
					value := string(buffer)
					res.RealName = &value
				case "twitter_link":
					value := string(buffer)
					res.TwitterLink = &value
				default:
					log.WithContext(ctx).Errorf("unknown form name: %s", part.FormName())
				}
			} else {
				// Process for the field that have a files

				filename := part.FileName()
				filePart := new(bytes.Buffer)
				if _, err := io.Copy(filePart, part); err != nil {
					return userdatas.MakeInternalServerError(errors.Errorf("failed to write data to %q: %w", filename, err))
				}

				switch part.FormName() {
				case "avatar_image":
					res.AvatarImage = &userdatas.UserImageUploadPayload{
						Content:  filePart.Bytes(),
						Filename: &filename,
					}
				case "cover_photo":
					res.CoverPhoto = &userdatas.UserImageUploadPayload{
						Content:  filePart.Bytes(),
						Filename: &filename,
					}
				default:
					log.WithContext(ctx).Errorf("unknown form name: %s", part.FormName())
				}
				log.WithContext(ctx).Debugf("Multipart process image: %q", filename)
			}
		}

		response = &res

		*p = response
		return nil
	}
}
