// Code generated by goa v3.15.0, DO NOT EDIT.
//
// Score HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	score "github.com/pastelnetwork/gonode/walletnode/api/gen/score"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeGetAggregatedChallengesScoresResponse returns an encoder for responses
// returned by the Score getAggregatedChallengesScores endpoint.
func EncodeGetAggregatedChallengesScoresResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.([]*score.ChallengesScores)
		enc := encoder(ctx, w)
		body := NewGetAggregatedChallengesScoresResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetAggregatedChallengesScoresRequest returns a decoder for requests
// sent to the Score getAggregatedChallengesScores endpoint.
func DecodeGetAggregatedChallengesScoresRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			pid string
			key string
			err error
		)
		pid = r.URL.Query().Get("pid")
		if pid == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("pid", "query string"))
		}
		key = r.Header.Get("Authorization")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("key", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetAggregatedChallengesScoresPayload(pid, key)
		if strings.Contains(payload.Key, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Key, " ", 2)[1]
			payload.Key = cred
		}

		return payload, nil
	}
}

// EncodeGetAggregatedChallengesScoresError returns an encoder for errors
// returned by the getAggregatedChallengesScores Score endpoint.
func EncodeGetAggregatedChallengesScoresError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en goa.GoaErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.GoaErrorName() {
		case "Unauthorized":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body any
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewGetAggregatedChallengesScoresUnauthorizedResponseBody(res)
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
				body = NewGetAggregatedChallengesScoresBadRequestResponseBody(res)
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
				body = NewGetAggregatedChallengesScoresNotFoundResponseBody(res)
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
				body = NewGetAggregatedChallengesScoresInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// marshalScoreChallengesScoresToChallengesScoresResponse builds a value of
// type *ChallengesScoresResponse from a value of type *score.ChallengesScores.
func marshalScoreChallengesScoresToChallengesScoresResponse(v *score.ChallengesScores) *ChallengesScoresResponse {
	res := &ChallengesScoresResponse{
		NodeID:                    v.NodeID,
		IPAddress:                 v.IPAddress,
		StorageChallengeScore:     v.StorageChallengeScore,
		HealthCheckChallengeScore: v.HealthCheckChallengeScore,
	}

	return res
}