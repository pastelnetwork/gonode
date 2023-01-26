// Code generated by goa v3.7.6, DO NOT EDIT.
//
// sense HTTP client encoders and decoders
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

	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	senseviews "github.com/pastelnetwork/gonode/walletnode/api/gen/sense/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// BuildUploadImageRequest instantiates a HTTP request object with method and
// path set to call the "sense" service "uploadImage" endpoint
func (c *Client) BuildUploadImageRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: UploadImageSensePath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("sense", "uploadImage", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeUploadImageRequest returns an encoder for requests sent to the sense
// uploadImage server.
func EncodeUploadImageRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*sense.UploadImagePayload)
		if !ok {
			return goahttp.ErrInvalidType("sense", "uploadImage", "*sense.UploadImagePayload", v)
		}
		if err := encoder(req).Encode(p); err != nil {
			return goahttp.ErrEncodingError("sense", "uploadImage", err)
		}
		return nil
	}
}

// NewSenseUploadImageEncoder returns an encoder to encode the multipart
// request for the "sense" service "uploadImage" endpoint.
func NewSenseUploadImageEncoder(encoderFn SenseUploadImageEncoderFunc) func(r *http.Request) goahttp.Encoder {
	return func(r *http.Request) goahttp.Encoder {
		body := &bytes.Buffer{}
		mw := multipart.NewWriter(body)
		return goahttp.EncodingFunc(func(v interface{}) error {
			p := v.(*sense.UploadImagePayload)
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
// sense uploadImage endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeUploadImageResponse may return the following errors:
//   - "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
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
				return nil, goahttp.ErrDecodingError("sense", "uploadImage", err)
			}
			p := NewUploadImageImageCreated(&body)
			view := "default"
			vres := &senseviews.Image{Projected: p, View: view}
			if err = senseviews.ValidateImage(vres); err != nil {
				return nil, goahttp.ErrValidationError("sense", "uploadImage", err)
			}
			res := sense.NewImage(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body UploadImageBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("sense", "uploadImage", err)
			}
			err = ValidateUploadImageBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "uploadImage", err)
			}
			return nil, NewUploadImageBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body UploadImageInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("sense", "uploadImage", err)
			}
			err = ValidateUploadImageInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "uploadImage", err)
			}
			return nil, NewUploadImageInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("sense", "uploadImage", resp.StatusCode, string(body))
		}
	}
}

// BuildStartProcessingRequest instantiates a HTTP request object with method
// and path set to call the "sense" service "startProcessing" endpoint
func (c *Client) BuildStartProcessingRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		imageID string
	)
	{
		p, ok := v.(*sense.StartProcessingPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("sense", "startProcessing", "*sense.StartProcessingPayload", v)
		}
		imageID = p.ImageID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: StartProcessingSensePath(imageID)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("sense", "startProcessing", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeStartProcessingRequest returns an encoder for requests sent to the
// sense startProcessing server.
func EncodeStartProcessingRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*sense.StartProcessingPayload)
		if !ok {
			return goahttp.ErrInvalidType("sense", "startProcessing", "*sense.StartProcessingPayload", v)
		}
		{
			head := p.AppPastelidPassphrase
			req.Header.Set("app_pastelid_passphrase", head)
		}
		body := NewStartProcessingRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("sense", "startProcessing", err)
		}
		return nil
	}
}

// DecodeStartProcessingResponse returns a decoder for responses returned by
// the sense startProcessing endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeStartProcessingResponse may return the following errors:
//   - "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
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
				return nil, goahttp.ErrDecodingError("sense", "startProcessing", err)
			}
			p := NewStartProcessingResultViewCreated(&body)
			view := "default"
			vres := &senseviews.StartProcessingResult{Projected: p, View: view}
			if err = senseviews.ValidateStartProcessingResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("sense", "startProcessing", err)
			}
			res := sense.NewStartProcessingResult(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body StartProcessingBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("sense", "startProcessing", err)
			}
			err = ValidateStartProcessingBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "startProcessing", err)
			}
			return nil, NewStartProcessingBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body StartProcessingInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("sense", "startProcessing", err)
			}
			err = ValidateStartProcessingInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "startProcessing", err)
			}
			return nil, NewStartProcessingInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("sense", "startProcessing", resp.StatusCode, string(body))
		}
	}
}

