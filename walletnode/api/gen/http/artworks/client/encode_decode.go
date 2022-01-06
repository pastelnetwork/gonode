// Code generated by goa v3.4.3, DO NOT EDIT.
//
// artworks HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	artworks "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	artworksviews "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks/views"
	goahttp "goa.design/goa/v3/http"
)

// BuildRegisterRequest instantiates a HTTP request object with method and path
// set to call the "artworks" service "register" endpoint
func (c *Client) BuildRegisterRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RegisterArtworksPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "register", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeRegisterRequest returns an encoder for requests sent to the artworks
// register server.
func EncodeRegisterRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*artworks.RegisterPayload)
		if !ok {
			return goahttp.ErrInvalidType("artworks", "register", "*artworks.RegisterPayload", v)
		}
		body := NewRegisterRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("artworks", "register", err)
		}
		return nil
	}
}

// DecodeRegisterResponse returns a decoder for responses returned by the
// artworks register endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeRegisterResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeRegisterResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusCreated:
			var (
				body RegisterResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "register", err)
			}
			p := NewRegisterResultViewCreated(&body)
			view := "default"
			vres := &artworksviews.RegisterResult{Projected: p, View: view}
			if err = artworksviews.ValidateRegisterResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("artworks", "register", err)
			}
			res := artworks.NewRegisterResult(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body RegisterBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "register", err)
			}
			err = ValidateRegisterBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "register", err)
			}
			return nil, NewRegisterBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "register", err)
			}
			err = ValidateRegisterInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "register", err)
			}
			return nil, NewRegisterInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "register", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTaskStateRequest instantiates a HTTP request object with method
// and path set to call the "artworks" service "registerTaskState" endpoint
func (c *Client) BuildRegisterTaskStateRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*artworks.RegisterTaskStatePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("artworks", "registerTaskState", "*artworks.RegisterTaskStatePayload", v)
		}
		taskID = p.TaskID
	}
	scheme := c.scheme
	switch c.scheme {
	case "http":
		scheme = "ws"
	case "https":
		scheme = "wss"
	}
	u := &url.URL{Scheme: scheme, Host: c.host, Path: RegisterTaskStateArtworksPath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "registerTaskState", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTaskStateResponse returns a decoder for responses returned by
// the artworks registerTaskState endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeRegisterTaskStateResponse may return the following errors:
//	- "NotFound" (type *goa.ServiceError): http.StatusNotFound
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeRegisterTaskStateResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body RegisterTaskStateResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTaskState", err)
			}
			res := NewRegisterTaskStateTaskStateOK(&body)
			return res, nil
		case http.StatusNotFound:
			var (
				body RegisterTaskStateNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterTaskStateInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "registerTaskState", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTaskRequest instantiates a HTTP request object with method and
// path set to call the "artworks" service "registerTask" endpoint
func (c *Client) BuildRegisterTaskRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*artworks.RegisterTaskPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("artworks", "registerTask", "*artworks.RegisterTaskPayload", v)
		}
		taskID = p.TaskID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RegisterTaskArtworksPath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "registerTask", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTaskResponse returns a decoder for responses returned by the
// artworks registerTask endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeRegisterTaskResponse may return the following errors:
//	- "NotFound" (type *goa.ServiceError): http.StatusNotFound
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeRegisterTaskResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body RegisterTaskResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTask", err)
			}
			p := NewRegisterTaskTaskOK(&body)
			view := "default"
			vres := &artworksviews.Task{Projected: p, View: view}
			if err = artworksviews.ValidateTask(vres); err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTask", err)
			}
			res := artworks.NewTask(vres)
			return res, nil
		case http.StatusNotFound:
			var (
				body RegisterTaskNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTask", err)
			}
			err = ValidateRegisterTaskNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTask", err)
			}
			return nil, NewRegisterTaskNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterTaskInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTask", err)
			}
			err = ValidateRegisterTaskInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTask", err)
			}
			return nil, NewRegisterTaskInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "registerTask", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTasksRequest instantiates a HTTP request object with method and
