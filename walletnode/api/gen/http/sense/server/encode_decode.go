// Code generated by goa v3.4.3, DO NOT EDIT.
//
// sense HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"io"
	"net/http"
	"unicode/utf8"

	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	senseviews "github.com/pastelnetwork/gonode/walletnode/api/gen/sense/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeUploadImageResponse returns an encoder for responses returned by the
// sense uploadImage endpoint.
func EncodeUploadImageResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*senseviews.Image)
		enc := encoder(ctx, w)
		body := NewUploadImageResponseBody(res.Projected)
		w.WriteHeader(http.StatusCreated)
		return enc.Encode(body)
	}
}

// DecodeUploadImageRequest returns a decoder for requests sent to the sense
// uploadImage endpoint.
func DecodeUploadImageRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var payload *sense.UploadImagePayload
		if err := decoder(r).Decode(&payload); err != nil {
			return nil, goa.DecodePayloadError(err.Error())
		}

		return payload, nil
	}
}

// NewSenseUploadImageDecoder returns a decoder to decode the multipart request
// for the "sense" service "uploadImage" endpoint.
func NewSenseUploadImageDecoder(mux goahttp.Muxer, senseUploadImageDecoderFn SenseUploadImageDecoderFunc) func(r *http.Request) goahttp.Decoder {
	return func(r *http.Request) goahttp.Decoder {
		return goahttp.EncodingFunc(func(v interface{}) error {
			mr, merr := r.MultipartReader()
			if merr != nil {
				return merr
			}
			p := v.(**sense.UploadImagePayload)
			if err := senseUploadImageDecoderFn(mr, p); err != nil {
				return err
			}
			return nil
		})
	}
}

// EncodeUploadImageError returns an encoder for errors returned by the
// uploadImage sense endpoint.
func EncodeUploadImageError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
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
				body = NewUploadImageBadRequestResponseBody(res)
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
				body = NewUploadImageInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeActionDetailsResponse returns an encoder for responses returned by the
// sense actionDetails endpoint.
func EncodeActionDetailsResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*senseviews.ActionDetailResult)
		enc := encoder(ctx, w)
		body := NewActionDetailsResponseBody(res.Projected)
		w.WriteHeader(http.StatusCreated)
		return enc.Encode(body)
	}
}

// DecodeActionDetailsRequest returns a decoder for requests sent to the sense
// actionDetails endpoint.
func DecodeActionDetailsRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body ActionDetailsRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateActionDetailsRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			imageID string

			params = mux.Vars(r)
		)
		imageID = params["image_id"]
		if utf8.RuneCountInString(imageID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("imageID", imageID, utf8.RuneCountInString(imageID), 8, true))
		}
		if utf8.RuneCountInString(imageID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("imageID", imageID, utf8.RuneCountInString(imageID), 8, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewActionDetailsPayload(&body, imageID)

		return payload, nil
	}
}

// EncodeActionDetailsError returns an encoder for errors returned by the
// actionDetails sense endpoint.
func EncodeActionDetailsError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
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
				body = NewActionDetailsBadRequestResponseBody(res)
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
				body = NewActionDetailsInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeStartProcessingResponse returns an encoder for responses returned by
// the sense startProcessing endpoint.
func EncodeStartProcessingResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*senseviews.StartProcessingResult)
		enc := encoder(ctx, w)
		body := NewStartProcessingResponseBody(res.Projected)
		w.WriteHeader(http.StatusCreated)
		return enc.Encode(body)
	}
}

// DecodeStartProcessingRequest returns a decoder for requests sent to the
// sense startProcessing endpoint.
func DecodeStartProcessingRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body StartProcessingRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateStartProcessingRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			imageID               string
			appPastelidPassphrase string

			params = mux.Vars(r)
		)
		imageID = params["image_id"]
		if utf8.RuneCountInString(imageID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("imageID", imageID, utf8.RuneCountInString(imageID), 8, true))
		}
		if utf8.RuneCountInString(imageID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("imageID", imageID, utf8.RuneCountInString(imageID), 8, false))
		}
		appPastelidPassphrase = r.Header.Get("app_pastelid_passphrase")
		if appPastelidPassphrase == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("app_pastelid_passphrase", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewStartProcessingPayload(&body, imageID, appPastelidPassphrase)

		return payload, nil
	}
}

// EncodeStartProcessingError returns an encoder for errors returned by the
// startProcessing sense endpoint.
func EncodeStartProcessingError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
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
				body = NewStartProcessingBadRequestResponseBody(res)
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
				body = NewStartProcessingInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// DecodeRegisterTaskStateRequest returns a decoder for requests sent to the
// sense registerTaskState endpoint.
func DecodeRegisterTaskStateRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["taskId"]
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskID", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewRegisterTaskStatePayload(taskID)

		return payload, nil
	}
}

// EncodeRegisterTaskStateError returns an encoder for errors returned by the
// registerTaskState sense endpoint.
func EncodeRegisterTaskStateError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "NotFound":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRegisterTaskStateNotFoundResponseBody(res)
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
				body = NewRegisterTaskStateInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}
