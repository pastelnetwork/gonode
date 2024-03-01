// Code generated by goa v3.14.0, DO NOT EDIT.
//
// StorageChallenge HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	storagechallenge "github.com/pastelnetwork/gonode/walletnode/api/gen/storage_challenge"
	storagechallengeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/storage_challenge/views"
	goahttp "goa.design/goa/v3/http"
)

// BuildGetSummaryStatsRequest instantiates a HTTP request object with method
// and path set to call the "StorageChallenge" service "getSummaryStats"
// endpoint
func (c *Client) BuildGetSummaryStatsRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetSummaryStatsStorageChallengePath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("StorageChallenge", "getSummaryStats", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetSummaryStatsRequest returns an encoder for requests sent to the
// StorageChallenge getSummaryStats server.
func EncodeGetSummaryStatsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*storagechallenge.GetSummaryStatsPayload)
		if !ok {
			return goahttp.ErrInvalidType("StorageChallenge", "getSummaryStats", "*storagechallenge.GetSummaryStatsPayload", v)
		}
		{
			head := p.Key
			req.Header.Set("Authorization", head)
		}
		values := req.URL.Query()
		if p.From != nil {
			values.Add("from", *p.From)
		}
		if p.To != nil {
			values.Add("to", *p.To)
		}
		values.Add("pid", p.Pid)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetSummaryStatsResponse returns a decoder for responses returned by
// the StorageChallenge getSummaryStats endpoint. restoreBody controls whether
// the response body should be restored after having been read.
// DecodeGetSummaryStatsResponse may return the following errors:
//   - "Unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//   - "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//   - "NotFound" (type *goa.ServiceError): http.StatusNotFound
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
func DecodeGetSummaryStatsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetSummaryStatsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("StorageChallenge", "getSummaryStats", err)
			}
			p := NewGetSummaryStatsSummaryStatsResultOK(&body)
			view := "default"
			vres := &storagechallengeviews.SummaryStatsResult{Projected: p, View: view}
			if err = storagechallengeviews.ValidateSummaryStatsResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("StorageChallenge", "getSummaryStats", err)
			}
			res := storagechallenge.NewSummaryStatsResult(vres)
			return res, nil
		case http.StatusUnauthorized:
			var (
				body GetSummaryStatsUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("StorageChallenge", "getSummaryStats", err)
			}
			err = ValidateGetSummaryStatsUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("StorageChallenge", "getSummaryStats", err)
			}
			return nil, NewGetSummaryStatsUnauthorized(&body)
		case http.StatusBadRequest:
			var (
				body GetSummaryStatsBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("StorageChallenge", "getSummaryStats", err)
			}
			err = ValidateGetSummaryStatsBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("StorageChallenge", "getSummaryStats", err)
			}
			return nil, NewGetSummaryStatsBadRequest(&body)
		case http.StatusNotFound:
			var (
				body GetSummaryStatsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("StorageChallenge", "getSummaryStats", err)
			}
			err = ValidateGetSummaryStatsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("StorageChallenge", "getSummaryStats", err)
			}
			return nil, NewGetSummaryStatsNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetSummaryStatsInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("StorageChallenge", "getSummaryStats", err)
			}
			err = ValidateGetSummaryStatsInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("StorageChallenge", "getSummaryStats", err)
			}
			return nil, NewGetSummaryStatsInternalServerError(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("StorageChallenge", "getSummaryStats", resp.StatusCode, string(body))
		}
	}
}

// unmarshalSCSummaryStatsResponseBodyToStoragechallengeviewsSCSummaryStatsView
// builds a value of type *storagechallengeviews.SCSummaryStatsView from a
// value of type *SCSummaryStatsResponseBody.
func unmarshalSCSummaryStatsResponseBodyToStoragechallengeviewsSCSummaryStatsView(v *SCSummaryStatsResponseBody) *storagechallengeviews.SCSummaryStatsView {
	res := &storagechallengeviews.SCSummaryStatsView{
		TotalChallengesIssued:                    v.TotalChallengesIssued,
		TotalChallengesProcessed:                 v.TotalChallengesProcessed,
		TotalChallengesVerifiedByChallenger:      v.TotalChallengesVerifiedByChallenger,
		TotalChallengesVerifiedByObservers:       v.TotalChallengesVerifiedByObservers,
		NoOfSlowResponsesObservedByObservers:     v.NoOfSlowResponsesObservedByObservers,
		NoOfInvalidSignaturesObservedByObservers: v.NoOfInvalidSignaturesObservedByObservers,
		NoOfInvalidEvaluationObservedByObservers: v.NoOfInvalidEvaluationObservedByObservers,
	}

	return res
}
