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
	// Fetches self-healing reports
	GetDetailedLogs(context.Context, *GetDetailedLogsPayload) (res *SelfHealingReports, err error)
	// Fetches metrics data over a specified time range
	GetSummaryStats(context.Context, *GetSummaryStatsPayload) (res *MetricsResult, err error)
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
var MethodNames = [2]string{"getDetailedLogs", "getSummaryStats"}

type ChallengeTicket struct {
	TxID        *string
	TicketType  *string
	MissingKeys []string
	DataHash    []byte
	Recipient   *string
}

// GetDetailedLogsPayload is the payload type of the metrics service
// getDetailedLogs method.
type GetDetailedLogsPayload struct {
	// PastelID of the user to fetch self-healing reports for
	Pid string
	// Specific event ID to fetch reports for
	EventID *string
	// Number of reports to fetch
	Count *int
	// Passphrase of the owner's PastelID
	Key string
}

// GetSummaryStatsPayload is the payload type of the metrics service
// getSummaryStats method.
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

// MetricsResult is the result type of the metrics service getSummaryStats
// method.
type MetricsResult struct {
	// Self-healing trigger stats
	ShTriggerMetrics []*SHTriggerStats
	// Self-healing execution stats
	ShExecutionMetrics *SHExecutionStats
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

// Self-healing execution stats
type SHExecutionStats struct {
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
	// Total number of reconstructions required with hash mismatch
	TotalReconstructionRequiredHashMismatch *int
	// Total number of files healed
	TotalFilesHealed int
	// Total number of file healings that failed
	TotalFileHealingFailed int
}

// Self-healing trigger stats
type SHTriggerStats struct {
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

type SelfHealingReport struct {
	// Map of message type to SelfHealingMessages
	Messages []*SelfHealingMessageKV
}

type SelfHealingReportKV struct {
	// Challenge ID
	EventID *string
	// Self-healing report
	Report *SelfHealingReport
}

// SelfHealingReports is the result type of the metrics service getDetailedLogs
// method.
type SelfHealingReports struct {
	// Map of challenge ID to SelfHealingReport
	Reports []*SelfHealingReportKV
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
	res := &MetricsResult{}
	if vres.ShTriggerMetrics != nil {
		res.ShTriggerMetrics = make([]*SHTriggerStats, len(vres.ShTriggerMetrics))
		for i, val := range vres.ShTriggerMetrics {
			res.ShTriggerMetrics[i] = transformMetricsviewsSHTriggerStatsViewToSHTriggerStats(val)
		}
	}
	if vres.ShExecutionMetrics != nil {
		res.ShExecutionMetrics = transformMetricsviewsSHExecutionStatsViewToSHExecutionStats(vres.ShExecutionMetrics)
	}
	return res
}

// newMetricsResultView projects result type MetricsResult to projected type
// MetricsResultView using the "default" view.
func newMetricsResultView(res *MetricsResult) *metricsviews.MetricsResultView {
	vres := &metricsviews.MetricsResultView{}
	if res.ShTriggerMetrics != nil {
		vres.ShTriggerMetrics = make([]*metricsviews.SHTriggerStatsView, len(res.ShTriggerMetrics))
		for i, val := range res.ShTriggerMetrics {
			vres.ShTriggerMetrics[i] = transformSHTriggerStatsToMetricsviewsSHTriggerStatsView(val)
		}
	} else {
		vres.ShTriggerMetrics = []*metricsviews.SHTriggerStatsView{}
	}
	if res.ShExecutionMetrics != nil {
		vres.ShExecutionMetrics = transformSHExecutionStatsToMetricsviewsSHExecutionStatsView(res.ShExecutionMetrics)
	}
	return vres
}

// transformMetricsviewsSHTriggerStatsViewToSHTriggerStats builds a value of
// type *SHTriggerStats from a value of type *metricsviews.SHTriggerStatsView.
func transformMetricsviewsSHTriggerStatsViewToSHTriggerStats(v *metricsviews.SHTriggerStatsView) *SHTriggerStats {
	if v == nil {
		return nil
	}
	res := &SHTriggerStats{
		TriggerID:              *v.TriggerID,
		NodesOffline:           *v.NodesOffline,
		ListOfNodes:            *v.ListOfNodes,
		TotalFilesIdentified:   *v.TotalFilesIdentified,
		TotalTicketsIdentified: *v.TotalTicketsIdentified,
	}

	return res
}

// transformMetricsviewsSHExecutionStatsViewToSHExecutionStats builds a value
// of type *SHExecutionStats from a value of type
// *metricsviews.SHExecutionStatsView.
func transformMetricsviewsSHExecutionStatsViewToSHExecutionStats(v *metricsviews.SHExecutionStatsView) *SHExecutionStats {
	if v == nil {
		return nil
	}
	res := &SHExecutionStats{
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
		TotalReconstructionRequiredHashMismatch:               v.TotalReconstructionRequiredHashMismatch,
		TotalFilesHealed:                                      *v.TotalFilesHealed,
		TotalFileHealingFailed:                                *v.TotalFileHealingFailed,
	}

	return res
}

// transformSHTriggerStatsToMetricsviewsSHTriggerStatsView builds a value of
// type *metricsviews.SHTriggerStatsView from a value of type *SHTriggerStats.
func transformSHTriggerStatsToMetricsviewsSHTriggerStatsView(v *SHTriggerStats) *metricsviews.SHTriggerStatsView {
	res := &metricsviews.SHTriggerStatsView{
		TriggerID:              &v.TriggerID,
		NodesOffline:           &v.NodesOffline,
		ListOfNodes:            &v.ListOfNodes,
		TotalFilesIdentified:   &v.TotalFilesIdentified,
		TotalTicketsIdentified: &v.TotalTicketsIdentified,
	}

	return res
}

// transformSHExecutionStatsToMetricsviewsSHExecutionStatsView builds a value
// of type *metricsviews.SHExecutionStatsView from a value of type
// *SHExecutionStats.
func transformSHExecutionStatsToMetricsviewsSHExecutionStatsView(v *SHExecutionStats) *metricsviews.SHExecutionStatsView {
	res := &metricsviews.SHExecutionStatsView{
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
		TotalReconstructionRequiredHashMismatch:               v.TotalReconstructionRequiredHashMismatch,
		TotalFilesHealed:                                      &v.TotalFilesHealed,
		TotalFileHealingFailed:                                &v.TotalFileHealingFailed,
	}

	return res
}
