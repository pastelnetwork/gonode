package artworksearch

// ArtSearchQueryField represents artsearch query-able field
type ArtSearchQueryField string

// List of types of activation ticket.
const (
	ArtSearchArtistName ArtSearchQueryField = "artist_name"
	ArtSearchArtTitle   ArtSearchQueryField = "art_title"
	ArtSearchSeries     ArtSearchQueryField = "series"
	ArtSearchDescr      ArtSearchQueryField = "descr"
	ArtSearchKeyword    ArtSearchQueryField = "keyword"
)

// ArtSearchQueryFields returns all query-able field names for an artwork
var ArtSearchQueryFields = []string{
	ArtSearchArtistName.String(),
	ArtSearchArtTitle.String(),
	ArtSearchSeries.String(),
	ArtSearchDescr.String(),
	ArtSearchKeyword.String(),
}

// String returns the string val of ArtSearchQueryField
func (f ArtSearchQueryField) String() string {
	return string(f)
}

// ArtSearchRequest represents artwork search payload
type ArtSearchRequest struct {
	// Artist PastelID or special value; mine
	Artist *string
	// Number of results to be return
	Limit int
	// Query is search query entered by user
	Query string
	// Search in Name of the artist
	ArtistName bool
	// Search in Title of artwork
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
}
