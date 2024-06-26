// Code generated by goa v3.15.0, DO NOT EDIT.
//
// Score HTTP client types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	score "github.com/pastelnetwork/gonode/walletnode/api/gen/score"
	goa "goa.design/goa/v3/pkg"
)

// GetAggregatedChallengesScoresResponseBody is the type of the "Score" service
// "getAggregatedChallengesScores" endpoint HTTP response body.
type GetAggregatedChallengesScoresResponseBody []*ChallengesScoresResponse

// GetAggregatedChallengesScoresUnauthorizedResponseBody is the type of the
// "Score" service "getAggregatedChallengesScores" endpoint HTTP response body
// for the "Unauthorized" error.
type GetAggregatedChallengesScoresUnauthorizedResponseBody struct {
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

// GetAggregatedChallengesScoresBadRequestResponseBody is the type of the
// "Score" service "getAggregatedChallengesScores" endpoint HTTP response body
// for the "BadRequest" error.
type GetAggregatedChallengesScoresBadRequestResponseBody struct {
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

// GetAggregatedChallengesScoresNotFoundResponseBody is the type of the "Score"
// service "getAggregatedChallengesScores" endpoint HTTP response body for the
// "NotFound" error.
type GetAggregatedChallengesScoresNotFoundResponseBody struct {
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

// GetAggregatedChallengesScoresInternalServerErrorResponseBody is the type of
// the "Score" service "getAggregatedChallengesScores" endpoint HTTP response
// body for the "InternalServerError" error.
type GetAggregatedChallengesScoresInternalServerErrorResponseBody struct {
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

// ChallengesScoresResponse is used to define fields on response body types.
type ChallengesScoresResponse struct {
	// Specific node id
	NodeID *string `form:"node_id,omitempty" json:"node_id,omitempty" xml:"node_id,omitempty"`
	// IPAddress of the node
	IPAddress *string `form:"ip_address,omitempty" json:"ip_address,omitempty" xml:"ip_address,omitempty"`
	// Total accumulated SC challenge score
	StorageChallengeScore *float64 `form:"storage_challenge_score,omitempty" json:"storage_challenge_score,omitempty" xml:"storage_challenge_score,omitempty"`
	// Total accumulated HC challenge score
	HealthCheckChallengeScore *float64 `form:"health_check_challenge_score,omitempty" json:"health_check_challenge_score,omitempty" xml:"health_check_challenge_score,omitempty"`
}

// NewGetAggregatedChallengesScoresChallengesScoresOK builds a "Score" service
// "getAggregatedChallengesScores" endpoint result from a HTTP "OK" response.
func NewGetAggregatedChallengesScoresChallengesScoresOK(body []*ChallengesScoresResponse) []*score.ChallengesScores {
	v := make([]*score.ChallengesScores, len(body))
	for i, val := range body {
		v[i] = unmarshalChallengesScoresResponseToScoreChallengesScores(val)
	}

	return v
}

// NewGetAggregatedChallengesScoresUnauthorized builds a Score service
// getAggregatedChallengesScores endpoint Unauthorized error.
func NewGetAggregatedChallengesScoresUnauthorized(body *GetAggregatedChallengesScoresUnauthorizedResponseBody) *goa.ServiceError {
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

// NewGetAggregatedChallengesScoresBadRequest builds a Score service
// getAggregatedChallengesScores endpoint BadRequest error.
func NewGetAggregatedChallengesScoresBadRequest(body *GetAggregatedChallengesScoresBadRequestResponseBody) *goa.ServiceError {
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

// NewGetAggregatedChallengesScoresNotFound builds a Score service
// getAggregatedChallengesScores endpoint NotFound error.
func NewGetAggregatedChallengesScoresNotFound(body *GetAggregatedChallengesScoresNotFoundResponseBody) *goa.ServiceError {
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

// NewGetAggregatedChallengesScoresInternalServerError builds a Score service
// getAggregatedChallengesScores endpoint InternalServerError error.
func NewGetAggregatedChallengesScoresInternalServerError(body *GetAggregatedChallengesScoresInternalServerErrorResponseBody) *goa.ServiceError {
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

// ValidateGetAggregatedChallengesScoresUnauthorizedResponseBody runs the
// validations defined on
// getAggregatedChallengesScores_Unauthorized_response_body
func ValidateGetAggregatedChallengesScoresUnauthorizedResponseBody(body *GetAggregatedChallengesScoresUnauthorizedResponseBody) (err error) {
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

// ValidateGetAggregatedChallengesScoresBadRequestResponseBody runs the
// validations defined on getAggregatedChallengesScores_BadRequest_response_body
func ValidateGetAggregatedChallengesScoresBadRequestResponseBody(body *GetAggregatedChallengesScoresBadRequestResponseBody) (err error) {
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

// ValidateGetAggregatedChallengesScoresNotFoundResponseBody runs the
// validations defined on getAggregatedChallengesScores_NotFound_response_body
func ValidateGetAggregatedChallengesScoresNotFoundResponseBody(body *GetAggregatedChallengesScoresNotFoundResponseBody) (err error) {
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

// ValidateGetAggregatedChallengesScoresInternalServerErrorResponseBody runs
// the validations defined on
// getAggregatedChallengesScores_InternalServerError_response_body
func ValidateGetAggregatedChallengesScoresInternalServerErrorResponseBody(body *GetAggregatedChallengesScoresInternalServerErrorResponseBody) (err error) {
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

// ValidateChallengesScoresResponse runs the validations defined on
// ChallengesScoresResponse
func ValidateChallengesScoresResponse(body *ChallengesScoresResponse) (err error) {
	if body.NodeID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("node_id", "body"))
	}
	if body.StorageChallengeScore == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("storage_challenge_score", "body"))
	}
	if body.HealthCheckChallengeScore == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("health_check_challenge_score", "body"))
	}
	return
}
