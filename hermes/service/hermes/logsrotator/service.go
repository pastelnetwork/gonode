package logsrotator

import (
	"github.com/pastelnetwork/gonode/hermes/service"
)

type logRotationService struct {
	appLogMap map[string]string
}

// NewLogRotationService returns a new log-rotation service
func NewLogRotationService() (service.SvcInterface, error) {
	lrs := &logRotationService{
		appLogMap: make(map[string]string),
	}

	lrs.appLogMap["supernode"] = "~/.pastel/supernode.log"
	lrs.appLogMap["hermes"] = "~/.pastel/hermes.log"

	return lrs, nil
}
