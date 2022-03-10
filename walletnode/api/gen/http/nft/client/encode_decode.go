// Code generated by goa v3.4.3, DO NOT EDIT.
//
// nft HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	nft "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	nftviews "github.com/pastelnetwork/gonode/walletnode/api/gen/nft/views"
	goahttp "goa.design/goa/v3/http"
)

// BuildRegisterRequest instantiates a HTTP request object with method and path
// set to call the "nft" service "register" endpoint
func (c *Client) BuildRegisterRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RegisterNftPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "register", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeRegisterRequest returns an encoder for requests sent to the nft
// register server.
func EncodeRegisterRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*nft.RegisterPayload)
		if !ok {
			return goahttp.ErrInvalidType("nft", "register", "*nft.RegisterPayload", v)
		}
		body := NewRegisterRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("nft", "register", err)
		}
		return nil
	}
}

// DecodeRegisterResponse returns a decoder for responses returned by the nft
// register endpoint. restoreBody controls whether the response body should be
// restored after having been read.
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
				return nil, goahttp.ErrDecodingError("nft", "register", err)
			}
			p := NewRegisterResultViewCreated(&body)
			view := "default"
			vres := &nftviews.RegisterResult{Projected: p, View: view}
			if err = nftviews.ValidateRegisterResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("nft", "register", err)
			}
			res := nft.NewRegisterResult(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body RegisterBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "register", err)
			}
			err = ValidateRegisterBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "register", err)
			}
			return nil, NewRegisterBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "register", err)
			}
			err = ValidateRegisterInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "register", err)
			}
			return nil, NewRegisterInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "register", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTaskStateRequest instantiates a HTTP request object with method
// and path set to call the "nft" service "registerTaskState" endpoint
func (c *Client) BuildRegisterTaskStateRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*nft.RegisterTaskStatePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("nft", "registerTaskState", "*nft.RegisterTaskStatePayload", v)
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
	u := &url.URL{Scheme: scheme, Host: c.host, Path: RegisterTaskStateNftPath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "registerTaskState", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTaskStateResponse returns a decoder for responses returned by
// the nft registerTaskState endpoint. restoreBody controls whether the
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
				return nil, goahttp.ErrDecodingError("nft", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTaskState", err)
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
				return nil, goahttp.ErrDecodingError("nft", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterTaskStateInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "registerTaskState", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTaskRequest instantiates a HTTP request object with method and
// path set to call the "nft" service "registerTask" endpoint
func (c *Client) BuildRegisterTaskRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*nft.RegisterTaskPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("nft", "registerTask", "*nft.RegisterTaskPayload", v)
		}
		taskID = p.TaskID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RegisterTaskNftPath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "registerTask", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTaskResponse returns a decoder for responses returned by the
// nft registerTask endpoint. restoreBody controls whether the response body
// should be restored after having been read.
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
				return nil, goahttp.ErrDecodingError("nft", "registerTask", err)
			}
			p := NewRegisterTaskTaskOK(&body)
			view := "default"
			vres := &nftviews.Task{Projected: p, View: view}
			if err = nftviews.ValidateTask(vres); err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTask", err)
			}
			res := nft.NewTask(vres)
			return res, nil
		case http.StatusNotFound:
			var (
				body RegisterTaskNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "registerTask", err)
			}
			err = ValidateRegisterTaskNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTask", err)
			}
			return nil, NewRegisterTaskNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterTaskInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "registerTask", err)
			}
			err = ValidateRegisterTaskInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTask", err)
			}
			return nil, NewRegisterTaskInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "registerTask", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTasksRequest instantiates a HTTP request object with method and
