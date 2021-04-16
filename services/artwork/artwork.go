package artwork

// Artwork represents artwork data.
type Artwork struct {
	Image            []byte
	Name             string
	Description      *string
	Keywords         *string
	SeriesName       *string
	IssuedCopies     int
	YoutubeURL       *string
	ArtistPastelID   string
	ArtistName       string
	ArtistWebsiteURL *string
	SpendableAddress string
	NetworkFee       float32
}
