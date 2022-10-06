// Code generated by goa v3.7.6, DO NOT EDIT.
//
// cascade HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	cascade "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	cascadeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// BuildUploadAssetRequest instantiates a HTTP request object with method and
// path set to call the "cascade" service "uploadAsset" endpoint
func (c *Client) BuildUploadAssetRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: UploadAssetCascadePath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("cascade", "uploadAsset", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeUploadAssetRequest returns an encoder for requests sent to the cascade
// uploadAsset server.
func EncodeUploadAssetRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*cascade.UploadAssetPayload)
		if !ok {
			return goahttp.ErrInvalidType("cascade", "uploadAsset", "*cascade.UploadAssetPayload", v)
		}
		if err := encoder(req).Encode(p); err != nil {
			return goahttp.ErrEncodingError("cascade", "uploadAsset", err)
		}
		return nil
	}
}

// NewCascadeUploadAssetEncoder returns an encoder to encode the multipart
// request for the "cascade" service "uploadAsset" endpoint.
func NewCascadeUploadAssetEncoder(encoderFn CascadeUploadAssetEncoderFunc) func(r *http.Request) goahttp.Encoder {
	return func(r *http.Request) goahttp.Encoder {
		body := &bytes.Buffer{}
		mw := multipart.NewWriter(body)
		return goahttp.EncodingFunc(func(v interface{}) error {
			p := v.(*cascade.UploadAssetPayload)
			if err := encoderFn(mw, p); err != nil {
				return err
			}
			r.Body = ioutil.NopCloser(body)
			r.Header.Set("Content-Type", mw.FormDataContentType())
			return mw.Close()
		})
	}
}

// DecodeUploadAssetResponse returns a decoder for responses returned by the
// cascade uploadAsset endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeUploadAssetResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeUploadAssetResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body UploadAssetResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "uploadAsset", err)
			}
			p := NewUploadAssetAssetCreated(&body)
			view := "default"
			vres := &cascadeviews.Asset{Projected: p, View: view}
			if err = cascadeviews.ValidateAsset(vres); err != nil {
				return nil, goahttp.ErrValidationError("cascade", "uploadAsset", err)
			}
			res := cascade.NewAsset(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body UploadAssetBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "uploadAsset", err)
			}
			err = ValidateUploadAssetBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "uploadAsset", err)
			}
			return nil, NewUploadAssetBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body UploadAssetInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "uploadAsset", err)
			}
			err = ValidateUploadAssetInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "uploadAsset", err)
			}
			return nil, NewUploadAssetInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("cascade", "uploadAsset", resp.StatusCode, string(body))
		}
	}
}

// BuildStartProcessingRequest instantiates a HTTP request object with method
// and path set to call the "cascade" service "startProcessing" endpoint
func (c *Client) BuildStartProcessingRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		fileID string
	)
	{
		p, ok := v.(*cascade.StartProcessingPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("cascade", "startProcessing", "*cascade.StartProcessingPayload", v)
		}
		fileID = p.FileID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: StartProcessingCascadePath(fileID)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("cascade", "startProcessing", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeStartProcessingRequest returns an encoder for requests sent to the
// cascade startProcessing server.
func EncodeStartProcessingRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*cascade.StartProcessingPayload)
		if !ok {
			return goahttp.ErrInvalidType("cascade", "startProcessing", "*cascade.StartProcessingPayload", v)
		}
		{
			head := p.AppPastelidPassphrase
			req.Header.Set("app_pastelid_passphrase", head)
		}
		body := NewStartProcessingRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("cascade", "startProcessing", err)
		}
		return nil
	}
}

