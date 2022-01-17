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

	Method("startProcessing", func() {
		Description("Start processing the image")
		Meta("swagger:summary", "Starts processing the image")

		Payload(func() {
			Extend(StartProcessingPayload)
		})
		Result(StartProcessingResult)

		HTTP(func() {
			POST("/start/{image_id}")
			Params(func() {
				Param("image_id", String)
			})
			Header("app_pastelid_passphrase")

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})

	Method("registerTaskState", func() {
		Description("Streams the state of the registration process.")
		Meta("swagger:summary", "Streams state by task ID")

		Payload(func() {
			Extend(RegisterTaskPayload)
		})
		StreamingResult(NftRegisterTaskState)

		HTTP(func() {
			GET("/start/{taskId}/state")
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
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
		Pattern(`^[a-zA-Z0-9]`)
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
		Pattern(`^[a-zA-Z0-9\/+]`)
		Example("bTwvO6UZSvFqHb9qmXbAOg2VmupmP70wfhYsAvwFfeC61cuoL9KIXZtbdQ/Ek8FVNoTpCY5BuxcA6lNjkIOBh4w9/RWtuqF16IaJhAnZ4JbZm1MiDCGcf7x0UU/GNRSk6rNAHlPYkOPdhkha+JjCwD4A")
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

// StartProcessingPayload - Payload for starting processing
var StartProcessingPayload = Type("StartProcessingPayload", func() {
	Description("Start Processing Payload")
	Attribute("image_id", String, func() {
		Description("Uploaded image ID")
		MinLength(8)
		MaxLength(8)
		Example("VK7mpAqZ")
	})
	Attribute("burn_txid", String, func() {
		Description("Burn transaction ID")
		MinLength(64)
		MaxLength(64)
		Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
	})
	Attribute("app_pastelid", String, func() {
		Meta("struct:field:name", "AppPastelID")
		Description("App PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("app_pastelid_passphrase", String, func() {
		Meta("struct:field:name", "AppPastelidPassphrase")
		Description("Passphrase of the App PastelID")
		Example("qwerasdf1234")
	})

	Required("image_id", "burn_txid", "app_pastelid", "app_pastelid_passphrase")
})

// StartProcessingResult - result of starting processing
var StartProcessingResult = ResultType("application/sense.start-processing", func() {
	TypeName("startProcessingResult")
	Attributes(func() {
		Attribute("task_id", String, func() {
			Description("Task ID of processing task")
			MinLength(8)
			MaxLength(8)
			Example("VK7mpAqZ")
		})
	})
	Required("task_id")
})
