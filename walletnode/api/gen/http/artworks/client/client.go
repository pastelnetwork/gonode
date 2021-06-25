// Code generated by goa v3.3.1, DO NOT EDIT.
//
// artworks client HTTP transport
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
	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the artworks service endpoint HTTP clients.
type Client struct {
	// Register Doer is the HTTP client used to make requests to the register
	// endpoint.
	RegisterDoer goahttp.Doer

	// RegisterTaskState Doer is the HTTP client used to make requests to the
	// registerTaskState endpoint.
	RegisterTaskStateDoer goahttp.Doer

	// RegisterTask Doer is the HTTP client used to make requests to the
	// registerTask endpoint.
	RegisterTaskDoer goahttp.Doer

	// RegisterTasks Doer is the HTTP client used to make requests to the
	// registerTasks endpoint.
	RegisterTasksDoer goahttp.Doer

	// UploadImage Doer is the HTTP client used to make requests to the uploadImage
	// endpoint.
	UploadImageDoer goahttp.Doer

	// Download Doer is the HTTP client used to make requests to the download
	// endpoint.
	DownloadDoer goahttp.Doer

	// DownloadTaskStateEndpoint Doer is the HTTP client used to make requests to
	// the downloadTaskState endpoint.
	DownloadTaskStateEndpointDoer goahttp.Doer

	// DowloadTask Doer is the HTTP client used to make requests to the dowloadTask
	// endpoint.
	DowloadTaskDoer goahttp.Doer

	// DownloadTasks Doer is the HTTP client used to make requests to the
	// downloadTasks endpoint.
	DownloadTasksDoer goahttp.Doer

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

// ArtworksUploadImageEncoderFunc is the type to encode multipart request for
// the "artworks" service "uploadImage" endpoint.
type ArtworksUploadImageEncoderFunc func(*multipart.Writer, *artworks.UploadImagePayload) error

// NewClient instantiates HTTP clients for all the artworks service servers.
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
		RegisterDoer:                  doer,
		RegisterTaskStateDoer:         doer,
		RegisterTaskDoer:              doer,
		RegisterTasksDoer:             doer,
		UploadImageDoer:               doer,
		DownloadDoer:                  doer,
		DownloadTaskStateEndpointDoer: doer,
		DowloadTaskDoer:               doer,
		DownloadTasksDoer:             doer,
		CORSDoer:                      doer,
		RestoreResponseBody:           restoreBody,
		scheme:                        scheme,
		host:                          host,
		decoder:                       dec,
		encoder:                       enc,
		dialer:                        dialer,
		configurer:                    cfn,
	}
}

// Register returns an endpoint that makes HTTP requests to the artworks
// service register server.
func (c *Client) Register() goa.Endpoint {
	var (
		encodeRequest  = EncodeRegisterRequest(c.encoder)
		decodeResponse = DecodeRegisterResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RegisterDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "register", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTaskState returns an endpoint that makes HTTP requests to the
// artworks service registerTaskState server.
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
		ctx, cancel = context.WithCancel(ctx)
		conn, resp, err := c.dialer.DialContext(ctx, req.URL.String(), req.Header)
		if err != nil {
			if resp != nil {
				return decodeResponse(resp)
			}
			return nil, goahttp.ErrRequestError("artworks", "registerTaskState", err)
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

// RegisterTask returns an endpoint that makes HTTP requests to the artworks
// service registerTask server.
func (c *Client) RegisterTask() goa.Endpoint {
	var (
		decodeResponse = DecodeRegisterTaskResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterTaskRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RegisterTaskDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "registerTask", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTasks returns an endpoint that makes HTTP requests to the artworks
// service registerTasks server.
func (c *Client) RegisterTasks() goa.Endpoint {
	var (
		decodeResponse = DecodeRegisterTasksResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterTasksRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RegisterTasksDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "registerTasks", err)
		}
		return decodeResponse(resp)
	}
}

// UploadImage returns an endpoint that makes HTTP requests to the artworks
// service uploadImage server.
func (c *Client) UploadImage(artworksUploadImageEncoderFn ArtworksUploadImageEncoderFunc) goa.Endpoint {
	var (
		encodeRequest  = EncodeUploadImageRequest(NewArtworksUploadImageEncoder(artworksUploadImageEncoderFn))
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
			return nil, goahttp.ErrRequestError("artworks", "uploadImage", err)
		}
		return decodeResponse(resp)
	}
}

// Download returns an endpoint that makes HTTP requests to the artworks
// service download server.
func (c *Client) Download() goa.Endpoint {
	var (
		encodeRequest  = EncodeDownloadRequest(c.encoder)
		decodeResponse = DecodeDownloadResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
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
			return nil, goahttp.ErrRequestError("artworks", "download", err)
		}
		return decodeResponse(resp)
	}
}

// DownloadTaskStateEndpoint returns an endpoint that makes HTTP requests to
// the artworks service downloadTaskState server.
func (c *Client) DownloadTaskStateEndpoint() goa.Endpoint {
	var (
		decodeResponse = DecodeDownloadTaskStateEndpointResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDownloadTaskStateEndpointRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		conn, resp, err := c.dialer.DialContext(ctx, req.URL.String(), req.Header)
		if err != nil {
			if resp != nil {
				return decodeResponse(resp)
			}
			return nil, goahttp.ErrRequestError("artworks", "downloadTaskState", err)
		}
		if c.configurer.DownloadTaskStateEndpointFn != nil {
			conn = c.configurer.DownloadTaskStateEndpointFn(conn, cancel)
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
		stream := &DownloadTaskStateClientStream{conn: conn}
		return stream, nil
	}
}

// DowloadTask returns an endpoint that makes HTTP requests to the artworks
// service dowloadTask server.
func (c *Client) DowloadTask() goa.Endpoint {
	var (
		decodeResponse = DecodeDowloadTaskResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDowloadTaskRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DowloadTaskDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "dowloadTask", err)
		}
		return decodeResponse(resp)
	}
}

// DownloadTasks returns an endpoint that makes HTTP requests to the artworks
// service downloadTasks server.
func (c *Client) DownloadTasks() goa.Endpoint {
	var (
		decodeResponse = DecodeDownloadTasksResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDownloadTasksRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DownloadTasksDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "downloadTasks", err)
		}
		return decodeResponse(resp)
	}
}
