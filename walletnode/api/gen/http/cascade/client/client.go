// Code generated by goa v3.13.0, DO NOT EDIT.
//
// cascade client HTTP transport
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
	cascade "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the cascade service endpoint HTTP clients.
type Client struct {
	// UploadAsset Doer is the HTTP client used to make requests to the uploadAsset
	// endpoint.
	UploadAssetDoer goahttp.Doer

	// StartProcessing Doer is the HTTP client used to make requests to the
	// startProcessing endpoint.
	StartProcessingDoer goahttp.Doer

	// RegisterTaskState Doer is the HTTP client used to make requests to the
	// registerTaskState endpoint.
	RegisterTaskStateDoer goahttp.Doer

	// GetTaskHistory Doer is the HTTP client used to make requests to the
	// getTaskHistory endpoint.
	GetTaskHistoryDoer goahttp.Doer

	// Download Doer is the HTTP client used to make requests to the download
	// endpoint.
	DownloadDoer goahttp.Doer

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

// CascadeUploadAssetEncoderFunc is the type to encode multipart request for
// the "cascade" service "uploadAsset" endpoint.
type CascadeUploadAssetEncoderFunc func(*multipart.Writer, *cascade.UploadAssetPayload) error

// NewClient instantiates HTTP clients for all the cascade service servers.
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
		UploadAssetDoer:       doer,
		StartProcessingDoer:   doer,
		RegisterTaskStateDoer: doer,
		GetTaskHistoryDoer:    doer,
		DownloadDoer:          doer,
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

// UploadAsset returns an endpoint that makes HTTP requests to the cascade
// service uploadAsset server.
func (c *Client) UploadAsset(cascadeUploadAssetEncoderFn CascadeUploadAssetEncoderFunc) goa.Endpoint {
	var (
		encodeRequest  = EncodeUploadAssetRequest(NewCascadeUploadAssetEncoder(cascadeUploadAssetEncoderFn))
		decodeResponse = DecodeUploadAssetResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildUploadAssetRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UploadAssetDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("cascade", "uploadAsset", err)
		}
		return decodeResponse(resp)
	}
}

// StartProcessing returns an endpoint that makes HTTP requests to the cascade
// service startProcessing server.
func (c *Client) StartProcessing() goa.Endpoint {
	var (
		encodeRequest  = EncodeStartProcessingRequest(c.encoder)
		decodeResponse = DecodeStartProcessingResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
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
			return nil, goahttp.ErrRequestError("cascade", "startProcessing", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTaskState returns an endpoint that makes HTTP requests to the
// cascade service registerTaskState server.
func (c *Client) RegisterTaskState() goa.Endpoint {
	var (
		decodeResponse = DecodeRegisterTaskStateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildRegisterTaskStateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		conn, resp, err := c.dialer.DialContext(ctx, req.URL.String(), req.Header)
		if err != nil {
			if resp != nil {
				return decodeResponse(resp)
			}
			return nil, goahttp.ErrRequestError("cascade", "registerTaskState", err)
		}
		if c.configurer.RegisterTaskStateFn != nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			conn = c.configurer.RegisterTaskStateFn(conn, cancel)
		}
		go func() {
			<-ctx.Done()
			conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client closing connection"),
				time.Now().UTC().Add(time.Second),
			)
			conn.Close()
		}()
		stream := &RegisterTaskStateClientStream{conn: conn}
		return stream, nil
	}
}

// GetTaskHistory returns an endpoint that makes HTTP requests to the cascade
// service getTaskHistory server.
func (c *Client) GetTaskHistory() goa.Endpoint {
	var (
		decodeResponse = DecodeGetTaskHistoryResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildGetTaskHistoryRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetTaskHistoryDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("cascade", "getTaskHistory", err)
		}
		return decodeResponse(resp)
	}
}

// Download returns an endpoint that makes HTTP requests to the cascade service
// download server.
func (c *Client) Download() goa.Endpoint {
	var (
		encodeRequest  = EncodeDownloadRequest(c.encoder)
		decodeResponse = DecodeDownloadResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildDownloadRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DownloadDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("cascade", "download", err)
		}
		return decodeResponse(resp)
	}
}
