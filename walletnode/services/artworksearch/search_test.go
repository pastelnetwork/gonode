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
			NFTTicketData: pastel.NFTTicket{
				AppTicketData: pastel.AppTicket{
					AuthorPastelID: "author-id",
					BlockTxID:      "block-txid",
					CreatorName:    "Alan Majchrowicz",
					NFTTitle:       "Lake Superior Sky III",
					TotalCopies:    10,
					CreatorWrittenStatement: `Lakes Why settle for blank walls, when you can
					 transform them into stunning vista points. Explore Lake Superiorfrom imaginative
					  scenic abstracts to sublime beach landscapes captured on camera.`,
					NFTKeywordSet: "Michigan,Midwest,Peninsula,Great Lakes,Lakeview",
					NFTSeriesName: "Science Art Lake",
				},
			},
			TotalCopies: 10,
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
				Query: "Explore Lake Superiorfrom ",
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
				Query:    "Lakes ",
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
