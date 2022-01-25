package artworksearch

import (
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/sahilm/fuzzy"
)

// RegTicketSearch is a helper to aggregate RegTicket search
type RegTicketSearch struct {
	*pastel.RegTicket
	Thumbnail         []byte
	ThumbnailSecondry []byte
	MaxScore          int
	Matches           []Match
	MatchIndex        int
}

// Match represents a matched string.
type Match struct {
	// The matched string.
	Str string
	// The field that is matched
	FieldType string
	// The indexes of matched characters. Useful for highlighting matches.
	MatchedIndexes []int
	// Score used to rank matches
	Score int
}

// getSearchableFields checks Request & returns list of data to be searched along with
// a Map to record index of each field type in the list as fuzzy search returns index
func (rs *RegTicketSearch) getSearchableFields(req *ArtSearchRequest) (data []string, mapper map[int]ArtSearchQueryField) {
	fieldIdxMapper := make(map[int]ArtSearchQueryField)
	idx := 0
	if req.ArtistName {
		data = append(data, rs.RegTicketData.NFTTicketData.AppTicketData.CreatorName)
		fieldIdxMapper[idx] = ArtSearchArtistName
		idx++
	}
	if req.ArtTitle {
		data = append(data, rs.RegTicketData.NFTTicketData.AppTicketData.NFTTitle)
		fieldIdxMapper[idx] = ArtSearchArtTitle
		idx++
	}
	if req.Descr {
		data = append(data, rs.RegTicketData.NFTTicketData.AppTicketData.CreatorWrittenStatement)
		fieldIdxMapper[idx] = ArtSearchDescr
		idx++
	}
	if req.Keyword {
		data = append(data, rs.RegTicketData.NFTTicketData.AppTicketData.NFTKeywordSet)
		fieldIdxMapper[idx] = ArtSearchKeyword
		idx++
	}
	if req.Series {
		data = append(data, rs.RegTicketData.NFTTicketData.AppTicketData.NFTSeriesName)
		fieldIdxMapper[idx] = ArtSearchSeries
	}

	return data, fieldIdxMapper
}

// Search does fuzzy search  on the reg ticket data for the query
func (rs *RegTicketSearch) Search(req *ArtSearchRequest) (srch *RegTicketSearch, matched bool) {
	data, mapper := rs.getSearchableFields(req)
	matches := fuzzy.Find(req.Query, data)
	for _, match := range matches {
		if match.Score <= 0 {
			continue
		}

		if match.Score > rs.MaxScore {
			rs.MaxScore = match.Score
		}
		rs.Matches = append(rs.Matches, Match{
			Str:            match.Str,
			FieldType:      string(mapper[match.Index]),
			MatchedIndexes: match.MatchedIndexes,
			Score:          match.Score,
		})
	}

	return rs, len(rs.Matches) > 0
}
