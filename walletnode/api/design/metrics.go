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
		Path("/metrics")
	})

	Error("Unauthorized", ErrorResult) // Assuming ErrorResult is defined in your design
	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("getMetrics", func() {
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

		Result(MetricsResult) // Define MetricsResult based on your metrics data structure

		HTTP(func() {
			GET("/")
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

var MetricsResult = ResultType("application/vnd.metrics.result", func() {
	Description("Structure representing the metrics data")

	Attributes(func() {
		Attribute("sc_metrics", Bytes, "SCMetrics represents serialized metrics data")
		Attribute("sh_trigger_metrics", ArrayOf(SHTriggerMetric), "Self-healing trigger metrics")
		Attribute("sh_execution_metrics", SHExecutionMetrics, "Self-healing execution metrics")
	})

	Required("sc_metrics", "sh_trigger_metrics", "sh_execution_metrics")
})

var SHTriggerMetric = Type("SHTriggerMetric", func() {
	Description("Self-healing trigger metric")

	Attribute("trigger_id", String, "Unique identifier for the trigger")
	Attribute("nodes_offline", Int, "Number of nodes offline")
	Attribute("list_of_nodes", String, "Comma-separated list of offline nodes")
	Attribute("total_files_identified", Int, "Total number of files identified for self-healing")
	Attribute("total_tickets_identified", Int, "Total number of tickets identified for self-healing")

	Required("trigger_id", "nodes_offline", "list_of_nodes", "total_files_identified", "total_tickets_identified")
})

var SHExecutionMetrics = Type("SHExecutionMetrics", func() {
	Description("Self-healing execution metrics")

	Attribute("total_challenges_issued", Int, "Total number of challenges issued")
	Attribute("total_challenges_rejected", Int, "Total number of challenges rejected")
	Attribute("total_challenges_accepted", Int, "Total number of challenges accepted")
	Attribute("total_challenges_failed", Int, "Total number of challenges failed")
	Attribute("total_challenges_successful", Int, "Total number of challenges successful")
	Attribute("total_files_healed", Int, "Total number of files healed")
	Attribute("total_file_healing_failed", Int, "Total number of file healings that failed")

	Required("total_challenges_issued", "total_challenges_rejected", "total_challenges_accepted", "total_challenges_failed", "total_challenges_successful", "total_files_healed", "total_file_healing_failed")
})
