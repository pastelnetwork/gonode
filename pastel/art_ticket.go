package pastel

// ArtTicket represents pastel art ticket.
type ArtTicket struct {
	AuthorPastelID string `json:"author_pastel_id"`
	BlockTxID      string `json:"block_tx_id"`
	BlockNum       int    `json:"block_num"`
	ImageHash      string `json:"image_hash"`

	ArtistName      string `json:"artist_name"`
	ArtistSite      string `json:"artist_site"`
	ArtistStatement string `json:"artist_statement"`

	ArtworkTitle      string   `json:"artwork_title"`
	ArtworkSeriesName string   `json:"artwork_series_name"`
	ArtworkURLs       []string `json:"artwork_urls"`
	ArtworkKeywords   []string `json:"artwork_keywords"`

	TotalCopies int `json:"total_copies"`

	Fingerprints  []float64 `json:"fingerprints"`
	ThumbnailHash string    `json:"thumbnail_hash"`

	IpfsHash string `json:"ipfs_hash"`
	IpfsName string `json:"ipfs_name"`
}
