// Code generated by goa v3.12.4, DO NOT EDIT.
//
// cascade endpoints
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package cascade

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "cascade" service endpoints.
type Endpoints struct {
	UploadAsset       goa.Endpoint
	StartProcessing   goa.Endpoint
	RegisterTaskState goa.Endpoint
	GetTaskHistory    goa.Endpoint
	Download          goa.Endpoint
}

// RegisterTaskStateEndpointInput holds both the payload and the server stream
// of the "registerTaskState" method.
type RegisterTaskStateEndpointInput struct {
	// Payload is the method payload.
	Payload *RegisterTaskStatePayload
	// Stream is the server stream used by the "registerTaskState" method to send
	// data.
	Stream RegisterTaskStateServerStream
}

// NewEndpoints wraps the methods of the "cascade" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		UploadAsset:       NewUploadAssetEndpoint(s),
		StartProcessing:   NewStartProcessingEndpoint(s, a.APIKeyAuth),
		RegisterTaskState: NewRegisterTaskStateEndpoint(s),
		GetTaskHistory:    NewGetTaskHistoryEndpoint(s),
		Download:          NewDownloadEndpoint(s, a.APIKeyAuth),
	}
}

// Use applies the given middleware to all the "cascade" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.UploadAsset = m(e.UploadAsset)
	e.StartProcessing = m(e.StartProcessing)
	e.RegisterTaskState = m(e.RegisterTaskState)
	e.GetTaskHistory = m(e.GetTaskHistory)
	e.Download = m(e.Download)
}

// NewUploadAssetEndpoint returns an endpoint function that calls the method
// "uploadAsset" of service "cascade".
func NewUploadAssetEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*UploadAssetPayload)
		res, err := s.UploadAsset(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedAsset(res, "default")
		return vres, nil
	}
}

// NewStartProcessingEndpoint returns an endpoint function that calls the
// method "startProcessing" of service "cascade".
func NewStartProcessingEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*StartProcessingPayload)
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
		res, err := s.StartProcessing(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStartProcessingResult(res, "default")
		return vres, nil
	}
}

// NewRegisterTaskStateEndpoint returns an endpoint function that calls the
// method "registerTaskState" of service "cascade".
func NewRegisterTaskStateEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		ep := req.(*RegisterTaskStateEndpointInput)
		return nil, s.RegisterTaskState(ctx, ep.Payload, ep.Stream)
	}
}

// NewGetTaskHistoryEndpoint returns an endpoint function that calls the method
// "getTaskHistory" of service "cascade".
func NewGetTaskHistoryEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*GetTaskHistoryPayload)
		return s.GetTaskHistory(ctx, p)
	}
}

// NewDownloadEndpoint returns an endpoint function that calls the method
// "download" of service "cascade".
func NewDownloadEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*DownloadPayload)
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
		return s.Download(ctx, p)
	}
}
