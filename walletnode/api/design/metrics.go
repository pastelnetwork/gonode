package design

import (
	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("metrics", func() {
	Description("OpenAPI Metrics service")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/openapi/metrics")
	})

	Error("UnAuthorized", ErrorResult)
	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("selfHealing", func() {
		Description("returns the self-healing metrics")
		Meta("swagger:summary", "Returns Self Healing Metrics")

		Payload(MetricsPayload)
		Security(APIKeyAuth)
		Result(SelfHealingMetricsResult)

		HTTP(func() {
			GET("/selfhealing")
			Params(func() {
				Param("pid")
			})
			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})
})

// MetricsPayload is payload to get metrics for self-healing & storage challenges
var MetricsPayload = Type("MetricsPayload", func() {
	Attribute("pid", String, func() {
		Meta("struct:field:name", "Pid")
		Description("Owner's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	APIKey("api_key", "key", String, func() {
		Description("Passphrase of the owner's PastelID")
		Example("Basic abcdef12345")
	})
	Required("pid", "key")
})

// SelfHealingMetricsResult is self-healing metrics result
var SelfHealingMetricsResult = ResultType("application/vnd.selfhealing.metrics", func() {
	TypeName("SelfHealingMetrics")
	Attributes(func() {
		Attribute("total_tickets_send_for_self_healing", Int, func() {
			Description("Tickets send for self-healing")
			Example(1000)
		})
		Attribute("estimated_missing_keys", Int, func() {
			Description("Total estimated missing keys")
			Example(1000)
		})

		Attribute("tickets_required_self_healing", Int, func() {
			Description("Tickets required self healing")
			Example(100)
		})

		Attribute("tickets_self_healed_successfully", Int, func() {
			Description("Tickets self-healed successfully")
			Example(100)
		})

		Attribute("tickets_verified_successfully", Int, func() {
			Description("Tickets verified successfully")
			Example(100)
		})
	})
})
