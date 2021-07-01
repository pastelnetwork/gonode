package pastel

// FIXME
type RegisterTicket struct {
	Version   int     `json:"version"`    // 1
	Author    []byte  `json:"author"`     // PastelID of the author (artist)
	Blocknum  int     `json:"blocknum"`   // block number when the ticket was created - this is to map the ticket to the MNs that should process it
	BlockHash []byte  `json:"block_hash"` // hash of the top block when the ticket was created - this is to map the ticket to the MNs that should process it
	Copies    int     `json:"copies"`     // number of copies
	Royalty   float64 `json:"royalty"`    // (not yet supported by cNode) how much artist should get on all future resales
	Green     string  `json:"green"`      // address for Green NFT payment (not yet supported by cNode)

	AppTicket []byte `json:"app_ticket"` // cNode DOES NOT parse this part!!!!

	ArtistName                     string `json:"artist_name"`
	ArtworkTitle                   string `json:"artwork_title"`
	ArtworkSeriesName              string `json:"artwork_series_name"`
	ArtworkKeywordSet              string `json:"artwork_keyword_set"`
	ArtistWebsite                  string `json:"artist_website"`
	ArtistWrittenStatement         string `json:"artist_written_statement"`
	ArtworkCreationVideoYoutubeURL string `json:"artwork_creation_video_youtube_url"`

	ThumbnailHash []byte `json:"thumbnail_hash"` // hash of the thumbnail !!!!SHA3-256!!!!
	DataHash      []byte `json:"data_hash"`      // hash of the image (or any other asset) that this ticket represents !!!!SHA3-256!!!!

	FingerprintsHash      []byte `json:"fingerprints_hash"`      // hash of the fingerprint !!!!SHA3-256!!!!
	Fingerprints          []byte `json:"fingerprints"`           // compressed fingerprint
	FingerprintsSignature []byte `json:"fingerprints_signature"` // signature on raw image fingerprint

	RqIds   []string `json:"rq_ids"`   // raptorq symbol identifiers -  !!!!SHA3-256 of symbol block!!!!
	RqCoti  int64    `json:"rq_coti"`  // raptorq CommonOTI
	RqSsoti int64    `json:"rq_ssoti"` // raptorq SchemeSpecificOTI

	RarenessScore int `json:"rareness_score"` // 0 to 1000
	NSFWScore     int `json:"nsfw_score"`     // 0 to 1000
	SeenScore     int `json:"seen_score"`     // 0 to 1000
}

// FIXME
type TradeTicket struct {
	Type      string `json:"type"`      // "trade",
	PastelID  string `json:"pastelID"`  // PastelID of the buyer
	SellTXID  string `json:"sell_txid"` // txid with sale ticket
	BueTXID   string `json:"buy_txid"`  // txid with buy ticket
	ArtTXID   string `json:"art_txid"`  // txid with either 1) art activation ticket or 2) trade ticket in it
	Price     string `json:"price"`
	Reserved  string `json:"reserved"`
	Signature string `json:"signature"`
}
