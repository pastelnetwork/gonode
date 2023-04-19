// Code generated by goa v3.7.6, DO NOT EDIT.
//
// collection client HTTP transport
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the collection service endpoint HTTP clients.
type Client struct {
	// RegisterCollection Doer is the HTTP client used to make requests to the
	// registerCollection endpoint.
	RegisterCollectionDoer goahttp.Doer

	// RegisterTaskState Doer is the HTTP client used to make requests to the
	// registerTaskState endpoint.
	RegisterTaskStateDoer goahttp.Doer

	// GetTaskHistory Doer is the HTTP client used to make requests to the
	// getTaskHistory endpoint.
	GetTaskHistoryDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme     string
	host       string
	encoder    func(*http.Request) goahttp.Encoder
	decoder    func(*http.Response) goahttp.Decoder
	dialer     goahttp.Dialer
	configurer *ConnConfigurer
}

// NewClient instantiates HTTP clients for all the collection service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
	dialer goahttp.Dialer,
	cfn *ConnConfigurer,
) *Client {
	if cfn == nil {
		cfn = &ConnConfigurer{}
	}
	return &Client{
		RegisterCollectionDoer: doer,
		RegisterTaskStateDoer:  doer,
		GetTaskHistoryDoer:     doer,
		CORSDoer:               doer,
		RestoreResponseBody:    restoreBody,
		scheme:                 scheme,
		host:                   host,
		decoder:                dec,
		encoder:                enc,
		dialer:                 dialer,
		configurer:             cfn,
	}
}

// RegisterCollection returns an endpoint that makes HTTP requests to the
// collection service registerCollection server.
func (c *Client) RegisterCollection() goa.Endpoint {
	var (
		encodeRequest  = EncodeRegisterCollectionRequest(c.encoder)
		decodeResponse = DecodeRegisterCollectionResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterCollectionRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RegisterCollectionDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("collection", "registerCollection", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTaskState returns an endpoint that makes HTTP requests to the
// collection service registerTaskState server.
func (c *Client) RegisterTaskState() goa.Endpoint {
	var (
		decodeResponse = DecodeRegisterTaskStateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterTaskStateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		var cancel context.CancelFunc
		_, cancel = context.WithCancel(ctx)
		defer cancel()
		conn, resp, err := c.dialer.DialContext(ctx, req.URL.String(), req.Header)
		if err != nil {
			if resp != nil {
				return decodeResponse(resp)
			}
			return nil, goahttp.ErrRequestError("collection", "registerTaskState", err)
		}
		if c.configurer.RegisterTaskStateFn != nil {
			conn = c.configurer.RegisterTaskStateFn(conn, cancel)
		}
		go func() {
			<-ctx.Done()
			conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client closing connection"),
				time.Now().Add(time.Second),
			)
			conn.Close()
		}()
		stream := &RegisterTaskStateClientStream{conn: conn}
		return stream, nil
	}
}

// GetTaskHistory returns an endpoint that makes HTTP requests to the
// collection service getTaskHistory server.
func (c *Client) GetTaskHistory() goa.Endpoint {
	var (
		decodeResponse = DecodeGetTaskHistoryResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetTaskHistoryRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetTaskHistoryDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("collection", "getTaskHistory", err)
		}
		return decodeResponse(resp)
	}
}