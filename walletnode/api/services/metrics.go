package services

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/metrics/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
	metricsregister "github.com/pastelnetwork/gonode/walletnode/services/metrics"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

// MetricsAPIHandler - MetricsAPIHandler service
type MetricsAPIHandler struct {
	*Common
	metricsService *metricsregister.Service
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
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

// APIKeyAuth implements the authorization logic for the APIKey security scheme.
func (service *MetricsAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// GetChallengeReports returns the self-healing challenge reports
func (service *MetricsAPIHandler) GetChallengeReports(ctx context.Context, p *metrics.GetChallengeReportsPayload) (*metrics.SelfHealingChallengeReports, error) {
	if p.Count == nil && p.ChallengeID == nil {
		return nil, metrics.MakeBadRequest(fmt.Errorf("count or challenge_id is required"))
	}

	if p.Count != nil && p.ChallengeID != nil {
		return nil, metrics.MakeBadRequest(fmt.Errorf("only one of count or challenge_id is allowed"))
	}

	req := metricsregister.SHChallengesRequest{
		PastelID:   p.Pid,
		Passphrase: p.Key,
	}

	if p.Count != nil {
		req.Count = *p.Count
	}

	if p.ChallengeID != nil {
		req.ChallengeID = *p.ChallengeID
	}

	reports, err := service.metricsService.GetSelfHealingChallengeReports(ctx, req)
	if err != nil {
		return nil, metrics.MakeInternalServerError(fmt.Errorf("failed to get challenge reports: %w", err))
	}

	return toSHChallengeReport(reports), nil
}

// GetMetrics returns the metrics data over a specified time range
func (service *MetricsAPIHandler) GetMetrics(ctx context.Context, p *metrics.GetMetricsPayload) (*metrics.MetricsResult, error) {
	var from, to *time.Time

	if p.From != nil {
		fromTime, err := time.Parse(time.RFC3339, *p.From)
		if err != nil {
			return nil, metrics.MakeBadRequest(fmt.Errorf("invalid from time format: %w", err))
		}

		from = &fromTime
	}

	if p.To != nil {
		toTime, err := time.Parse(time.RFC3339, *p.To)
		if err != nil {
			return nil, metrics.MakeBadRequest(fmt.Errorf("invalid to time format: %w", err))
		}
		to = &toTime
	}

	req := metricsregister.GetMetricsRequest{
		From:       from,
		To:         to,
		PastelID:   p.Pid,
		Passphrase: p.Key,
	}

	res, err := service.metricsService.GetMetrics(ctx, req)
	if err != nil {
		return nil, metrics.MakeInternalServerError(fmt.Errorf("failed to get metrics: %w", err))
	}

	// Convert SHTriggerMetrics from slice of SHTriggerMetric to slice of pointers to SHTriggerMetric for the result
	var shTriggerMetrics []*metrics.SHTriggerMetric
	for _, metric := range res.SHTriggerMetrics {
		m := metric
		shTriggerMetrics = append(shTriggerMetrics, &metrics.SHTriggerMetric{
			TriggerID:              m.TriggerID,
			NodesOffline:           m.NodesOffline,
			ListOfNodes:            m.ListOfNodes,
			TotalFilesIdentified:   m.TotalFilesIdentified,
			TotalTicketsIdentified: m.TotalTicketsIdentified,
		})
	}

	return &metrics.MetricsResult{
		ScMetrics:        res.SCMetrics,
		ShTriggerMetrics: shTriggerMetrics,
		ShExecutionMetrics: &metrics.SHExecutionMetrics{
			TotalChallengesIssued:     res.SHExecutionMetrics.TotalChallengesIssued,
			TotalChallengesRejected:   res.SHExecutionMetrics.TotalChallengesRejected,
			TotalChallengesAccepted:   res.SHExecutionMetrics.TotalChallengesAccepted,
			TotalChallengesFailed:     res.SHExecutionMetrics.TotalChallengesFailed,
			TotalChallengesSuccessful: res.SHExecutionMetrics.TotalChallengesSuccessful,
			TotalFilesHealed:          res.SHExecutionMetrics.TotalFilesHealed,
			TotalFileHealingFailed:    res.SHExecutionMetrics.TotalFileHealingFailed,
		},
	}, nil
}

// NewMetricsAPIHandler returns the swagger OpenAPI implementation.
func NewMetricsAPIHandler(srvc *metricsregister.Service) *MetricsAPIHandler {
	return &MetricsAPIHandler{
		metricsService: srvc,
	}
}
