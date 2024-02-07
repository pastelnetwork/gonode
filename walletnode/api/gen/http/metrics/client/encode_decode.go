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
	"fmt"
	"io"
	"net/http"
	"net/url"

	metrics "github.com/pastelnetwork/gonode/walletnode/api/gen/metrics"
	metricsviews "github.com/pastelnetwork/gonode/walletnode/api/gen/metrics/views"
	goahttp "goa.design/goa/v3/http"
)

// BuildGetChallengeReportsRequest instantiates a HTTP request object with
// method and path set to call the "metrics" service "getChallengeReports"
// endpoint
func (c *Client) BuildGetChallengeReportsRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetChallengeReportsMetricsPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("metrics", "getChallengeReports", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetChallengeReportsRequest returns an encoder for requests sent to the
// metrics getChallengeReports server.
func EncodeGetChallengeReportsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*metrics.GetChallengeReportsPayload)
		if !ok {
			return goahttp.ErrInvalidType("metrics", "getChallengeReports", "*metrics.GetChallengeReportsPayload", v)
		}
		{
			head := p.Key
			req.Header.Set("Authorization", head)
		}
		values := req.URL.Query()
		values.Add("pid", p.Pid)
		if p.ChallengeID != nil {
			values.Add("challenge_id", *p.ChallengeID)
		}
		if p.Count != nil {
			values.Add("count", fmt.Sprintf("%v", *p.Count))
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetChallengeReportsResponse returns a decoder for responses returned
// by the metrics getChallengeReports endpoint. restoreBody controls whether
// the response body should be restored after having been read.
// DecodeGetChallengeReportsResponse may return the following errors:
//   - "Unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//   - "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//   - "NotFound" (type *goa.ServiceError): http.StatusNotFound
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
func DecodeGetChallengeReportsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
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
				body GetChallengeReportsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getChallengeReports", err)
			}
			res := NewGetChallengeReportsSelfHealingChallengeReportsOK(&body)
			return res, nil
		case http.StatusUnauthorized:
			var (
				body GetChallengeReportsUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getChallengeReports", err)
			}
			err = ValidateGetChallengeReportsUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getChallengeReports", err)
			}
			return nil, NewGetChallengeReportsUnauthorized(&body)
		case http.StatusBadRequest:
			var (
				body GetChallengeReportsBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getChallengeReports", err)
			}
			err = ValidateGetChallengeReportsBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getChallengeReports", err)
			}
			return nil, NewGetChallengeReportsBadRequest(&body)
		case http.StatusNotFound:
			var (
				body GetChallengeReportsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getChallengeReports", err)
			}
			err = ValidateGetChallengeReportsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getChallengeReports", err)
			}
			return nil, NewGetChallengeReportsNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetChallengeReportsInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getChallengeReports", err)
			}
			err = ValidateGetChallengeReportsInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getChallengeReports", err)
			}
			return nil, NewGetChallengeReportsInternalServerError(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("metrics", "getChallengeReports", resp.StatusCode, string(body))
		}
	}
}

// BuildGetMetricsRequest instantiates a HTTP request object with method and
// path set to call the "metrics" service "getMetrics" endpoint
func (c *Client) BuildGetMetricsRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetMetricsMetricsPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("metrics", "getMetrics", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetMetricsRequest returns an encoder for requests sent to the metrics
// getMetrics server.
func EncodeGetMetricsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*metrics.GetMetricsPayload)
		if !ok {
			return goahttp.ErrInvalidType("metrics", "getMetrics", "*metrics.GetMetricsPayload", v)
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

// DecodeGetMetricsResponse returns a decoder for responses returned by the
// metrics getMetrics endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeGetMetricsResponse may return the following errors:
//   - "Unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//   - "BadRequest" (type *goa.ServiceError): http.StatusBadRequest
//   - "NotFound" (type *goa.ServiceError): http.StatusNotFound
//   - "InternalServerError" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
func DecodeGetMetricsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
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
				body GetMetricsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getMetrics", err)
			}
			p := NewGetMetricsMetricsResultOK(&body)
			view := "default"
			vres := &metricsviews.MetricsResult{Projected: p, View: view}
			if err = metricsviews.ValidateMetricsResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getMetrics", err)
			}
			res := metrics.NewMetricsResult(vres)
			return res, nil
		case http.StatusUnauthorized:
			var (
				body GetMetricsUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getMetrics", err)
			}
			err = ValidateGetMetricsUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getMetrics", err)
			}
			return nil, NewGetMetricsUnauthorized(&body)
		case http.StatusBadRequest:
			var (
				body GetMetricsBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getMetrics", err)
			}
			err = ValidateGetMetricsBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getMetrics", err)
			}
			return nil, NewGetMetricsBadRequest(&body)
		case http.StatusNotFound:
			var (
				body GetMetricsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getMetrics", err)
			}
			err = ValidateGetMetricsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getMetrics", err)
			}
			return nil, NewGetMetricsNotFound(&body)
		case http.StatusInternalServerError:
			var (
				body GetMetricsInternalServerErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("metrics", "getMetrics", err)
			}
			err = ValidateGetMetricsInternalServerErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("metrics", "getMetrics", err)
			}
			return nil, NewGetMetricsInternalServerError(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("metrics", "getMetrics", resp.StatusCode, string(body))
		}
	}
}