// path set to call the "nft" service "registerTasks" endpoint
func (c *Client) BuildRegisterTasksRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: RegisterTasksNftPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "registerTasks", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTasksResponse returns a decoder for responses returned by the
// nft registerTasks endpoint. restoreBody controls whether the response body
// should be restored after having been read.
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
				return nil, goahttp.ErrDecodingError("nft", "registerTasks", err)
			}
			p := NewRegisterTasksTaskCollectionOK(body)
			view := "tiny"
			vres := nftviews.TaskCollection{Projected: p, View: view}
			if err = nftviews.ValidateTaskCollection(vres); err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTasks", err)
			}
			res := nft.NewTaskCollection(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body RegisterTasksInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "registerTasks", err)
			}
			err = ValidateRegisterTasksInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "registerTasks", err)
			}
			return nil, NewRegisterTasksInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "registerTasks", resp.StatusCode, string(body))
		}
	}
}

// BuildUploadImageRequest instantiates a HTTP request object with method and
// path set to call the "nft" service "uploadImage" endpoint
func (c *Client) BuildUploadImageRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: UploadImageNftPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "uploadImage", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeUploadImageRequest returns an encoder for requests sent to the nft
// uploadImage server.
func EncodeUploadImageRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*nft.UploadImagePayload)
		if !ok {
			return goahttp.ErrInvalidType("nft", "uploadImage", "*nft.UploadImagePayload", v)
		}
		if err := encoder(req).Encode(p); err != nil {
			return goahttp.ErrEncodingError("nft", "uploadImage", err)
		}
		return nil
	}
}

// NewNftUploadImageEncoder returns an encoder to encode the multipart request
// for the "nft" service "uploadImage" endpoint.
func NewNftUploadImageEncoder(encoderFn NftUploadImageEncoderFunc) func(r *http.Request) goahttp.Encoder {
	return func(r *http.Request) goahttp.Encoder {
		body := &bytes.Buffer{}
		mw := multipart.NewWriter(body)
		return goahttp.EncodingFunc(func(v interface{}) error {
			p := v.(*nft.UploadImagePayload)
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
// nft uploadImage endpoint. restoreBody controls whether the response body
// should be restored after having been read.
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
				return nil, goahttp.ErrDecodingError("nft", "uploadImage", err)
			}
			p := NewUploadImageImageCreated(&body)
			view := "default"
			vres := &nftviews.Image{Projected: p, View: view}
			if err = nftviews.ValidateImage(vres); err != nil {
				return nil, goahttp.ErrValidationError("nft", "uploadImage", err)
			}
			res := nft.NewImage(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body UploadImageBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "uploadImage", err)
			}
			err = ValidateUploadImageBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "uploadImage", err)
			}
			return nil, NewUploadImageBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body UploadImageInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "uploadImage", err)
			}
			err = ValidateUploadImageInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "uploadImage", err)
			}
			return nil, NewUploadImageInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "uploadImage", resp.StatusCode, string(body))
		}
	}
}

// BuildNftSearchRequest instantiates a HTTP request object with method and
// path set to call the "nft" service "nftSearch" endpoint
func (c *Client) BuildNftSearchRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	scheme := c.scheme
	switch c.scheme {
	case "http":
		scheme = "ws"
	case "https":
		scheme = "wss"
	}
	u := &url.URL{Scheme: scheme, Host: c.host, Path: NftSearchNftPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "nftSearch", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeNftSearchRequest returns an encoder for requests sent to the nft
// nftSearch server.
func EncodeNftSearchRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*nft.NftSearchPayload)
		if !ok {
			return goahttp.ErrInvalidType("nft", "nftSearch", "*nft.NftSearchPayload", v)
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
		values.Add("creator_name", fmt.Sprintf("%v", p.CreatorName))
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
		if p.IsLikelyDupe != nil {
			values.Add("is_likely_dupe", fmt.Sprintf("%v", *p.IsLikelyDupe))
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
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeNftSearchResponse returns a decoder for responses returned by the nft
// nftSearch endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeNftSearchResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeNftSearchResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body NftSearchResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "nftSearch", err)
			}
			err = ValidateNftSearchResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "nftSearch", err)
			}
			res := NewNftSearchResultOK(&body)
			return res, nil
		case http.StatusBadRequest:
			var (
				body NftSearchBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "nftSearch", err)
			}
			err = ValidateNftSearchBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "nftSearch", err)
			}
			return nil, NewNftSearchBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body NftSearchInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "nftSearch", err)
			}
			err = ValidateNftSearchInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "nftSearch", err)
			}
			return nil, NewNftSearchInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "nftSearch", resp.StatusCode, string(body))
		}
	}
}

