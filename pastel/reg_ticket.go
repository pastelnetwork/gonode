package pastel

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/probe"
)

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

// FingerAndScores represents structure of app ticket
type FingerAndScores struct {
	FingerprintData []byte
	Fingerprints    probe.Fingerprints
	RarenessScore   int
	NSFWScore       int
	SeenScore       int
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

// TicketSignatures represents signatures from parties
type TicketSignatures struct {
	Artist map[string]string `json:"artist,omitempty"`
	Mn1    map[string]string `json:"mn1,omitempty"`
	Mn2    map[string]string `json:"mn2,omitempty"`
	Mn3    map[string]string `json:"mn3,omitempty"`
}

// RegisterArtRequest represents request to register an art
// "ticket" "{signatures}" "jXYqZNPj21RVnwxnEJ654wEdzi7GZTZ5LAdiotBmPrF7pDMkpX1JegDMQZX55WZLkvy9fxNpZcbBJuE8QYUqBF" "passphrase", "key1", "key2", 100)")
type RegisterArtRequest struct {
	Ticket      *ArtTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Pasphase    string
	Key1        string
	Key2        string
	Fee         int64
}

// GetRegisterArtFeeRequest represents a request to get registration fee
type GetRegisterArtFeeRequest struct {
	Ticket      *ArtTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Passphrase  string
	Key1        string
	Key2        string
	Fee         int64
	ImgSizeInMb int64
}

type internalArtTicket struct {
	Version   int    `json:"version"`
	Author    []byte `json:"author"`
	BlockNum  int    `json:"blocknum"`
	BlockHash []byte `json:"block_hash"`
	Copies    int    `json:"copies"`
	Royalty   int    `json:"royalty"`
	Green     string `json:"green_address"`
	AppTicket []byte `json:"app_ticket"`
}

// EncodeArtTicket encodes  ArtTicket into byte array
func EncodeArtTicket(ticket *ArtTicket) ([]byte, error) {
	appTicket, err := json.Marshal(ticket.AppTicketData)
	if err != nil {
		return nil, errors.Errorf("failed to marshal app ticket data: %w", err)
	}

	// ArtTicket is Pastel Art Ticket
	artTicket := internalArtTicket{
		Version:   ticket.Version,
		Author:    ticket.Author,
		BlockNum:  ticket.BlockNum,
		BlockHash: ticket.BlockHash,
		Copies:    ticket.Copies,
		Royalty:   ticket.Royalty,
		Green:     ticket.Green,
		AppTicket: appTicket,
	}

	b, err := json.Marshal(artTicket)
	if err != nil {
		return nil, errors.Errorf("failed to marshal art ticket: %w", err)
	}

	return b, nil
}

// DecodeArtTicket decoded byte array into ArtTicket
func DecodeArtTicket(b []byte) (*ArtTicket, error) {
	res := internalArtTicket{}
	err := json.Unmarshal(b, &res)

	if err != nil {
		return nil, errors.Errorf("failed to unmarshal art ticket: %w", err)
	}

	appTicket := AppTicket{}

	err = json.Unmarshal(res.AppTicket, &appTicket)
	if err != nil {
		return nil, errors.Errorf("failed to unmarshal app ticket data: %w", err)
	}

	return &ArtTicket{
		Version:       res.Version,
		Author:        res.Author,
		BlockNum:      res.BlockNum,
		BlockHash:     res.BlockHash,
		Copies:        res.Copies,
		Royalty:       res.Royalty,
		Green:         res.Green,
		AppTicketData: appTicket,
	}, nil

}

// EncodeSignatures encodes TicketSignatures into byte array
func EncodeSignatures(signatures TicketSignatures) ([]byte, error) {
	// reset signatures of Mn1 if any
	signatures.Mn1 = nil

	b, err := json.Marshal(signatures)

	if err != nil {
		return nil, errors.Errorf("failed to marshal signatures: %w", err)
	}

	return b, nil
}
