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
