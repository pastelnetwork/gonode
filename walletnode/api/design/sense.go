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
			Extend(SenseImageUploadPayload)
		})
		Result(SenseImageUploadResult)

		HTTP(func() {
			POST("/upload")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})

	Method("startTask", func() {
		Description("Start Action")
		Meta("swagger:summary", "Start Action Task")

		Payload(func() {
			Extend(SenseStartActionPayload)
		})
		Result(SenseStartActionResult)

		HTTP(func() {
			POST("/start/{task_id}")
			Params(func() {
				Param("task_id", String)
			})

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})
})

// SenseImageUploadPayload - Payload for uploading an image
var SenseImageUploadPayload = Type("SenseImageUploadPayload", func() {
	Description("Image upload payload")
	Attribute("file", Bytes, func() {
		Meta("struct:field:name", "Bytes")
		Description("File to upload")
	})
	Attribute("filename", String, func() {
		Meta("swagger:example", "false")
		Description("For internal use")
	})
	Required("file")
})

// SenseImageUploadResult - result of the uploadImage method
var SenseImageUploadResult = ResultType("application/sense.upload-image", func() {
	TypeName("imageUploadResult")
	Attributes(func() {
		Attribute("task_id", String, func() {
			Description("Task ID")
			MinLength(64)
			MaxLength(64)
			Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
		})
	})
	Required("task_id")
})

// SenseStartActionPayload - Payload for starting an action
var SenseStartActionPayload = Type("SenseStartActionPayload", func() {
	Description("Start Action Payload")
	Attribute("task_id", String, func() {
		Description("Task ID")
		MinLength(64)
		MaxLength(64)
		Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
	})
	Attribute("app_pastelid", String, func() {
		Meta("struct:field:name", "PastelID")
		Description("3rd party app's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("action_data_hash", Bytes, func() {
		Meta("struct:field:name", "action_data_hash")
		Description("Hash (SHA3-256) of the Action Data")
	})
	Attribute("action_data_signature", Bytes, func() {
		Meta("struct:field:name", "action_data_signature")
		Description("The signature of the Action Data")
	})

	Required("app_pastelid", "action_data_hash", "action_data_signature")
})

// SenseStartActionResult - result of the Start Action method
var SenseStartActionResult = ResultType("application/sense.start-action", func() {
	TypeName("startActionDataResult")
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
