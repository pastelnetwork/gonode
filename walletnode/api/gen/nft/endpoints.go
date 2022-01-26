// Code generated by goa v3.4.3, DO NOT EDIT.
//
// nft endpoints
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package nft

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "nft" service endpoints.
type Endpoints struct {
	Register          goa.Endpoint
	RegisterTaskState goa.Endpoint
	RegisterTask      goa.Endpoint
	RegisterTasks     goa.Endpoint
	UploadImage       goa.Endpoint
	NftSearch         goa.Endpoint
	NftGet            goa.Endpoint
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

// NftSearchEndpointInput holds both the payload and the server stream of the
// "nftSearch" method.
type NftSearchEndpointInput struct {
	// Payload is the method payload.
	Payload *NftSearchPayload
	// Stream is the server stream used by the "nftSearch" method to send data.
	Stream NftSearchServerStream
}

// NewEndpoints wraps the methods of the "nft" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Register:          NewRegisterEndpoint(s),
		RegisterTaskState: NewRegisterTaskStateEndpoint(s),
		RegisterTask:      NewRegisterTaskEndpoint(s),
		RegisterTasks:     NewRegisterTasksEndpoint(s),
		UploadImage:       NewUploadImageEndpoint(s),
		NftSearch:         NewNftSearchEndpoint(s),
		NftGet:            NewNftGetEndpoint(s),
		Download:          NewDownloadEndpoint(s, a.APIKeyAuth),
	}
}

// Use applies the given middleware to all the "nft" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Register = m(e.Register)
	e.RegisterTaskState = m(e.RegisterTaskState)
	e.RegisterTask = m(e.RegisterTask)
	e.RegisterTasks = m(e.RegisterTasks)
	e.UploadImage = m(e.UploadImage)
	e.NftSearch = m(e.NftSearch)
	e.NftGet = m(e.NftGet)
	e.Download = m(e.Download)
}

// NewRegisterEndpoint returns an endpoint function that calls the method
// "register" of service "nft".
func NewRegisterEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RegisterPayload)
		res, err := s.Register(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedRegisterResult(res, "default")
		return vres, nil
	}
}

// NewRegisterTaskStateEndpoint returns an endpoint function that calls the
// method "registerTaskState" of service "nft".
func NewRegisterTaskStateEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ep := req.(*RegisterTaskStateEndpointInput)
		return nil, s.RegisterTaskState(ctx, ep.Payload, ep.Stream)
	}
}

// NewRegisterTaskEndpoint returns an endpoint function that calls the method
// "registerTask" of service "nft".
func NewRegisterTaskEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RegisterTaskPayload)
		res, err := s.RegisterTask(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedTask(res, "default")
		return vres, nil
	}
}

// NewRegisterTasksEndpoint returns an endpoint function that calls the method
// "registerTasks" of service "nft".
func NewRegisterTasksEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		res, err := s.RegisterTasks(ctx)
		if err != nil {
			return nil, err
		}
		vres := NewViewedTaskCollection(res, "tiny")
		return vres, nil
	}
}

// NewUploadImageEndpoint returns an endpoint function that calls the method
// "uploadImage" of service "nft".
func NewUploadImageEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UploadImagePayload)
		res, err := s.UploadImage(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedImage(res, "default")
		return vres, nil
	}
}

// NewNftSearchEndpoint returns an endpoint function that calls the method
// "nftSearch" of service "nft".
func NewNftSearchEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ep := req.(*NftSearchEndpointInput)
		return nil, s.NftSearch(ctx, ep.Payload, ep.Stream)
	}
}

// NewNftGetEndpoint returns an endpoint function that calls the method
// "nftGet" of service "nft".
func NewNftGetEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*NftGetPayload)
		return s.NftGet(ctx, p)
	}
}

// NewDownloadEndpoint returns an endpoint function that calls the method
// "download" of service "nft".
func NewDownloadEndpoint(s Service, authAPIKeyFn security.AuthAPIKeyFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*NftDownloadPayload)
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
