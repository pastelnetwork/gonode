package pastel

import (
	"bytes"
	a85 "encoding/ascii85"
	"encoding/json"

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
	Type           string           `json:"type"`
	Version        int              `json:"version"`
	Signatures     TicketSignatures `json:"signatures"`
	Key1           string           `json:"key1"`
	Key2           string           `json:"key2"`
	CreatorHeight  int              `json:"creator_height"`
	TotalCopies    int              `json:"total_copies"`
	Royalty        float64          `json:"royalty"`
	RoyaltyAddress string           `json:"royalty_address"`
	Green          bool             `json:"green"`
	GreenAddress   string           `json:"green_address"`
	StorageFee     int              `json:"storage_fee"`
	NFTTicket      []byte           `json:"nft_ticket"`
	NFTTicketData  NFTTicket        `json:"-"`
}

// NFTTicket is Pastel NFT Ticket
type NFTTicket struct {
	Version       int       `json:"nft_ticket_version"`
	Author        string    `json:"author"`
	BlockNum      int       `json:"blocknum"`
	BlockHash     string    `json:"block_hash"`
	Copies        int       `json:"copies"`
	Royalty       float64   `json:"royalty"`
	Green         bool      `json:"green"`
	AppTicket     []byte    `json:"app_ticket"`
	AppTicketData AppTicket `json:"-"`
}

// AlternativeNSFWScore represents the not-safe-for-work of an image
type AlternativeNSFWScore struct {
	Drawing float32 `json:"drawing"`
	Hentai  float32 `json:"hentai"`
	Neutral float32 `json:"neutral"`
	Porn    float32 `json:"porn"`
	Sexy    float32 `json:"sexy"`
}

// FingerAndScores represents structure of app ticket
type FingerAndScores struct {
	DupeDectectionSystemVersion string               `json:"dupe_dectection_system_version"`
	HashOfCandidateImageFile    []byte               `json:"hash_of_candidate_image_file"`
	OverallAverageRarenessScore float32              `json:"overall_average_rareness_score"`
	IsLikelyDupe                bool                 `json:"is_rare_on_internet"`
	IsRareOnInternet            bool                 `json:"is_likely_dupe"`
	NumberOfPagesOfResults      uint32               `json:"number_of_pages_of_results"`
	MatchesFoundOnFirstPage     uint32               `json:"matches_found_on_first_page"`
	URLOfFirstMatchInPage       string               `json:"url_of_first_match_in_page"`
	OpenNSFWScore               float32              `json:"open_nsfw_score"`
	ZstdCompressedFingerprint   []byte               `json:"zstd_compressed_fingerprint"`
	AlternativeNSFWScore        AlternativeNSFWScore `json:"alternative_nsfw_score"`
	ImageHashes                 ImageHashes          `json:"image_hashes"`
}

// AppTicket represents pastel App ticket.
type AppTicket struct {
	AuthorPastelID string `json:"author_pastel_id"`
	BlockTxID      string `json:"block_tx_id"`
	BlockNum       int    `json:"block_num"`

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

	FingerprintsHash      []byte `json:"fingerprints_hash"`
	FingerprintsSignature []byte `json:"fingerprints_signature"`

	DupeDetectionSystemVer  string              `json:"dupe_detection_system_version"`
	MatchesFoundOnFirstPage int                 `json:"matches_found_on_first_page"`
	NumberOfResultPages     int                 `json:"number_of_pages_of_results"`
	FirstMatchURL           string              `json:"url_of_first_match_in_page"`
	PastelRarenessScore     float32             `json:"pastel_rareness_score"`
	InternetRarenessScore   float32             `json:"internet_rareness_score"`
	OpenNSFWScore           float32             `json:"open_nsfw_score"`
	AlternateNSFWScores     AlternateNSFWScores `json:"alternate_nsfw_scores"`
	ImageHashes             ImageHashes         `json:"image_hashes"`
	IsLikelyDupe            bool                `json:"is_likely_dupe"`
	IsRareOnInternet        bool                `json:"is_rare_on_internet"`

	RQIDs []string `json:"rq_ids"`
	RQOti []byte   `json:"rq_oti"`
}

// AlternateNSFWScores represents alternate NSFW scores from dupe detection service
type AlternateNSFWScores struct {
	Drawing float32 `json:"drawing"`
	Hentai  float32 `json:"hentai"`
	Neutral float32 `json:"neutral"`
	Porn    float32 `json:"porn"`
	Sexy    float32 `json:"sexy"`
}

// ImageHashes represents image hashes from dupe detection service
type ImageHashes struct {
	PDQHash        string `json:"pdq_hash"`
	PerceptualHash string `json:"perceptual_hash"`
	AverageHash    string `json:"average_hash"`
	DifferenceHash string `json:"difference_hash"`
	NeuralHash     string `json:"neuralhash_hash"`
}

// GetRegisterNFTFeeRequest represents a request to get registration fee
type GetRegisterNFTFeeRequest struct {
	Ticket      *NFTTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Passphrase  string
	Key1        string
	Key2        string
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
	AppTicket []byte  `json:"app_ticket"`
}

// EncodeNFTTicket encodes  NFTTicket into byte array
func EncodeNFTTicket(ticket *NFTTicket) ([]byte, error) {
	appTicketBytes, err := json.Marshal(ticket.AppTicketData)
	if err != nil {
		return nil, errors.Errorf("marshal app ticket: %w", err)
	}

	appTicket := make([]byte, a85.MaxEncodedLen(len(appTicketBytes)))
	_ = a85.Encode(appTicket, appTicketBytes)

	// NFTTicket is Pastel Art Ticket
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
	b85Decode := func(in []byte) []byte {
		decodedBytes := make([]byte, len(in))
		nDecodedBytes, _, _ := a85.Decode(decodedBytes, in, true)
		decodedBytes = decodedBytes[:nDecodedBytes]

		//ascii85 adds /x00 null bytes at the end
		return bytes.Trim(decodedBytes, "\x00")
	}

	res := internalNFTTicket{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Errorf("unmarshal nft ticket: %w", err)
	}

	appDecodedBytes := b85Decode(res.AppTicket)
	appTicket := AppTicket{}
	err = json.Unmarshal(appDecodedBytes, &appTicket)
	if err != nil {
		return nil, errors.Errorf("unmarshal nft ticket: %w", err)
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

// TicketSignatures represents signatures from parties
type TicketSignatures struct {
	Creator map[string]string `json:"creator,omitempty"`
	Mn1     map[string]string `json:"mn1,omitempty"`
	Mn2     map[string]string `json:"mn2,omitempty"`
	Mn3     map[string]string `json:"mn3,omitempty"`
}

// RegisterNFTRequest represents request to register an art
// "ticket" "{signatures}" "jXYqZNPj21RVnwxnEJ654wEdzi7GZTZ5LAdiotBmPrF7pDMkpX1JegDMQZX55WZLkvy9fxNpZcbBJuE8QYUqBF" "passphrase", "key1", "key2", 100)")
type RegisterNFTRequest struct {
	Ticket      *NFTTicket
	Signatures  *TicketSignatures
	Mn1PastelID string
	Pasphase    string
	Key1        string
	Key2        string
	Fee         int64
}

// EncodeSignatures encodes TicketSignatures into byte array
func EncodeSignatures(signatures TicketSignatures) ([]byte, error) {
	// reset signatures of Mn1 if any
	signatures.Mn1 = nil

	b, err := json.Marshal(signatures)

	if err != nil {
		return nil, errors.Errorf("marshal signatures: %w", err)
	}

	return b, nil
}
