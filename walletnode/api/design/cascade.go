package design

import (
	"time"

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("cascade", func() {
	Description("OpenAPI Cascade service")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/openapi/cascade")
	})

	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("uploadAsset", func() {
		Description("Upload the asset file")
		Meta("swagger:summary", "Uploads Action Data")

		Payload(func() {
			Extend(AssetUploadPayload)
		})
		Result(AssetUploadResult)

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

		Security(APIKeyAuth)

		Payload(func() {
			Extend(StartCascadeProcessingPayload)
		})
		Result(StartProcessingResult)

		HTTP(func() {
			POST("/start/{file_id}")
			Params(func() {
				Param("file_id", String)
			})
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

	Method("getTaskHistory", func() {
		Description("Gets the history of the task's states.")
		Meta("swagger:summary", "Get history of states as a json string with a list of state objects.")

		Payload(func() {
			Extend(RegisterTaskPayload)
		})
		Result(ArrayOf(TaskHistory))

		HTTP(func() {
			GET("/{taskId}/history")
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("download", func() {
		Description("Download cascade Artifact.")
		Meta("swagger:summary", "Downloads cascade artifact")

		Security(APIKeyAuth)

		Payload(DownloadPayload)
		//Result(DownloadResult)
		Result(Any)
		HTTP(func() {
			GET("/download")
			//SkipResponseBodyEncodeDecode()

			Param("txid")
			Param("pid")

			// Header("key:Authorization") // Provide the key in Authorization header (default)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})
})

// AssetUploadPayload represents a payload for uploading asset
var AssetUploadPayload = Type("AssetUploadPayload", func() {
	Description("Asset upload payload")
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

// AssetUploadResult is asset file upload result.
var AssetUploadResult = ResultType("application/vnd.cascade.upload-file", func() {
	TypeName("Asset")
	Attributes(func() {
		Attribute("file_id", String, func() {
			Description("Uploaded file ID")
			MinLength(8)
			MaxLength(8)
			Example("VK7mpAqZ")
		})
		Attribute("expires_in", String, func() {
			Description("File expiration")
			Format(FormatDateTime)
			Example(time.RFC3339)
		})

		Attribute("total_estimated_fee", Float64, func() {
			Description("Estimated fee")
			Minimum(0.00001)
			Default(1)
			Example(100)
		})

		Attribute("required_preburn_amount", Float64, func() {
			Description("The amount that's required to be preburned")
			Minimum(0.00001)
			Default(1)
			Example(20)
		})
	})
	Required("file_id", "expires_in", "total_estimated_fee")
})

// StartCascadeProcessingPayload - Payload for starting processing
var StartCascadeProcessingPayload = Type("StartCascadeProcessingPayload", func() {
	Description("Start Processing Payload")
	Attribute("file_id", String, func() {
		Description("Uploaded asset file ID")
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
	Attribute("make_publicly_accessible", Boolean, func() {
		Meta("struct:field:name", "MakePubliclyAccessible")
		Description("To make it publicly accessible")
		Example(false)
		Default(false)

	})
	APIKey("api_key", "key", String, func() {
		Description("Passphrase of the owner's PastelID")
		Example("Basic abcdef12345")
	})

	Required("file_id", "burn_txid", "app_pastelid", "key")
})
