// Code generated by goa v3.15.0, DO NOT EDIT.
//
// StorageChallenge HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	storagechallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/storage_challenge"
	storagechallengeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/storage_challenge/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeGetSummaryStatsResponse returns an encoder for responses returned by
// the StorageChallenge getSummaryStats endpoint.
func EncodeGetSummaryStatsResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res := v.(*storagechallengeviews.SummaryStatsResult)
		enc := encoder(ctx, w)
		body := NewGetSummaryStatsResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetSummaryStatsRequest returns a decoder for requests sent to the
// StorageChallenge getSummaryStats endpoint.
func DecodeGetSummaryStatsRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			from *string
			to   *string
			pid  string
			key  string
			err  error
		)
		fromRaw := r.URL.Query().Get("from")
		if fromRaw != "" {
			from = &fromRaw
		}
		if from != nil {
			err = goa.MergeErrors(err, goa.ValidateFormat("from", *from, goa.FormatDateTime))
		}
		toRaw := r.URL.Query().Get("to")
		if toRaw != "" {
			to = &toRaw
		}
		if to != nil {
			err = goa.MergeErrors(err, goa.ValidateFormat("to", *to, goa.FormatDateTime))
		}
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
		payload := NewGetSummaryStatsPayload(from, to, pid, key)
		if strings.Contains(payload.Key, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Key, " ", 2)[1]
			payload.Key = cred
		}

		return payload, nil
	}
}

// EncodeGetSummaryStatsError returns an encoder for errors returned by the
// getSummaryStats StorageChallenge endpoint.
func EncodeGetSummaryStatsError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
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
				body = NewGetSummaryStatsUnauthorizedResponseBody(res)
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
				body = NewGetSummaryStatsBadRequestResponseBody(res)
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
				body = NewGetSummaryStatsNotFoundResponseBody(res)
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
				body = NewGetSummaryStatsInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetDetailedLogsResponse returns an encoder for responses returned by
// the StorageChallenge getDetailedLogs endpoint.
func EncodeGetDetailedLogsResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.([]*storagechallenge.StorageMessage)
		enc := encoder(ctx, w)
		body := NewGetDetailedLogsResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetDetailedLogsRequest returns a decoder for requests sent to the
// StorageChallenge getDetailedLogs endpoint.
func DecodeGetDetailedLogsRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		var (
			pid         string
			challengeID string
			key         string
			err         error
		)
		pid = r.URL.Query().Get("pid")
		if pid == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("pid", "query string"))
		}
		challengeID = r.URL.Query().Get("challenge_id")
		if challengeID == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("challenge_id", "query string"))
		}
		key = r.Header.Get("Authorization")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("key", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetDetailedLogsPayload(pid, challengeID, key)
		if strings.Contains(payload.Key, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Key, " ", 2)[1]
			payload.Key = cred
		}

		return payload, nil
	}
}

// EncodeGetDetailedLogsError returns an encoder for errors returned by the
// getDetailedLogs StorageChallenge endpoint.
func EncodeGetDetailedLogsError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
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
				body = NewGetDetailedLogsUnauthorizedResponseBody(res)
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
				body = NewGetDetailedLogsBadRequestResponseBody(res)
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
				body = NewGetDetailedLogsNotFoundResponseBody(res)
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
				body = NewGetDetailedLogsInternalServerErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// marshalStoragechallengeviewsSCSummaryStatsViewToSCSummaryStatsResponseBody
// builds a value of type *SCSummaryStatsResponseBody from a value of type
// *storagechallengeviews.SCSummaryStatsView.
func marshalStoragechallengeviewsSCSummaryStatsViewToSCSummaryStatsResponseBody(v *storagechallengeviews.SCSummaryStatsView) *SCSummaryStatsResponseBody {
	res := &SCSummaryStatsResponseBody{
		TotalChallengesIssued:                    *v.TotalChallengesIssued,
		TotalChallengesProcessedByRecipient:      *v.TotalChallengesProcessedByRecipient,
		TotalChallengesEvaluatedByChallenger:     *v.TotalChallengesEvaluatedByChallenger,
		TotalChallengesVerified:                  *v.TotalChallengesVerified,
		NoOfSlowResponsesObservedByObservers:     *v.NoOfSlowResponsesObservedByObservers,
		NoOfInvalidSignaturesObservedByObservers: *v.NoOfInvalidSignaturesObservedByObservers,
		NoOfInvalidEvaluationObservedByObservers: *v.NoOfInvalidEvaluationObservedByObservers,
	}

	return res
}

