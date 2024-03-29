// Code generated by goa v3.15.0, DO NOT EDIT.
//
// HealthCheckChallenge service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package healthcheckchallenge

import (
	"context"

	healthcheckchallengeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/health_check_challenge/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// HealthCheck Challenge service for to return health check related data
type Service interface {
	// Fetches summary stats data over a specified time range
	GetSummaryStats(context.Context, *GetSummaryStatsPayload) (res *HcSummaryStatsResult, err error)
	// Fetches health-check-challenge reports
	GetDetailedLogs(context.Context, *GetDetailedLogsPayload) (res []*HcDetailedLogsMessage, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// APIName is the name of the API as defined in the design.
const APIName = "walletnode"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "1.0"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "HealthCheckChallenge"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"getSummaryStats", "getDetailedLogs"}

// GetDetailedLogsPayload is the payload type of the HealthCheckChallenge
// service getDetailedLogs method.
type GetDetailedLogsPayload struct {
	// PastelID of the user to fetch health-check-challenge reports for
	Pid string
	// ChallengeID of the health check challenge to fetch their logs
	ChallengeID *string
	// Passphrase of the owner's PastelID
	Key string
}

// GetSummaryStatsPayload is the payload type of the HealthCheckChallenge
// service getSummaryStats method.
type GetSummaryStatsPayload struct {
	// Start time for the metrics data range
	From *string
	// End time for the metrics data range
	To *string
	// PastelID of the user to fetch metrics for
	Pid string
	// Passphrase of the owner's PastelID
	Key string
}

// Data of challenge
type HCChallengeData struct {
	// Block
	Block *int32
	// Merkelroot
	Merkelroot *string
	// Timestamp
	Timestamp string
}

// Data of evaluation
type HCEvaluationData struct {
	// Block
	Block *int32
	// Merkelroot
	Merkelroot *string
	// Timestamp
	Timestamp string
	// IsVerified
	IsVerified bool
}

// Data of Observer's evaluation
type HCObserverEvaluationData struct {
	// Block
	Block *int32
	// Merkelroot
	Merkelroot *string
	// IsChallengeTimestampOK
	IsChallengeTimestampOK bool
	// IsProcessTimestampOK
	IsProcessTimestampOK bool
	// IsEvaluationTimestampOK
	IsEvaluationTimestampOK bool
	// IsRecipientSignatureOK
	IsRecipientSignatureOK bool
	// IsChallengerSignatureOK
	IsChallengerSignatureOK bool
	// IsEvaluationResultOK
	IsEvaluationResultOK bool
	// Timestamp
	Timestamp string
}

// Data of response
type HCResponseData struct {
	// Block
	Block *int32
	// Merkelroot
	Merkelroot *string
	// Timestamp
	Timestamp string
}

// HealthCheck-Challenge SummaryStats
type HCSummaryStats struct {
	// Total number of challenges issued
	TotalChallengesIssued int
	// Total number of challenges processed by the recipient node
	TotalChallengesProcessedByRecipient int
	// Total number of challenges evaluated by the challenger node
	TotalChallengesEvaluatedByChallenger int
	// Total number of challenges verified by observers
	TotalChallengesVerified int
	// challenges failed due to slow-responses evaluated by observers
	NoOfSlowResponsesObservedByObservers int
	// challenges failed due to invalid signatures evaluated by observers
	NoOfInvalidSignaturesObservedByObservers int
	// challenges failed due to invalid evaluation evaluated by observers
	NoOfInvalidEvaluationObservedByObservers int
}

// HealthCheck challenge message data
type HcDetailedLogsMessage struct {
	// ID of the challenge
	ChallengeID string
	// type of the message
	MessageType string
	// ID of the sender's node
	SenderID string
	// signature of the sender
	SenderSignature *string
	// ID of the challenger
	ChallengerID string
	// Challenge data
	Challenge *HCChallengeData
	// List of observer IDs
	Observers []string
	// ID of the recipient
	RecipientID string
	// Response data
	Response *HCResponseData
	// Challenger evaluation data
	ChallengerEvaluation *HCEvaluationData
	// Observer evaluation data
	ObserverEvaluation *HCObserverEvaluationData
}

// HcSummaryStatsResult is the result type of the HealthCheckChallenge service
// getSummaryStats method.
type HcSummaryStatsResult struct {
	// HCSummaryStats represents health check challenge summary of metrics stats
	HcSummaryStats *HCSummaryStats
}

// MakeUnauthorized builds a goa.ServiceError from an error.
func MakeUnauthorized(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "Unauthorized", false, false, false)
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "BadRequest", false, false, false)
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "NotFound", false, false, false)
}

