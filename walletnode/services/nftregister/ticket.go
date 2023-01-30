package nftregister

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
)

// NftRegistrationRequest represents nft registration request.
type NftRegistrationRequest struct {
	Image                     *files.File               `json:"image"`
	Name                      string                    `json:"name"`
	Description               *string                   `json:"description"`
	Keywords                  *string                   `json:"keywords"`
	SeriesName                *string                   `json:"series_name"`
	IssuedCopies              int                       `json:"issued_copies"`
	YoutubeURL                *string                   `json:"youtube_url"`
	CreatorPastelID           string                    `json:"creator_pastel_id"`
	CreatorPastelIDPassphrase string                    `json:"creator_pastel_id_passphrase"`
	CreatorName               string                    `json:"creator_name"`
	CreatorWebsiteURL         *string                   `json:"creator_website_url"`
	SpendableAddress          string                    `json:"spendable_address"`
	MaximumFee                float64                   `json:"maximum_fee"`
	Green                     bool                      `json:"green"`
	Royalty                   float64                   `json:"royalty"`
	Thumbnail                 files.ThumbnailCoordinate `json:"thumbnail_coordinate"`
	MakePubliclyAccessible    bool                      `json:"make_publicly_accessible"`
}

// FromNftRegisterPayload converts from one to another
func FromNftRegisterPayload(payload *nft.RegisterPayload) *NftRegistrationRequest {
	thumbnail := files.ThumbnailCoordinate{
		TopLeftX:     payload.ThumbnailCoordinate.TopLeftX,
		TopLeftY:     payload.ThumbnailCoordinate.TopLeftY,
		BottomRightX: payload.ThumbnailCoordinate.BottomRightX,
		BottomRightY: payload.ThumbnailCoordinate.BottomRightY,
	}

	return &NftRegistrationRequest{
		Name:                      payload.Name,
		Description:               payload.Description,
		Keywords:                  payload.Keywords,
		SeriesName:                payload.SeriesName,
		IssuedCopies:              payload.IssuedCopies,
		YoutubeURL:                payload.YoutubeURL,
		CreatorPastelID:           payload.CreatorPastelID,
		CreatorPastelIDPassphrase: payload.CreatorPastelIDPassphrase,
		CreatorName:               payload.CreatorName,
		CreatorWebsiteURL:         payload.CreatorWebsiteURL,
		SpendableAddress:          payload.SpendableAddress,
		MaximumFee:                payload.MaximumFee,
		Green:                     payload.Green,
		Royalty:                   payload.Royalty,
		Thumbnail:                 thumbnail,
		MakePubliclyAccessible:    payload.MakePubliclyAccessible,
	}
}

// ToNftRegisterTicket converts from one to another
func ToNftRegisterTicket(request *NftRegistrationRequest) *nft.NftRegisterPayload {
	thumbnail := nft.Thumbnailcoordinate{
		TopLeftX:     request.Thumbnail.TopLeftX,
		TopLeftY:     request.Thumbnail.TopLeftY,
		BottomRightY: request.Thumbnail.BottomRightY,
		BottomRightX: request.Thumbnail.BottomRightX,
	}
	return &nft.NftRegisterPayload{
		Name:                      request.Name,
		Description:               request.Description,
		Keywords:                  request.Keywords,
		SeriesName:                request.SeriesName,
		IssuedCopies:              request.IssuedCopies,
		YoutubeURL:                request.YoutubeURL,
		CreatorPastelID:           request.CreatorPastelID,
		CreatorPastelIDPassphrase: request.CreatorPastelIDPassphrase,
		CreatorName:               request.CreatorName,
		CreatorWebsiteURL:         request.CreatorWebsiteURL,
		SpendableAddress:          request.SpendableAddress,
		MaximumFee:                request.MaximumFee,
		Green:                     request.Green,
		Royalty:                   request.Royalty,
		ThumbnailCoordinate:       &thumbnail,
	}
}