// marshalStoragechallengeStorageMessageToStorageMessageResponse builds a value
// of type *StorageMessageResponse from a value of type
// *storagechallenge.StorageMessage.
func marshalStoragechallengeStorageMessageToStorageMessageResponse(v *storagechallenge.StorageMessage) *StorageMessageResponse {
	res := &StorageMessageResponse{
		ChallengeID:     v.ChallengeID,
		MessageType:     v.MessageType,
		SenderID:        v.SenderID,
		SenderSignature: v.SenderSignature,
		ChallengerID:    v.ChallengerID,
		RecipientID:     v.RecipientID,
	}
	if v.Challenge != nil {
		res.Challenge = marshalStoragechallengeChallengeDataToChallengeDataResponse(v.Challenge)
	}
	if v.Observers != nil {
		res.Observers = make([]string, len(v.Observers))
		for i, val := range v.Observers {
			res.Observers[i] = val
		}
	} else {
		res.Observers = []string{}
	}
	if v.Response != nil {
		res.Response = marshalStoragechallengeResponseDataToResponseDataResponse(v.Response)
	}
	if v.ChallengerEvaluation != nil {
		res.ChallengerEvaluation = marshalStoragechallengeEvaluationDataToEvaluationDataResponse(v.ChallengerEvaluation)
	}
	if v.ObserverEvaluation != nil {
		res.ObserverEvaluation = marshalStoragechallengeObserverEvaluationDataToObserverEvaluationDataResponse(v.ObserverEvaluation)
	}

	return res
}

// marshalStoragechallengeChallengeDataToChallengeDataResponse builds a value
// of type *ChallengeDataResponse from a value of type
// *storagechallenge.ChallengeData.
func marshalStoragechallengeChallengeDataToChallengeDataResponse(v *storagechallenge.ChallengeData) *ChallengeDataResponse {
	if v == nil {
		return nil
	}
	res := &ChallengeDataResponse{
		Block:      v.Block,
		Merkelroot: v.Merkelroot,
		Timestamp:  v.Timestamp,
		FileHash:   v.FileHash,
		StartIndex: v.StartIndex,
		EndIndex:   v.EndIndex,
	}

	return res
}

// marshalStoragechallengeResponseDataToResponseDataResponse builds a value of
// type *ResponseDataResponse from a value of type
// *storagechallenge.ResponseData.
func marshalStoragechallengeResponseDataToResponseDataResponse(v *storagechallenge.ResponseData) *ResponseDataResponse {
	if v == nil {
		return nil
	}
	res := &ResponseDataResponse{
		Block:      v.Block,
		Merkelroot: v.Merkelroot,
		Hash:       v.Hash,
		Timestamp:  v.Timestamp,
	}

	return res
}

// marshalStoragechallengeEvaluationDataToEvaluationDataResponse builds a value
// of type *EvaluationDataResponse from a value of type
// *storagechallenge.EvaluationData.
func marshalStoragechallengeEvaluationDataToEvaluationDataResponse(v *storagechallenge.EvaluationData) *EvaluationDataResponse {
	if v == nil {
		return nil
	}
	res := &EvaluationDataResponse{
		Block:      v.Block,
		Merkelroot: v.Merkelroot,
		Timestamp:  v.Timestamp,
		Hash:       v.Hash,
		IsVerified: v.IsVerified,
	}

	return res
}

// marshalStoragechallengeObserverEvaluationDataToObserverEvaluationDataResponse
// builds a value of type *ObserverEvaluationDataResponse from a value of type
// *storagechallenge.ObserverEvaluationData.
func marshalStoragechallengeObserverEvaluationDataToObserverEvaluationDataResponse(v *storagechallenge.ObserverEvaluationData) *ObserverEvaluationDataResponse {
	if v == nil {
		return nil
	}
	res := &ObserverEvaluationDataResponse{
		Block:                   v.Block,
		Merkelroot:              v.Merkelroot,
		IsChallengeTimestampOK:  v.IsChallengeTimestampOK,
		IsProcessTimestampOK:    v.IsProcessTimestampOK,
		IsEvaluationTimestampOK: v.IsEvaluationTimestampOK,
		IsRecipientSignatureOK:  v.IsRecipientSignatureOK,
		IsChallengerSignatureOK: v.IsChallengerSignatureOK,
		IsEvaluationResultOK:    v.IsEvaluationResultOK,
		Reason:                  v.Reason,
		TrueHash:                v.TrueHash,
		Timestamp:               v.Timestamp,
	}

	return res
}
