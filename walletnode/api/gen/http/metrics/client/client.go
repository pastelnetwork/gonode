// Code generated by goa v3.14.0, DO NOT EDIT.
//
// metrics client HTTP transport
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the metrics service endpoint HTTP clients.
type Client struct {
	// SelfHealing Doer is the HTTP client used to make requests to the selfHealing
	// endpoint.
	SelfHealingDoer goahttp.Doer

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

// NewClient instantiates HTTP clients for all the metrics service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		SelfHealingDoer:     doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// SelfHealing returns an endpoint that makes HTTP requests to the metrics
// service selfHealing server.
func (c *Client) SelfHealing() goa.Endpoint {
	var (
		encodeRequest  = EncodeSelfHealingRequest(c.encoder)
		decodeResponse = DecodeSelfHealingResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildSelfHealingRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.SelfHealingDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("metrics", "selfHealing", err)
		}
		return decodeResponse(resp)
	}
}
