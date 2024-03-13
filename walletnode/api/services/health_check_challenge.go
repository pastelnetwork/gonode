package services

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	healthCheckChallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/health_check_challenge"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/health_check_challenge/server"
	healthCheckChallengeRegister "github.com/pastelnetwork/gonode/walletnode/services/healthcheckchallenge"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

// HealthCheckChallengeAPIHandler - HealthCheckChallengeAPIHandler service
type HealthCheckChallengeAPIHandler struct {
	*Common
	healthCheckChallengeService *healthCheckChallengeRegister.Service
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *HealthCheckChallengeAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := healthCheckChallenge.NewEndpoints(service)

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
func (service *HealthCheckChallengeAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// GetSummaryStats returns the stats over a specified time range
func (service *HealthCheckChallengeAPIHandler) GetSummaryStats(ctx context.Context, p *healthCheckChallenge.GetSummaryStatsPayload) (*healthCheckChallenge.HcSummaryStatsResult, error) {
	var from, to *time.Time

	if p.From != nil {
		fromTime, err := time.Parse(time.RFC3339, *p.From)
		if err != nil {
			return nil, healthCheckChallenge.MakeBadRequest(fmt.Errorf("invalid from time format: %w", err))
		}

		from = &fromTime
	}

	if p.To != nil {
		toTime, err := time.Parse(time.RFC3339, *p.To)
		if err != nil {
			return nil, healthCheckChallenge.MakeBadRequest(fmt.Errorf("invalid to time format: %w", err))
		}
		to = &toTime
	}

	req := healthCheckChallengeRegister.GetSummaryStats{
		From:       from,
		To:         to,
		PastelID:   p.Pid,
		Passphrase: p.Key,
	}

	res, err := service.healthCheckChallengeService.GetSummaryStats(ctx, req)
	if err != nil {
		return nil, healthCheckChallenge.MakeInternalServerError(fmt.Errorf("failed to get metrics: %w", err))
	}

	return &healthCheckChallenge.HcSummaryStatsResult{
		HcSummaryStats: &healthCheckChallenge.HCSummaryStats{
			TotalChallengesIssued:                    res.HCSummaryStats.TotalChallenges,
			TotalChallengesProcessedByRecipient:      res.HCSummaryStats.TotalChallengesProcessed,
			TotalChallengesEvaluatedByChallenger:     res.HCSummaryStats.TotalChallengesEvaluatedByChallenger,
			TotalChallengesVerified:                  res.HCSummaryStats.TotalChallengesVerified,
			NoOfSlowResponsesObservedByObservers:     res.HCSummaryStats.SlowResponsesObservedByObservers,
			NoOfInvalidSignaturesObservedByObservers: res.HCSummaryStats.InvalidSignaturesObservedByObservers,
			NoOfInvalidEvaluationObservedByObservers: res.HCSummaryStats.InvalidEvaluationObservedByObservers,
		},
	}, nil
}

// NewHealthCheckChallengeAPIHandler returns the swagger OpenAPI implementation.
func NewHealthCheckChallengeAPIHandler(srvc *healthCheckChallengeRegister.Service) *HealthCheckChallengeAPIHandler {
	return &HealthCheckChallengeAPIHandler{
		healthCheckChallengeService: srvc,
	}
}
