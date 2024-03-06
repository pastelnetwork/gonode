package design

import cors "goa.design/plugins/v3/cors/dsl"

import (

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports
)

var _ = Service("StorageChallenge", func() {
	Description("Storage Challenge service for to return storage-challenge related data")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/storage_challenges")
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

		Result(SCSummaryStatsRes)

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
		Description("Fetches storage-challenge reports")
		Meta("swagger:summary", "Fetches storage-challenge reports")

		Security(APIKeyAuth)

		Payload(func() {
			Attribute("pid", String, "PastelID of the user to fetch storage-challenge reports for", func() {
				Example("jXYJud3rm...")
			})

			Attribute("challenge_id", String, "ChallengeID of the storage challenge to fetch their logs", func() {
				Example("jXYJ")
			})

			APIKey("api_key", "key", String, func() {
				Description("Passphrase of the owner's PastelID")
				Example("Basic abcdef12345")
			})

			Required("pid", "key")
		})

		Result(ArrayOf(SCDetailedLogs))

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

// SCDetailedLogs defines the storage challenge message data type
var SCDetailedLogs = ResultType("application/vnd.storage.message", func() {
	Description("Storage challenge message data")
	Attributes(func() {
		Attribute("challenge_id", String, "ID of the challenge")
		Attribute("message_type", String, "type of the message")
		Attribute("sender_id", String, "ID of the sender's node")
		Attribute("sender_signature", String, "signature of the sender")
		Attribute("challenger_id", String, "ID of the challenger")
		Attribute("challenge", SCChallengeData, "Challenge data")
		Attribute("observers", ArrayOf(String), "List of observer IDs")
		Attribute("recipient_id", String, "ID of the recipient")
		Attribute("response", SCResponseData, "Response data")
		Attribute("challenger_evaluation", SCEvaluationData, "Challenger evaluation data")
		Attribute("observer_evaluation", SCObserverEvaluationData, "Observer evaluation data")
		Required("challenge_id", "message_type", "sender_id", "challenger_id", "observers", "recipient_id")
	})
})

// SCChallengeData defines the challenge data type
var SCChallengeData = Type("ChallengeData", func() {
	Description("Data of challenge")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("timestamp", String, "Timestamp")
	Attribute("file_hash", String, "File hash")
	Attribute("start_index", Int, "Start index")
	Attribute("end_index", Int, "End index")
	Required("timestamp", "file_hash", "start_index", "end_index")
})

// SCResponseData defines the response data type
var SCResponseData = Type("ResponseData", func() {
	Description("Data of response")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("hash", String, "Hash")
	Attribute("timestamp", String, "Timestamp")
	Required("timestamp")
})

// SCEvaluationData defines the evaluation data type
var SCEvaluationData = Type("EvaluationData", func() {
	Description("Data of evaluation")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("timestamp", String, "Timestamp")
	Attribute("hash", String, "Hash")
	Attribute("is_verified", Boolean, "IsVerified")
	Required("timestamp", "hash", "is_verified")
})

// SCObserverEvaluationData defines the observer evaluation data type
var SCObserverEvaluationData = Type("ObserverEvaluationData", func() {
	Description("Data of Observer's evaluation")
	Attribute("block", Int32, "Block")
	Attribute("merkelroot", String, "Merkelroot")
	Attribute("is_challenge_timestamp_ok", Boolean, "IsChallengeTimestampOK")
	Attribute("is_process_timestamp_ok", Boolean, "IsProcessTimestampOK")
	Attribute("is_evaluation_timestamp_ok", Boolean, "IsEvaluationTimestampOK")
	Attribute("is_recipient_signature_ok", Boolean, "IsRecipientSignatureOK")
	Attribute("is_challenger_signature_ok", Boolean, "IsChallengerSignatureOK")
	Attribute("is_evaluation_result_ok", Boolean, "IsEvaluationResultOK")
	Attribute("reason", String, "Reason")
	Attribute("true_hash", String, "TrueHash")
	Attribute("timestamp", String, "Timestamp")
	Required("timestamp", "is_challenge_timestamp_ok", "is_process_timestamp_ok", "is_evaluation_timestamp_ok", "is_recipient_signature_ok", "is_challenger_signature_ok", "is_evaluation_result_ok", "true_hash", "timestamp")
})

// SCSummaryStatsRes is the result type for the getSummaryStats method
var SCSummaryStatsRes = ResultType("application/vnd.summary_stats.result", func() {
	Description("Structure representing the metrics data")

	Attributes(func() {
		Attribute("sc_summary_stats", SCSummaryStats, "SCSummaryStats represents storage challenge summary of metrics stats")
	})

	Required("sc_summary_stats")
})

// SCSummaryStats is the result type for the storage-challenge summary stats
var SCSummaryStats = Type("SCSummaryStats", func() {
	Description("Storage-Challenge SummaryStats")

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
