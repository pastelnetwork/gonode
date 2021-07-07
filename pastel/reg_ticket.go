package pastel

// Refer https://pastel.wiki/en/Architecture/Components/TicketStructures

// RegTickets is a collection of RegTicket
type RegTickets []RegTicket

// RegTicket represents pastel registration ticket.
type RegTicket struct {
	Height        int           `json:"height"`
	TXID          string        `json:"txid"`
	RegTicketData RegTicketData `json:"ticket"`
}

// RegTicketData is Pastel Registration ticket structure
type RegTicketData struct {
	Type          string    `json:"type"`
	ArtistHeight  int       `json:"artist_height"`
	Key1          string    `json:"key1"`
	Key2          string    `json:"key2"`
	IsGreen       bool      `json:"is_green"`
	StorageFee    int       `json:"storage_fee"`
	TotalCopes    int       `json:"total_copes"`
	Royalty       int       `json:"royalty"`
	Version       int       `json:"version"`
	ArtTicket     []byte    `json:"art_ticket"`
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

	ArtistName             string `json:"artist_name"`
	ArtistWebsite          string `json:"artist_website"`
	ArtistWrittenStatement string `json:"artist_written_statement"`

	ArtworkTitle                   string `json:"artwork_title"`
	ArtworkSeriesName              string `json:"artwork_series_name"`
	ArtworkCreationVideoYoutubeURL string `json:"artwork_creation_video_youtube_url"`
	ArtworkKeywordSet              string `json:"artwork_keyword_set"`
	TotalCopies                    int    `json:"total_copies"`

	ThumbnailHash []byte `json:"thumbnail_hash"`
	DataHash      []byte `json:"data_hash"`

	Fingerprints          []byte `json:"fingerprints"`
	FingerprintsHash      []byte `json:"fingerprints_hash"`
	FingerprintsSignature []byte `json:"fingerprints_signature"`

	RarenessScore int `json:"rareness_score"`
	NSFWScore     int `json:"nsfw_score"`
	SeenScore     int `json:"seen_score"`

	RQIDs   []string `json:"rq_ids"`
	RQCoti  int64    `json:"rq_coti"`
	RQSsoti int64    `json:"rq_ssoti"`
}

type TicketSignatures struct {
	Artist map[string][]byte `json:"artist,omitempty"`
	Mn1    map[string][]byte `json:"mn1,omitempty"`
	Mn2    map[string][]byte `json:"mn2,omitempty"`
	Mn3    map[string][]byte `json:"mn3,omitempty"`
}

// Command
// "ticket" "{signatures}" "jXYqZNPj21RVnwxnEJ654wEdzi7GZTZ5LAdiotBmPrF7pDMkpX1JegDMQZX55WZLkvy9fxNpZcbBJuE8QYUqBF" "passphrase", "key1", "key2", 100)")
type RegisterArtRequest struct {
	Ticket      *ArtTicket
	Signatures  *TicketSignatures
	Mn1PastelId string
	Pasphase    string
	Key1        string
	Key2        string
	Fee         int64
}

type GetRegisterArtFeeRequest struct {
	Ticket      *ArtTicket
	Signatures  *TicketSignatures
	Mn1PastelId string
	Pasphase    string
	Key1        string
	Key2        string
	Fee         int64
	ImgSizeInMb int64
}
