package artworksearch

import (
	"testing"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/stretchr/testify/assert"
)

func TestGetSearchableFields(t *testing.T) {
	t.Parallel()
	regTicket := &pastel.RegTicket{
		Height: 0,
		TXID:   "id",
		RegTicketData: pastel.RegTicketData{
			ArtTicketData: pastel.ArtTicket{
				AppTicketData: pastel.AppTicket{
					AuthorPastelID: "author-id",
					BlockTxID:      "block-txid",
					ArtistName:     "Alan Majchrowicz",
					ArtworkTitle:   "Lake Superior Sky III",
					TotalCopies:    10,
					ArtistWrittenStatement: `Why settle for blank walls, when you can
					 transform them into stunning vista points. Explore from imaginative
					  scenic abstracts to sublime beach landscapes captured on camera.`,
					ArtworkKeywordSet: "Michigan,Midwest,Peninsula,Great Lakes,Lakeview",
					ArtworkSeriesName: "Science Art Lake",
				},
			},
			TotalCopes: 10,
		},
	}

	tests := map[string]struct {
		search   *RegTicketSearch
		req      *ArtSearchRequest
		matching bool
		matches  int
	}{
		"artist-name-match": {
			search: &RegTicketSearch{RegTicket: regTicket},
			req: &ArtSearchRequest{
				Query:      "alan",
				ArtTitle:   true,
				ArtistName: true,
			},
			matching: true,
			matches:  1,
		},
		"art-title-match": {
			search: &RegTicketSearch{RegTicket: regTicket},
			req: &ArtSearchRequest{
				Query:    "super",
				ArtTitle: true,
			},
			matching: true,
			matches:  1,
		},
		"keyword-match": {
			search: &RegTicketSearch{RegTicket: regTicket},
			req: &ArtSearchRequest{
				Query:   "Lakes",
				Keyword: true,
			},
			matching: true,
			matches:  1,
		},
		"descr-match": {
			search: &RegTicketSearch{RegTicket: regTicket},
			req: &ArtSearchRequest{
				Query: "beach",
				Descr: true,
			},
			matching: true,
			matches:  1,
		},
		"keyword-match-false": {
			search: &RegTicketSearch{RegTicket: regTicket},
			req: &ArtSearchRequest{
				Query:   "Lakes",
				Keyword: false, // keyword shouldn't be searched
			},
			matching: false,
			matches:  0,
		},
		"title,descr match": {
			search: &RegTicketSearch{RegTicket: regTicket},
			req: &ArtSearchRequest{
				Query:    "Lakes",
				ArtTitle: true,
				Descr:    true,
			},
			matching: true,
			matches:  2,
		},
		"title,descr, keyword match": {
			search: &RegTicketSearch{RegTicket: regTicket},
			req: &ArtSearchRequest{
				Query:    "lAkEs",
				ArtTitle: true,
				Descr:    true,
				Keyword:  true,
			},
			matching: true,
			matches:  3,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			srch, isMatched := tc.search.Search(tc.req)
			assert.Equal(t, tc.matching, isMatched)
			assert.Equal(t, tc.matches, len(srch.Matches))
		})
	}
}