// DecodeStartProcessingResponse returns a decoder for responses returned by
// the cascade startProcessing endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeStartProcessingResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeStartProcessingResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body StartProcessingResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "startProcessing", err)
			}
			p := NewStartProcessingResultViewCreated(&body)
			view := "default"
			vres := &cascadeviews.StartProcessingResult{Projected: p, View: view}
			if err = cascadeviews.ValidateStartProcessingResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("cascade", "startProcessing", err)
			}
			res := cascade.NewStartProcessingResult(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body StartProcessingBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "startProcessing", err)
			}
			err = ValidateStartProcessingBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "startProcessing", err)
			}
			return nil, NewStartProcessingBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body StartProcessingInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "startProcessing", err)
			}
			err = ValidateStartProcessingInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "startProcessing", err)
			}
			return nil, NewStartProcessingInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("cascade", "startProcessing", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTaskStateRequest instantiates a HTTP request object with method
// and path set to call the "cascade" service "registerTaskState" endpoint
func (c *Client) BuildRegisterTaskStateRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*cascade.RegisterTaskStatePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("cascade", "registerTaskState", "*cascade.RegisterTaskStatePayload", v)
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
	u := &url.URL{Scheme: scheme, Host: c.host, Path: RegisterTaskStateCascadePath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("cascade", "registerTaskState", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTaskStateResponse returns a decoder for responses returned by
// the cascade registerTaskState endpoint. restoreBody controls whether the
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
				return nil, goahttp.ErrDecodingError("cascade", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "registerTaskState", err)
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
				return nil, goahttp.ErrDecodingError("cascade", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterTaskStateInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("cascade", "registerTaskState", resp.StatusCode, string(body))
		}
	}
}

// BuildGetTaskHistoryRequest instantiates a HTTP request object with method
// and path set to call the "cascade" service "getTaskHistory" endpoint
func (c *Client) BuildGetTaskHistoryRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*cascade.GetTaskHistoryPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("cascade", "getTaskHistory", "*cascade.GetTaskHistoryPayload", v)
		}
		taskID = p.TaskID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetTaskHistoryCascadePath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("cascade", "getTaskHistory", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetTaskHistoryResponse returns a decoder for responses returned by the
// cascade getTaskHistory endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeGetTaskHistoryResponse may return the following errors:
//	- "NotFound" (type *goa.ServiceError): http.StatusNotFound
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeGetTaskHistoryResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body GetTaskHistoryResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "getTaskHistory", err)
			}
			for _, e := range body {
				if e != nil {
					if err2 := ValidateTaskHistoryResponse(e); err2 != nil {
						err = goa.MergeErrors(err, err2)
					}
				}
			}
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "getTaskHistory", err)
			}
			res := NewGetTaskHistoryTaskHistoryOK(body)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetTaskHistoryNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "getTaskHistory", err)
			}
			err = ValidateGetTaskHistoryNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "getTaskHistory", err)
			}
			return nil, NewGetTaskHistoryNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetTaskHistoryInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "getTaskHistory", err)
			}
			err = ValidateGetTaskHistoryInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "getTaskHistory", err)
			}
			return nil, NewGetTaskHistoryInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("cascade", "getTaskHistory", resp.StatusCode, string(body))
		}
	}
}

// BuildDownloadRequest instantiates a HTTP request object with method and path
// set to call the "cascade" service "download" endpoint
func (c *Client) BuildDownloadRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: DownloadCascadePath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("cascade", "download", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeDownloadRequest returns an encoder for requests sent to the cascade
// download server.
func EncodeDownloadRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*cascade.DownloadPayload)
		if !ok {
			return goahttp.ErrInvalidType("cascade", "download", "*cascade.DownloadPayload", v)
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
// cascade download endpoint. restoreBody controls whether the response body
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
				return nil, goahttp.ErrDecodingError("cascade", "download", err)
			}
			err = ValidateDownloadResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "download", err)
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
				return nil, goahttp.ErrDecodingError("cascade", "download", err)
			}
			err = ValidateDownloadNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "download", err)
			}
			return nil, NewDownloadNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body DownloadInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "download", err)
			}
			err = ValidateDownloadInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "download", err)
			}
			return nil, NewDownloadInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("cascade", "download", resp.StatusCode, string(body))
		}
	}
}

// unmarshalTaskHistoryResponseToCascadeTaskHistory builds a value of type
// *cascade.TaskHistory from a value of type *TaskHistoryResponse.
func unmarshalTaskHistoryResponseToCascadeTaskHistory(v *TaskHistoryResponse) *cascade.TaskHistory {
	res := &cascade.TaskHistory{
		Timestamp: v.Timestamp,
		Status:    *v.Status,
		Message:   v.Message,
	}

	return res
}
