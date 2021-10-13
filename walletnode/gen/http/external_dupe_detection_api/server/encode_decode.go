// Code generated by goa v3.5.2, DO NOT EDIT.
//
// external_dupe_detection_api HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"io"
	"net/http"

	externaldupedetectionapiviews "github.com/pastelnetwork/gonode/walletnode/gen/external_dupe_detection_api/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeInitiateSubmissionResponse returns an encoder for responses returned
// by the external_dupe_detection_api initiate_submission endpoint.
func EncodeInitiateSubmissionResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*externaldupedetectionapiviews.Externaldupedetetionapiinitiatesubmissionresult)
		enc := encoder(ctx, w)
		body := NewInitiateSubmissionResponseBody(res.Projected)
		w.WriteHeader(http.StatusCreated)
		return enc.Encode(body)
	}
}

// DecodeInitiateSubmissionRequest returns a decoder for requests sent to the
// external_dupe_detection_api initiate_submission endpoint.
func DecodeInitiateSubmissionRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body InitiateSubmissionRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateInitiateSubmissionRequestBody(&body)
		if err != nil {
			return nil, err
		}
		payload := NewInitiateSubmissionExternalDupeDetetionAPIInitiateSubmission(&body)

		return payload, nil
	}
}

// EncodeInitiateSubmissionError returns an encoder for errors returned by the
// initiate_submission external_dupe_detection_api endpoint.
func EncodeInitiateSubmissionError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
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
				body = NewInitiateSubmissionBadRequestResponseBody(res)
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
				body = NewInitiateSubmissionInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}
