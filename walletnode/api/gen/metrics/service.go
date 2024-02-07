// Code generated by goa v3.14.0, DO NOT EDIT.
//
// metrics service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package metrics

import (
	"context"

	metricsviews "github.com/pastelnetwork/gonode/walletnode/api/gen/metrics/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Metrics service for fetching data over a specified time range
type Service interface {
	// Fetches self-healing challenge reports
	GetChallengeReports(context.Context, *GetChallengeReportsPayload) (res *SelfHealingChallengeReports, err error)
	// Fetches metrics data over a specified time range
	GetMetrics(context.Context, *GetMetricsPayload) (res *MetricsResult, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "metrics"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"getChallengeReports", "getMetrics"}

type ChallengeTicket struct {
	TxID        *string
	TicketType  *string
	MissingKeys []string
	DataHash    []byte
	Recipient   *string
}

// GetChallengeReportsPayload is the payload type of the metrics service
// getChallengeReports method.
type GetChallengeReportsPayload struct {
	// PastelID of the user to fetch challenge reports for
	Pid string
	// Specific challenge ID to fetch reports for
	ChallengeID *string
	// Number of reports to fetch
	Count *int
	// Passphrase of the owner's PastelID
	Key string
}

// GetMetricsPayload is the payload type of the metrics service getMetrics
// method.
type GetMetricsPayload struct {
	// Start time for the metrics data range
	From *string
	// End time for the metrics data range
	To *string
	// PastelID of the user to fetch metrics for
	Pid string
	// Passphrase of the owner's PastelID
	Key string
}

// MetricsResult is the result type of the metrics service getMetrics method.
type MetricsResult struct {
	// SCMetrics represents serialized metrics data
	ScMetrics []byte
	// Self-healing trigger metrics
	ShTriggerMetrics []*SHTriggerMetric
	// Self-healing execution metrics
	ShExecutionMetrics *SHExecutionMetrics
}

type RespondedTicket struct {
	TxID                     *string
	TicketType               *string
	MissingKeys              []string
	ReconstructedFileHash    []byte
	SenseFileIds             []string
	RaptorQSymbols           []byte
	IsReconstructionRequired *bool
}

// Self-healing execution metrics
type SHExecutionMetrics struct {
	// Total number of challenges issued
	TotalChallengesIssued int
	// Total number of challenges acknowledged by the healer node
	TotalChallengesAcknowledged int
	// Total number of challenges rejected (healer node evaluated that
	// reconstruction is not required)
	TotalChallengesRejected int
	// Total number of challenges accepted (healer node evaluated that
	// reconstruction is required)
	TotalChallengesAccepted int
	// Total number of challenges verified
	TotalChallengeEvaluationsVerified int
	// Total number of reconstructions approved by verifier nodes
	TotalReconstructionRequiredEvaluationsApproved int
	// Total number of reconstructions not required approved by verifier nodes
	TotalReconstructionNotRequiredEvaluationsApproved int
	// Total number of challenge evaluations unverified by verifier nodes
	TotalChallengeEvaluationsUnverified int
	// Total number of reconstructions not approved by verifier nodes
	TotalReconstructionRequiredEvaluationsNotApproved int
	// Total number of reconstructions not required evaluation not approved by
	// verifier nodes
	TotalReconstructionsNotRequiredEvaluationsNotApproved int
	// Total number of files healed
	TotalFilesHealed int
	// Total number of file healings that failed
	TotalFileHealingFailed int
}

// Self-healing trigger metric
type SHTriggerMetric struct {
	// Unique identifier for the trigger
	TriggerID string
	// Number of nodes offline
	NodesOffline int
	// Comma-separated list of offline nodes
	ListOfNodes string
	// Total number of files identified for self-healing
	TotalFilesIdentified int
	// Total number of tickets identified for self-healing
	TotalTicketsIdentified int
}

type SelfHealingChallengeData struct {
	Block            *int32
	Merkelroot       *string
	Timestamp        *string
	ChallengeTickets []*ChallengeTicket
	NodesOnWatchlist *string
}

type SelfHealingChallengeReport struct {
	// Map of message type to SelfHealingMessages
	Messages []*SelfHealingMessageKV
}

type SelfHealingChallengeReportKV struct {
	// Challenge ID
	ChallengeID *string
	// Self-healing challenge report
	Report *SelfHealingChallengeReport
}

// SelfHealingChallengeReports is the result type of the metrics service
// getChallengeReports method.
type SelfHealingChallengeReports struct {
	// Map of challenge ID to SelfHealingChallengeReport
	Reports []*SelfHealingChallengeReportKV
}

type SelfHealingMessage struct {
	TriggerID       *string
	MessageType     *string
	Data            *SelfHealingMessageData
	SenderID        *string
	SenderSignature []byte
}

type SelfHealingMessageData struct {
	ChallengerID *string
	RecipientID  *string
	Challenge    *SelfHealingChallengeData
	Response     *SelfHealingResponseData
	Verification *SelfHealingVerificationData
}

type SelfHealingMessageKV struct {
	// Message type
	MessageType *string
	// Self-healing messages
	Messages []*SelfHealingMessage
}

type SelfHealingResponseData struct {
	ChallengeID     *string
	Block           *int32
	Merkelroot      *string
	Timestamp       *string
	RespondedTicket *RespondedTicket
	Verifiers       []string
}

type SelfHealingVerificationData struct {
	ChallengeID    *string
	Block          *int32
	Merkelroot     *string
	Timestamp      *string
	VerifiedTicket *VerifiedTicket
	VerifiersData  map[string][]byte
}

type VerifiedTicket struct {
	TxID                     *string
	TicketType               *string
	MissingKeys              []string
	ReconstructedFileHash    []byte
	IsReconstructionRequired *bool
	RaptorQSymbols           []byte
	SenseFileIds             []string
	IsVerified               *bool
	Message                  *string
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

// NewMetricsResult initializes result type MetricsResult from viewed result
// type MetricsResult.
func NewMetricsResult(vres *metricsviews.MetricsResult) *MetricsResult {
	return newMetricsResult(vres.Projected)
}

// NewViewedMetricsResult initializes viewed result type MetricsResult from
// result type MetricsResult using the given view.
func NewViewedMetricsResult(res *MetricsResult, view string) *metricsviews.MetricsResult {
	p := newMetricsResultView(res)
	return &metricsviews.MetricsResult{Projected: p, View: "default"}
}

// newMetricsResult converts projected type MetricsResult to service type
// MetricsResult.
func newMetricsResult(vres *metricsviews.MetricsResultView) *MetricsResult {
	res := &MetricsResult{
		ScMetrics: vres.ScMetrics,
	}
	if vres.ShTriggerMetrics != nil {
		res.ShTriggerMetrics = make([]*SHTriggerMetric, len(vres.ShTriggerMetrics))
		for i, val := range vres.ShTriggerMetrics {
			res.ShTriggerMetrics[i] = transformMetricsviewsSHTriggerMetricViewToSHTriggerMetric(val)
		}
	}
	if vres.ShExecutionMetrics != nil {
		res.ShExecutionMetrics = transformMetricsviewsSHExecutionMetricsViewToSHExecutionMetrics(vres.ShExecutionMetrics)
	}
	return res
}

// newMetricsResultView projects result type MetricsResult to projected type
// MetricsResultView using the "default" view.
func newMetricsResultView(res *MetricsResult) *metricsviews.MetricsResultView {
	vres := &metricsviews.MetricsResultView{
		ScMetrics: res.ScMetrics,
	}
	if res.ShTriggerMetrics != nil {
		vres.ShTriggerMetrics = make([]*metricsviews.SHTriggerMetricView, len(res.ShTriggerMetrics))
		for i, val := range res.ShTriggerMetrics {
			vres.ShTriggerMetrics[i] = transformSHTriggerMetricToMetricsviewsSHTriggerMetricView(val)
		}
	} else {
		vres.ShTriggerMetrics = []*metricsviews.SHTriggerMetricView{}
	}
	if res.ShExecutionMetrics != nil {
		vres.ShExecutionMetrics = transformSHExecutionMetricsToMetricsviewsSHExecutionMetricsView(res.ShExecutionMetrics)
	}
	return vres
}

// transformMetricsviewsSHTriggerMetricViewToSHTriggerMetric builds a value of
// type *SHTriggerMetric from a value of type *metricsviews.SHTriggerMetricView.
func transformMetricsviewsSHTriggerMetricViewToSHTriggerMetric(v *metricsviews.SHTriggerMetricView) *SHTriggerMetric {
	if v == nil {
		return nil
	}
	res := &SHTriggerMetric{
		TriggerID:              *v.TriggerID,
		NodesOffline:           *v.NodesOffline,
		ListOfNodes:            *v.ListOfNodes,
		TotalFilesIdentified:   *v.TotalFilesIdentified,
		TotalTicketsIdentified: *v.TotalTicketsIdentified,
	}

	return res
}

// transformMetricsviewsSHExecutionMetricsViewToSHExecutionMetrics builds a
// value of type *SHExecutionMetrics from a value of type
// *metricsviews.SHExecutionMetricsView.
func transformMetricsviewsSHExecutionMetricsViewToSHExecutionMetrics(v *metricsviews.SHExecutionMetricsView) *SHExecutionMetrics {
	if v == nil {
		return nil
	}
	res := &SHExecutionMetrics{
		TotalChallengesIssued:                                 *v.TotalChallengesIssued,
		TotalChallengesAcknowledged:                           *v.TotalChallengesAcknowledged,
		TotalChallengesRejected:                               *v.TotalChallengesRejected,
		TotalChallengesAccepted:                               *v.TotalChallengesAccepted,
		TotalChallengeEvaluationsVerified:                     *v.TotalChallengeEvaluationsVerified,
		TotalReconstructionRequiredEvaluationsApproved:        *v.TotalReconstructionRequiredEvaluationsApproved,
		TotalReconstructionNotRequiredEvaluationsApproved:     *v.TotalReconstructionNotRequiredEvaluationsApproved,
		TotalChallengeEvaluationsUnverified:                   *v.TotalChallengeEvaluationsUnverified,
		TotalReconstructionRequiredEvaluationsNotApproved:     *v.TotalReconstructionRequiredEvaluationsNotApproved,
		TotalReconstructionsNotRequiredEvaluationsNotApproved: *v.TotalReconstructionsNotRequiredEvaluationsNotApproved,
		TotalFilesHealed:                                      *v.TotalFilesHealed,
		TotalFileHealingFailed:                                *v.TotalFileHealingFailed,
	}

	return res
}

// transformSHTriggerMetricToMetricsviewsSHTriggerMetricView builds a value of
// type *metricsviews.SHTriggerMetricView from a value of type *SHTriggerMetric.
func transformSHTriggerMetricToMetricsviewsSHTriggerMetricView(v *SHTriggerMetric) *metricsviews.SHTriggerMetricView {
	res := &metricsviews.SHTriggerMetricView{
		TriggerID:              &v.TriggerID,
		NodesOffline:           &v.NodesOffline,
		ListOfNodes:            &v.ListOfNodes,
		TotalFilesIdentified:   &v.TotalFilesIdentified,
		TotalTicketsIdentified: &v.TotalTicketsIdentified,
	}

	return res
}

// transformSHExecutionMetricsToMetricsviewsSHExecutionMetricsView builds a
// value of type *metricsviews.SHExecutionMetricsView from a value of type
// *SHExecutionMetrics.
func transformSHExecutionMetricsToMetricsviewsSHExecutionMetricsView(v *SHExecutionMetrics) *metricsviews.SHExecutionMetricsView {
	res := &metricsviews.SHExecutionMetricsView{
		TotalChallengesIssued:                                 &v.TotalChallengesIssued,
		TotalChallengesAcknowledged:                           &v.TotalChallengesAcknowledged,
		TotalChallengesRejected:                               &v.TotalChallengesRejected,
		TotalChallengesAccepted:                               &v.TotalChallengesAccepted,
		TotalChallengeEvaluationsVerified:                     &v.TotalChallengeEvaluationsVerified,
		TotalReconstructionRequiredEvaluationsApproved:        &v.TotalReconstructionRequiredEvaluationsApproved,
		TotalReconstructionNotRequiredEvaluationsApproved:     &v.TotalReconstructionNotRequiredEvaluationsApproved,
		TotalChallengeEvaluationsUnverified:                   &v.TotalChallengeEvaluationsUnverified,
		TotalReconstructionRequiredEvaluationsNotApproved:     &v.TotalReconstructionRequiredEvaluationsNotApproved,
		TotalReconstructionsNotRequiredEvaluationsNotApproved: &v.TotalReconstructionsNotRequiredEvaluationsNotApproved,
		TotalFilesHealed:                                      &v.TotalFilesHealed,
		TotalFileHealingFailed:                                &v.TotalFileHealingFailed,
	}

	return res
}
