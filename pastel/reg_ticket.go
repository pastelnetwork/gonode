package pastel

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
	Type          string           `json:"type"`
	ArtistHeight  int              `json:"artist_height"`
	Signatures    TicketSignatures `json:"signatures"`
	Key1          string           `json:"key1"`
	Key2          string           `json:"key2"`
	IsGreen       bool             `json:"is_green"`
	StorageFee    int              `json:"storage_fee"`
	TotalCopies   int              `json:"total_copies"`
	Royalty       int              `json:"royalty"`
	Version       int              `json:"version"`
	ArtTicket     []byte           `json:"art_ticket"`
	ArtTicketData ArtTicket        `json:"-"`
}

// ArtTicket is Pastel Art Ticket
type ArtTicket struct {
	Version       int       `json:"version"`
	Author        []byte    `json:"author"`
	BlockNum      int       `json:"blocknum"`
	BlockHash     []byte    `json:"block_hash"`
	Copies        int       `json:"copies"`
	Royalty       int       `json:"royalty"`
	Green         string    `json:"green_address"`
	AppTicket     []byte    `json:"app_ticket"`
	AppTicketData AppTicket `json:"-"`
}

// AppTicket represents pastel App ticket.
type AppTicket struct {
	AuthorPastelID string `json:"author_pastel_id"`
	BlockTxID      string `json:"block_tx_id"`
	BlockNum       int    `json:"block_num"`

	ArtistName             string `json:"artist_name"`
	ArtistWebsite          string `json:"artist_website"`
	ArtistWrittenStatement string `json:"artist_written_statement"`

	ArtworkTitle                   string `json:"artwork_title"`
	ArtworkSeriesName              string `json:"artwork_series_name"`
	ArtworkCreationVideoYoutubeURL string `json:"artwork_creation_video_youtube_url"`
	ArtworkKeywordSet              string `json:"artwork_keyword_set"`
	TotalCopies                    int    `json:"total_copies"`

	PreviewHash    []byte `json:"preview_hash"`
	Thumbnail1Hash []byte `json:"thumbnail1_hash"`
	Thumbnail2Hash []byte `json:"thumbnail2_hash"`

	DataHash []byte `json:"data_hash"`

	Fingerprints          []byte `json:"fingerprints"`
	FingerprintsHash      []byte `json:"fingerprints_hash"`
	FingerprintsSignature []byte `json:"fingerprints_signature"`

	RarenessScore int `json:"rareness_score"`
	NSFWScore     int `json:"nsfw_score"`
	SeenScore     int `json:"seen_score"`

	RQIDs []string `json:"rq_ids"`
	RQOti []byte   `json:"rq_oti"`
}

type TicketSignatures struct {
	Artist map[string]string `json:"artist,omitempty"`
	Mn1    map[string]string `json:"mn1,omitempty"`
	Mn2    map[string]string `json:"mn2,omitempty"`
	Mn3    map[string]string `json:"mn3,omitempty"`
}
