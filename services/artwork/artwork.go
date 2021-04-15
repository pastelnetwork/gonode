package artwork

type Artwork struct {
	Name             string
	Description      *string
	Keywords         *string
	SeriesName       *string
	IssuedCopies     int
	Image            *Image
	YoutubeURL       *string
	ArtistPastelID   string
	ArtistName       string
	ArtistWebsiteURL *string
	SpendableAddress string
	NetworkFee       float32
}
