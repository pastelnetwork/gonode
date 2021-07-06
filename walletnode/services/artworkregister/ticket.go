package artworkregister

import "github.com/pastelnetwork/gonode/common/service/artwork"

// Request represents artwork registration request.
type Request struct {
	Image                    *artwork.File               `json:"image"`
	Name                     string                      `json:"name"`
	Description              *string                     `json:"description"`
	Keywords                 *string                     `json:"keywords"`
	SeriesName               *string                     `json:"series_name"`
	IssuedCopies             int                         `json:"issued_copies"`
	YoutubeURL               *string                     `json:"youtube_url"`
	ArtistPastelID           string                      `json:"artist_pastel_id"`
	ArtistPastelIDPassphrase string                      `json:"artist_pastel_id_passphrase"`
	ArtistName               string                      `json:"artist_name"`
	ArtistWebsiteURL         *string                     `json:"artist_website_url"`
	SpendableAddress         string                      `json:"spendable_address"`
	MaximumFee               float64                     `json:"maximum_fee"`
	Green                    bool                        `json:"green"`
	Royalty                  float64                     `json:"royalty"`
	Thumbnail                artwork.ThumbnailCoordinate `json:"thumbnail_coordinate"`
}
