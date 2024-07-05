package common

import (
	"context"

	"github.com/pastelnetwork/gonode/pastel"
)

// ValidateUser validates user by id and password.
func ValidateUser(ctx context.Context, pc pastel.Client, id string, pass string) bool {
	_, err := pc.Sign(ctx, []byte("data"), id, pass, pastel.SignAlgorithmED448)
	return err == nil
}

// IsPastelIDTicketRegistered validates if the user has a valid registered pastelID ticket
func IsPastelIDTicketRegistered(ctx context.Context, pc pastel.Client, id string) bool {
	idTicket, err := pc.FindTicketByID(ctx, id)
	return err == nil && idTicket.PastelID == id
}

type AddTaskPayload struct {
	FileID                 string
	BurnTxid               *string
	BurnTxids              []string
	AppPastelID            string
	MakePubliclyAccessible bool
	SpendableAddress       *string
	Key                    string
}
