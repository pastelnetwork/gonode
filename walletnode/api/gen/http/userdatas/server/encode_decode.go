// Code generated by goa v3.5.3, DO NOT EDIT.
//
// userdatas HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"net/http"
	"unicode/utf8"

	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeCreateUserdataResponse returns an encoder for responses returned by
// the userdatas createUserdata endpoint.
func EncodeCreateUserdataResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*userdatas.UserdataProcessResult)
		enc := encoder(ctx, w)
		body := NewCreateUserdataResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeCreateUserdataRequest returns a decoder for requests sent to the
// userdatas createUserdata endpoint.
func DecodeCreateUserdataRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var payload *userdatas.CreateUserdataPayload
		if err := decoder(r).Decode(&payload); err != nil {
			return nil, goa.DecodePayloadError(err.Error())
		}

		return payload, nil
	}
}

// NewUserdatasCreateUserdataDecoder returns a decoder to decode the multipart
// request for the "userdatas" service "createUserdata" endpoint.
func NewUserdatasCreateUserdataDecoder(mux goahttp.Muxer, userdatasCreateUserdataDecoderFn UserdatasCreateUserdataDecoderFunc) func(r *http.Request) goahttp.Decoder {
	return func(r *http.Request) goahttp.Decoder {
		return goahttp.EncodingFunc(func(v interface{}) error {
			mr, merr := r.MultipartReader()
			if merr != nil {
				return merr
			}
			p := v.(**userdatas.CreateUserdataPayload)
			if err := userdatasCreateUserdataDecoderFn(mr, p); err != nil {
				return err
			}
			return nil
		})
	}
}

// EncodeCreateUserdataError returns an encoder for errors returned by the
// createUserdata userdatas endpoint.
func EncodeCreateUserdataError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "BadRequest":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewCreateUserdataBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "InternalServerError":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewCreateUserdataInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeUpdateUserdataResponse returns an encoder for responses returned by
// the userdatas updateUserdata endpoint.
func EncodeUpdateUserdataResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*userdatas.UserdataProcessResult)
		enc := encoder(ctx, w)
		body := NewUpdateUserdataResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeUpdateUserdataRequest returns a decoder for requests sent to the
// userdatas updateUserdata endpoint.
func DecodeUpdateUserdataRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var payload *userdatas.UpdateUserdataPayload
		if err := decoder(r).Decode(&payload); err != nil {
			return nil, goa.DecodePayloadError(err.Error())
		}

		return payload, nil
	}
}

// NewUserdatasUpdateUserdataDecoder returns a decoder to decode the multipart
// request for the "userdatas" service "updateUserdata" endpoint.
func NewUserdatasUpdateUserdataDecoder(mux goahttp.Muxer, userdatasUpdateUserdataDecoderFn UserdatasUpdateUserdataDecoderFunc) func(r *http.Request) goahttp.Decoder {
	return func(r *http.Request) goahttp.Decoder {
		return goahttp.EncodingFunc(func(v interface{}) error {
			mr, merr := r.MultipartReader()
			if merr != nil {
				return merr
			}
			p := v.(**userdatas.UpdateUserdataPayload)
			if err := userdatasUpdateUserdataDecoderFn(mr, p); err != nil {
				return err
			}
			return nil
		})
	}
}

// EncodeUpdateUserdataError returns an encoder for errors returned by the
// updateUserdata userdatas endpoint.
func EncodeUpdateUserdataError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "BadRequest":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUpdateUserdataBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "InternalServerError":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUpdateUserdataInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetUserdataResponse returns an encoder for responses returned by the
// userdatas getUserdata endpoint.
func EncodeGetUserdataResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*userdatas.UserSpecifiedData)
		enc := encoder(ctx, w)
		body := NewGetUserdataResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetUserdataRequest returns a decoder for requests sent to the
// userdatas getUserdata endpoint.
func DecodeGetUserdataRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			pastelid string
			err      error

			params = mux.Vars(r)
		)
		pastelid = params["pastelid"]
		err = goa.MergeErrors(err, goa.ValidatePattern("pastelid", pastelid, "^[a-zA-Z0-9]+$"))
		if utf8.RuneCountInString(pastelid) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("pastelid", pastelid, utf8.RuneCountInString(pastelid), 86, true))
		}
		if utf8.RuneCountInString(pastelid) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("pastelid", pastelid, utf8.RuneCountInString(pastelid), 86, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetUserdataPayload(pastelid)

		return payload, nil
	}
}

// EncodeGetUserdataError returns an encoder for errors returned by the
// getUserdata userdatas endpoint.
func EncodeGetUserdataError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "BadRequest":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetUserdataBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "NotFound":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetUserdataNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetUserdataInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// unmarshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload
// builds a value of type *userdatas.UserImageUploadPayload from a value of
// type *UserImageUploadPayloadRequestBody.
func unmarshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(v *UserImageUploadPayloadRequestBody) *userdatas.UserImageUploadPayload {
	if v == nil {
		return nil
	}
	res := &userdatas.UserImageUploadPayload{
		Content:  v.Content,
		Filename: v.Filename,
	}

	return res
}

// marshalUserdatasUserImageUploadPayloadToUserImageUploadPayloadResponseBody
// builds a value of type *UserImageUploadPayloadResponseBody from a value of
// type *userdatas.UserImageUploadPayload.
func marshalUserdatasUserImageUploadPayloadToUserImageUploadPayloadResponseBody(v *userdatas.UserImageUploadPayload) *UserImageUploadPayloadResponseBody {
	if v == nil {
		return nil
	}
	res := &UserImageUploadPayloadResponseBody{
		Content:  v.Content,
		Filename: v.Filename,
	}

	return res
}