// unmarshalSelfHealingChallengeReportKVResponseBodyToMetricsSelfHealingChallengeReportKV
// builds a value of type *metrics.SelfHealingChallengeReportKV from a value of
// type *SelfHealingChallengeReportKVResponseBody.
func unmarshalSelfHealingChallengeReportKVResponseBodyToMetricsSelfHealingChallengeReportKV(v *SelfHealingChallengeReportKVResponseBody) *metrics.SelfHealingChallengeReportKV {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingChallengeReportKV{
		ChallengeID: v.ChallengeID,
	}
	if v.Report != nil {
		res.Report = unmarshalSelfHealingChallengeReportResponseBodyToMetricsSelfHealingChallengeReport(v.Report)
	}

	return res
}

// unmarshalSelfHealingChallengeReportResponseBodyToMetricsSelfHealingChallengeReport
// builds a value of type *metrics.SelfHealingChallengeReport from a value of
// type *SelfHealingChallengeReportResponseBody.
func unmarshalSelfHealingChallengeReportResponseBodyToMetricsSelfHealingChallengeReport(v *SelfHealingChallengeReportResponseBody) *metrics.SelfHealingChallengeReport {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingChallengeReport{}
	if v.Messages != nil {
		res.Messages = make([]*metrics.SelfHealingMessageKV, len(v.Messages))
		for i, val := range v.Messages {
			res.Messages[i] = unmarshalSelfHealingMessageKVResponseBodyToMetricsSelfHealingMessageKV(val)
		}
	}

	return res
}

// unmarshalSelfHealingMessageKVResponseBodyToMetricsSelfHealingMessageKV
// builds a value of type *metrics.SelfHealingMessageKV from a value of type
// *SelfHealingMessageKVResponseBody.
func unmarshalSelfHealingMessageKVResponseBodyToMetricsSelfHealingMessageKV(v *SelfHealingMessageKVResponseBody) *metrics.SelfHealingMessageKV {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingMessageKV{
		MessageType: v.MessageType,
	}
	if v.Messages != nil {
		res.Messages = make([]*metrics.SelfHealingMessage, len(v.Messages))
		for i, val := range v.Messages {
			res.Messages[i] = unmarshalSelfHealingMessageResponseBodyToMetricsSelfHealingMessage(val)
		}
	}

	return res
}

// unmarshalSelfHealingMessageResponseBodyToMetricsSelfHealingMessage builds a
// value of type *metrics.SelfHealingMessage from a value of type
// *SelfHealingMessageResponseBody.
func unmarshalSelfHealingMessageResponseBodyToMetricsSelfHealingMessage(v *SelfHealingMessageResponseBody) *metrics.SelfHealingMessage {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingMessage{
		TriggerID:       v.TriggerID,
		MessageType:     v.MessageType,
		SenderID:        v.SenderID,
		SenderSignature: v.SenderSignature,
	}
	if v.Data != nil {
		res.Data = unmarshalSelfHealingMessageDataResponseBodyToMetricsSelfHealingMessageData(v.Data)
	}

	return res
}

// unmarshalSelfHealingMessageDataResponseBodyToMetricsSelfHealingMessageData
// builds a value of type *metrics.SelfHealingMessageData from a value of type
// *SelfHealingMessageDataResponseBody.
func unmarshalSelfHealingMessageDataResponseBodyToMetricsSelfHealingMessageData(v *SelfHealingMessageDataResponseBody) *metrics.SelfHealingMessageData {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingMessageData{
		ChallengerID: v.ChallengerID,
		RecipientID:  v.RecipientID,
	}
	if v.Challenge != nil {
		res.Challenge = unmarshalSelfHealingChallengeDataResponseBodyToMetricsSelfHealingChallengeData(v.Challenge)
	}
	if v.Response != nil {
		res.Response = unmarshalSelfHealingResponseDataResponseBodyToMetricsSelfHealingResponseData(v.Response)
	}
	if v.Verification != nil {
		res.Verification = unmarshalSelfHealingVerificationDataResponseBodyToMetricsSelfHealingVerificationData(v.Verification)
	}

	return res
}

