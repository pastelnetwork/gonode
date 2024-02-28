package design

import (

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("metrics", func() {
	Description("Metrics service for fetching data over a specified time range")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/self_healing")
	})

	Error("Unauthorized", ErrorResult) // Assuming ErrorResult is defined in your design
	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("getDetailedLogs", func() {
		Description("Fetches self-healing reports")
		Meta("swagger:summary", "Fetches self-healing reports")

		Security(APIKeyAuth)

		Payload(func() {
			Attribute("pid", String, "PastelID of the user to fetch self-healing reports for", func() {
				Example("jXYJud3rm...")
			})
			Attribute("event_id", String, "Specific event ID to fetch reports for", func() {
				Example("event-123")
			})
			Attribute("count", Int, "Number of reports to fetch", func() {
				Example(10)
			})
			APIKey("api_key", "key", String, func() {
				Description("Passphrase of the owner's PastelID")
				Example("Basic abcdef12345")
			})

			Required("pid", "key")
		})

		Result(SelfHealingReports)

		HTTP(func() {
			GET("/detailed_logs")
			Param("pid")
			Param("event_id")
			Param("count")

			Response("Unauthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("getSummaryStats", func() {
		Description("Fetches metrics data over a specified time range")
		Meta("swagger:summary", "Fetches metrics data")

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

		Result(SummaryStats)

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
})

// SummaryStats is the result type for the getSummaryStats method
var SummaryStats = ResultType("application/vnd.metrics.result", func() {
	Description("Structure representing the metrics data")

	Attributes(func() {
		Attribute("sh_trigger_metrics", ArrayOf(SHTriggerStats), "Self-healing trigger stats")
		Attribute("sh_execution_metrics", SHExecutionStats, "Self-healing execution stats")
	})

	Required("sh_trigger_metrics", "sh_execution_metrics")
})

// SHTriggerStats is the result type for the self-healing trigger stats
var SHTriggerStats = Type("SHTriggerStats", func() {
	Description("Self-healing trigger stats")

	Attribute("trigger_id", String, "Unique identifier for the trigger")
	Attribute("nodes_offline", Int, "Number of nodes offline")
	Attribute("list_of_nodes", String, "Comma-separated list of offline nodes")
	Attribute("total_files_identified", Int, "Total number of files identified for self-healing")
	Attribute("total_tickets_identified", Int, "Total number of tickets identified for self-healing")

	Required("trigger_id", "nodes_offline", "list_of_nodes", "total_files_identified", "total_tickets_identified")
})

// SHExecutionStats is the result type for the self-healing execution stats
var SHExecutionStats = Type("SHExecutionStats", func() {
	Description("Self-healing execution stats")

	Attribute("total_challenges_issued", Int, "Total number of challenges issued")
	Attribute("total_challenges_acknowledged", Int, "Total number of challenges acknowledged by the healer node")
	Attribute("total_challenges_rejected", Int, "Total number of challenges rejected (healer node evaluated that reconstruction is not required)")
	Attribute("total_challenges_accepted", Int, "Total number of challenges accepted (healer node evaluated that reconstruction is required)")

	Attribute("total_challenge_evaluations_verified", Int, "Total number of challenges verified")
	Attribute("total_reconstruction_required_evaluations_approved", Int, "Total number of reconstructions approved by verifier nodes")
	Attribute("total_reconstruction_not_required_evaluations_approved", Int, "Total number of reconstructions not required approved by verifier nodes")
	Attribute("total_challenge_evaluations_unverified", Int, "Total number of challenge evaluations unverified by verifier nodes")
	Attribute("total_reconstruction_required_evaluations_not_approved", Int, "Total number of reconstructions not approved by verifier nodes")
	Attribute("total_reconstructions_not_required_evaluations_not_approved", Int, "Total number of reconstructions not required evaluation not approved by verifier nodes")
	Attribute("total_reconstruction_required_hash_mismatch", Int, "Total number of reconstructions required with hash mismatch")
	Attribute("total_files_healed", Int, "Total number of files healed")
	Attribute("total_file_healing_failed", Int, "Total number of file healings that failed")

	Required("total_challenges_issued", "total_challenges_acknowledged", "total_challenges_rejected",
		"total_challenges_accepted", "total_challenge_evaluations_verified", "total_reconstruction_required_evaluations_approved",
		"total_reconstruction_not_required_evaluations_approved", "total_challenge_evaluations_unverified",
		"total_reconstruction_required_evaluations_not_approved", "total_reconstructions_not_required_evaluations_not_approved",
		"total_files_healed", "total_file_healing_failed")
})

// SCMetrics is the result type for the storage-challenge metrics
var SCMetrics = Type("SCMetrics", func() {
	Description("Storage-Challenge Metrics")

	Attribute("total_challenges_issued", Int, "Total number of challenges issued")
	Attribute("total_challenges_processed", Int, "Total number of challenges processed by the recipient node")
	Attribute("total_challenges_verified_by_challenger", Int, "Total number of challenges verified by the challenger node")
	Attribute("total_challenges_verified_by_observers", Int, "Total number of challenges verified by observers")
	Attribute("slow_response_observed_by_observers", Int, "challenges failed due to slow-responses evaluated by observers")
	Attribute("invalid_signatures_observed_by_observers", Int, "challenges failed due to invalid signatures evaluated by observers")
	Attribute("invalid_evaluation_observed_by_observers", Int, "challenges failed due to invalid evaluation evaluated by observers")

	Required("total_challenges_issued", "total_challenges_processed", "total_challenges_verified_by_challenger",
		"total_challenges_verified_by_observers", "slow_response_observed_by_observers", "invalid_signatures_observed_by_observers",
		"invalid_evaluation_observed_by_observers")
})

// SelfHealingReports is the result type for the getSelfHealingReports method
var SelfHealingReports = Type("SelfHealingReports", func() {
	Description("Self-healing challenge reports")
	Attribute("reports", ArrayOf(SelfHealingReportKV), "Map of challenge ID to SelfHealingReport")
})

// SelfHealingReportKV is the result type for the self-healing challenge report
var SelfHealingReportKV = Type("SelfHealingReportKV", func() {
	Attribute("event_id", String, "Challenge ID")
	Attribute("report", SelfHealingReport, "Self-healing report")
})

// SelfHealingReport is the result type for the self-healing report
var SelfHealingReport = Type("SelfHealingReport", func() {
	Attribute("messages", ArrayOf(SelfHealingMessageKV), "Map of message type to SelfHealingMessages")
})

// SelfHealingMessageKV is the result type for the self-healing message
var SelfHealingMessageKV = Type("SelfHealingMessageKV", func() {
	Attribute("message_type", String, "Message type")
	Attribute("messages", ArrayOf(SelfHealingMessage), "Self-healing messages")
})

// SelfHealingMessage is the result type for the self-healing message
var SelfHealingMessage = Type("SelfHealingMessage", func() {
	Attribute("trigger_id", String)
	Attribute("message_type", String)
	Attribute("data", SelfHealingMessageData)
	Attribute("sender_id", String)
	Attribute("sender_signature", Bytes)
})

// SelfHealingMessageData is the result type for the self-healing message data
var SelfHealingMessageData = Type("SelfHealingMessageData", func() {
	Attribute("challenger_id", String)
	Attribute("recipient_id", String)
	Attribute("challenge", SelfHealingChallengeData)
	Attribute("response", SelfHealingResponseData)
	Attribute("verification", SelfHealingVerificationData)
})

// SelfHealingChallengeData is the result type for the self-healing challenge data
var SelfHealingChallengeData = Type("SelfHealingChallengeData", func() {
	Attribute("block", Int32)
	Attribute("merkelroot", String)
	Attribute("timestamp", String) // Goa does not directly support time.Time, use string and format as RFC3339
	Attribute("challenge_tickets", ArrayOf(ChallengeTicket))
	Attribute("nodes_on_watchlist", String)
})

// ChallengeTicket is the result type for the challenge ticket
var ChallengeTicket = Type("ChallengeTicket", func() {
	Attribute("tx_id", String)
	Attribute("ticket_type", String) // Assuming TicketType is an enum or similar in Go, represented as String here
	Attribute("missing_keys", ArrayOf(String))
	Attribute("data_hash", Bytes)
	Attribute("recipient", String)
})

// SelfHealingResponseData is the result type for the self-healing response data
var SelfHealingResponseData = Type("SelfHealingResponseData", func() {
	Attribute("challenge_id", String)
	Attribute("block", Int32)
	Attribute("merkelroot", String)
	Attribute("timestamp", String) // Use string for time.Time
	Attribute("responded_ticket", RespondedTicket)
	Attribute("verifiers", ArrayOf(String))
})

// RespondedTicket is the result type for the responded ticket
var RespondedTicket = Type("RespondedTicket", func() {
	Attribute("tx_id", String)
	Attribute("ticket_type", String) // Assuming TicketType is an enum or similar in Go
	Attribute("missing_keys", ArrayOf(String))
	Attribute("reconstructed_file_hash", Bytes)
	Attribute("sense_file_ids", ArrayOf(String))
	Attribute("raptor_q_symbols", Bytes)
	Attribute("is_reconstruction_required", Boolean)
})

// SelfHealingVerificationData is the result type for the self-healing verification data
var SelfHealingVerificationData = Type("SelfHealingVerificationData", func() {
	Attribute("challenge_id", String)
	Attribute("block", Int32)
	Attribute("merkelroot", String)
	Attribute("timestamp", String) // Use string for time.Time
	Attribute("verified_ticket", VerifiedTicket)
	Attribute("verifiers_data", MapOf(String, Bytes)) // Goa supports MapOf for simple key-value pairs
})

// VerifiedTicket is the result type for the verified ticket
var VerifiedTicket = Type("VerifiedTicket", func() {
	Attribute("tx_id", String)
	Attribute("ticket_type", String)
	Attribute("missing_keys", ArrayOf(String))
	Attribute("reconstructed_file_hash", Bytes)
	Attribute("is_reconstruction_required", Boolean)
	Attribute("raptor_q_symbols", Bytes)
	Attribute("sense_file_ids", ArrayOf(String))
	Attribute("is_verified", Boolean)
	Attribute("message", String)
})
