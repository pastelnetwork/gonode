package pastel

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
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
	Blocknum      int       `json:"blocknum"`
	DataHash      []byte    `json:"data_hash"`
	Copies        int       `json:"copies"`
	Royalty       int       `json:"royalty"`
	Green         string    `json:"green"`
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

type internalArtTicket struct {
	Version   int    `json:"version"`
	Author    []byte `json:"author"`
	Blocknum  int    `json:"blocknum"`
	DataHash  []byte `json:"data_hash"`
	Copies    int    `json:"copies"`
	Royalty   int    `json:"royalty"`
	Green     string `json:"green"`
	AppTicket []byte `json:"app_ticket"`
}

func EncodeArtTicket(ticket *ArtTicket) ([]byte, error) {
	appTicket, err := json.Marshal(ticket.AppTicketData)
	if err != nil {
		return nil, errors.Errorf("failed to marshal app ticket data: %w", err)
	}

	// ArtTicket is Pastel Art Ticket
	artTicket := internalArtTicket{
		Version:   ticket.Version,
		Author:    ticket.Author,
		Blocknum:  ticket.Blocknum,
		DataHash:  ticket.DataHash,
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
		Blocknum:      res.Blocknum,
		DataHash:      res.DataHash,
		Copies:        res.Copies,
		Royalty:       res.Royalty,
		Green:         res.Green,
		AppTicketData: appTicket,
	}, nil

}

func EncodeSignatures(signatures TicketSignatures) ([]byte, error) {
	// reset signatures of Mn1 if any
	signatures.Mn1 = nil

	b, err := json.Marshal(signatures)

	if err != nil {
		return nil, errors.Errorf("failed to marshal signatures: %w", err)
	}

	return b, nil
}
