// Code generated by goa v3.15.0, DO NOT EDIT.
//
// HealthCheckChallenge HTTP client types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	healthcheckchallengeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/health_check_challenge/views"
	goa "goa.design/goa/v3/pkg"
)

// GetSummaryStatsResponseBody is the type of the "HealthCheckChallenge"
// service "getSummaryStats" endpoint HTTP response body.
type GetSummaryStatsResponseBody struct {
	// HCSummaryStats represents health check challenge summary of metrics stats
	HcSummaryStats *HCSummaryStatsResponseBody `form:"hc_summary_stats,omitempty" json:"hc_summary_stats,omitempty" xml:"hc_summary_stats,omitempty"`
}

// GetSummaryStatsUnauthorizedResponseBody is the type of the
// "HealthCheckChallenge" service "getSummaryStats" endpoint HTTP response body
// for the "Unauthorized" error.
type GetSummaryStatsUnauthorizedResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// GetSummaryStatsBadRequestResponseBody is the type of the
// "HealthCheckChallenge" service "getSummaryStats" endpoint HTTP response body
// for the "BadRequest" error.
type GetSummaryStatsBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// GetSummaryStatsNotFoundResponseBody is the type of the
// "HealthCheckChallenge" service "getSummaryStats" endpoint HTTP response body
// for the "NotFound" error.
type GetSummaryStatsNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// GetSummaryStatsInternalServerErrorResponseBody is the type of the
// "HealthCheckChallenge" service "getSummaryStats" endpoint HTTP response body
// for the "InternalServerError" error.
type GetSummaryStatsInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// HCSummaryStatsResponseBody is used to define fields on response body types.
type HCSummaryStatsResponseBody struct {
	// Total number of challenges issued
	TotalChallengesIssued *int `form:"total_challenges_issued,omitempty" json:"total_challenges_issued,omitempty" xml:"total_challenges_issued,omitempty"`
	// Total number of challenges processed by the recipient node
	TotalChallengesProcessedByRecipient *int `form:"total_challenges_processed_by_recipient,omitempty" json:"total_challenges_processed_by_recipient,omitempty" xml:"total_challenges_processed_by_recipient,omitempty"`
	// Total number of challenges evaluated by the challenger node
	TotalChallengesEvaluatedByChallenger *int `form:"total_challenges_evaluated_by_challenger,omitempty" json:"total_challenges_evaluated_by_challenger,omitempty" xml:"total_challenges_evaluated_by_challenger,omitempty"`
	// Total number of challenges verified by observers
	TotalChallengesVerified *int `form:"total_challenges_verified,omitempty" json:"total_challenges_verified,omitempty" xml:"total_challenges_verified,omitempty"`
	// challenges failed due to slow-responses evaluated by observers
	NoOfSlowResponsesObservedByObservers *int `form:"no_of_slow_responses_observed_by_observers,omitempty" json:"no_of_slow_responses_observed_by_observers,omitempty" xml:"no_of_slow_responses_observed_by_observers,omitempty"`
	// challenges failed due to invalid signatures evaluated by observers
	NoOfInvalidSignaturesObservedByObservers *int `form:"no_of_invalid_signatures_observed_by_observers,omitempty" json:"no_of_invalid_signatures_observed_by_observers,omitempty" xml:"no_of_invalid_signatures_observed_by_observers,omitempty"`
	// challenges failed due to invalid evaluation evaluated by observers
	NoOfInvalidEvaluationObservedByObservers *int `form:"no_of_invalid_evaluation_observed_by_observers,omitempty" json:"no_of_invalid_evaluation_observed_by_observers,omitempty" xml:"no_of_invalid_evaluation_observed_by_observers,omitempty"`
}

// NewGetSummaryStatsHcSummaryStatsResultOK builds a "HealthCheckChallenge"
// service "getSummaryStats" endpoint result from a HTTP "OK" response.
func NewGetSummaryStatsHcSummaryStatsResultOK(body *GetSummaryStatsResponseBody) *healthcheckchallengeviews.HcSummaryStatsResultView {
	v := &healthcheckchallengeviews.HcSummaryStatsResultView{}
	v.HcSummaryStats = unmarshalHCSummaryStatsResponseBodyToHealthcheckchallengeviewsHCSummaryStatsView(body.HcSummaryStats)

	return v
}

// NewGetSummaryStatsUnauthorized builds a HealthCheckChallenge service
// getSummaryStats endpoint Unauthorized error.
func NewGetSummaryStatsUnauthorized(body *GetSummaryStatsUnauthorizedResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewGetSummaryStatsBadRequest builds a HealthCheckChallenge service
// getSummaryStats endpoint BadRequest error.
func NewGetSummaryStatsBadRequest(body *GetSummaryStatsBadRequestResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewGetSummaryStatsNotFound builds a HealthCheckChallenge service
// getSummaryStats endpoint NotFound error.
func NewGetSummaryStatsNotFound(body *GetSummaryStatsNotFoundResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewGetSummaryStatsInternalServerError builds a HealthCheckChallenge service
// getSummaryStats endpoint InternalServerError error.
func NewGetSummaryStatsInternalServerError(body *GetSummaryStatsInternalServerErrorResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// ValidateGetSummaryStatsUnauthorizedResponseBody runs the validations defined
// on getSummaryStats_Unauthorized_response_body
func ValidateGetSummaryStatsUnauthorizedResponseBody(body *GetSummaryStatsUnauthorizedResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateGetSummaryStatsBadRequestResponseBody runs the validations defined
// on getSummaryStats_BadRequest_response_body
func ValidateGetSummaryStatsBadRequestResponseBody(body *GetSummaryStatsBadRequestResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateGetSummaryStatsNotFoundResponseBody runs the validations defined on
// getSummaryStats_NotFound_response_body
func ValidateGetSummaryStatsNotFoundResponseBody(body *GetSummaryStatsNotFoundResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateGetSummaryStatsInternalServerErrorResponseBody runs the validations
// defined on getSummaryStats_InternalServerError_response_body
func ValidateGetSummaryStatsInternalServerErrorResponseBody(body *GetSummaryStatsInternalServerErrorResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateHCSummaryStatsResponseBody runs the validations defined on
// HCSummaryStatsResponseBody
func ValidateHCSummaryStatsResponseBody(body *HCSummaryStatsResponseBody) (err error) {
	if body.TotalChallengesIssued == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_issued", "body"))
	}
	if body.TotalChallengesProcessedByRecipient == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_processed_by_recipient", "body"))
	}
	if body.TotalChallengesEvaluatedByChallenger == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_evaluated_by_challenger", "body"))
	}
	if body.TotalChallengesVerified == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_verified", "body"))
	}
	if body.NoOfSlowResponsesObservedByObservers == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("no_of_slow_responses_observed_by_observers", "body"))
	}
	if body.NoOfInvalidSignaturesObservedByObservers == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("no_of_invalid_signatures_observed_by_observers", "body"))
	}
	if body.NoOfInvalidEvaluationObservedByObservers == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("no_of_invalid_evaluation_observed_by_observers", "body"))
	}
	return
}