// unmarshalSelfHealingChallengeDataResponseBodyToMetricsSelfHealingChallengeData
// builds a value of type *metrics.SelfHealingChallengeData from a value of
// type *SelfHealingChallengeDataResponseBody.
func unmarshalSelfHealingChallengeDataResponseBodyToMetricsSelfHealingChallengeData(v *SelfHealingChallengeDataResponseBody) *metrics.SelfHealingChallengeData {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingChallengeData{
		Block:            v.Block,
		Merkelroot:       v.Merkelroot,
		Timestamp:        v.Timestamp,
		NodesOnWatchlist: v.NodesOnWatchlist,
	}
	if v.ChallengeTickets != nil {
		res.ChallengeTickets = make([]*metrics.ChallengeTicket, len(v.ChallengeTickets))
		for i, val := range v.ChallengeTickets {
			res.ChallengeTickets[i] = unmarshalChallengeTicketResponseBodyToMetricsChallengeTicket(val)
		}
	}

	return res
}

// unmarshalChallengeTicketResponseBodyToMetricsChallengeTicket builds a value
// of type *metrics.ChallengeTicket from a value of type
// *ChallengeTicketResponseBody.
func unmarshalChallengeTicketResponseBodyToMetricsChallengeTicket(v *ChallengeTicketResponseBody) *metrics.ChallengeTicket {
	if v == nil {
		return nil
	}
	res := &metrics.ChallengeTicket{
		TxID:       v.TxID,
		TicketType: v.TicketType,
		DataHash:   v.DataHash,
		Recipient:  v.Recipient,
	}
	if v.MissingKeys != nil {
		res.MissingKeys = make([]string, len(v.MissingKeys))
		for i, val := range v.MissingKeys {
			res.MissingKeys[i] = val
		}
	}

	return res
}

// unmarshalSelfHealingResponseDataResponseBodyToMetricsSelfHealingResponseData
// builds a value of type *metrics.SelfHealingResponseData from a value of type
// *SelfHealingResponseDataResponseBody.
func unmarshalSelfHealingResponseDataResponseBodyToMetricsSelfHealingResponseData(v *SelfHealingResponseDataResponseBody) *metrics.SelfHealingResponseData {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingResponseData{
		ChallengeID: v.ChallengeID,
		Block:       v.Block,
		Merkelroot:  v.Merkelroot,
		Timestamp:   v.Timestamp,
	}
	if v.RespondedTicket != nil {
		res.RespondedTicket = unmarshalRespondedTicketResponseBodyToMetricsRespondedTicket(v.RespondedTicket)
	}
	if v.Verifiers != nil {
		res.Verifiers = make([]string, len(v.Verifiers))
		for i, val := range v.Verifiers {
			res.Verifiers[i] = val
		}
	}

	return res
}

// unmarshalRespondedTicketResponseBodyToMetricsRespondedTicket builds a value
// of type *metrics.RespondedTicket from a value of type
// *RespondedTicketResponseBody.
func unmarshalRespondedTicketResponseBodyToMetricsRespondedTicket(v *RespondedTicketResponseBody) *metrics.RespondedTicket {
	if v == nil {
		return nil
	}
	res := &metrics.RespondedTicket{
		TxID:                     v.TxID,
		TicketType:               v.TicketType,
		ReconstructedFileHash:    v.ReconstructedFileHash,
		RaptorQSymbols:           v.RaptorQSymbols,
		IsReconstructionRequired: v.IsReconstructionRequired,
	}
	if v.MissingKeys != nil {
		res.MissingKeys = make([]string, len(v.MissingKeys))
		for i, val := range v.MissingKeys {
			res.MissingKeys[i] = val
		}
	}
	if v.SenseFileIds != nil {
		res.SenseFileIds = make([]string, len(v.SenseFileIds))
		for i, val := range v.SenseFileIds {
			res.SenseFileIds[i] = val
		}
	}

	return res
}

