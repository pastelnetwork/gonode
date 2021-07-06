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
	Type          string     `json:"type"`
	ArtistHeight  int        `json:"artist_height"`
	Signatures    Signatures `json:"signatures"`
	Key1          string     `json:"key1"`
	Key2          string     `json:"key2"`
	IsGreen       bool       `json:"is_green"`
	StorageFee    int        `json:"storage_fee"`
	TotalCopies   int        `json:"total_copies"`
	Royalty       int        `json:"royalty"`
	Version       int        `json:"version"`
	ArtTicket     []byte     `json:"art_ticket"`
	ArtTicketData ArtTicket  `json:"-"`
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

	RQIDs   []string `json:"rq_ids"`
	RQCoti  int64    `json:"rq_coti"`
	RQSsoti int64    `json:"rq_ssoti"`
}

// Signature represents signatures element of art ticket
type Signature struct {
	Signature []byte `json:"signature"`
	PubKey    []byte `json:"pubkey"`
}

// Signature represents signalture element of signatures
type Signatures struct {
	SignatureAuthor Signature `json:"signature_author"`
	Signature1      Signature `json:"signature_1"`
	Signature2      Signature `json:"signature_2"`
	Signature3      Signature `json:"signature_3"`
}