// MakeInternalServerError builds a goa.ServiceError from an error.
func MakeInternalServerError(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "InternalServerError", false, false, false)
}

// NewHcSummaryStatsResult initializes result type HcSummaryStatsResult from
// viewed result type HcSummaryStatsResult.
func NewHcSummaryStatsResult(vres *healthcheckchallengeviews.HcSummaryStatsResult) *HcSummaryStatsResult {
	return newHcSummaryStatsResult(vres.Projected)
}

// NewViewedHcSummaryStatsResult initializes viewed result type
// HcSummaryStatsResult from result type HcSummaryStatsResult using the given
// view.
func NewViewedHcSummaryStatsResult(res *HcSummaryStatsResult, view string) *healthcheckchallengeviews.HcSummaryStatsResult {
	p := newHcSummaryStatsResultView(res)
	return &healthcheckchallengeviews.HcSummaryStatsResult{Projected: p, View: "default"}
}

// newHcSummaryStatsResult converts projected type HcSummaryStatsResult to
// service type HcSummaryStatsResult.
func newHcSummaryStatsResult(vres *healthcheckchallengeviews.HcSummaryStatsResultView) *HcSummaryStatsResult {
	res := &HcSummaryStatsResult{}
	if vres.HcSummaryStats != nil {
		res.HcSummaryStats = transformHealthcheckchallengeviewsHCSummaryStatsViewToHCSummaryStats(vres.HcSummaryStats)
	}
	return res
}

// newHcSummaryStatsResultView projects result type HcSummaryStatsResult to
// projected type HcSummaryStatsResultView using the "default" view.
func newHcSummaryStatsResultView(res *HcSummaryStatsResult) *healthcheckchallengeviews.HcSummaryStatsResultView {
	vres := &healthcheckchallengeviews.HcSummaryStatsResultView{}
	if res.HcSummaryStats != nil {
		vres.HcSummaryStats = transformHCSummaryStatsToHealthcheckchallengeviewsHCSummaryStatsView(res.HcSummaryStats)
	}
	return vres
}

// transformHealthcheckchallengeviewsHCSummaryStatsViewToHCSummaryStats builds
// a value of type *HCSummaryStats from a value of type
// *healthcheckchallengeviews.HCSummaryStatsView.
func transformHealthcheckchallengeviewsHCSummaryStatsViewToHCSummaryStats(v *healthcheckchallengeviews.HCSummaryStatsView) *HCSummaryStats {
	if v == nil {
		return nil
	}
	res := &HCSummaryStats{
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

// transformHCSummaryStatsToHealthcheckchallengeviewsHCSummaryStatsView builds
// a value of type *healthcheckchallengeviews.HCSummaryStatsView from a value
// of type *HCSummaryStats.
func transformHCSummaryStatsToHealthcheckchallengeviewsHCSummaryStatsView(v *HCSummaryStats) *healthcheckchallengeviews.HCSummaryStatsView {
	res := &healthcheckchallengeviews.HCSummaryStatsView{
		TotalChallengesIssued:                    &v.TotalChallengesIssued,
		TotalChallengesProcessedByRecipient:      &v.TotalChallengesProcessedByRecipient,
		TotalChallengesEvaluatedByChallenger:     &v.TotalChallengesEvaluatedByChallenger,
		TotalChallengesVerified:                  &v.TotalChallengesVerified,
		NoOfSlowResponsesObservedByObservers:     &v.NoOfSlowResponsesObservedByObservers,
		NoOfInvalidSignaturesObservedByObservers: &v.NoOfInvalidSignaturesObservedByObservers,
		NoOfInvalidEvaluationObservedByObservers: &v.NoOfInvalidEvaluationObservedByObservers,
	}

	return res
}
