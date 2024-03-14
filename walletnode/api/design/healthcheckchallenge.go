package design

import (
	cors "goa.design/plugins/v3/cors/dsl"

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
)

//revive:enable:dot-imports

var _ = Service("HealthCheckChallenge", func() {
	Description("HealthCheck Challenge service for to return health check related data")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/healthcheck_challenge")
	})

	Error("Unauthorized", ErrorResult) // Assuming ErrorResult is defined in your design
	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("getSummaryStats", func() {
		Description("Fetches summary stats data over a specified time range")
		Meta("swagger:summary", "Fetches summary stats")

		Security(APIKeyAuth)

		Payload(func() {
			Attribute("from", String, func() {
				Description("Start time for the metrics data range")
				Format(FormatDateTime)
				Example("2023-01-01T00:00:00Z")
			})
			Attribute("to", String, func() {
				Description("End time for the metrics data range")
				Format(FormatDateTime)
				Example("2023-01-02T00:00:00Z")
			})
			Attribute("pid", String, func() {
				Description("PastelID of the user to fetch metrics for")
				Example("jXYJud3rm...")
			})

			APIKey("api_key", "key", String, func() {
				Description("Passphrase of the owner's PastelID")
				Example("Basic abcdef12345")
			})

			Required("pid", "key")
		})

		Result(HCSummaryStatsRes)

		HTTP(func() {
			GET("/summary_stats")
			Param("from")
			Param("to")
			Param("pid")

			Response("Unauthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("getDetailedLogs", func() {
		Description("Fetches health-check-challenge reports")
		Meta("swagger:summary", "Fetches health-check-challenge reports")

		Security(APIKeyAuth)

		Payload(func() {
			Attribute("pid", String, "PastelID of the user to fetch health-check-challenge reports for", func() {
				Example("jXYJud3rm...")
			})

			Attribute("challenge_id", String, "ChallengeID of the health check challenge to fetch their logs", func() {
				Example("jXYJ")
			})

			APIKey("api_key", "key", String, func() {
				Description("Passphrase of the owner's PastelID")
				Example("Basic abcdef12345")
			})

			Required("pid", "key")
		})

		Result(ArrayOf(HCDetailedLogs))

		HTTP(func() {
			GET("/detailed_logs")
			Param("pid")
			Param("challenge_id")

			Response("Unauthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})
})

// HCSummaryStatsRes is the result type for the getSummaryStats method
var HCSummaryStatsRes = ResultType("application/vnd.hc_summary_stats.result", func() {
	Description("Structure representing the metrics data")

	Attributes(func() {
		Attribute("hc_summary_stats", HCSummaryStats, "HCSummaryStats represents health check challenge summary of metrics stats")
	})

	Required("hc_summary_stats")
})

// HCSummaryStats is the result type for the healthcheck-challenge summary stats
var HCSummaryStats = Type("HCSummaryStats", func() {
	Description("HealthCheck-Challenge SummaryStats")

	Attribute("total_challenges_issued", Int, "Total number of challenges issued")
	Attribute("total_challenges_processed_by_recipient", Int, "Total number of challenges processed by the recipient node")
	Attribute("total_challenges_evaluated_by_challenger", Int, "Total number of challenges evaluated by the challenger node")
	Attribute("total_challenges_verified", Int, "Total number of challenges verified by observers")
	Attribute("no_of_slow_responses_observed_by_observers", Int, "challenges failed due to slow-responses evaluated by observers")
	Attribute("no_of_invalid_signatures_observed_by_observers", Int, "challenges failed due to invalid signatures evaluated by observers")
	Attribute("no_of_invalid_evaluation_observed_by_observers", Int, "challenges failed due to invalid evaluation evaluated by observers")

	Required("total_challenges_issued", "total_challenges_processed_by_recipient", "total_challenges_evaluated_by_challenger",
		"total_challenges_verified", "no_of_slow_responses_observed_by_observers", "no_of_invalid_signatures_observed_by_observers",
		"no_of_invalid_evaluation_observed_by_observers")
})

// HCDetailedLogs defines the health check challenge message data type
var HCDetailedLogs = ResultType("application/vnd.hc_detailed_logs.message", func() {
	Description("HealthCheck challenge message data")
	Attributes(func() {
		Attribute("challenge_id", String, "ID of the challenge")
		Attribute("message_type", String, "type of the message")
		Attribute("sender_id", String, "ID of the sender's node")
		Attribute("sender_signature", String, "signature of the sender")
		Attribute("challenger_id", String, "ID of the challenger")
		Attribute("challenge", HCChallengeData, "Challenge data")
		Attribute("observers", ArrayOf(String), "List of observer IDs")
		Attribute("recipient_id", String, "ID of the recipient")
		Attribute("response", HCResponseData, "Response data")
		Attribute("challenger_evaluation", HCEvaluationData, "Challenger evaluation data")
		Attribute("observer_evaluation", HCObserverEvaluationData, "Observer evaluation data")
		Required("challenge_id", "message_type", "sender_id", "challenger_id", "observers", "recipient_id")
	})
})

// HCChallengeData defines the challenge data type
var HCChallengeData = Type("HCChallengeData", func() {
	Description("Data of challenge")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("timestamp", String, "Timestamp")
	Required("timestamp")
})

// HCResponseData defines the response data type
var HCResponseData = Type("HCResponseData", func() {
	Description("Data of response")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("timestamp", String, "Timestamp")
	Required("timestamp")
})

// HCEvaluationData defines the evaluation data type
var HCEvaluationData = Type("HCEvaluationData", func() {
	Description("Data of evaluation")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("timestamp", String, "Timestamp")
	Attribute("is_verified", Boolean, "IsVerified")
	Required("timestamp", "is_verified")
})

// HCObserverEvaluationData defines the observer evaluation data type
var HCObserverEvaluationData = Type("HCObserverEvaluationData", func() {
	Description("Data of Observer's evaluation")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("is_challenge_timestamp_ok", Boolean, "IsChallengeTimestampOK")
	Attribute("is_process_timestamp_ok", Boolean, "IsProcessTimestampOK")
	Attribute("is_evaluation_timestamp_ok", Boolean, "IsEvaluationTimestampOK")
	Attribute("is_recipient_signature_ok", Boolean, "IsRecipientSignatureOK")
	Attribute("is_challenger_signature_ok", Boolean, "IsChallengerSignatureOK")
	Attribute("is_evaluation_result_ok", Boolean, "IsEvaluationResultOK")
	Attribute("timestamp", String, "Timestamp")
	Required("timestamp", "is_challenge_timestamp_ok", "is_process_timestamp_ok", "is_evaluation_timestamp_ok", "is_recipient_signature_ok", "is_challenger_signature_ok", "is_evaluation_result_ok", "timestamp")
})
