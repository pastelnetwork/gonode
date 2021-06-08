package pastel

// Ticket is Pastel NFT ticket structure
type Ticket struct {
	Type          string    `json:"type"`
	ArtTicket     []byte    `json:"art_ticket"`
	ArtistHeight  int       `json:"artist_height"`
	Key1          string    `json:"key1"`
	Key2          string    `json:"key2"`
	StorageFee    int       `json:"storage_fee"`
	TotalCopes    int       `json:"total_copes"`
	ArtTicketData ArtTicket `json:"-"`
}

// ArtTicket is Pastel Art Ticket
type ArtTicket struct {
	Version       int       `json:"version"`
	Blocknum      int       `json:"blocknum"`
	Author        []byte    `json:"author"`
	DataHash      []byte    `json:"data_hash"`
	Copies        int       `json:"copies"`
	Reserved      []byte    `json:"reserved"`
	AppTicket     []byte    `json:"app_ticket"`
	AppTicketData AppTicket `json:"-"`
}

// AppTicket represents pastel App ticket.
type AppTicket struct {
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
