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
	Attribute("total_challenges_processed", Int, "Total number of challenges processed by the recipient node")
	Attribute("total_challenges_verified_by_challenger", Int, "Total number of challenges verified by the challenger node")
	Attribute("total_challenges_verified_by_observers", Int, "Total number of challenges verified by observers")
	Attribute("no_of_slow_responses_observed_by_observers", Int, "challenges failed due to slow-responses evaluated by observers")
	Attribute("no_of_invalid_signatures_observed_by_observers", Int, "challenges failed due to invalid signatures evaluated by observers")
	Attribute("no_of_invalid_evaluation_observed_by_observers", Int, "challenges failed due to invalid evaluation evaluated by observers")

	Required("total_challenges_issued", "total_challenges_processed", "total_challenges_verified_by_challenger",
		"total_challenges_verified_by_observers", "no_of_slow_responses_observed_by_observers", "no_of_invalid_signatures_observed_by_observers",
		"no_of_invalid_evaluation_observed_by_observers")
})