// BuildNftGetRequest instantiates a HTTP request object with method and path
// set to call the "nft" service "nftGet" endpoint
func (c *Client) BuildNftGetRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		txid string
	)
	{
		p, ok := v.(*nft.NftGetPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("nft", "nftGet", "*nft.NftGetPayload", v)
		}
		txid = p.Txid
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: NftGetNftPath(txid)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "nftGet", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeNftGetRequest returns an encoder for requests sent to the nft nftGet
// server.
func EncodeNftGetRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*nft.NftGetPayload)
		if !ok {
			return goahttp.ErrInvalidType("nft", "nftGet", "*nft.NftGetPayload", v)
		}
		body := NewNftGetRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("nft", "nftGet", err)
		}
		return nil
	}
}

// DecodeNftGetResponse returns a decoder for responses returned by the nft
// nftGet endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeNftGetResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "NotFound" (type *goa.ServiceError): http.StatusNotFound
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeNftGetResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body NftGetResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "nftGet", err)
			}
			err = ValidateNftGetResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "nftGet", err)
			}
			res := NewNftGetNftDetailOK(&body)
			return res, nil
		case http.StatusBadRequest:
			var (
				body NftGetBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "nftGet", err)
			}
			err = ValidateNftGetBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "nftGet", err)
			}
			return nil, NewNftGetBadRequest(&body)
		case http.StatusNotFound:
			var (
				body NftGetNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "nftGet", err)
			}
			err = ValidateNftGetNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "nftGet", err)
			}
			return nil, NewNftGetNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body NftGetInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "nftGet", err)
			}
			err = ValidateNftGetInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "nftGet", err)
			}
			return nil, NewNftGetInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "nftGet", resp.StatusCode, string(body))
		}
	}
}

// BuildDownloadRequest instantiates a HTTP request object with method and path
// set to call the "nft" service "download" endpoint
func (c *Client) BuildDownloadRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: DownloadNftPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("nft", "download", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeDownloadRequest returns an encoder for requests sent to the nft
// download server.
func EncodeDownloadRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*nft.NftDownloadPayload)
		if !ok {
			return goahttp.ErrInvalidType("nft", "download", "*nft.NftDownloadPayload", v)
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

// DecodeDownloadResponse returns a decoder for responses returned by the nft
// download endpoint. restoreBody controls whether the response body should be
// restored after having been read.
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
				return nil, goahttp.ErrDecodingError("nft", "download", err)
			}
			err = ValidateDownloadResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "download", err)
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
				return nil, goahttp.ErrDecodingError("nft", "download", err)
			}
			err = ValidateDownloadNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "download", err)
			}
			return nil, NewDownloadNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body DownloadInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("nft", "download", err)
			}
			err = ValidateDownloadInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("nft", "download", err)
			}
			return nil, NewDownloadInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("nft", "download", resp.StatusCode, string(body))
		}
	}
}

