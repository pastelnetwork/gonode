package services

import (
	"context"
	"fmt"
	"strings"
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

// GetDetailedLogs returns the detailed self-healing reports
func (service *MetricsAPIHandler) GetDetailedLogs(ctx context.Context, p *metrics.GetDetailedLogsPayload) (*metrics.SelfHealingReports, error) {
	var req metricsregister.SHReportRequest

	if p.EventID != nil {
		req = metricsregister.SHReportRequest{
			EventID:    *p.EventID,
			PastelID:   p.Pid,
			Passphrase: p.Key,
		}
	} else {
		req = metricsregister.SHReportRequest{
			Count:      10,
			PastelID:   p.Pid,
			Passphrase: p.Key,
		}
	}

	reports, err := service.metricsService.GetDetailedLogs(ctx, req)
	if err != nil {
		return nil, metrics.MakeInternalServerError(fmt.Errorf("failed to get self-healing detailed logs: %w", err))
	}

	return toSHReport(reports), nil
}

// GetSummaryStats returns the stats over a specified time range
func (service *MetricsAPIHandler) GetSummaryStats(ctx context.Context, p *metrics.GetSummaryStatsPayload) (*metrics.MetricsResult, error) {
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

	req := metricsregister.GetSummaryStats{
		From:       from,
		To:         to,
		PastelID:   p.Pid,
		Passphrase: p.Key,
	}

	res, err := service.metricsService.GetSummaryStats(ctx, req)
	if err != nil {
		return nil, metrics.MakeInternalServerError(fmt.Errorf("failed to get summary stats: %w", err))
	}

	// Convert SHTriggerMetrics from slice of SHTriggerMetric to slice of pointers to SHTriggerMetric for the result
	var shTriggerMetrics []*metrics.SHTriggerStats
	for _, metric := range res.SHTriggerMetrics {
		m := metric

		noOfNodesOffline := len(strings.Split(m.ListOfNodes, "-"))

		shTriggerMetrics = append(shTriggerMetrics, &metrics.SHTriggerStats{
			TriggerID:              m.TriggerID,
			NodesOffline:           noOfNodesOffline,
			ListOfNodes:            m.ListOfNodes,
			TotalFilesIdentified:   m.TotalFilesIdentified,
			TotalTicketsIdentified: m.TotalTicketsIdentified,
		})
	}

	return &metrics.MetricsResult{
		SelfHealingTriggerEventsStats: shTriggerMetrics,
		SelfHealingExecutionEventsStats: &metrics.SHExecutionStats{
			TotalSelfHealingEventsIssued:       res.SHExecutionMetrics.TotalChallengesIssued,
			TotalSelfHealingEventsAcknowledged: res.SHExecutionMetrics.TotalChallengesAcknowledged,
			TotalSelfHealingEventsRejected:     res.SHExecutionMetrics.TotalChallengesRejected,
			TotalSelfHealingEventsAccepted:     res.SHExecutionMetrics.TotalChallengesAccepted,

			TotalSelfHealingEventsEvaluationsVerified:             res.SHExecutionMetrics.TotalChallengeEvaluationsVerified,
			TotalReconstructionRequiredEvaluationsApproved:        res.SHExecutionMetrics.TotalReconstructionsApproved,
			TotalReconstructionNotRequiredEvaluationsApproved:     res.SHExecutionMetrics.TotalReconstructionsNotRquiredApproved,
			TotalSelfHealingEventsEvaluationsUnverified:           res.SHExecutionMetrics.TotalChallengeEvaluationsUnverified,
			TotalReconstructionRequiredEvaluationsNotApproved:     res.SHExecutionMetrics.TotalReconstructionsNotApproved,
			TotalReconstructionsNotRequiredEvaluationsNotApproved: res.SHExecutionMetrics.TotalReconstructionsNotRequiredEvaluationNotApproved,

			TotalFilesHealed:       res.SHExecutionMetrics.TotalFilesHealed,
			TotalFileHealingFailed: res.SHExecutionMetrics.TotalFileHealingFailed,

			TotalReconstructionRequiredHashMismatch: &res.SHExecutionMetrics.TotalReconstructionRequiredHashMismatch,
		},
	}, nil
}

// NewMetricsAPIHandler returns the swagger OpenAPI implementation.
func NewMetricsAPIHandler(srvc *metricsregister.Service) *MetricsAPIHandler {
	return &MetricsAPIHandler{
		metricsService: srvc,
	}
}