// path set to call the "artworks" service "registerTasks" endpoint
func (c *Client) BuildRegisterTasksRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RegisterTasksArtworksPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "registerTasks", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTasksResponse returns a decoder for responses returned by the
// artworks registerTasks endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeRegisterTasksResponse may return the following errors:
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeRegisterTasksResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body RegisterTasksResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTasks", err)
			}
			p := NewRegisterTasksTaskCollectionOK(body)
			view := "tiny"
			vres := artworksviews.TaskCollection{Projected: p, View: view}
			if err = artworksviews.ValidateTaskCollection(vres); err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTasks", err)
			}
			res := artworks.NewTaskCollection(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body RegisterTasksInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "registerTasks", err)
			}
			err = ValidateRegisterTasksInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "registerTasks", err)
			}
			return nil, NewRegisterTasksInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "registerTasks", resp.StatusCode, string(body))
		}
	}
}

// BuildUploadImageRequest instantiates a HTTP request object with method and
// path set to call the "artworks" service "uploadImage" endpoint
func (c *Client) BuildUploadImageRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: UploadImageArtworksPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "uploadImage", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeUploadImageRequest returns an encoder for requests sent to the
// artworks uploadImage server.
func EncodeUploadImageRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*artworks.UploadImagePayload)
		if !ok {
			return goahttp.ErrInvalidType("artworks", "uploadImage", "*artworks.UploadImagePayload", v)
		}
		if err := encoder(req).Encode(p); err != nil {
			return goahttp.ErrEncodingError("artworks", "uploadImage", err)
		}
		return nil
	}
}

// NewArtworksUploadImageEncoder returns an encoder to encode the multipart
// request for the "artworks" service "uploadImage" endpoint.
func NewArtworksUploadImageEncoder(encoderFn ArtworksUploadImageEncoderFunc) func(r *http.Request) goahttp.Encoder {
	return func(r *http.Request) goahttp.Encoder {
		body := &bytes.Buffer{}
		mw := multipart.NewWriter(body)
		return goahttp.EncodingFunc(func(v interface{}) error {
			p := v.(*artworks.UploadImagePayload)
			if err := encoderFn(mw, p); err != nil {
				return err
			}
			r.Body = ioutil.NopCloser(body)
			r.Header.Set("Content-Type", mw.FormDataContentType())
			return mw.Close()
		})
	}
}

// DecodeUploadImageResponse returns a decoder for responses returned by the
// artworks uploadImage endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeUploadImageResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeUploadImageResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusCreated:
			var (
				body UploadImageResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "uploadImage", err)
			}
			p := NewUploadImageImageCreated(&body)
			view := "default"
			vres := &artworksviews.Image{Projected: p, View: view}
			if err = artworksviews.ValidateImage(vres); err != nil {
				return nil, goahttp.ErrValidationError("artworks", "uploadImage", err)
			}
			res := artworks.NewImage(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body UploadImageBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "uploadImage", err)
			}
			err = ValidateUploadImageBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "uploadImage", err)
			}
			return nil, NewUploadImageBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body UploadImageInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "uploadImage", err)
			}
			err = ValidateUploadImageInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "uploadImage", err)
			}
			return nil, NewUploadImageInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "uploadImage", resp.StatusCode, string(body))
		}
	}
}