// marshalNftThumbnailcoordinateToThumbnailcoordinateRequestBody builds a value
// of type *ThumbnailcoordinateRequestBody from a value of type
// *nft.Thumbnailcoordinate.
func marshalNftThumbnailcoordinateToThumbnailcoordinateRequestBody(v *nft.Thumbnailcoordinate) *ThumbnailcoordinateRequestBody {
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

// marshalThumbnailcoordinateRequestBodyToNftThumbnailcoordinate builds a value
// of type *nft.Thumbnailcoordinate from a value of type
// *ThumbnailcoordinateRequestBody.
func marshalThumbnailcoordinateRequestBodyToNftThumbnailcoordinate(v *ThumbnailcoordinateRequestBody) *nft.Thumbnailcoordinate {
	if v == nil {
		return nil
	}
	res := &nft.Thumbnailcoordinate{
		TopLeftX:     v.TopLeftX,
		TopLeftY:     v.TopLeftY,
		BottomRightX: v.BottomRightX,
		BottomRightY: v.BottomRightY,
	}

	return res
}

// unmarshalTaskStateResponseBodyToNftviewsTaskStateView builds a value of type
// *nftviews.TaskStateView from a value of type *TaskStateResponseBody.
func unmarshalTaskStateResponseBodyToNftviewsTaskStateView(v *TaskStateResponseBody) *nftviews.TaskStateView {
	if v == nil {
		return nil
	}
	res := &nftviews.TaskStateView{
		Date:   v.Date,
		Status: v.Status,
	}

	return res
}

// unmarshalNftRegisterPayloadResponseBodyToNftviewsNftRegisterPayloadView
// builds a value of type *nftviews.NftRegisterPayloadView from a value of type
// *NftRegisterPayloadResponseBody.
func unmarshalNftRegisterPayloadResponseBodyToNftviewsNftRegisterPayloadView(v *NftRegisterPayloadResponseBody) *nftviews.NftRegisterPayloadView {
	res := &nftviews.NftRegisterPayloadView{
		Name:                      v.Name,
		Description:               v.Description,
		Keywords:                  v.Keywords,
		SeriesName:                v.SeriesName,
		IssuedCopies:              v.IssuedCopies,
		YoutubeURL:                v.YoutubeURL,
		CreatorPastelID:           v.CreatorPastelID,
		CreatorPastelIDPassphrase: v.CreatorPastelIDPassphrase,
		CreatorName:               v.CreatorName,
		CreatorWebsiteURL:         v.CreatorWebsiteURL,
		SpendableAddress:          v.SpendableAddress,
		MaximumFee:                v.MaximumFee,
		Royalty:                   v.Royalty,
		Green:                     v.Green,
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = unmarshalThumbnailcoordinateResponseBodyToNftviewsThumbnailcoordinateView(v.ThumbnailCoordinate)
	}

	return res
}

// unmarshalThumbnailcoordinateResponseBodyToNftviewsThumbnailcoordinateView
// builds a value of type *nftviews.ThumbnailcoordinateView from a value of
// type *ThumbnailcoordinateResponseBody.
func unmarshalThumbnailcoordinateResponseBodyToNftviewsThumbnailcoordinateView(v *ThumbnailcoordinateResponseBody) *nftviews.ThumbnailcoordinateView {
	if v == nil {
		return nil
	}
	res := &nftviews.ThumbnailcoordinateView{
		TopLeftX:     v.TopLeftX,
		TopLeftY:     v.TopLeftY,
		BottomRightX: v.BottomRightX,
		BottomRightY: v.BottomRightY,
	}

	return res
}

// unmarshalTaskResponseToNftviewsTaskView builds a value of type
// *nftviews.TaskView from a value of type *TaskResponse.
func unmarshalTaskResponseToNftviewsTaskView(v *TaskResponse) *nftviews.TaskView {
	res := &nftviews.TaskView{
		ID:     v.ID,
		Status: v.Status,
		Txid:   v.Txid,
	}
	if v.States != nil {
		res.States = make([]*nftviews.TaskStateView, len(v.States))
		for i, val := range v.States {
			res.States[i] = unmarshalTaskStateResponseToNftviewsTaskStateView(val)
		}
	}
	res.Ticket = unmarshalNftRegisterPayloadResponseToNftviewsNftRegisterPayloadView(v.Ticket)

	return res
}

// unmarshalTaskStateResponseToNftviewsTaskStateView builds a value of type
// *nftviews.TaskStateView from a value of type *TaskStateResponse.
func unmarshalTaskStateResponseToNftviewsTaskStateView(v *TaskStateResponse) *nftviews.TaskStateView {
	if v == nil {
		return nil
	}
	res := &nftviews.TaskStateView{
		Date:   v.Date,
		Status: v.Status,
	}

	return res
}

// unmarshalNftRegisterPayloadResponseToNftviewsNftRegisterPayloadView builds a
// value of type *nftviews.NftRegisterPayloadView from a value of type
// *NftRegisterPayloadResponse.
func unmarshalNftRegisterPayloadResponseToNftviewsNftRegisterPayloadView(v *NftRegisterPayloadResponse) *nftviews.NftRegisterPayloadView {
	res := &nftviews.NftRegisterPayloadView{
		Name:                      v.Name,
		Description:               v.Description,
		Keywords:                  v.Keywords,
		SeriesName:                v.SeriesName,
		IssuedCopies:              v.IssuedCopies,
		YoutubeURL:                v.YoutubeURL,
		CreatorPastelID:           v.CreatorPastelID,
		CreatorPastelIDPassphrase: v.CreatorPastelIDPassphrase,
		CreatorName:               v.CreatorName,
		CreatorWebsiteURL:         v.CreatorWebsiteURL,
		SpendableAddress:          v.SpendableAddress,
		MaximumFee:                v.MaximumFee,
		Royalty:                   v.Royalty,
		Green:                     v.Green,
	}
	if v.ThumbnailCoordinate != nil {
		res.ThumbnailCoordinate = unmarshalThumbnailcoordinateResponseToNftviewsThumbnailcoordinateView(v.ThumbnailCoordinate)
	}

	return res
}

// unmarshalThumbnailcoordinateResponseToNftviewsThumbnailcoordinateView builds
// a value of type *nftviews.ThumbnailcoordinateView from a value of type
// *ThumbnailcoordinateResponse.
func unmarshalThumbnailcoordinateResponseToNftviewsThumbnailcoordinateView(v *ThumbnailcoordinateResponse) *nftviews.ThumbnailcoordinateView {
	if v == nil {
		return nil
	}
	res := &nftviews.ThumbnailcoordinateView{
		TopLeftX:     v.TopLeftX,
		TopLeftY:     v.TopLeftY,
		BottomRightX: v.BottomRightX,
		BottomRightY: v.BottomRightY,
	}

	return res
}

// unmarshalNftSummaryResponseBodyToNftNftSummary builds a value of type
// *nft.NftSummary from a value of type *NftSummaryResponseBody.
func unmarshalNftSummaryResponseBodyToNftNftSummary(v *NftSummaryResponseBody) *nft.NftSummary {
	res := &nft.NftSummary{
		Thumbnail1:        v.Thumbnail1,
		Thumbnail2:        v.Thumbnail2,
		Txid:              *v.Txid,
		Title:             *v.Title,
		Description:       *v.Description,
		Keywords:          v.Keywords,
		SeriesName:        v.SeriesName,
		Copies:            *v.Copies,
		YoutubeURL:        v.YoutubeURL,
		CreatorPastelID:   *v.CreatorPastelID,
		CreatorName:       *v.CreatorName,
		CreatorWebsiteURL: v.CreatorWebsiteURL,
		NsfwScore:         v.NsfwScore,
		RarenessScore:     v.RarenessScore,
		IsLikelyDupe:      v.IsLikelyDupe,
	}

	return res
}

// unmarshalFuzzyMatchResponseBodyToNftFuzzyMatch builds a value of type
// *nft.FuzzyMatch from a value of type *FuzzyMatchResponseBody.
func unmarshalFuzzyMatchResponseBodyToNftFuzzyMatch(v *FuzzyMatchResponseBody) *nft.FuzzyMatch {
	res := &nft.FuzzyMatch{
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
