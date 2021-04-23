package api

import (
	"time"

	"github.com/pastelnetwork/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/walletnode/services/artwork/register"
	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

func fromRegisterPayload(payload *artworks.RegisterPayload) *register.Ticket {
	return &register.Ticket{
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
		NetworkFee:               payload.NetworkFee,
	}
}

func toArtworkTicket(ticket *register.Ticket) *artworks.ArtworkTicket {
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
		NetworkFee:               ticket.NetworkFee,
	}
}

func toArtworkStates(msgs []*state.Message) []*artworks.TaskState {
	var states []*artworks.TaskState

	for _, msg := range msgs {
		states = append(states, &artworks.TaskState{
			Date:   msg.CreatedAt.Format(time.RFC3339),
			Status: msg.Status.String(),
		})
	}
	return states
}
