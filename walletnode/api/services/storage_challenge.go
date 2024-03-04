package services

import (
	"context"
	"fmt"
	storagechallengeregister "github.com/pastelnetwork/gonode/walletnode/services/storagechallenge"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/storage_challenge/server"
	storagechallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/storage_challenge"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

// StorageChallengeAPIHandler - StorageChallengeAPIHandler service
type StorageChallengeAPIHandler struct {
	*Common
	storageChallengeService *storagechallengeregister.Service
}

// Mount onfigures the mux to serve the OpenAPI enpoints.
func (service *StorageChallengeAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := storagechallenge.NewEndpoints(service)

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
func (service *StorageChallengeAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// GetSummaryStats returns the stats over a specified time range
func (service *StorageChallengeAPIHandler) GetSummaryStats(ctx context.Context, p *storagechallenge.GetSummaryStatsPayload) (*storagechallenge.SummaryStatsResult, error) {
	var from, to *time.Time

	if p.From != nil {
		fromTime, err := time.Parse(time.RFC3339, *p.From)
		if err != nil {
			return nil, storagechallenge.MakeBadRequest(fmt.Errorf("invalid from time format: %w", err))
		}

		from = &fromTime
	}

	if p.To != nil {
		toTime, err := time.Parse(time.RFC3339, *p.To)
		if err != nil {
			return nil, storagechallenge.MakeBadRequest(fmt.Errorf("invalid to time format: %w", err))
		}
		to = &toTime
	}

	req := storagechallengeregister.GetSummaryStats{
		From:       from,
		To:         to,
		PastelID:   p.Pid,
		Passphrase: p.Key,
	}

	res, err := service.storageChallengeService.GetSummaryStats(ctx, req)
	if err != nil {
		return nil, storagechallenge.MakeInternalServerError(fmt.Errorf("failed to get metrics: %w", err))
	}

	return &storagechallenge.SummaryStatsResult{
		ScSummaryStats: &storagechallenge.SCSummaryStats{
			TotalChallengesIssued:                    res.SCSummaryStats.TotalChallenges,
			TotalChallengesProcessed:                 res.SCSummaryStats.TotalChallengesProcessed,
			TotalChallengesVerifiedByChallenger:      res.SCSummaryStats.TotalChallengesVerifiedByChallenger,
			TotalChallengesVerifiedByObservers:       res.SCSummaryStats.TotalChallengesVerifiedByObservers,
			NoOfSlowResponsesObservedByObservers:     res.SCSummaryStats.SlowResponsesObservedByObservers,
			NoOfInvalidSignaturesObservedByObservers: res.SCSummaryStats.InvalidSignaturesObservedByObservers,
			NoOfInvalidEvaluationObservedByObservers: res.SCSummaryStats.InvalidEvaluationObservedByObservers,
		},
	}, nil
}

// NewStorageChallengeAPIHandler returns the swagger OpenAPI implementation.
func NewStorageChallengeAPIHandler(srvc *storagechallengeregister.Service) *StorageChallengeAPIHandler {
	return &StorageChallengeAPIHandler{
		storageChallengeService: srvc,
	}
}
