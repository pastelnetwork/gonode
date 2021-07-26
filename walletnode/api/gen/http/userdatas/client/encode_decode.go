// Code generated by goa v3.3.1, DO NOT EDIT.
//
// userdatas HTTP client encoders and decoders
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

	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	goahttp "goa.design/goa/v3/http"
)

// BuildProcessUserdataRequest instantiates a HTTP request object with method
// and path set to call the "userdatas" service "processUserdata" endpoint
func (c *Client) BuildProcessUserdataRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ProcessUserdataUserdatasPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("userdatas", "processUserdata", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeProcessUserdataRequest returns an encoder for requests sent to the
// userdatas processUserdata server.
func EncodeProcessUserdataRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*userdatas.ProcessUserdataPayload)
		if !ok {
			return goahttp.ErrInvalidType("userdatas", "processUserdata", "*userdatas.ProcessUserdataPayload", v)
		}
		if err := encoder(req).Encode(p); err != nil {
			return goahttp.ErrEncodingError("userdatas", "processUserdata", err)
		}
		return nil
	}
}

// NewUserdatasProcessUserdataEncoder returns an encoder to encode the
// multipart request for the "userdatas" service "processUserdata" endpoint.
func NewUserdatasProcessUserdataEncoder(encoderFn UserdatasProcessUserdataEncoderFunc) func(r *http.Request) goahttp.Encoder {
	return func(r *http.Request) goahttp.Encoder {
		body := &bytes.Buffer{}
		mw := multipart.NewWriter(body)
		return goahttp.EncodingFunc(func(v interface{}) error {
			p := v.(*userdatas.ProcessUserdataPayload)
			if err := encoderFn(mw, p); err != nil {
				return err
			}
			r.Body = ioutil.NopCloser(body)
			r.Header.Set("Content-Type", mw.FormDataContentType())
			return mw.Close()
		})
	}
}

// DecodeProcessUserdataResponse returns a decoder for responses returned by
// the userdatas processUserdata endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeProcessUserdataResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeProcessUserdataResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body ProcessUserdataResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("userdatas", "processUserdata", err)
			}
			err = ValidateProcessUserdataResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("userdatas", "processUserdata", err)
			}
			res := NewProcessUserdataUserdataProcessResultCreated(&body)
			return res, nil
		case http.StatusBadRequest:
			var (
				body ProcessUserdataBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("userdatas", "processUserdata", err)
			}
			err = ValidateProcessUserdataBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("userdatas", "processUserdata", err)
			}
			return nil, NewProcessUserdataBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body ProcessUserdataInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("userdatas", "processUserdata", err)
			}
			err = ValidateProcessUserdataInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("userdatas", "processUserdata", err)
			}
			return nil, NewProcessUserdataInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("userdatas", "processUserdata", resp.StatusCode, string(body))
		}
	}
}

// BuildUserdataGetRequest instantiates a HTTP request object with method and
// path set to call the "userdatas" service "userdataGet" endpoint
func (c *Client) BuildUserdataGetRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		pastelid string
	)
	{
		p, ok := v.(*userdatas.UserdataGetPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("userdatas", "userdataGet", "*userdatas.UserdataGetPayload", v)
		}
		pastelid = p.Pastelid
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: UserdataGetUserdatasPath(pastelid)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("userdatas", "userdataGet", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeUserdataGetResponse returns a decoder for responses returned by the
// userdatas userdataGet endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeUserdataGetResponse may return the following errors:
//	- "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//	- "NotFound" (type *goa.ServiceError): http.StatusNotFound
//	- "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeUserdataGetResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body UserdataGetResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("userdatas", "userdataGet", err)
			}
			err = ValidateUserdataGetResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("userdatas", "userdataGet", err)
			}
			res := NewUserdataGetUserSpecifiedDataOK(&body)
			return res, nil
		case http.StatusBadRequest:
			var (
				body UserdataGetBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("userdatas", "userdataGet", err)
			}
			err = ValidateUserdataGetBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("userdatas", "userdataGet", err)
			}
			return nil, NewUserdataGetBadRequest(&body)
		case http.StatusNotFound:
			var (
				body UserdataGetNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("userdatas", "userdataGet", err)
			}
			err = ValidateUserdataGetNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("userdatas", "userdataGet", err)
			}
			return nil, NewUserdataGetNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body UserdataGetInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("userdatas", "userdataGet", err)
			}
			err = ValidateUserdataGetInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("userdatas", "userdataGet", err)
			}
			return nil, NewUserdataGetInternalServerError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("userdatas", "userdataGet", resp.StatusCode, string(body))
		}
	}
}

// marshalUserdatasUserImageUploadPayloadToUserImageUploadPayloadRequestBody
// builds a value of type *UserImageUploadPayloadRequestBody from a value of
// type *userdatas.UserImageUploadPayload.
func marshalUserdatasUserImageUploadPayloadToUserImageUploadPayloadRequestBody(v *userdatas.UserImageUploadPayload) *UserImageUploadPayloadRequestBody {
	if v == nil {
		return nil
	}
	res := &UserImageUploadPayloadRequestBody{
		Content:  v.Content,
		Filename: v.Filename,
	}

	return res
}

// marshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload
// builds a value of type *userdatas.UserImageUploadPayload from a value of
// type *UserImageUploadPayloadRequestBody.
func marshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(v *UserImageUploadPayloadRequestBody) *userdatas.UserImageUploadPayload {
	if v == nil {
		return nil
	}
	res := &userdatas.UserImageUploadPayload{
		Content:  v.Content,
		Filename: v.Filename,
	}

	return res
}

// unmarshalUserImageUploadPayloadResponseBodyToUserdatasUserImageUploadPayload
// builds a value of type *userdatas.UserImageUploadPayload from a value of
// type *UserImageUploadPayloadResponseBody.
func unmarshalUserImageUploadPayloadResponseBodyToUserdatasUserImageUploadPayload(v *UserImageUploadPayloadResponseBody) *userdatas.UserImageUploadPayload {
	if v == nil {
		return nil
	}
	res := &userdatas.UserImageUploadPayload{
		Content:  v.Content,
		Filename: v.Filename,
	}

	return res
}
