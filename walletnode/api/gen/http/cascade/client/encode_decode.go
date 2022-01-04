// Code generated by goa v3.5.3, DO NOT EDIT.
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
)

// BuildUploadImageRequest instantiates a HTTP request object with method and
// path set to call the "cascade" service "uploadImage" endpoint
func (c *Client) BuildUploadImageRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: UploadImageCascadePath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("cascade", "uploadImage", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeUploadImageRequest returns an encoder for requests sent to the cascade
// uploadImage server.
func EncodeUploadImageRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*cascade.UploadImagePayload)
		if !ok {
			return goahttp.ErrInvalidType("cascade", "uploadImage", "*cascade.UploadImagePayload", v)
		}
		if err := encoder(req).Encode(p); err != nil {
			return goahttp.ErrEncodingError("cascade", "uploadImage", err)
		}
		return nil
	}
}

// NewCascadeUploadImageEncoder returns an encoder to encode the multipart
// request for the "cascade" service "uploadImage" endpoint.
func NewCascadeUploadImageEncoder(encoderFn CascadeUploadImageEncoderFunc) func(r *http.Request) goahttp.Encoder {
	return func(r *http.Request) goahttp.Encoder {
		body := &bytes.Buffer{}
		mw := multipart.NewWriter(body)
		return goahttp.EncodingFunc(func(v interface{}) error {
			p := v.(*cascade.UploadImagePayload)
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
// cascade uploadImage endpoint. restoreBody controls whether the response body
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
				return nil, goahttp.ErrDecodingError("cascade", "uploadImage", err)
			}
			p := NewUploadImageImageCreated(&body)
			view := "default"
			vres := &cascadeviews.Image{Projected: p, View: view}
			if err = cascadeviews.ValidateImage(vres); err != nil {
				return nil, goahttp.ErrValidationError("cascade", "uploadImage", err)
			}
			res := cascade.NewImage(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body UploadImageBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "uploadImage", err)
			}
			err = ValidateUploadImageBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "uploadImage", err)
			}
			return nil, NewUploadImageBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body UploadImageInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "uploadImage", err)
			}
			err = ValidateUploadImageInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "uploadImage", err)
			}
			return nil, NewUploadImageInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("cascade", "uploadImage", resp.StatusCode, string(body))
		}
	}
}

// BuildActionDetailsRequest instantiates a HTTP request object with method and
// path set to call the "cascade" service "actionDetails" endpoint
func (c *Client) BuildActionDetailsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		imageID string
	)
	{
		p, ok := v.(*cascade.ActionDetailsPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("cascade", "actionDetails", "*cascade.ActionDetailsPayload", v)
		}
		imageID = p.ImageID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ActionDetailsCascadePath(imageID)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("cascade", "actionDetails", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeActionDetailsRequest returns an encoder for requests sent to the
// cascade actionDetails server.
func EncodeActionDetailsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*cascade.ActionDetailsPayload)
		if !ok {
			return goahttp.ErrInvalidType("cascade", "actionDetails", "*cascade.ActionDetailsPayload", v)
		}
		body := NewActionDetailsRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("cascade", "actionDetails", err)
		}
		return nil
	}
}

// DecodeActionDetailsResponse returns a decoder for responses returned by the
// cascade actionDetails endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeActionDetailsResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeActionDetailsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body ActionDetailsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "actionDetails", err)
			}
			p := NewActionDetailsActionDetailResultCreated(&body)
			view := "default"
			vres := &cascadeviews.ActionDetailResult{Projected: p, View: view}
			if err = cascadeviews.ValidateActionDetailResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("cascade", "actionDetails", err)
			}
			res := cascade.NewActionDetailResult(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body ActionDetailsBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "actionDetails", err)
			}
			err = ValidateActionDetailsBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "actionDetails", err)
			}
			return nil, NewActionDetailsBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body ActionDetailsInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("cascade", "actionDetails", err)
			}
			err = ValidateActionDetailsInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("cascade", "actionDetails", err)
			}
			return nil, NewActionDetailsInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("cascade", "actionDetails", resp.StatusCode, string(body))
		}
	}
}

// BuildStartProcessingRequest instantiates a HTTP request object with method
// and path set to call the "cascade" service "startProcessing" endpoint
func (c *Client) BuildStartProcessingRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		imageID string
	)
	{
		p, ok := v.(*cascade.StartProcessingPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("cascade", "startProcessing", "*cascade.StartProcessingPayload", v)
		}
		imageID = p.ImageID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: StartProcessingCascadePath(imageID)}
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
