package services

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/artworks/server"

	mdlserver "github.com/pastelnetwork/gonode/walletnode/api/gen/http/userdatas/server"
	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
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
			log.WithContext(ctx).Debugf("Upload image %q", filename)

			image := service.register.Storage.NewFile()
			if err := image.SetFormatFromExtension(filepath.Ext(filename)); err != nil {
				return artworks.MakeBadRequest(err)
			}
			//service.register.Storage.AddFile(image)
			filename = image.Name()
			log.WithContext(ctx).Debugf("Upload image new name %q", filename)

			fl, err := image.Create()
			if err != nil {
				return artworks.MakeInternalServerError(errors.Errorf("failed to create temp file %q: %w", filename, err))
			}
			defer fl.Close()

			if _, err := io.Copy(fl, part); err != nil {
				return artworks.MakeInternalServerError(errors.Errorf("failed to write data to %q: %w", filename, err))
			}
			log.WithContext(ctx).Debugf("Uploaded image to %q", filename)
			//log.WithContext(ctx).Debugf("Storage files %v", service.register.Storage.Files())

			res.Filename = &filename
		}

		*p = &res
		return nil
	}
}

// UserdatasCreateUserdataDecoderFunc implements the multipart decoder for service "userdatas" endpoint "/create".
// The decoder must populate the argument p after encoding.
func UserdatasCreateUserdataDecoderFunc(ctx context.Context, _ *Userdata) mdlserver.UserdatasCreateUserdataDecoderFunc {
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
					res.ArtistPastelID = string(buffer)
				case "artist_pastelid_passphrase":
					res.ArtistPastelIDPassphrase = string(buffer)
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
					res.Realname = &value
				case "twitter_link":
					value := string(buffer)
					res.TwitterLink = &value
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
func UserdatasUpdateUserdataDecoderFunc(ctx context.Context, _ *Userdata) mdlserver.UserdatasUpdateUserdataDecoderFunc {
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
					res.ArtistPastelID = string(buffer)
				case "artist_pastelid_passphrase":
					res.ArtistPastelIDPassphrase = string(buffer)
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
					res.Realname = &value
				case "twitter_link":
					value := string(buffer)
					res.TwitterLink = &value
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
				}
				log.WithContext(ctx).Debugf("Multipart process image: %q", filename)
			}
		}

		response = &res

		*p = response
		return nil
	}
}
