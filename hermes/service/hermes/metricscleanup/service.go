package metricscleanup

import (
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/hermes/service"
)

type metricsCleanupService struct {
	historyDB queries.LocalStoreInterface
}

// NewMetricsCleanupService returns a new metric cleanup service
func NewMetricsCleanupService(hDB queries.LocalStoreInterface) (service.SvcInterface, error) {
	return &metricsCleanupService{
		historyDB: hDB,
	}, nil
}
