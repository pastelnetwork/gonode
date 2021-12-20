package design

import (

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("sense", func() {
	Description("OpenAPI Sense service")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/openapi/sense")
	})

	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("uploadImage", func() {
		Description("Upload the image")
		Meta("swagger:summary", "Uploads Action Data (Image)")

		Payload(func() {
			Extend(ImageUploadPayload)
		})
		Result(ImageUploadResult)

		HTTP(func() {
			POST("/upload")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})

	Method("actionDetails", func() {
		Description("Provide action details")
		Meta("swagger:summary", "Provides action details")

		Payload(func() {
			Extend(ActionDetailsPayload)
		})
		Result(ActionDetailsResult)

		HTTP(func() {
			POST("/details/{image_id}")
			Params(func() {
				Param("image_id", String)
			})

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})
})

// ActionDetailsPayload - Payload for posting action details
var ActionDetailsPayload = Type("ActionDetailsPayload", func() {
	Description("Provide Action Details Payload")
	Attribute("image_id", String, func() {
		Description("Uploaded image ID")
		MinLength(8)
		MaxLength(8)
		Example("VK7mpAqZ")
	})
	Attribute("app_pastelid", String, func() {
		Meta("struct:field:name", "PastelID")
		Description("3rd party app's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXZqaS48TT6LFjxnf9P68hZNQVCFqY631FPz4CtM6VugDi5yLNB51ccy17kgKKCyFuL4qadQsXHH3QwHVuPVyY")
	})
	Attribute("action_data_hash", String, func() {
		Meta("struct:field:name", "action_data_hash")
		Description("Hash (SHA3-256) of the Action Data")
		MinLength(64)
		MaxLength(64)
		Pattern(`^[a-fA-F0-9]`)
		Example("7ae3874ff2df92df38cce7586c08fe8f3687884edf3b0543f8d9420f4df31265")
	})
	Attribute("action_data_signature", String, func() {
		Meta("struct:field:name", "action_data_signature")
		Description("The signature (base64) of the Action Data")
		MinLength(152)
		MaxLength(152)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("be9XjyZuhKSlxoxRO4zj3s+hlzF8gWVCwaBFq44n92nONvclfXul/Bn8UgXOIbRGP6LNuzLfiU+ApSUxcaw7NuK0h1tXOmNTt78T/aNT9SiFbBx3wcDqZJtKNlI4a/p7wekDvGjlKfU8jn+RauYCvhEA")
	})

	Required("image_id", "app_pastelid", "action_data_hash", "action_data_signature")
})

// ActionDetailsResult - result of posting action details
var ActionDetailsResult = ResultType("application/sense.start-action", func() {
	TypeName("actionDetailResult")
	Attributes(func() {
		Attribute("estimated_fee", Float64, func() {
			Description("Estimated fee")
			Minimum(0.00001)
			Default(1)
			Example(100)
		})
	})
	Required("estimated_fee")
})
