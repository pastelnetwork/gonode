package services

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/state"
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
			Status: status.Type.String(),
		})
	}
	return states
}
