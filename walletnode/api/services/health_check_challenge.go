package services

import (
	"context"
	"fmt"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/walletnode/api"
	healthCheckChallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/health_check_challenge"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/health_check_challenge/server"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
	healthCheckChallengeRegister "github.com/pastelnetwork/gonode/walletnode/services/healthcheckchallenge"
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

// GetDetailedLogs returns the detailed health-check-challenge data logs
func (service *HealthCheckChallengeAPIHandler) GetDetailedLogs(ctx context.Context, p *healthCheckChallenge.GetDetailedLogsPayload) (hcDetailedLogMessages []*healthCheckChallenge.HcDetailedLogsMessage, err error) {
	var req healthCheckChallengeRegister.HCDetailedLogRequest

	if p.ChallengeID != nil {
		req = healthCheckChallengeRegister.HCDetailedLogRequest{
			ChallengeID: *p.ChallengeID,
			PastelID:    p.Pid,
			Passphrase:  p.Key,
		}
	} else {
		req = healthCheckChallengeRegister.HCDetailedLogRequest{
			Count:      50,
			PastelID:   p.Pid,
			Passphrase: p.Key,
		}
	}

	hcDetailedMessageDataList, err := service.healthCheckChallengeService.GetDetailedMessageDataList(ctx, req)
	if err != nil {
		return nil, metrics.MakeInternalServerError(fmt.Errorf("failed to get health-check-challenge detailed logs: %w", err))
	}

	for _, hcl := range hcDetailedMessageDataList {

		hcMsg := &healthCheckChallenge.HcDetailedLogsMessage{
			ChallengeID: hcl.ChallengeID,
			MessageType: hcl.MessageType.String(),
			SenderID:    hcl.Sender,

			ChallengerID: hcl.Data.ChallengerID,
			Observers:    hcl.Data.Observers,
			RecipientID:  hcl.Data.RecipientID,
		}

		if hcl.MessageType == types.HealthCheckChallengeMessageType {
			hcMsg.Challenge = &healthCheckChallenge.HCChallengeData{
				Timestamp: hcl.Data.Challenge.Timestamp.Format(time.RFC3339),
			}

			if hcl.Data.Challenge.Block != 0 {
				hcMsg.Challenge.Block = &hcl.Data.Challenge.Block
			}

			if hcl.Data.Challenge.Merkelroot != "" {
				hcMsg.Challenge.Merkelroot = &hcl.Data.Challenge.Merkelroot
			}
		}

		if hcl.MessageType == types.HealthCheckResponseMessageType {
			hcMsg.Response = &healthCheckChallenge.HCResponseData{
				Timestamp: hcl.Data.Response.Timestamp.Format(time.RFC3339),
			}

			if hcl.Data.Response.Block != 0 {
				hcMsg.Response.Block = &hcl.Data.Response.Block
			}

			if hcl.Data.Response.Merkelroot != "" {
				hcMsg.Response.Merkelroot = &hcl.Data.Response.Merkelroot
			}
		}

		if hcl.MessageType == types.HealthCheckEvaluationMessageType {
			hcMsg.ChallengerEvaluation = &healthCheckChallenge.HCEvaluationData{
				Timestamp:  hcl.Data.ChallengerEvaluation.Timestamp.Format(time.RFC3339),
				IsVerified: hcl.Data.ChallengerEvaluation.IsVerified,
			}

			if hcl.Data.ChallengerEvaluation.Block != 0 {
				hcMsg.ChallengerEvaluation.Block = &hcl.Data.ChallengerEvaluation.Block
			}

			if hcl.Data.ChallengerEvaluation.Merkelroot != "" {
				hcMsg.ChallengerEvaluation.Merkelroot = &hcl.Data.ChallengerEvaluation.Merkelroot
			}
		}

		if hcl.MessageType == types.HealthCheckAffirmationMessageType {
			hcMsg.ObserverEvaluation = &healthCheckChallenge.HCObserverEvaluationData{
				IsChallengeTimestampOK:  hcl.Data.ObserverEvaluation.IsChallengeTimestampOK,
				IsProcessTimestampOK:    hcl.Data.ObserverEvaluation.IsProcessTimestampOK,
				IsEvaluationTimestampOK: hcl.Data.ObserverEvaluation.IsEvaluationTimestampOK,
				IsRecipientSignatureOK:  hcl.Data.ObserverEvaluation.IsRecipientSignatureOK,
				IsChallengerSignatureOK: hcl.Data.ObserverEvaluation.IsChallengerSignatureOK,
				IsEvaluationResultOK:    hcl.Data.ObserverEvaluation.IsEvaluationResultOK,
				Timestamp:               hcl.Data.ObserverEvaluation.Timestamp.Format(time.RFC3339),
			}

			if hcl.Data.ObserverEvaluation.Block != 0 {
				hcMsg.ObserverEvaluation.Block = &hcl.Data.ObserverEvaluation.Block
			}

			if hcl.Data.ObserverEvaluation.Merkelroot != "" {
				hcMsg.ObserverEvaluation.Merkelroot = &hcl.Data.ObserverEvaluation.Merkelroot
			}
		}

		hcDetailedLogMessages = append(hcDetailedLogMessages, hcMsg)
	}

	return hcDetailedLogMessages, nil
}

// NewHealthCheckChallengeAPIHandler returns the swagger OpenAPI implementation.
func NewHealthCheckChallengeAPIHandler(srvc *healthCheckChallengeRegister.Service) *HealthCheckChallengeAPIHandler {
	return &HealthCheckChallengeAPIHandler{
		healthCheckChallengeService: srvc,
	}
}
