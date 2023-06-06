// Code generated by goa v3.7.6, DO NOT EDIT.
//
// nft client HTTP transport
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
	nft "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the nft service endpoint HTTP clients.
type Client struct {
	// Register Doer is the HTTP client used to make requests to the register
	// endpoint.
	RegisterDoer goahttp.Doer

	// RegisterTaskState Doer is the HTTP client used to make requests to the
	// registerTaskState endpoint.
	RegisterTaskStateDoer goahttp.Doer

	// GetTaskHistory Doer is the HTTP client used to make requests to the
	// getTaskHistory endpoint.
	GetTaskHistoryDoer goahttp.Doer

	// RegisterTask Doer is the HTTP client used to make requests to the
	// registerTask endpoint.
	RegisterTaskDoer goahttp.Doer

	// RegisterTasks Doer is the HTTP client used to make requests to the
	// registerTasks endpoint.
	RegisterTasksDoer goahttp.Doer

	// UploadImage Doer is the HTTP client used to make requests to the uploadImage
	// endpoint.
	UploadImageDoer goahttp.Doer

	// NftSearch Doer is the HTTP client used to make requests to the nftSearch
	// endpoint.
	NftSearchDoer goahttp.Doer

	// NftGet Doer is the HTTP client used to make requests to the nftGet endpoint.
	NftGetDoer goahttp.Doer

	// Download Doer is the HTTP client used to make requests to the download
	// endpoint.
	DownloadDoer goahttp.Doer

	// DdServiceOutputFileDetail Doer is the HTTP client used to make requests to
	// the ddServiceOutputFileDetail endpoint.
	DdServiceOutputFileDetailDoer goahttp.Doer

	// DdServiceOutputFile Doer is the HTTP client used to make requests to the
	// ddServiceOutputFile endpoint.
	DdServiceOutputFileDoer goahttp.Doer

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

// NftUploadImageEncoderFunc is the type to encode multipart request for the
// "nft" service "uploadImage" endpoint.
type NftUploadImageEncoderFunc func(*multipart.Writer, *nft.UploadImagePayload) error

// NewClient instantiates HTTP clients for all the nft service servers.
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
		GetTaskHistoryDoer:            doer,
		RegisterTaskDoer:              doer,
		RegisterTasksDoer:             doer,
		UploadImageDoer:               doer,
		NftSearchDoer:                 doer,
		NftGetDoer:                    doer,
		DownloadDoer:                  doer,
		DdServiceOutputFileDetailDoer: doer,
		DdServiceOutputFileDoer:       doer,
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

// Register returns an endpoint that makes HTTP requests to the nft service
// register server.
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
			return nil, goahttp.ErrRequestError("nft", "register", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTaskState returns an endpoint that makes HTTP requests to the nft
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
		ctx, cancel = context.WithCancel(ctx)
		conn, resp, err := c.dialer.DialContext(ctx, req.URL.String(), req.Header)
		if err != nil {
			if resp != nil {
				return decodeResponse(resp)
			}
			return nil, goahttp.ErrRequestError("nft", "registerTaskState", err)
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

// GetTaskHistory returns an endpoint that makes HTTP requests to the nft
// service getTaskHistory server.
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
			return nil, goahttp.ErrRequestError("nft", "getTaskHistory", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTask returns an endpoint that makes HTTP requests to the nft service
// registerTask server.
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
			return nil, goahttp.ErrRequestError("nft", "registerTask", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterTasks returns an endpoint that makes HTTP requests to the nft
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
			return nil, goahttp.ErrRequestError("nft", "registerTasks", err)
		}
		return decodeResponse(resp)
	}
}

// UploadImage returns an endpoint that makes HTTP requests to the nft service
// uploadImage server.
func (c *Client) UploadImage(nftUploadImageEncoderFn NftUploadImageEncoderFunc) goa.Endpoint {
	var (
		encodeRequest  = EncodeUploadImageRequest(NewNftUploadImageEncoder(nftUploadImageEncoderFn))
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
			return nil, goahttp.ErrRequestError("nft", "uploadImage", err)
		}
		return decodeResponse(resp)
	}
}

// NftSearch returns an endpoint that makes HTTP requests to the nft service
// nftSearch server.
func (c *Client) NftSearch() goa.Endpoint {
	var (
		encodeRequest  = EncodeNftSearchRequest(c.encoder)
		decodeResponse = DecodeNftSearchResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildNftSearchRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
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
			return nil, goahttp.ErrRequestError("nft", "nftSearch", err)
		}
		if c.configurer.NftSearchFn != nil {
			conn = c.configurer.NftSearchFn(conn, cancel)
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
		stream := &NftSearchClientStream{conn: conn}
		return stream, nil
	}
}

// NftGet returns an endpoint that makes HTTP requests to the nft service
// nftGet server.
func (c *Client) NftGet() goa.Endpoint {
	var (
		encodeRequest  = EncodeNftGetRequest(c.encoder)
		decodeResponse = DecodeNftGetResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildNftGetRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.NftGetDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("nft", "nftGet", err)
		}
		return decodeResponse(resp)
	}
}

// Download returns an endpoint that makes HTTP requests to the nft service
// download server.
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
			return nil, goahttp.ErrRequestError("nft", "download", err)
		}
		return decodeResponse(resp)
	}
}

// DdServiceOutputFileDetail returns an endpoint that makes HTTP requests to
// the nft service ddServiceOutputFileDetail server.
func (c *Client) DdServiceOutputFileDetail() goa.Endpoint {
	var (
		encodeRequest  = EncodeDdServiceOutputFileDetailRequest(c.encoder)
		decodeResponse = DecodeDdServiceOutputFileDetailResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDdServiceOutputFileDetailRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DdServiceOutputFileDetailDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("nft", "ddServiceOutputFileDetail", err)
		}
		return decodeResponse(resp)
	}
}

// DdServiceOutputFile returns an endpoint that makes HTTP requests to the nft
// service ddServiceOutputFile server.
func (c *Client) DdServiceOutputFile() goa.Endpoint {
	var (
		encodeRequest  = EncodeDdServiceOutputFileRequest(c.encoder)
		decodeResponse = DecodeDdServiceOutputFileResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDdServiceOutputFileRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DdServiceOutputFileDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("nft", "ddServiceOutputFile", err)
		}
		return decodeResponse(resp)
	}
}
