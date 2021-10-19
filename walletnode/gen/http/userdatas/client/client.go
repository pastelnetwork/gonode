// Code generated by goa v3.5.2, DO NOT EDIT.
//
// userdatas client HTTP transport
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"context"
	"mime/multipart"
	"net/http"

	userdatas "github.com/pastelnetwork/gonode/walletnode/gen/userdatas"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the userdatas service endpoint HTTP clients.
type Client struct {
	// CreateUserdata Doer is the HTTP client used to make requests to the
	// createUserdata endpoint.
	CreateUserdataDoer goahttp.Doer

	// UpdateUserdata Doer is the HTTP client used to make requests to the
	// updateUserdata endpoint.
	UpdateUserdataDoer goahttp.Doer

	// GetUserdata Doer is the HTTP client used to make requests to the getUserdata
	// endpoint.
	GetUserdataDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// UserdatasCreateUserdataEncoderFunc is the type to encode multipart request
// for the "userdatas" service "createUserdata" endpoint.
type UserdatasCreateUserdataEncoderFunc func(*multipart.Writer, *userdatas.CreateUserdataPayload) error

// UserdatasUpdateUserdataEncoderFunc is the type to encode multipart request
// for the "userdatas" service "updateUserdata" endpoint.
type UserdatasUpdateUserdataEncoderFunc func(*multipart.Writer, *userdatas.UpdateUserdataPayload) error

// NewClient instantiates HTTP clients for all the userdatas service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		CreateUserdataDoer:  doer,
		UpdateUserdataDoer:  doer,
		GetUserdataDoer:     doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// CreateUserdata returns an endpoint that makes HTTP requests to the userdatas
// service createUserdata server.
func (c *Client) CreateUserdata(userdatasCreateUserdataEncoderFn UserdatasCreateUserdataEncoderFunc) goa.Endpoint {
	var (
		encodeRequest  = EncodeCreateUserdataRequest(NewUserdatasCreateUserdataEncoder(userdatasCreateUserdataEncoderFn))
		decodeResponse = DecodeCreateUserdataResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildCreateUserdataRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.CreateUserdataDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("userdatas", "createUserdata", err)
		}
		return decodeResponse(resp)
	}
}

// UpdateUserdata returns an endpoint that makes HTTP requests to the userdatas
// service updateUserdata server.
func (c *Client) UpdateUserdata(userdatasUpdateUserdataEncoderFn UserdatasUpdateUserdataEncoderFunc) goa.Endpoint {
	var (
		encodeRequest  = EncodeUpdateUserdataRequest(NewUserdatasUpdateUserdataEncoder(userdatasUpdateUserdataEncoderFn))
		decodeResponse = DecodeUpdateUserdataResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUpdateUserdataRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UpdateUserdataDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("userdatas", "updateUserdata", err)
		}
		return decodeResponse(resp)
	}
}

// GetUserdata returns an endpoint that makes HTTP requests to the userdatas
// service getUserdata server.
func (c *Client) GetUserdata() goa.Endpoint {
	var (
		decodeResponse = DecodeGetUserdataResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetUserdataRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetUserdataDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("userdatas", "getUserdata", err)
		}
		return decodeResponse(resp)
	}
}