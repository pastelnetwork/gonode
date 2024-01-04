package services

import (
	"context"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/metrics/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
	mService "github.com/pastelnetwork/gonode/walletnode/services/metrics"
)

// MetricsAPIHandler - MetricsAPIHandler service
type MetricsAPIHandler struct {
	*Common
	metrics *mService.MetricsService
}

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *MetricsAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// Mount configures the mux to serve the OpenAPI enpoints.
func (service *MetricsAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := metrics.NewEndpoints(service)

	srv := server.New(
		endpoints,
		mux,
		goahttp.RequestDecoder,
		goahttp.ResponseEncoder,
		api.ErrorHandler,
		nil,
	)
	server.Mount(mux, srv)

	for _, m := range srv.Mounts {
		log.WithContext(ctx).Infof("%q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	return srv
}

// SelfHealing - returns the metrics for self-healing
func (service *MetricsAPIHandler) SelfHealing(ctx context.Context, p *metrics.MetricsPayload) (res *metrics.SelfHealingMetrics, err error) {
	shMetrics, err := service.metrics.GetSelfHealingMetrics(ctx, p.Pid, p.Key)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting self-healing metrics")
	}
	return &metrics.SelfHealingMetrics{
		TotalTicketsSendForSelfHealing: &shMetrics.SentTicketsForSelfHealing,
		EstimatedMissingKeys:           &shMetrics.EstimatedMissingKeys,
		TicketsRequiredSelfHealing:     &shMetrics.TicketsRequiredSelfHealing,
		TicketsSelfHealedSuccessfully:  &shMetrics.SuccessfullySelfHealedTickets,
		TicketsVerifiedSuccessfully:    &shMetrics.SuccessfullyVerifiedTickets,
	}, nil
}

// NewMetricsAPIHandler returns the metrics swagger OpenAPI implementation.
func NewMetricsAPIHandler(config *Config, metricsService *mService.MetricsService) *MetricsAPIHandler {
	return &MetricsAPIHandler{
		Common:  NewCommon(config),
		metrics: metricsService,
	}
}
