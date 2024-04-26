// Code generated by goa v3.15.0, DO NOT EDIT.
//
// Score endpoints
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package score

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "Score" service endpoints.
type Endpoints struct {
	GetAggregatedChallengesScores goa.Endpoint
}

// NewEndpoints wraps the methods of the "Score" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		GetAggregatedChallengesScores: NewGetAggregatedChallengesScoresEndpoint(s, a.APIKeyAuth),
	}
}

// Use applies the given middleware to all the "Score" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.GetAggregatedChallengesScores = m(e.GetAggregatedChallengesScores)
}

// NewGetAggregatedChallengesScoresEndpoint returns an endpoint function that
// calls the method "getAggregatedChallengesScores" of service "Score".
func NewGetAggregatedChallengesScoresEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetAggregatedChallengesScoresPayload)
		var err error
		sc := security.APIKeyScheme{
			Name:           "api_key",
			Scopes:         []string{},
			RequiredScopes: []string{},
		}
		ctx, err = authAPIKeyFn(ctx, p.Key, &sc)
		if err != nil {
			return nil, err
		}
		return s.GetAggregatedChallengesScores(ctx, p)
	}
}