// BuildRegisterTaskStateRequest instantiates a HTTP request object with method
// and path set to call the "sense" service "registerTaskState" endpoint
func (c *Client) BuildRegisterTaskStateRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*sense.RegisterTaskStatePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("sense", "registerTaskState", "*sense.RegisterTaskStatePayload", v)
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
	u := &url.URL{Scheme: scheme, Host: c.host, Path: RegisterTaskStateSensePath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("sense", "registerTaskState", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeRegisterTaskStateResponse returns a decoder for responses returned by
// the sense registerTaskState endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeRegisterTaskStateResponse may return the following errors:
//   - "NotFound" (type *goa.ServiceError): http.StatusNotFound
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
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
				return nil, goahttp.ErrDecodingError("sense", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "registerTaskState", err)
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
				return nil, goahttp.ErrDecodingError("sense", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body RegisterTaskStateInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("sense", "registerTaskState", err)
			}
			err = ValidateRegisterTaskStateInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "registerTaskState", err)
			}
			return nil, NewRegisterTaskStateInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("sense", "registerTaskState", resp.StatusCode, string(body))
		}
	}
}

// BuildGetTaskHistoryRequest instantiates a HTTP request object with method
// and path set to call the "sense" service "getTaskHistory" endpoint
func (c *Client) BuildGetTaskHistoryRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		taskID string
	)
	{
		p, ok := v.(*sense.GetTaskHistoryPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("sense", "getTaskHistory", "*sense.GetTaskHistoryPayload", v)
		}
		taskID = p.TaskID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetTaskHistorySensePath(taskID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("sense", "getTaskHistory", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetTaskHistoryResponse returns a decoder for responses returned by the
// sense getTaskHistory endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeGetTaskHistoryResponse may return the following errors:
//   - "NotFound" (type *goa.ServiceError): http.StatusNotFound
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
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
				return nil, goahttp.ErrDecodingError("sense", "getTaskHistory", err)
			}
			for _, e := range body {
				if e != nil {
					if err2 := ValidateTaskHistoryResponse(e); err2 != nil {
						err = goa.MergeErrors(err, err2)
					}
				}
			}
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "getTaskHistory", err)
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
				return nil, goahttp.ErrDecodingError("sense", "getTaskHistory", err)
			}
			err = ValidateGetTaskHistoryNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "getTaskHistory", err)
			}
			return nil, NewGetTaskHistoryNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetTaskHistoryInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("sense", "getTaskHistory", err)
			}
			err = ValidateGetTaskHistoryInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "getTaskHistory", err)
			}
			return nil, NewGetTaskHistoryInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("sense", "getTaskHistory", resp.StatusCode, string(body))
		}
	}
}

// BuildDownloadRequest instantiates a HTTP request object with method and path
// set to call the "sense" service "download" endpoint
func (c *Client) BuildDownloadRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: DownloadSensePath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("sense", "download", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeDownloadRequest returns an encoder for requests sent to the sense
// download server.
func EncodeDownloadRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*sense.DownloadPayload)
		if !ok {
			return goahttp.ErrInvalidType("sense", "download", "*sense.DownloadPayload", v)
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

// DecodeDownloadResponse returns a decoder for responses returned by the sense
// download endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeDownloadResponse may return the following errors:
//   - "NotFound" (type *goa.ServiceError): http.StatusNotFound
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
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
				return nil, goahttp.ErrDecodingError("sense", "download", err)
			}
			err = ValidateDownloadResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "download", err)
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
				return nil, goahttp.ErrDecodingError("sense", "download", err)
			}
			err = ValidateDownloadNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "download", err)
			}
			return nil, NewDownloadNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body DownloadInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("sense", "download", err)
			}
			err = ValidateDownloadInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("sense", "download", err)
			}
			return nil, NewDownloadInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("sense", "download", resp.StatusCode, string(body))
		}
	}
}

// unmarshalTaskHistoryResponseToSenseTaskHistory builds a value of type
// *sense.TaskHistory from a value of type *TaskHistoryResponse.
func unmarshalTaskHistoryResponseToSenseTaskHistory(v *TaskHistoryResponse) *sense.TaskHistory {
	res := &sense.TaskHistory{
		Timestamp: v.Timestamp,
		Status:    *v.Status,
		Message:   v.Message,
	}
	if v.Details != nil {
		res.Details = unmarshalDetailsResponseToSenseDetails(v.Details)
	}

	return res
}

// unmarshalDetailsResponseToSenseDetails builds a value of type *sense.Details
// from a value of type *DetailsResponse.
func unmarshalDetailsResponseToSenseDetails(v *DetailsResponse) *sense.Details {
	if v == nil {
		return nil
	}
	res := &sense.Details{
		Message: v.Message,
	}
	if v.Fields != nil {
		res.Fields = make(map[string]interface{}, len(v.Fields))
		for key, val := range v.Fields {
			tk := key
			tv := val
			res.Fields[tk] = tv
		}
	}

	return res
}
