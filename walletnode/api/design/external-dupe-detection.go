package design

import (

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import

	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("external_dupe_detection_api", func() {
	Description("API detect duplication for external images")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/external_dupe_detection_api")
	})

	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("initiate_submission", func() {
		Description("API initiate submit external dupe detection request")
		Meta("swagger:summary", "initiate submit external dupe detection request")
		Payload(ExternalDupeDetetionAPIInitiateSubmission)
		Result(ExternalDupeDetetionAPIInitiateSubmissionResult)

		HTTP(func() {
			POST("/initiate_submission")

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})
})

var ExternalDupeDetetionAPIInitiateSubmission = Type("ExternalDupeDetetionAPIInitiateSubmission", func() {
	Attribute("request_id", String, "SHA3-256 format, unique external dupe detection request id", func() {
		Meta("struct:field:name", "RequestID")
		MinLength(64)
		MaxLength(64)
		Example("4c278a281fe2d077b2bcd141b382a98e94afb2d1557c87d0013bcf23d18052c5")
	})
	Attribute("datetime_request_initiated", String, "SHA3-256 format, unique external dupe detection request id", func() {
		Meta("struct:field:name", "Datetime_Request_Initiated")
		Format(FormatDateTime)
		Example("2006-01-02T15:04:05Z07:00")
	})
	Attribute("pastel_service_llc_api_id", String, "Pastel service LLC API ID", func() {
		Meta("struct:field:name", "PastelServiceLLCAPIID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("pastel_service_llc_api_secret", String, "Pastel service LLC API Secret", func() {
		Meta("struct:field:name", "PastelServiceLLCAPISecret")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("max_cost_for_dupe_detection_in_usd", Float64, func() {
		Meta("struct:field:name", "MaxCoseForDupeDetectionInUSD")
		Minimum(0)
		Example(10)
	})
	Attribute("submitted_image_data_base64", String, "Base64 encoded of submitted image", func() {
		Meta("struct:field:name", "SubmittedImageDataBase64")
	})
	Attribute("sha3256_hash_of_submitted_image", String, "SHA3-256 hash of submitted image", func() {
		Meta("struct:field:name", "SHA3256HashOfSubmittedImage")
		MinLength(64)
		MaxLength(64)
		Example("4c278a281fe2d077b2bcd141b382a98e94afb2d1557c87d0013bcf23d18052c5")
	})

	Required("request_id", "datetime_request_initiated", "pastel_service_llc_api_secret", "max_cost_for_dupe_detection_in_usd", "submitted_image_data_base64", "sha3256_hash_of_submitted_image")
})

var ExternalDupeDetetionAPIInitiateSubmissionResult = ResultType("ExternalDupeDetetionAPIInitiateSubmissionResult", func() {
	Attribute("request_id", String, "SHA3-256 format, unique external dupe detection request id", func() {
		Meta("struct:field:name", "RequestID")
		MinLength(64)
		MaxLength(64)
		Example("4c278a281fe2d077b2bcd141b382a98e94afb2d1557c87d0013bcf23d18052c5")
	})

	Required("request_id")
})
