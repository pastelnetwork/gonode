package pastel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"

	"github.com/pastelnetwork/gonode/common/b85"

	"github.com/pastelnetwork/gonode/common/errors"
)

// Refer https://pastel.wiki/en/Architecture/Components/TicketStructures

const (
	// NFTTypeImage is NFT type "image"
	NFTTypeImage = "image"
)

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
	Type           string              `json:"type"`
	Version        int                 `json:"version"`
	Signatures     RegTicketSignatures `json:"signatures"`
	Label          string              `json:"label"`
	CreatorHeight  int                 `json:"creator_height"`
	TotalCopies    int                 `json:"total_copies"`
	Royalty        float64             `json:"royalty"`
	RoyaltyAddress string              `json:"royalty_address"`
	Green          bool                `json:"green"`
	GreenAddress   string              `json:"green_address"`
	StorageFee     int                 `json:"storage_fee"`
	NFTTicket      []byte              `json:"nft_ticket"`
	NFTTicketData  NFTTicket           `json:"-"`
}

// NFTTicket is Pastel NFT Request
type NFTTicket struct {
	Version       int       `json:"nft_ticket_version"`
	Author        string    `json:"author"`
	BlockNum      int       `json:"blocknum"`
	BlockHash     string    `json:"block_hash"`
	Copies        int       `json:"copies"`
	Royalty       float64   `json:"royalty"`
	Green         bool      `json:"green"`
	AppTicket     string    `json:"app_ticket"`
	AppTicketData AppTicket `json:"-"`
}

// AppTicket represents pastel App ticket.
type AppTicket struct {
	CreatorName             string `json:"creator_name"`
	CreatorWebsite          string `json:"creator_website"`
	CreatorWrittenStatement string `json:"creator_written_statement"`

	NFTTitle                   string `json:"nft_title"`
	NFTType                    string `json:"nft_type"`
	NFTSeriesName              string `json:"nft_series_name"`
	NFTCreationVideoYoutubeURL string `json:"nft_creation_video_youtube_url"`
	NFTKeywordSet              string `json:"nft_keyword_set"`
	TotalCopies                int    `json:"total_copies"`

	PreviewHash    []byte `json:"preview_hash"`
	Thumbnail1Hash []byte `json:"thumbnail1_hash"`
	Thumbnail2Hash []byte `json:"thumbnail2_hash"`

	DataHash []byte `json:"data_hash"`

	DDAndFingerprintsIc  uint32   `json:"dd_and_fingerprints_ic"`
	DDAndFingerprintsMax uint32   `json:"dd_and_fingerprints_max"`
	DDAndFingerprintsIDs []string `json:"dd_and_fingerprints_ids"`

	RQIc  uint32   `json:"rq_ic"`
	RQMax uint32   `json:"rq_max"`
	RQIDs []string `json:"rq_ids"`
	RQOti []byte   `json:"rq_oti"`

	OriginalFileSizeInBytes int    `json:"original_file_size_in_bytes"`
	FileType                string `json:"file_type"`
	MakePubliclyAccessible  bool   `json:"make_publicly_accessible"`
}

// GetRegisterNFTFeeRequest represents a request to get registration fee
type GetRegisterNFTFeeRequest struct {
	Ticket      *NFTTicket
	Signatures  *RegTicketSignatures
	Mn1PastelID string
	Passphrase  string
	Label       string
	Fee         int64
	ImgSizeInMb int64
}

type internalNFTTicket struct {
	Version   int     `json:"nft_ticket_version"`
	Author    string  `json:"author"`
	BlockNum  int     `json:"blocknum"`
	BlockHash string  `json:"block_hash"`
	Copies    int     `json:"copies"`
	Royalty   float64 `json:"royalty"`
	Green     bool    `json:"green"`
	AppTicket string  `json:"app_ticket"`
}

// EncodeNFTTicket encodes  NFTTicket into byte array
func EncodeNFTTicket(ticket *NFTTicket) ([]byte, error) {
	appTicketBytes, err := json.Marshal(ticket.AppTicketData)
	if err != nil {
		return nil, errors.Errorf("marshal app ticket: %w", err)
	}

	appTicket := b85.Encode(appTicketBytes)

	// NFTTicket is Pastel Art Request
	nftTicket := internalNFTTicket{
		Version:   ticket.Version,
		Author:    ticket.Author,
		BlockNum:  ticket.BlockNum,
		BlockHash: ticket.BlockHash,
		Copies:    ticket.Copies,
		Royalty:   ticket.Royalty,
		Green:     ticket.Green,
		AppTicket: appTicket,
	}

	b, err := json.Marshal(nftTicket)
	if err != nil {
		return nil, errors.Errorf("marshal nft ticket: %w", err)
	}

	return b, nil
}

// DecodeNFTTicket decoded byte array into ArtTicket
func DecodeNFTTicket(b []byte) (*NFTTicket, error) {
	res := internalNFTTicket{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Errorf("unmarshal nft ticket: %w", err)
	}

	appDecodedBytes, err := b85.Decode(res.AppTicket)
	if err != nil {
		log.Warnf("b85 decoding failed, trying to base64 decode - err: %v", err)
		appDecodedBytes, err = base64.StdEncoding.DecodeString(res.AppTicket)
		if err != nil {
			return nil, fmt.Errorf("b64 decode: %v", err)
		}
	}

	appTicket := AppTicket{}
	err = json.Unmarshal(appDecodedBytes, &appTicket)
	if err != nil {
		return nil, errors.Errorf("unmarshal app ticket: %w", err)
	}

	return &NFTTicket{
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

// RegTicketSignatures represents signatures from parties
type RegTicketSignatures struct {
	Creator map[string]string `json:"creator,omitempty"`
	Mn1     map[string]string `json:"mn1,omitempty"`
	Mn2     map[string]string `json:"mn2,omitempty"`
	Mn3     map[string]string `json:"mn3,omitempty"`
}

// RegisterNFTRequest represents request to register an art
// "ticket" "{signatures}" "jXYqZNPj21RVnwxnEJ654wEdzi7GZTZ5LAdiotBmPrF7pDMkpX1JegDMQZX55WZLkvy9fxNpZcbBJuE8QYUqBF" "passphrase", "label", 100)")
type RegisterNFTRequest struct {
	Ticket      *NFTTicket
	Signatures  *RegTicketSignatures
	Mn1PastelID string
	Passphrase  string
	Label       string
	Fee         int64
}

// EncodeRegSignatures encodes RegTicketSignatures into byte array
func EncodeRegSignatures(signatures RegTicketSignatures) ([]byte, error) {
	// reset signatures of Mn1 if any
	signatures.Mn1 = nil

	b, err := json.Marshal(signatures)

	if err != nil {
		return nil, errors.Errorf("marshal signatures: %w", err)
	}

	return b, nil
}
