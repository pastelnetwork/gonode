package nftsearch

import "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"

// NftSearchingQueryField represents artsearch query-able field
type NftSearchingQueryField string

// List of types of activation ticket.
const (
	//NftSearchArtistName is field
	NftSearchArtistName NftSearchingQueryField = "creator_name"
	//NftSearchArtTitle is field
	NftSearchArtTitle NftSearchingQueryField = "art_title"
	//NftSearchSeries is field
	NftSearchSeries NftSearchingQueryField = "series"
	//NftSearchDescr is field
	NftSearchDescr NftSearchingQueryField = "descr"
	//NftSearchKeyword is field
	NftSearchKeyword NftSearchingQueryField = "keyword"
)

// NftSearchQueryFields returns all query-able field names for an nft
var NftSearchQueryFields = []string{
	NftSearchArtistName.String(),
	NftSearchArtTitle.String(),
	NftSearchSeries.String(),
	NftSearchDescr.String(),
	NftSearchKeyword.String(),
}

// String returns the string val of NftSearchingQueryField
func (f NftSearchingQueryField) String() string {
	return string(f)
}

// NftSearchingRequest represents nft search payload
type NftSearchingRequest struct {
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
	// Search in Nft series name
	Series bool
	// Search in Artist written statement
	Descr bool
	// Search Keywords that Artist assigns to Nft
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
	MinNsfwScore float64
	// Maximum nsfw score
	MaxNsfwScore float64
	// Minimum rareness score
	MinRarenessScore float64
	// Maximum rareness score
	MaxRarenessScore float64
	// Is Likely Dupe
	IsLikelyDupe bool
	//UserPastelID
	UserPastelID string
	//UserPastelID
	UserPassphrase string
}

// FromNftSearchRequest from one to another
func FromNftSearchRequest(req *nft.NftSearchPayload) *NftSearchingRequest {
	rq := &NftSearchingRequest{
		Artist:     req.Artist,
		Limit:      req.Limit,
		Query:      req.Query,
		ArtistName: req.CreatorName,
		ArtTitle:   req.ArtTitle,
		Series:     req.Series,
		Descr:      req.Descr,
		Keyword:    req.Keyword,
		MinBlock:   req.MinBlock,
		MaxBlock:   req.MaxBlock,
		MinCopies:  req.MinCopies,
		MaxCopies:  req.MaxCopies,
	}

	if req.UserPastelid != nil {
		rq.UserPastelID = *req.UserPastelid
	}

	if req.UserPassphrase != nil {
		rq.UserPassphrase = *req.UserPassphrase
	}

	if req.MinNsfwScore != nil {
		rq.MinNsfwScore = *req.MinNsfwScore
	}

	if req.MaxNsfwScore != nil {
		rq.MaxNsfwScore = *req.MaxNsfwScore
	}

	if req.MinRarenessScore != nil {
		rq.MinRarenessScore = *req.MinRarenessScore
	}

	if req.MaxRarenessScore != nil {
		rq.MaxRarenessScore = *req.MaxRarenessScore
	}

	if req.IsLikelyDupe != nil {
		rq.IsLikelyDupe = *req.IsLikelyDupe
	}

	return rq
}
