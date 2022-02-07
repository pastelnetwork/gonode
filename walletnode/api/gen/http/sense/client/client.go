// Code generated by goa v3.4.3, DO NOT EDIT.
//
// sense client HTTP transport
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the sense service endpoint HTTP clients.
type Client struct {
	// UploadImage Doer is the HTTP client used to make requests to the uploadImage
	// endpoint.
	UploadImageDoer goahttp.Doer

	// ActionDetails Doer is the HTTP client used to make requests to the
	// actionDetails endpoint.
	ActionDetailsDoer goahttp.Doer

	// StartProcessing Doer is the HTTP client used to make requests to the
	// startProcessing endpoint.
	StartProcessingDoer goahttp.Doer

	// RegisterTaskState Doer is the HTTP client used to make requests to the
	// registerTaskState endpoint.
	RegisterTaskStateDoer goahttp.Doer

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

// SenseUploadImageEncoderFunc is the type to encode multipart request for the
// "sense" service "uploadImage" endpoint.
type SenseUploadImageEncoderFunc func(*multipart.Writer, *sense.UploadImagePayload) error

// NewClient instantiates HTTP clients for all the sense service servers.
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
		UploadImageDoer:       doer,
		ActionDetailsDoer:     doer,
		StartProcessingDoer:   doer,
		RegisterTaskStateDoer: doer,
		CORSDoer:              doer,
		RestoreResponseBody:   restoreBody,
		scheme:                scheme,
		host:                  host,
		decoder:               dec,
		encoder:               enc,
		dialer:                dialer,
		configurer:            cfn,
	}
}

// UploadImage returns an endpoint that makes HTTP requests to the sense
// service uploadImage server.
func (c *Client) UploadImage(senseUploadImageEncoderFn SenseUploadImageEncoderFunc) goa.Endpoint {
	var (
		encodeRequest  = EncodeUploadImageRequest(NewSenseUploadImageEncoder(senseUploadImageEncoderFn))
		decodeResponse = DecodeUploadImageResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUploadImageRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UploadImageDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("sense", "uploadImage", err)
		}
		return decodeResponse(resp)
	}
}

// ActionDetails returns an endpoint that makes HTTP requests to the sense
// service actionDetails server.
func (c *Client) ActionDetails() goa.Endpoint {
	var (
		encodeRequest  = EncodeActionDetailsRequest(c.encoder)
		decodeResponse = DecodeActionDetailsResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildActionDetailsRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ActionDetailsDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("sense", "actionDetails", err)
		}
		return decodeResponse(resp)
	}
}

// StartProcessing returns an endpoint that makes HTTP requests to the sense
// service startProcessing server.
func (c *Client) StartProcessing() goa.Endpoint {
	var (
		encodeRequest  = EncodeStartProcessingRequest(c.encoder)
		decodeResponse = DecodeStartProcessingResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildStartProcessingRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.StartProcessingDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("sense", "startProcessing", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTaskState returns an endpoint that makes HTTP requests to the sense
// service registerTaskState server.
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
			return nil, goahttp.ErrRequestError("sense", "registerTaskState", err)
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
