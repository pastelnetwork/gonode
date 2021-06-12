package services

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"

	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
)

func fromRegisterPayload(payload *artworks.RegisterPayload) *artworkregister.Ticket {
	return &artworkregister.Ticket{
		Name:                     payload.Name,
		Description:              payload.Description,
		Keywords:                 payload.Keywords,
		SeriesName:               payload.SeriesName,
		IssuedCopies:             payload.IssuedCopies,
		YoutubeURL:               payload.YoutubeURL,
		ArtistPastelID:           payload.ArtistPastelID,
		ArtistPastelIDPassphrase: payload.ArtistPastelIDPassphrase,
		ArtistName:               payload.ArtistName,
		ArtistWebsiteURL:         payload.ArtistWebsiteURL,
		SpendableAddress:         payload.SpendableAddress,
		MaximumFee:               payload.MaximumFee,
	}
}

func toArtworkTicket(ticket *artworkregister.Ticket) *artworks.ArtworkTicket {
	return &artworks.ArtworkTicket{
		Name:                     ticket.Name,
		Description:              ticket.Description,
		Keywords:                 ticket.Keywords,
		SeriesName:               ticket.SeriesName,
		IssuedCopies:             ticket.IssuedCopies,
		YoutubeURL:               ticket.YoutubeURL,
		ArtistPastelID:           ticket.ArtistPastelID,
		ArtistPastelIDPassphrase: ticket.ArtistPastelIDPassphrase,
		ArtistName:               ticket.ArtistName,
		ArtistWebsiteURL:         ticket.ArtistWebsiteURL,
		SpendableAddress:         ticket.SpendableAddress,
		MaximumFee:               ticket.MaximumFee,
	}
}

func toArtworkStates(statuses []*state.Status) []*artworks.TaskState {
	var states []*artworks.TaskState

	for _, status := range statuses {
		states = append(states, &artworks.TaskState{
			Date:   status.CreatedAt.Format(time.RFC3339),
			Status: status.String(),
		})
	}
	return states
}

func toArtSearchResult(srch *artworksearch.RegTicketSearch) *artworks.ArtworkSearchResult {
	ticketData := srch.RegTicketData.ArtTicketData.AppTicketData
	res := &artworks.ArtworkSearchResult{
		Artwork: &artworks.ArtworkTicket{
			Name: ticketData.ArtworkTitle,

			IssuedCopies:     srch.RegTicketData.ArtTicketData.Copies,
			ArtistName:       ticketData.ArtistName,
			YoutubeURL:       &ticketData.ArtworkCreationVideoYoutubeURL,
			ArtistPastelID:   ticketData.AuthorPastelID,
			ArtistWebsiteURL: &ticketData.ArtistWebsite,
			Description:      &ticketData.ArtistWrittenStatement,
			Keywords:         &ticketData.ArtworkKeywordSet,
			SeriesName:       &ticketData.ArtworkSeriesName,
		},
		Image:      srch.Thumbnail,
		MatchIndex: srch.MatchIndex,
	}

	res.Matches = []*artworks.FuzzyMatch{}
	for _, match := range res.Matches {
		res.Matches = append(res.Matches, &artworks.FuzzyMatch{
			Score:          match.Score,
			Str:            match.Str,
			FieldType:      match.FieldType,
			MatchedIndexes: match.MatchedIndexes,
		})
	}

	return res
}

func fromArtSearchRequest(req *artworks.ArtSearchPayload) *artworksearch.ArtSearchRequest {
	return &artworksearch.ArtSearchRequest{
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
}
