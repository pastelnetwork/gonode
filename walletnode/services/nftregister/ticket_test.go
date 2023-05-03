package nftregister

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
)

func TestFromNftRegisterPayload(t *testing.T) {
	des := "description"
	key := "key"
	series := "series"
	url := "url"
	issuedCopies := 5
	green := true
	royalty := 10.4

	type args struct {
		payload *nft.RegisterPayload
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				payload: &nft.RegisterPayload{
					Name:                      "name",
					Description:               &des,
					Keywords:                  &key,
					SeriesName:                &series,
					IssuedCopies:              &issuedCopies,
					YoutubeURL:                &url,
					CreatorPastelID:           "pastelid",
					CreatorPastelIDPassphrase: "passphrase",
					CreatorName:               "creator",
					CreatorWebsiteURL:         &url,
					SpendableAddress:          "spendaddr",
					MaximumFee:                12,
					Green:                     &green,
					Royalty:                   &royalty,
					ThumbnailCoordinate: &nft.Thumbnailcoordinate{
						TopLeftX:     0,
						TopLeftY:     0,
						BottomRightY: 640,
						BottomRightX: 480,
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			gotReq := FromNftRegisterPayload(tc.args.payload)
			assert.Equal(t, gotReq.CreatorName, tc.args.payload.CreatorName)
			assert.Equal(t, gotReq.CreatorPastelID, tc.args.payload.CreatorPastelID)
			assert.Equal(t, gotReq.CreatorPastelIDPassphrase, tc.args.payload.CreatorPastelIDPassphrase)
			assert.Equal(t, gotReq.Description, tc.args.payload.Description)
			assert.Equal(t, gotReq.Keywords, tc.args.payload.Keywords)
			assert.Equal(t, gotReq.CreatorWebsiteURL, tc.args.payload.CreatorWebsiteURL)
			assert.Equal(t, gotReq.MaximumFee, tc.args.payload.MaximumFee)
			assert.Equal(t, gotReq.SeriesName, tc.args.payload.SeriesName)
			assert.Equal(t, gotReq.Thumbnail.BottomRightX, tc.args.payload.ThumbnailCoordinate.BottomRightX)
			assert.Equal(t, gotReq.Thumbnail.BottomRightY, tc.args.payload.ThumbnailCoordinate.BottomRightY)
			assert.Equal(t, gotReq.Thumbnail.TopLeftY, tc.args.payload.ThumbnailCoordinate.TopLeftY)

		})
	}
}

func TestToNftRegisterTicket(t *testing.T) {
	des := "description"
	key := "key"
	series := "series"
	url := "url"
	issuedCopies := 5
	green := true
	royalty := 10.4

	type args struct {
		payload *NftRegistrationRequest
	}

	testCases := map[string]struct {
		args    args
		wantErr error
	}{
		"success": {
			args: args{
				payload: &NftRegistrationRequest{
					Name:                      "name",
					Description:               &des,
					Keywords:                  &key,
					SeriesName:                &series,
					IssuedCopies:              &issuedCopies,
					YoutubeURL:                &url,
					CreatorPastelID:           "pastelid",
					CreatorPastelIDPassphrase: "passphrase",
					CreatorName:               "creator",
					CreatorWebsiteURL:         &url,
					SpendableAddress:          "spendaddr",
					MaximumFee:                12,
					Green:                     &green,
					Royalty:                   &royalty,
					Thumbnail: files.ThumbnailCoordinate{
						TopLeftX:     0,
						TopLeftY:     0,
						BottomRightY: 640,
						BottomRightX: 480,
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("testCase-%v", name), func(t *testing.T) {
			t.Parallel()

			gotReq := ToNftRegisterTicket(tc.args.payload)
			assert.Equal(t, gotReq.CreatorName, tc.args.payload.CreatorName)
			assert.Equal(t, gotReq.CreatorPastelID, tc.args.payload.CreatorPastelID)
			assert.Equal(t, gotReq.CreatorPastelIDPassphrase, tc.args.payload.CreatorPastelIDPassphrase)
			assert.Equal(t, gotReq.Description, tc.args.payload.Description)
			assert.Equal(t, gotReq.Keywords, tc.args.payload.Keywords)
			assert.Equal(t, gotReq.CreatorWebsiteURL, tc.args.payload.CreatorWebsiteURL)
			assert.Equal(t, gotReq.MaximumFee, tc.args.payload.MaximumFee)
			assert.Equal(t, gotReq.SeriesName, tc.args.payload.SeriesName)
			assert.Equal(t, gotReq.ThumbnailCoordinate.BottomRightX, tc.args.payload.Thumbnail.BottomRightX)
			assert.Equal(t, gotReq.ThumbnailCoordinate.BottomRightY, tc.args.payload.Thumbnail.BottomRightY)
			assert.Equal(t, gotReq.ThumbnailCoordinate.TopLeftY, tc.args.payload.Thumbnail.TopLeftY)

		})
	}
}