// BuildArtSearchRequest instantiates a HTTP request object with method and
// path set to call the "artworks" service "artSearch" endpoint
func (c *Client) BuildArtSearchRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	scheme := c.scheme
	switch c.scheme {
	case "http":
		scheme = "ws"
	case "https":
		scheme = "wss"
	}
	u := &url.URL{Scheme: scheme, Host: c.host, Path: ArtSearchArtworksPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "artSearch", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeArtSearchRequest returns an encoder for requests sent to the artworks
// artSearch server.
func EncodeArtSearchRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*artworks.ArtSearchPayload)
		if !ok {
			return goahttp.ErrInvalidType("artworks", "artSearch", "*artworks.ArtSearchPayload", v)
		}
		if p.UserPastelid != nil {
			head := *p.UserPastelid
			req.Header.Set("user_pastelid", head)
		}
		if p.UserPassphrase != nil {
			head := *p.UserPassphrase
			req.Header.Set("user_passphrase", head)
		}
		values := req.URL.Query()
		if p.Artist != nil {
			values.Add("artist", *p.Artist)
		}
		values.Add("limit", fmt.Sprintf("%v", p.Limit))
		values.Add("query", p.Query)
		values.Add("artist_name", fmt.Sprintf("%v", p.ArtistName))
		values.Add("art_title", fmt.Sprintf("%v", p.ArtTitle))
		values.Add("series", fmt.Sprintf("%v", p.Series))
		values.Add("descr", fmt.Sprintf("%v", p.Descr))
		values.Add("keyword", fmt.Sprintf("%v", p.Keyword))
		if p.MinCopies != nil {
			values.Add("min_copies", fmt.Sprintf("%v", *p.MinCopies))
		}
		if p.MaxCopies != nil {
			values.Add("max_copies", fmt.Sprintf("%v", *p.MaxCopies))
		}
		values.Add("min_block", fmt.Sprintf("%v", p.MinBlock))
		if p.MaxBlock != nil {
			values.Add("max_block", fmt.Sprintf("%v", *p.MaxBlock))
		}
		if p.MinRarenessScore != nil {
			values.Add("min_rareness_score", fmt.Sprintf("%v", *p.MinRarenessScore))
		}
		if p.MaxRarenessScore != nil {
			values.Add("max_rareness_score", fmt.Sprintf("%v", *p.MaxRarenessScore))
		}
		if p.MinNsfwScore != nil {
			values.Add("min_nsfw_score", fmt.Sprintf("%v", *p.MinNsfwScore))
		}
		if p.MaxNsfwScore != nil {
			values.Add("max_nsfw_score", fmt.Sprintf("%v", *p.MaxNsfwScore))
		}
		if p.MinInternetRarenessScore != nil {
			values.Add("min_internet_rareness_score", fmt.Sprintf("%v", *p.MinInternetRarenessScore))
		}
		if p.MaxInternetRarenessScore != nil {
			values.Add("max_internet_rareness_score", fmt.Sprintf("%v", *p.MaxInternetRarenessScore))
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeArtSearchResponse returns a decoder for responses returned by the
// artworks artSearch endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeArtSearchResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeArtSearchResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ArtSearchResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "artSearch", err)
			}
			err = ValidateArtSearchResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "artSearch", err)
			}
			res := NewArtSearchArtworkSearchResultOK(&body)
			return res, nil
		case http.StatusBadRequest:
			var (
				body ArtSearchBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "artSearch", err)
			}
			err = ValidateArtSearchBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "artSearch", err)
			}
			return nil, NewArtSearchBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body ArtSearchInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "artSearch", err)
			}
			err = ValidateArtSearchInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "artSearch", err)
			}
			return nil, NewArtSearchInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "artSearch", resp.StatusCode, string(body))
		}
	}
}

// BuildArtworkGetRequest instantiates a HTTP request object with method and
// path set to call the "artworks" service "artworkGet" endpoint
func (c *Client) BuildArtworkGetRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		txid string
	)
	{
		p, ok := v.(*artworks.ArtworkGetPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("artworks", "artworkGet", "*artworks.ArtworkGetPayload", v)
		}
		txid = p.Txid
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ArtworkGetArtworksPath(txid)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "artworkGet", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeArtworkGetRequest returns an encoder for requests sent to the artworks
// artworkGet server.
func EncodeArtworkGetRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*artworks.ArtworkGetPayload)
		if !ok {
			return goahttp.ErrInvalidType("artworks", "artworkGet", "*artworks.ArtworkGetPayload", v)
		}
		body := NewArtworkGetRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("artworks", "artworkGet", err)
		}
		return nil
	}
}

