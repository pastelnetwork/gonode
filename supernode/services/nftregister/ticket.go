package nftregister

// Ticket represents NFT registration ticket.
type Ticket struct {
	Image                     []byte
	Name                      string
	Description               *string
	Keywords                  *string
	SeriesName                *string
	IssuedCopies              int
	YoutubeURL                *string
	CreatorPastelID           string
	CreatorPastelIDPassphrase string
	ArtistName                string
	ArtistWebsiteURL          *string
	SpendableAddress          string
	MaximumFee                float64
}
