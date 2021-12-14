package design

import (

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("openapi", func() {
	Description("OpenAPI Sensei service")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/openapi")
	})

	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("uploadImage", func() {
		Description("Upload the image")
		Meta("swagger:summary", "Uploads an image")

		Payload(func() {
			Extend(SenseImageUploadPayload)
		})
		Result(SenseImageUploadResult)

		HTTP(func() {
			POST("/sense/upload")
			MultipartRequest()

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
var SenseImageUploadResult = ResultType("application/json", func() {
	TypeName("Image")
	Attributes(func() {
		Attribute("txid", String, func() {
			Description("txid")
			MinLength(64)
			MaxLength(64)
			Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
		})
	})
	Required("txid")
})
