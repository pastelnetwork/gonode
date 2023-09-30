// Code generated by goa v3.13.1, DO NOT EDIT.
//
// userdatas client
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package userdatas

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "userdatas" service client.
type Client struct {
	CreateUserdataEndpoint goa.Endpoint
	UpdateUserdataEndpoint goa.Endpoint
	GetUserdataEndpoint    goa.Endpoint
}

// NewClient initializes a "userdatas" service client given the endpoints.
func NewClient(createUserdata, updateUserdata, getUserdata goa.Endpoint) *Client {
	return &Client{
		CreateUserdataEndpoint: createUserdata,
		UpdateUserdataEndpoint: updateUserdata,
		GetUserdataEndpoint:    getUserdata,
	}
}

// CreateUserdata calls the "createUserdata" endpoint of the "userdatas"
// service.
// CreateUserdata may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) CreateUserdata(ctx context.Context, p *CreateUserdataPayload) (res *UserdataProcessResult, err error) {
	var ires any
	ires, err = c.CreateUserdataEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UserdataProcessResult), nil
}

// UpdateUserdata calls the "updateUserdata" endpoint of the "userdatas"
// service.
// UpdateUserdata may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) UpdateUserdata(ctx context.Context, p *UpdateUserdataPayload) (res *UserdataProcessResult, err error) {
	var ires any
	ires, err = c.UpdateUserdataEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UserdataProcessResult), nil
}

// GetUserdata calls the "getUserdata" endpoint of the "userdatas" service.
// GetUserdata may return the following errors:
//   - "BadRequest" (type *goa.ServiceError)
//   - "NotFound" (type *goa.ServiceError)
//   - "InternalServerError" (type *goa.ServiceError)
//   - error: internal error
func (c *Client) GetUserdata(ctx context.Context, p *GetUserdataPayload) (res *UserSpecifiedData, err error) {
	var ires any
	ires, err = c.GetUserdataEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UserSpecifiedData), nil
}
