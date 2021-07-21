// Code generated by goa v3.3.1, DO NOT EDIT.
//
// userdatas HTTP client CLI support package
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	goa "goa.design/goa/v3/pkg"
)

// BuildProcessUserdataPayload builds the payload for the userdatas
// processUserdata endpoint from CLI flags.
func BuildProcessUserdataPayload(userdatasProcessUserdataBody string) (*userdatas.ProcessUserdataPayload, error) {
	var err error
	var body ProcessUserdataRequestBody
	{
		err = json.Unmarshal([]byte(userdatasProcessUserdataBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"artist_pastelid\": \"jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS\",\n      \"artist_pastelid_passphrase\": \"qwerasdf1234\",\n      \"avatar_image\": {\n         \"file\": \"TmloaWwgcmVydW0gZXQu\"\n      },\n      \"biography\": \"I\\'m a digital artist based in Paris, France. ...\",\n      \"categories\": \"Voluptas voluptates totam aut sed quo consequatur.\",\n      \"cover_photo\": {\n         \"file\": \"TmloaWwgcmVydW0gZXQu\"\n      },\n      \"facebook_link\": \"https://www.facebook.com/Williams_Scottish\",\n      \"location\": \"New York, US\",\n      \"native_currency\": \"USD\",\n      \"primary_language\": \"English\",\n      \"realname\": \"Williams Scottish\",\n      \"twitter_link\": \"https://www.twitter.com/@Williams_Scottish\"\n   }'")
		}
		if body.Realname != nil {
			if utf8.RuneCountInString(*body.Realname) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.realname", *body.Realname, utf8.RuneCountInString(*body.Realname), 256, false))
			}
		}
		if body.FacebookLink != nil {
			if utf8.RuneCountInString(*body.FacebookLink) > 128 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.facebook_link", *body.FacebookLink, utf8.RuneCountInString(*body.FacebookLink), 128, false))
			}
		}
		if body.TwitterLink != nil {
			if utf8.RuneCountInString(*body.TwitterLink) > 128 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.twitter_link", *body.TwitterLink, utf8.RuneCountInString(*body.TwitterLink), 128, false))
			}
		}
		if body.NativeCurrency != nil {
			if utf8.RuneCountInString(*body.NativeCurrency) < 3 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.native_currency", *body.NativeCurrency, utf8.RuneCountInString(*body.NativeCurrency), 3, true))
			}
		}
		if body.NativeCurrency != nil {
			if utf8.RuneCountInString(*body.NativeCurrency) > 3 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.native_currency", *body.NativeCurrency, utf8.RuneCountInString(*body.NativeCurrency), 3, false))
			}
		}
		if body.Location != nil {
			if utf8.RuneCountInString(*body.Location) > 256 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.location", *body.Location, utf8.RuneCountInString(*body.Location), 256, false))
			}
		}
		if body.PrimaryLanguage != nil {
			if utf8.RuneCountInString(*body.PrimaryLanguage) > 30 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.primary_language", *body.PrimaryLanguage, utf8.RuneCountInString(*body.PrimaryLanguage), 30, false))
			}
		}
		if body.Biography != nil {
			if utf8.RuneCountInString(*body.Biography) > 1024 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.biography", *body.Biography, utf8.RuneCountInString(*body.Biography), 1024, false))
			}
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.artist_pastelid", body.ArtistPastelID, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(body.ArtistPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.artist_pastelid", body.ArtistPastelID, utf8.RuneCountInString(body.ArtistPastelID), 86, true))
		}
		if utf8.RuneCountInString(body.ArtistPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.artist_pastelid", body.ArtistPastelID, utf8.RuneCountInString(body.ArtistPastelID), 86, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &userdatas.ProcessUserdataPayload{
		Realname:                 body.Realname,
		FacebookLink:             body.FacebookLink,
		TwitterLink:              body.TwitterLink,
		NativeCurrency:           body.NativeCurrency,
		Location:                 body.Location,
		PrimaryLanguage:          body.PrimaryLanguage,
		Categories:               body.Categories,
		Biography:                body.Biography,
		ArtistPastelID:           body.ArtistPastelID,
		ArtistPastelIDPassphrase: body.ArtistPastelIDPassphrase,
	}
	if body.AvatarImage != nil {
		v.AvatarImage = marshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(body.AvatarImage)
	}
	if body.CoverPhoto != nil {
		v.CoverPhoto = marshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(body.CoverPhoto)
	}

	return v, nil
}
