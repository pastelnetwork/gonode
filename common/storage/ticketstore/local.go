package ticketstore

import (
	"context"
)

// TicketStorageInterface is interface for queries ticket sqlite store
type TicketStorageInterface interface {
	CloseTicketDB(ctx context.Context)

	FilesQueries
	ActivationAttemptsQueries
	RegistrationAttemptsQueries
	MultiVolCascadeTicketMapQueries
}
