// Code generated by goa v3.13.1, DO NOT EDIT.
//
// collection HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	collection "github.com/pastelnetwork/gonode/walletnode/api/gen/collection"
	collectionviews "github.com/pastelnetwork/gonode/walletnode/api/gen/collection/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeRegisterCollectionResponse returns an encoder for responses returned
// by the collection registerCollection endpoint.
func EncodeRegisterCollectionResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res := v.(*collectionviews.RegisterCollectionResponse)
		enc := encoder(ctx, w)
		body := NewRegisterCollectionResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeRegisterCollectionRequest returns a decoder for requests sent to the
// collection registerCollection endpoint.
func DecodeRegisterCollectionRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			body RegisterCollectionRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateRegisterCollectionRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			key *string
		)
		keyRaw := r.Header.Get("Authorization")
		if keyRaw != "" {
			key = &keyRaw
		}
		payload := NewRegisterCollectionPayload(&body, key)
		if payload.Key != nil {
			if strings.Contains(*payload.Key, " ") {
				// Remove authorization scheme prefix (e.g. "Bearer")
				cred := strings.SplitN(*payload.Key, " ", 2)[1]
				payload.Key = &cred
			}
		}

		return payload, nil
	}
}

// EncodeRegisterCollectionError returns an encoder for errors returned by the
// registerCollection collection endpoint.
func EncodeRegisterCollectionError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en goa.GoaErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.GoaErrorName() {
		case "UnAuthorized":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewRegisterCollectionUnAuthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "BadRequest":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewRegisterCollectionBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewRegisterCollectionNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewRegisterCollectionInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// DecodeRegisterTaskStateRequest returns a decoder for requests sent to the
// collection registerTaskState endpoint.
func DecodeRegisterTaskStateRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			taskID string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["taskId"]
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewRegisterTaskStatePayload(taskID)

		return payload, nil
	}
}

// EncodeRegisterTaskStateError returns an encoder for errors returned by the
// registerTaskState collection endpoint.
func EncodeRegisterTaskStateError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en goa.GoaErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.GoaErrorName() {
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewRegisterTaskStateNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewRegisterTaskStateInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetTaskHistoryResponse returns an encoder for responses returned by
// the collection getTaskHistory endpoint.
func EncodeGetTaskHistoryResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.([]*collection.TaskHistory)
		enc := encoder(ctx, w)
		body := NewGetTaskHistoryResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetTaskHistoryRequest returns a decoder for requests sent to the
// collection getTaskHistory endpoint.
func DecodeGetTaskHistoryRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			taskID string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["taskId"]
		if utf8.RuneCountInString(taskID) < 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, true))
		}
		if utf8.RuneCountInString(taskID) > 8 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("taskId", taskID, utf8.RuneCountInString(taskID), 8, false))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetTaskHistoryPayload(taskID)

		return payload, nil
	}
}

// EncodeGetTaskHistoryError returns an encoder for errors returned by the
// getTaskHistory collection endpoint.
func EncodeGetTaskHistoryError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en goa.GoaErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.GoaErrorName() {
		case "NotFound":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewGetTaskHistoryNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "InternalServerError":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewGetTaskHistoryInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// marshalCollectionTaskHistoryToTaskHistoryResponse builds a value of type
// *TaskHistoryResponse from a value of type *collection.TaskHistory.
func marshalCollectionTaskHistoryToTaskHistoryResponse(v *collection.TaskHistory) *TaskHistoryResponse {
	res := &TaskHistoryResponse{
		Timestamp: v.Timestamp,
		Status:    v.Status,
		Message:   v.Message,
	}
	if v.Details != nil {
		res.Details = marshalCollectionDetailsToDetailsResponse(v.Details)
	}

	return res
}

// marshalCollectionDetailsToDetailsResponse builds a value of type
// *DetailsResponse from a value of type *collection.Details.
func marshalCollectionDetailsToDetailsResponse(v *collection.Details) *DetailsResponse {
	if v == nil {
		return nil
	}
	res := &DetailsResponse{
		Message: v.Message,
	}
	if v.Fields != nil {
		res.Fields = make(map[string]any, len(v.Fields))
		for key, val := range v.Fields {
			tk := key
			tv := val
			res.Fields[tk] = tv
		}
	}

	return res
}
