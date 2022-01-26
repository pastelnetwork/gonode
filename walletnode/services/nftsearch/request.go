package nftsearch

import "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"

// NftSearchQueryField represents artsearch query-able field
type NftSearchQueryField string

// List of types of activation ticket.
const (
	NftSearchArtistName NftSearchQueryField = "artist_name"
	NftSearchArtTitle   NftSearchQueryField = "art_title"
	NftSearchSeries     NftSearchQueryField = "series"
	NftSearchDescr      NftSearchQueryField = "descr"
	NftSearchKeyword    NftSearchQueryField = "keyword"
)

// NftSearchQueryFields returns all query-able field names for an nft
var NftSearchQueryFields = []string{
	NftSearchArtistName.String(),
	NftSearchArtTitle.String(),
	NftSearchSeries.String(),
	NftSearchDescr.String(),
	NftSearchKeyword.String(),
}

// String returns the string val of NftSearchQueryField
func (f NftSearchQueryField) String() string {
	return string(f)
}

// NftSearchRequest represents nft search payload
type NftSearchRequest struct {
	// Artist PastelID or special value; mine
	Artist *string
	// Number of results to be return
	Limit int
	// Query is search query entered by user
	Query string
	// Search in Name of the artist
	ArtistName bool
	// Search in Title of nft
	ArtTitle bool
	// Search in Artwork series name
	Series bool
	// Search in Artist written statement
	Descr bool
	// Search Keywords that Artist assigns to Artwork
	Keyword bool
	// Minimum blocknum
	MinBlock int
	// Maximum blocknum
	MaxBlock *int
	// Minimum number of created copies
	MinCopies *int
	// Maximum number of created copies
	MaxCopies *int
	// Minimum nsfw score
	MinNsfwScore *float64
	// Maximum nsfw score
	MaxNsfwScore *float64
	// Minimum rareness score
	MinRarenessScore *float64
	// Maximum rareness score
	MaxRarenessScore *float64
	// MinInternetRarenessScore
	MinInternetRarenessScore *float64
	// MaxInternetRarenessScore
	MaxInternetRarenessScore *float64
	//UserPastelID
	UserPastelID string
	//UserPastelID
	UserPassphrase string
}

func FromNftSearchRequest(req *nft.NftSearchPayload) *NftSearchRequest {
	rq := &NftSearchRequest{
		Artist:           req.Artist,
		Limit:            req.Limit,
		Query:            req.Query,
		ArtistName:       req.ArtistName,
		ArtTitle:         req.ArtTitle,
		Series:           req.Series,
		Descr:            req.Descr,
		Keyword:          req.Keyword,
		MinBlock:         req.MinBlock,
		MaxBlock:         req.MaxBlock,
		MinCopies:        req.MinCopies,
		MaxCopies:        req.MaxCopies,
		MinNsfwScore:     req.MinNsfwScore,
		MaxNsfwScore:     req.MaxNsfwScore,
		MinRarenessScore: req.MinRarenessScore,
		MaxRarenessScore: req.MaxRarenessScore,
	}

	if req.UserPastelid != nil {
		rq.UserPastelID = *req.UserPastelid
	}

	if req.UserPassphrase != nil {
		rq.UserPassphrase = *req.UserPassphrase
	}

	return rq
}