// DecodeArtworkGetResponse returns a decoder for responses returned by the
// artworks artworkGet endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeArtworkGetResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "NotFound" (type *goa.ServiceError): http.StatusNotFound
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeArtworkGetResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ArtworkGetResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "artworkGet", err)
			}
			err = ValidateArtworkGetResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "artworkGet", err)
			}
			res := NewArtworkGetArtworkDetailOK(&body)
			return res, nil
		case http.StatusBadRequest:
			var (
				body ArtworkGetBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "artworkGet", err)
			}
			err = ValidateArtworkGetBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "artworkGet", err)
			}
			return nil, NewArtworkGetBadRequest(&body)
		case http.StatusNotFound:
			var (
				body ArtworkGetNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "artworkGet", err)
			}
			err = ValidateArtworkGetNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "artworkGet", err)
			}
			return nil, NewArtworkGetNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body ArtworkGetInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "artworkGet", err)
			}
			err = ValidateArtworkGetInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "artworkGet", err)
			}
			return nil, NewArtworkGetInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "artworkGet", resp.StatusCode, string(body))
		}
	}
}

// BuildDownloadRequest instantiates a HTTP request object with method and path
// set to call the "artworks" service "download" endpoint
func (c *Client) BuildDownloadRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: DownloadArtworksPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("artworks", "download", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeDownloadRequest returns an encoder for requests sent to the artworks
// download server.
func EncodeDownloadRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*artworks.ArtworkDownloadPayload)
		if !ok {
			return goahttp.ErrInvalidType("artworks", "download", "*artworks.ArtworkDownloadPayload", v)
		}
		{
			head := p.Key
			req.Header.Set("Authorization", head)
		}
		values := req.URL.Query()
		values.Add("txid", p.Txid)
		values.Add("pid", p.Pid)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeDownloadResponse returns a decoder for responses returned by the
// artworks download endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeDownloadResponse may return the following errors:
//	- "NotFound" (type *goa.ServiceError): http.StatusNotFound
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeDownloadResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body DownloadResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "download", err)
			}
			err = ValidateDownloadResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "download", err)
			}
			res := NewDownloadResultOK(&body)
			return res, nil
		case http.StatusNotFound:
			var (
				body DownloadNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "download", err)
			}
			err = ValidateDownloadNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "download", err)
			}
			return nil, NewDownloadNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body DownloadInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("artworks", "download", err)
			}
			err = ValidateDownloadInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("artworks", "download", err)
			}
			return nil, NewDownloadInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("artworks", "download", resp.StatusCode, string(body))
		}
	}
}

// marshalArtworksThumbnailcoordinateToThumbnailcoordinateRequestBody builds a
// value of type *ThumbnailcoordinateRequestBody from a value of type
// *artworks.Thumbnailcoordinate.
func marshalArtworksThumbnailcoordinateToThumbnailcoordinateRequestBody(v *artworks.Thumbnailcoordinate) *ThumbnailcoordinateRequestBody {
	if v == nil {
		return nil
	}
	res := &ThumbnailcoordinateRequestBody{
		TopLeftX:     v.TopLeftX,
		TopLeftY:     v.TopLeftY,
		BottomRightX: v.BottomRightX,
		BottomRightY: v.BottomRightY,
	}

	return res
}

