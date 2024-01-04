// Code generated by goa v3.14.0, DO NOT EDIT.
//
// metrics HTTP client encoders and decoders
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

	metrics "github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
	metricsviews "github.com/pastelnetwork/gonode/walletnode/api/gen/metrics/views"
	goahttp "goa.design/goa/v3/http"
)

// BuildSelfHealingRequest instantiates a HTTP request object with method and
// path set to call the "metrics" service "selfHealing" endpoint
func (c *Client) BuildSelfHealingRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: SelfHealingMetricsPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("metrics", "selfHealing", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeSelfHealingRequest returns an encoder for requests sent to the metrics
// selfHealing server.
func EncodeSelfHealingRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*metrics.MetricsPayload)
		if !ok {
			return goahttp.ErrInvalidType("metrics", "selfHealing", "*metrics.MetricsPayload", v)
		}
		{
			head := p.Key
			req.Header.Set("Authorization", head)
		}
		values := req.URL.Query()
		values.Add("pid", p.Pid)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeSelfHealingResponse returns a decoder for responses returned by the
// metrics selfHealing endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeSelfHealingResponse may return the following errors:
//   - "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
func DecodeSelfHealingResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
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
		case http.StatusCreated:
			var (
				body SelfHealingResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "selfHealing", err)
			}
			p := NewSelfHealingMetricsViewCreated(&body)
			view := "default"
			vres := &metricsviews.SelfHealingMetrics{Projected: p, View: view}
			if err = metricsviews.ValidateSelfHealingMetrics(vres); err != nil {
				return nil, goahttp.ErrValidationError("metrics", "selfHealing", err)
			}
			res := metrics.NewSelfHealingMetrics(vres)
			return res, nil
		case http.StatusBadRequest:
			var (
				body SelfHealingBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "selfHealing", err)
			}
			err = ValidateSelfHealingBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "selfHealing", err)
			}
			return nil, NewSelfHealingBadRequest(&body)
		case http.StatusInternalServerError:
			var (
				body SelfHealingInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "selfHealing", err)
			}
			err = ValidateSelfHealingInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "selfHealing", err)
			}
			return nil, NewSelfHealingInternalServerError(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("metrics", "selfHealing", resp.StatusCode, string(body))
		}
	}
}
