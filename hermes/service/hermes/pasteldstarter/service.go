package pasteldstarter

import (
	"github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/synchronizer"
	"github.com/pastelnetwork/gonode/pastel"
)

type restartPastelDService struct {
	pastelClient      pastel.Client
	sync              *synchronizer.Synchronizer
	currentBlockCount int32
}

// NewRestartPastelDService returns a new restart pastel-d service
func NewRestartPastelDService(pastelClient pastel.Client, s *synchronizer.Synchronizer) (service.SvcInterface, error) {
	return &restartPastelDService{
		pastelClient: pastelClient,
		sync:         s,
	}, nil
}
