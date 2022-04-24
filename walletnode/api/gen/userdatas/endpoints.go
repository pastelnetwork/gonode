// Code generated by goa v3.6.2, DO NOT EDIT.
//
// userdatas endpoints
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package userdatas

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "userdatas" service endpoints.
type Endpoints struct {
	CreateUserdata goa.Endpoint
	UpdateUserdata goa.Endpoint
	GetUserdata    goa.Endpoint
}

// NewEndpoints wraps the methods of the "userdatas" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		CreateUserdata: NewCreateUserdataEndpoint(s),
		UpdateUserdata: NewUpdateUserdataEndpoint(s),
		GetUserdata:    NewGetUserdataEndpoint(s),
	}
}

// Use applies the given middleware to all the "userdatas" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.CreateUserdata = m(e.CreateUserdata)
	e.UpdateUserdata = m(e.UpdateUserdata)
	e.GetUserdata = m(e.GetUserdata)
}

// NewCreateUserdataEndpoint returns an endpoint function that calls the method
// "createUserdata" of service "userdatas".
func NewCreateUserdataEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*CreateUserdataPayload)
		return s.CreateUserdata(ctx, p)
	}
}

// NewUpdateUserdataEndpoint returns an endpoint function that calls the method
// "updateUserdata" of service "userdatas".
func NewUpdateUserdataEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UpdateUserdataPayload)
		return s.UpdateUserdata(ctx, p)
	}
}

// NewGetUserdataEndpoint returns an endpoint function that calls the method
// "getUserdata" of service "userdatas".
func NewGetUserdataEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetUserdataPayload)
		return s.GetUserdata(ctx, p)
	}
}
