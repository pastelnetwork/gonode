package services

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/http/score/server"
	scorechallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/score"
	scoreregister "github.com/pastelnetwork/gonode/walletnode/services/score"
	"math"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/security"
)

// ScoreAPIHandler - ScoreAPIHandler service
type ScoreAPIHandler struct {
	*Common
	scoreService *scoreregister.Service
}

// Mount configures the mux to serve the OpenAPI endpoints.
func (service *ScoreAPIHandler) Mount(ctx context.Context, mux goahttp.Muxer) goahttp.Server {
	endpoints := scorechallenge.NewEndpoints(service)

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
func (service *ScoreAPIHandler) APIKeyAuth(ctx context.Context, _ string, _ *security.APIKeyScheme) (context.Context, error) {
	return ctx, nil
}

// GetAggregatedChallengesScores returns the aggregated challenges scores
func (service *ScoreAPIHandler) GetAggregatedChallengesScores(ctx context.Context, p *scorechallenge.GetAggregatedChallengesScoresPayload) ([]*scorechallenge.ChallengesScores, error) {

	req := scoreregister.ChallengesScoreRequest{
		PastelID:   p.Pid,
		Passphrase: p.Key,
	}

	res, err := service.scoreService.GetAggregatedChallengeScore(ctx, req)
	if err != nil {
		return nil, scorechallenge.MakeInternalServerError(fmt.Errorf("failed to get aggregated challenges scores: %w", err))
	}

	var challengesScores []*scorechallenge.ChallengesScores
	for _, cs := range res {

		score := &scorechallenge.ChallengesScores{
			NodeID:                    cs.NodeID,
			StorageChallengeScore:     math.Round(cs.StorageChallengeScore*100) / 100,
			HealthCheckChallengeScore: math.Round(cs.HealthCheckChallengeScore*100) / 100,
		}

		if cs.IPAddress != "" {
			score.IPAddress = &cs.IPAddress
		}

		challengesScores = append(challengesScores, score)
	}

	return challengesScores, nil
}

// NewScoreAPIHandler returns the swagger OpenAPI implementation.
func NewScoreAPIHandler(srvc *scoreregister.Service) *ScoreAPIHandler {
	return &ScoreAPIHandler{
		scoreService: srvc,
	}
}
