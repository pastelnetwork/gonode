// Code generated by goa v3.15.0, DO NOT EDIT.
//
// HealthCheckChallenge views
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// HcSummaryStatsResult is the viewed result type that is projected based on a
// view.
type HcSummaryStatsResult struct {
	// Type to project
	Projected *HcSummaryStatsResultView
	// View to render
	View string
}

// HcSummaryStatsResultView is a type that runs validations on a projected type.
type HcSummaryStatsResultView struct {
	// HCSummaryStats represents health check challenge summary of metrics stats
	HcSummaryStats *HCSummaryStatsView
}

// HCSummaryStatsView is a type that runs validations on a projected type.
type HCSummaryStatsView struct {
	// Total number of challenges issued
	TotalChallengesIssued *int
	// Total number of challenges processed by the recipient node
	TotalChallengesProcessedByRecipient *int
	// Total number of challenges evaluated by the challenger node
	TotalChallengesEvaluatedByChallenger *int
	// Total number of challenges verified by observers
	TotalChallengesVerified *int
	// challenges failed due to slow-responses evaluated by observers
	NoOfSlowResponsesObservedByObservers *int
	// challenges failed due to invalid signatures evaluated by observers
	NoOfInvalidSignaturesObservedByObservers *int
	// challenges failed due to invalid evaluation evaluated by observers
	NoOfInvalidEvaluationObservedByObservers *int
}

var (
	// HcSummaryStatsResultMap is a map indexing the attribute names of
	// HcSummaryStatsResult by view name.
	HcSummaryStatsResultMap = map[string][]string{
		"default": {
			"hc_summary_stats",
		},
	}
)

// ValidateHcSummaryStatsResult runs the validations defined on the viewed
// result type HcSummaryStatsResult.
func ValidateHcSummaryStatsResult(result *HcSummaryStatsResult) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateHcSummaryStatsResultView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
	}
	return
}

// ValidateHcSummaryStatsResultView runs the validations defined on
// HcSummaryStatsResultView using the "default" view.
func ValidateHcSummaryStatsResultView(result *HcSummaryStatsResultView) (err error) {
	if result.HcSummaryStats == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hc_summary_stats", "result"))
	}
	if result.HcSummaryStats != nil {
		if err2 := ValidateHCSummaryStatsView(result.HcSummaryStats); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateHCSummaryStatsView runs the validations defined on
// HCSummaryStatsView.
func ValidateHCSummaryStatsView(result *HCSummaryStatsView) (err error) {
	if result.TotalChallengesIssued == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_issued", "result"))
	}
	if result.TotalChallengesProcessedByRecipient == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_processed_by_recipient", "result"))
	}
	if result.TotalChallengesEvaluatedByChallenger == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_evaluated_by_challenger", "result"))
	}
	if result.TotalChallengesVerified == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total_challenges_verified", "result"))
	}
	if result.NoOfSlowResponsesObservedByObservers == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("no_of_slow_responses_observed_by_observers", "result"))
	}
	if result.NoOfInvalidSignaturesObservedByObservers == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("no_of_invalid_signatures_observed_by_observers", "result"))
	}
	if result.NoOfInvalidEvaluationObservedByObservers == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("no_of_invalid_evaluation_observed_by_observers", "result"))
	}
	return
}
