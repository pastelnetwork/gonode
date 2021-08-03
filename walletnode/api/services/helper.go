package services

import (
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"

	"github.com/pastelnetwork/gonode/common/service/task/state"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
)

func fromRegisterPayload(payload *artworks.RegisterPayload) *artworkregister.Request {
	thumbnail := artwork.ThumbnailCoordinate{
		TopLeftX:     payload.ThumbnailCoordinate.TopLeftX,
		TopLeftY:     payload.ThumbnailCoordinate.TopLeftY,
		BottomRightX: payload.ThumbnailCoordinate.BottomRightX,
		BottomRightY: payload.ThumbnailCoordinate.BottomRightY,
	}

	return &artworkregister.Request{
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
		Green:                    payload.Green,
		Royalty:                  payload.Royalty,
		Thumbnail:                thumbnail,
	}
}

func toArtworkTicket(ticket *artworkregister.Request) *artworks.ArtworkTicket {
	thumbnail := artworks.Thumbnailcoordinate{
		TopLeftX:     ticket.Thumbnail.TopLeftX,
		TopLeftY:     ticket.Thumbnail.TopLeftY,
		BottomRightY: ticket.Thumbnail.BottomRightX,
		BottomRightX: ticket.Thumbnail.BottomRightY,
	}
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
		Green:                    ticket.Green,
		Royalty:                  ticket.Royalty,
		ThumbnailCoordinate:      &thumbnail,
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
		Artwork: &artworks.ArtworkSummary{
			Txid:      srch.TXID,
			Thumbnail: srch.Thumbnail,
			Title:     ticketData.ArtworkTitle,

			Copies:           srch.RegTicketData.ArtTicketData.Copies,
			ArtistName:       ticketData.ArtistName,
			YoutubeURL:       &ticketData.ArtworkCreationVideoYoutubeURL,
			ArtistPastelID:   ticketData.AuthorPastelID,
			ArtistWebsiteURL: &ticketData.ArtistWebsite,
			Description:      ticketData.ArtistWrittenStatement,
			Keywords:         &ticketData.ArtworkKeywordSet,
			SeriesName:       &ticketData.ArtworkSeriesName,
		},

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

func toArtworkDetail(ticket *pastel.RegTicket) *artworks.ArtworkDetail {
	return &artworks.ArtworkDetail{
		Txid:             ticket.TXID,
		Title:            ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkTitle,
		Copies:           ticket.RegTicketData.ArtTicketData.Copies,
		ArtistName:       ticket.RegTicketData.ArtTicketData.AppTicketData.ArtistName,
		YoutubeURL:       &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkCreationVideoYoutubeURL,
		ArtistPastelID:   ticket.RegTicketData.ArtTicketData.AppTicketData.AuthorPastelID,
		ArtistWebsiteURL: &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtistWebsite,
		Description:      ticket.RegTicketData.ArtTicketData.AppTicketData.ArtistWrittenStatement,
		Keywords:         &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkKeywordSet,
		SeriesName:       &ticket.RegTicketData.ArtTicketData.AppTicketData.ArtworkSeriesName,
		IsGreen:          ticket.RegTicketData.IsGreen,
		Royalty:          float64(ticket.RegTicketData.Royalty),
		// FIXME
		// RarenessScore:    ticket.RegTicketData.ArtTicketData.AppTicketData.RarenessScore,
		// NsfwScore:        ticket.RegTicketData.ArtTicketData.AppTicketData.NSFWScore,
		// SeenScore:        ticket.RegTicketData.ArtTicketData.AppTicketData.SeenScore,
		Version:    &ticket.RegTicketData.ArtTicketData.Version,
		StorageFee: &ticket.RegTicketData.StorageFee,
	}
}

func fromDownloadPayload(payload *artworks.DownloadPayload) *artworkdownload.Ticket {
	return &artworkdownload.Ticket{
		Txid:               payload.Txid,
		PastelID:           payload.Pid,
		PastelIDPassphrase: payload.Key,
	}
}
