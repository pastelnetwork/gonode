package services

import (
	"context"

	"github.com/pastelnetwork/go-commons/log"
	tickets "github.com/pastelnetwork/walletnode/rest/gen/tickets"
)

type ticketsService struct{}

// Add a new ticket and return its ID.
func (service *ticketsService) Add(ctx context.Context, p *tickets.AddPayload) (res string, err error) {
	log.WithContext(ctx).Debug("tickets.add")
	return
}

// List all stored tickets.
func (service *ticketsService) List(ctx context.Context) (res tickets.StoredTicketCollection, err error) {
	log.WithContext(ctx).Debug("tickets.list")
	return
}

// NewTickets returns the tickets service implementation.
func NewTickets() tickets.Service {
	return &ticketsService{}
}