// unmarshalSelfHealingVerificationDataResponseBodyToMetricsSelfHealingVerificationData
// builds a value of type *metrics.SelfHealingVerificationData from a value of
// type *SelfHealingVerificationDataResponseBody.
func unmarshalSelfHealingVerificationDataResponseBodyToMetricsSelfHealingVerificationData(v *SelfHealingVerificationDataResponseBody) *metrics.SelfHealingVerificationData {
	if v == nil {
		return nil
	}
	res := &metrics.SelfHealingVerificationData{
		ChallengeID: v.ChallengeID,
		Block:       v.Block,
		Merkelroot:  v.Merkelroot,
		Timestamp:   v.Timestamp,
	}
	if v.VerifiedTicket != nil {
		res.VerifiedTicket = unmarshalVerifiedTicketResponseBodyToMetricsVerifiedTicket(v.VerifiedTicket)
	}
	if v.VerifiersData != nil {
		res.VerifiersData = make(map[string][]byte, len(v.VerifiersData))
		for key, val := range v.VerifiersData {
			tk := key
			tv := val
			res.VerifiersData[tk] = tv
		}
	}

	return res
}

// unmarshalVerifiedTicketResponseBodyToMetricsVerifiedTicket builds a value of
// type *metrics.VerifiedTicket from a value of type
// *VerifiedTicketResponseBody.
func unmarshalVerifiedTicketResponseBodyToMetricsVerifiedTicket(v *VerifiedTicketResponseBody) *metrics.VerifiedTicket {
	if v == nil {
		return nil
	}
	res := &metrics.VerifiedTicket{
		TxID:                     v.TxID,
		TicketType:               v.TicketType,
		ReconstructedFileHash:    v.ReconstructedFileHash,
		IsReconstructionRequired: v.IsReconstructionRequired,
		RaptorQSymbols:           v.RaptorQSymbols,
		IsVerified:               v.IsVerified,
		Message:                  v.Message,
	}
	if v.MissingKeys != nil {
		res.MissingKeys = make([]string, len(v.MissingKeys))
		for i, val := range v.MissingKeys {
			res.MissingKeys[i] = val
		}
	}
	if v.SenseFileIds != nil {
		res.SenseFileIds = make([]string, len(v.SenseFileIds))
		for i, val := range v.SenseFileIds {
			res.SenseFileIds[i] = val
		}
	}

	return res
}

// unmarshalSHTriggerMetricResponseBodyToMetricsviewsSHTriggerMetricView builds
// a value of type *metricsviews.SHTriggerMetricView from a value of type
// *SHTriggerMetricResponseBody.
func unmarshalSHTriggerMetricResponseBodyToMetricsviewsSHTriggerMetricView(v *SHTriggerMetricResponseBody) *metricsviews.SHTriggerMetricView {
	res := &metricsviews.SHTriggerMetricView{
		TriggerID:              v.TriggerID,
		NodesOffline:           v.NodesOffline,
		ListOfNodes:            v.ListOfNodes,
		TotalFilesIdentified:   v.TotalFilesIdentified,
		TotalTicketsIdentified: v.TotalTicketsIdentified,
	}

	return res
}

// unmarshalSHExecutionMetricsResponseBodyToMetricsviewsSHExecutionMetricsView
// builds a value of type *metricsviews.SHExecutionMetricsView from a value of
// type *SHExecutionMetricsResponseBody.
func unmarshalSHExecutionMetricsResponseBodyToMetricsviewsSHExecutionMetricsView(v *SHExecutionMetricsResponseBody) *metricsviews.SHExecutionMetricsView {
	res := &metricsviews.SHExecutionMetricsView{
		TotalChallengesIssued:                                 v.TotalChallengesIssued,
		TotalChallengesAcknowledged:                           v.TotalChallengesAcknowledged,
		TotalChallengesRejected:                               v.TotalChallengesRejected,
		TotalChallengesAccepted:                               v.TotalChallengesAccepted,
		TotalChallengeEvaluationsVerified:                     v.TotalChallengeEvaluationsVerified,
		TotalReconstructionRequiredEvaluationsApproved:        v.TotalReconstructionRequiredEvaluationsApproved,
		TotalReconstructionNotRequiredEvaluationsApproved:     v.TotalReconstructionNotRequiredEvaluationsApproved,
		TotalChallengeEvaluationsUnverified:                   v.TotalChallengeEvaluationsUnverified,
		TotalReconstructionRequiredEvaluationsNotApproved:     v.TotalReconstructionRequiredEvaluationsNotApproved,
		TotalReconstructionsNotRequiredEvaluationsNotApproved: v.TotalReconstructionsNotRequiredEvaluationsNotApproved,
		TotalReconstructionRequiredHashMismatch:               v.TotalReconstructionRequiredHashMismatch,
		TotalFilesHealed:                                      v.TotalFilesHealed,
		TotalFileHealingFailed:                                v.TotalFileHealingFailed,
	}

	return res
}