// marshalThumbnailcoordinateRequestBodyToArtworksThumbnailcoordinate builds a
// value of type *artworks.Thumbnailcoordinate from a value of type
// *ThumbnailcoordinateRequestBody.
func marshalThumbnailcoordinateRequestBodyToArtworksThumbnailcoordinate(v *ThumbnailcoordinateRequestBody) *artworks.Thumbnailcoordinate {
	if v == nil {
		return nil
	}
	res := &artworks.Thumbnailcoordinate{
		TopLeftX:     v.TopLeftX,
		TopLeftY:     v.TopLeftY,
		BottomRightX: v.BottomRightX,
		BottomRightY: v.BottomRightY,
	}

	return res
}

// unmarshalTaskStateResponseBodyToArtworksviewsTaskStateView builds a value of
// type *artworksviews.TaskStateView from a value of type
// *TaskStateResponseBody.
func unmarshalTaskStateResponseBodyToArtworksviewsTaskStateView(v *TaskStateResponseBody) *artworksviews.TaskStateView {
	if v == nil {
		return nil
	}
	res := &artworksviews.TaskStateView{
		Date:   v.Date,
		Status: v.Status,
	}

	return res
}

// unmarshalArtworkTicketResponseBodyToArtworksviewsArtworkTicketView builds a
// value of type *artworksviews.ArtworkTicketView from a value of type
// *ArtworkTicketResponseBody.
func unmarshalArtworkTicketResponseBodyToArtworksviewsArtworkTicketView(v *ArtworkTicketResponseBody) *artworksviews.ArtworkTicketView {
	res := &artworksviews.ArtworkTicketView{
		Name:                     v.Name,
		Description:              v.Description,
		Keywords:                 v.Keywords,
		SeriesName:               v.SeriesName,
		IssuedCopies:             v.IssuedCopies,
		YoutubeURL:               v.YoutubeURL,
		ArtistPastelID:           v.ArtistPastelID,
		ArtistPastelIDPassphrase: v.ArtistPastelIDPassphrase,
		ArtistName:               v.ArtistName,
		ArtistWebsiteURL:         v.ArtistWebsiteURL,
		SpendableAddress:         v.SpendableAddress,
		MaximumFee:               v.MaximumFee,
		Royalty:                  v.Royalty,
		Green:                    v.Green,
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = unmarshalThumbnailcoordinateResponseBodyToArtworksviewsThumbnailcoordinateView(v.ThumbnailCoordinate)
	}

	return res
}

// unmarshalThumbnailcoordinateResponseBodyToArtworksviewsThumbnailcoordinateView
// builds a value of type *artworksviews.ThumbnailcoordinateView from a value
// of type *ThumbnailcoordinateResponseBody.
func unmarshalThumbnailcoordinateResponseBodyToArtworksviewsThumbnailcoordinateView(v *ThumbnailcoordinateResponseBody) *artworksviews.ThumbnailcoordinateView {
	if v == nil {
		return nil
	}
	res := &artworksviews.ThumbnailcoordinateView{
		TopLeftX:     v.TopLeftX,
		TopLeftY:     v.TopLeftY,
		BottomRightX: v.BottomRightX,
		BottomRightY: v.BottomRightY,
	}

	return res
}

// unmarshalTaskResponseToArtworksviewsTaskView builds a value of type
// *artworksviews.TaskView from a value of type *TaskResponse.
func unmarshalTaskResponseToArtworksviewsTaskView(v *TaskResponse) *artworksviews.TaskView {
	res := &artworksviews.TaskView{
		ID:     v.ID,
		Status: v.Status,
		Txid:   v.Txid,
	}
	if v.States != nil {
		res.States = make([]*artworksviews.TaskStateView, len(v.States))
		for i, val := range v.States {
			res.States[i] = unmarshalTaskStateResponseToArtworksviewsTaskStateView(val)
		}
	}
	res.Ticket = unmarshalArtworkTicketResponseToArtworksviewsArtworkTicketView(v.Ticket)

	return res
}

// unmarshalTaskStateResponseToArtworksviewsTaskStateView builds a value of
// type *artworksviews.TaskStateView from a value of type *TaskStateResponse.
func unmarshalTaskStateResponseToArtworksviewsTaskStateView(v *TaskStateResponse) *artworksviews.TaskStateView {
	if v == nil {
		return nil
	}
	res := &artworksviews.TaskStateView{
		Date:   v.Date,
		Status: v.Status,
	}

	return res
}

