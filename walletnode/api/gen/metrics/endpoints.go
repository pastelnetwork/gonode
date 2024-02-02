// Code generated by goa v3.14.0, DO NOT EDIT.
//
// metrics endpoints
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package metrics

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "metrics" service endpoints.
type Endpoints struct {
	GetMetrics goa.Endpoint
}

// NewEndpoints wraps the methods of the "metrics" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		GetMetrics: NewGetMetricsEndpoint(s, a.APIKeyAuth),
	}
}

// Use applies the given middleware to all the "metrics" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.GetMetrics = m(e.GetMetrics)
}

// NewGetMetricsEndpoint returns an endpoint function that calls the method
// "getMetrics" of service "metrics".
func NewGetMetricsEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetMetricsPayload)
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
		res, err := s.GetMetrics(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedMetricsResult(res, "default")
		return vres, nil
	}
}
