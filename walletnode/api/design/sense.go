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
		Meta("swagger:summary", "Uploads Action Data")

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
		StreamingResult(RegisterTaskState)

		HTTP(func() {
			GET("/start/{taskId}/state")
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("download", func() {
		Description("Download sense result; duplication detection results file.")
		Meta("swagger:summary", "Download sense result; duplication detection results file.")

		Security(APIKeyAuth)

		Payload(DownloadPayload)
		Result(DownloadResult)

		HTTP(func() {
			GET("/download")
			Param("txid")
			Param("pid")
			// Header("key:Authorization") // Provide the key in Authorization header (default)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})
})
