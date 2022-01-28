package nftregister

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
)

// NftRegistrationRequest represents nft registration request.
type NftRegistrationRequest struct {
	Image                    *files.File               `json:"image"`
	Name                     string                    `json:"name"`
	Description              *string                   `json:"description"`
	Keywords                 *string                   `json:"keywords"`
	SeriesName               *string                   `json:"series_name"`
	IssuedCopies             int                       `json:"issued_copies"`
	YoutubeURL               *string                   `json:"youtube_url"`
	ArtistPastelID           string                    `json:"artist_pastel_id"`
	ArtistPastelIDPassphrase string                    `json:"artist_pastel_id_passphrase"`
	ArtistName               string                    `json:"artist_name"`
	ArtistWebsiteURL         *string                   `json:"artist_website_url"`
	SpendableAddress         string                    `json:"spendable_address"`
	MaximumFee               float64                   `json:"maximum_fee"`
	Green                    bool                      `json:"green"`
	Royalty                  float64                   `json:"royalty"`
	Thumbnail                files.ThumbnailCoordinate `json:"thumbnail_coordinate"`
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
		Name:                     payload.Name,
		Description:              payload.Description,
		Keywords:                 payload.Keywords,
		SeriesName:               payload.SeriesName,
		IssuedCopies:             payload.IssuedCopies,
		YoutubeURL:               payload.YoutubeURL,
		ArtistPastelID:           payload.ArtistPastelID,
		ArtistPastelIDPassphrase: payload.ArtistPastelIDPassphrase,
		ArtistName:               payload.ArtistName,
		ArtistWebsiteURL:         payload.ArtistWebsiteURL,
		SpendableAddress:         payload.SpendableAddress,
		MaximumFee:               payload.MaximumFee,
		Green:                    payload.Green,
		Royalty:                  payload.Royalty,
		Thumbnail:                thumbnail,
	}
}

// ToNftRegisterTicket converts from one to another
func ToNftRegisterTicket(request *NftRegistrationRequest) *nft.NftRegisterPayload {
	thumbnail := nft.Thumbnailcoordinate{
		TopLeftX:     request.Thumbnail.TopLeftX,
		TopLeftY:     request.Thumbnail.TopLeftY,
		BottomRightY: request.Thumbnail.BottomRightX,
		BottomRightX: request.Thumbnail.BottomRightY,
	}
	return &nft.NftRegisterPayload{
		Name:                     request.Name,
		Description:              request.Description,
		Keywords:                 request.Keywords,
		SeriesName:               request.SeriesName,
		IssuedCopies:             request.IssuedCopies,
		YoutubeURL:               request.YoutubeURL,
		ArtistPastelID:           request.ArtistPastelID,
		ArtistPastelIDPassphrase: request.ArtistPastelIDPassphrase,
		ArtistName:               request.ArtistName,
		ArtistWebsiteURL:         request.ArtistWebsiteURL,
		SpendableAddress:         request.SpendableAddress,
		MaximumFee:               request.MaximumFee,
		Green:                    request.Green,
		Royalty:                  request.Royalty,
		ThumbnailCoordinate:      &thumbnail,
	}
}
