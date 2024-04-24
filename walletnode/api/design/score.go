package design

import (
	cors "goa.design/plugins/v3/cors/dsl"

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
)

//revive:enable:dot-imports

var _ = Service("Score", func() {
	Description("Score service for return score related to challenges")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/nodes")
	})

	Error("Unauthorized", ErrorResult) // Assuming ErrorResult is defined in your design
	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("getAggregatedChallengesScores", func() {
		Description("Fetches aggregated challenges score for SC and HC")
		Meta("swagger:summary", "Fetches aggregated challenges score for sc and hc")

		Security(APIKeyAuth)

		Payload(func() {
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

		Result(ArrayOf(ChallengesScores))

		HTTP(func() {
			GET("/challenges_score")
			Param("pid")

			Response("Unauthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})
})

// ChallengesScores is the combined accumulated score of HC and SC challenges
var ChallengesScores = Type("ChallengesScores", func() {
	Description("Combined accumulated scores for HC and SC challenges")

	Attribute("node_id", String, "Specific node id")
	Attribute("ip_address", String, "IPAddress of the node")
	Attribute("storage_challenge_score", Float64, "Total accumulated SC challenge score")
	Attribute("health_check_challenge_score", Float64, "Total accumulated HC challenge score")

	Required("node_id", "storage_challenge_score",
		"health_check_challenge_score")
})