// unmarshalArtworkTicketResponseToArtworksviewsArtworkTicketView builds a
// value of type *artworksviews.ArtworkTicketView from a value of type
// *ArtworkTicketResponse.
func unmarshalArtworkTicketResponseToArtworksviewsArtworkTicketView(v *ArtworkTicketResponse) *artworksviews.ArtworkTicketView {
	res := &artworksviews.ArtworkTicketView{
		Name:                     v.Name,
		Description:              v.Description,
		Keywords:                 v.Keywords,
		SeriesName:               v.SeriesName,
		IssuedCopies:             v.IssuedCopies,
		YoutubeURL:               v.YoutubeURL,
		ArtistPastelID:           v.ArtistPastelID,
		ArtistPastelIDPassphrase: v.ArtistPastelIDPassphrase,
		ArtistName:               v.ArtistName,
		ArtistWebsiteURL:         v.ArtistWebsiteURL,
		SpendableAddress:         v.SpendableAddress,
		MaximumFee:               v.MaximumFee,
		Royalty:                  v.Royalty,
		Green:                    v.Green,
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = unmarshalThumbnailcoordinateResponseToArtworksviewsThumbnailcoordinateView(v.ThumbnailCoordinate)
	}

	return res
}

// unmarshalThumbnailcoordinateResponseToArtworksviewsThumbnailcoordinateView
// builds a value of type *artworksviews.ThumbnailcoordinateView from a value
// of type *ThumbnailcoordinateResponse.
func unmarshalThumbnailcoordinateResponseToArtworksviewsThumbnailcoordinateView(v *ThumbnailcoordinateResponse) *artworksviews.ThumbnailcoordinateView {
	if v == nil {
		return nil
	}
	res := &artworksviews.ThumbnailcoordinateView{
		TopLeftX:     v.TopLeftX,
		TopLeftY:     v.TopLeftY,
		BottomRightX: v.BottomRightX,
		BottomRightY: v.BottomRightY,
	}

	return res
}

// unmarshalArtworkSummaryResponseBodyToArtworksArtworkSummary builds a value
// of type *artworks.ArtworkSummary from a value of type
// *ArtworkSummaryResponseBody.
func unmarshalArtworkSummaryResponseBodyToArtworksArtworkSummary(v *ArtworkSummaryResponseBody) *artworks.ArtworkSummary {
	res := &artworks.ArtworkSummary{
		Thumbnail1:       v.Thumbnail1,
		Thumbnail2:       v.Thumbnail2,
		Txid:             *v.Txid,
		Title:            *v.Title,
		Description:      *v.Description,
		Keywords:         v.Keywords,
		SeriesName:       v.SeriesName,
		Copies:           *v.Copies,
		YoutubeURL:       v.YoutubeURL,
		ArtistPastelID:   *v.ArtistPastelID,
		ArtistName:       *v.ArtistName,
		ArtistWebsiteURL: v.ArtistWebsiteURL,
	}

	return res
}

// unmarshalFuzzyMatchResponseBodyToArtworksFuzzyMatch builds a value of type
// *artworks.FuzzyMatch from a value of type *FuzzyMatchResponseBody.
func unmarshalFuzzyMatchResponseBodyToArtworksFuzzyMatch(v *FuzzyMatchResponseBody) *artworks.FuzzyMatch {
	res := &artworks.FuzzyMatch{
		Str:       v.Str,
		FieldType: v.FieldType,
		Score:     v.Score,
	}
	if v.MatchedIndexes != nil {
		res.MatchedIndexes = make([]int, len(v.MatchedIndexes))
		for i, val := range v.MatchedIndexes {
			res.MatchedIndexes[i] = val
		}
	}

	return res
}
