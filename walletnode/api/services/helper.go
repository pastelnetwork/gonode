package services

import (
	"time"

	"github.com/pastelnetwork/gonode/common/service/task/state"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
)

func fromRegisterPayload(payload *artworks.RegisterPayload) *artworkregister.Ticket {
	thumbnail := artwork.ThumbnailCoordinate{
		TopLeftX:     payload.ThumbnailCoordinate.TopLeftX,
		TopLeftY:     payload.ThumbnailCoordinate.TopLeftY,
		BottomRightX: payload.ThumbnailCoordinate.BottomRightX,
		BottomRightY: payload.ThumbnailCoordinate.BottomRightY,
	}

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
		Green:                    payload.Green,
		Royalty:                  payload.Royalty,
		Thumbnail:                thumbnail,
	}
}

func toArtworkTicket(ticket *artworkregister.Ticket) *artworks.ArtworkTicket {
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
