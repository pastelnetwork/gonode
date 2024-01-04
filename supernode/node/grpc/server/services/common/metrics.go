package common

import (
	"github.com/pastelnetwork/gonode/supernode/services/metrics"
)

// Metrics represents common grpc service for metrics.
type Metrics struct {
	*metrics.FetchMetricsService
}

// NewMetricsService returns a new Metrics instance.
func NewMetricsService(service *metrics.FetchMetricsService) *Metrics {
	return &Metrics{
		FetchMetricsService: service,
	}
}
