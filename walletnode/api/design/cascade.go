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

	Error("UnAuthorized", ErrorResult)
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

	Method("uploadAssetV2", func() {
		Description("Upload the asset file - This endpoint is for the new version of the upload endpoint that supports larger files as well.")
		Meta("swagger:summary", "Uploads Cascade File")

		Payload(func() {
			Extend(AssetUploadPayloadV2)
		})
		Result(AssetUploadResultV2)

		HTTP(func() {
			POST("/v2/upload")
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
			Response("UnAuthorized", StatusUnauthorized)
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
		Result(FileDownloadResult)

		HTTP(func() {
			GET("/download")
			//SkipResponseBodyEncodeDecode()

			Param("txid")
			Param("pid")

			Response("UnAuthorized", StatusUnauthorized)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("registrationDetails", func() {
		Description("Get the file registration details")
		Meta("swagger:summary", "Get the file registration details")

		Payload(func() {
			Extend(FileRegistrationDetailPayload)
		})
		Result(FileRegistrationDetailResult)

		HTTP(func() {
			GET("/registration_details/{file_id}")
			Params(func() {
				Param("file_id", String)
			})

			// Define error HTTP statuses.
			Response("UnAuthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
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

	Attribute("hash", String, func() {
		Meta("swagger:example", "false")
		Description("For internal use")
	})

	Attribute("size", Int64, func() {
		Meta("swagger:example", "false")
		Description("For internal use")
	})

	Required("file")
})

// AssetUploadPayload represents a payload for uploading asset
var AssetUploadPayloadV2 = Type("AssetUploadPayloadV2", func() {
	Description("Asset upload payload")
	Attribute("file", Bytes, func() {
		Meta("struct:field:name", "Bytes")
		Description("File to upload")
	})
	Attribute("filename", String, func() {
		Meta("swagger:example", "false")
		Description("-For internal use-")
	})

	Attribute("hash", String, func() {
		Meta("swagger:example", "false")
		Description("For internal use")
	})

	Attribute("size", Int64, func() {
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

var AssetUploadResultV2 = ResultType("application/vnd.cascade.upload-file-v2", func() {
	TypeName("AssetV2")
	Attributes(func() {
		Attribute("file_id", String, func() {
			Description("Uploaded file ID")
			MinLength(8)
			MaxLength(8)
			Example("VK7mpAqZ")
		})

		Attribute("total_estimated_fee", Float64, func() {
			Description("Estimated fee")
			Minimum(0.00001)
			Default(1)
			Example(100)
		})

		Attribute("required_preburn_transaction_amounts", ArrayOf(Float64), func() {
			Description("The amounts that's required to be preburned - one per transaction")
		})
	})

	Required("file_id", "total_estimated_fee")
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
	Attribute("spendable_address", String, func() {
		Meta("struct:field:name", "SpendableAddress")
		Description("Address to use for registration fee ")
		MinLength(35)
		MaxLength(35)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j")
	})
	APIKey("api_key", "key", String, func() {
		Description("Passphrase of the owner's PastelID")
		Example("Basic abcdef12345")
	})

	Required("file_id", "burn_txid", "app_pastelid", "key")
})

// FileRegistrationDetailPayload - Payload for registration detail
var FileRegistrationDetailPayload = Type("FileRegistrationDetailPayload", func() {
	Description("File registration details")
	Attribute("file_id", String, func() {
		Description("file ID")
		MaxLength(8)
		Example("VK7mpAqZ")
	})

	Required("file_id")
})

// FileRegistrationDetailResult is registration detail result.
var FileRegistrationDetailResult = ResultType("application/vnd.cascade.registration-detail", func() {
	TypeName("Registration")
	Attributes(func() {
		Attribute("files", ArrayOf(File), "List of files")
	})
	Required("files")
})

var File = Type("File", func() {
	Attribute("file_id", String, "File ID")
	Attribute("upload_timestamp", String, "Upload Timestamp in datetime format", func() {
		Format(FormatDateTime)
	})
	Attribute("path", String, "Path to the file")
	Attribute("file_index", String, "Index of the file")
	Attribute("base_file_id", String, "Base File ID")
	Attribute("task_id", String, "Task ID")
	Attribute("reg_txid", String, "Registration Transaction ID")
	Attribute("activation_txid", String, "Activation Transaction ID")
	Attribute("req_burn_txn_amount", Float64, "Required Burn Transaction Amount")
	Attribute("burn_txn_id", String, "Burn Transaction ID")
	Attribute("req_amount", Float64, "Required Amount")
	Attribute("is_concluded", Boolean, "Indicates if the process is concluded")
	Attribute("cascade_metadata_ticket_id", String, "Cascade Metadata Ticket ID")
	Attribute("uuid_key", String, "UUID Key")
	Attribute("hash_of_original_big_file", String, "Hash of the Original Big File")
	Attribute("name_of_original_big_file_with_ext", String, "Name of the Original Big File with Extension")
	Attribute("size_of_original_big_file", Float64, "Size of the Original Big File")
	Attribute("data_type_of_original_big_file", String, "Data Type of the Original Big File")
	Attribute("start_block", Int32, "Start Block")
	Attribute("done_block", Int, "Done Block")
	Attribute("registration_attempts", ArrayOf(RegistrationAttempt), "List of registration attempts")
	Attribute("activation_attempts", ArrayOf(ActivationAttempt), "List of activation attempts")
	Required("file_id", "task_id", "upload_timestamp", "base_file_id", "registration_attempts", "activation_attempts",
		"req_burn_txn_amount", "req_amount", "cascade_metadata_ticket_id", "hash_of_original_big_file", "name_of_original_big_file_with_ext",
		"size_of_original_big_file", "data_type_of_original_big_file")
})

var RegistrationAttempt = Type("RegistrationAttempt", func() {
	Attribute("id", Int, "ID")
	Attribute("file_id", String, "File ID")
	Attribute("reg_started_at", String, "Registration Started At in datetime format", func() {
		Format(FormatDateTime)
	})
	Attribute("processor_sns", String, "Processor SNS")
	Attribute("finished_at", String, "Finished At in datetime format", func() {
		Format(FormatDateTime)
	})
	Attribute("is_successful", Boolean, "Indicates if the registration was successful")
	Attribute("error_message", String, "Error Message")
	Required("id", "file_id", "reg_started_at", "finished_at")
})

var ActivationAttempt = Type("ActivationAttempt", func() {
	Attribute("id", Int, "ID")
	Attribute("file_id", String, "File ID")
	Attribute("activation_attempt_at", String, "Activation Attempt At in datetime format", func() {
		Format(FormatDateTime)
	})
	Attribute("is_successful", Boolean, "Indicates if the activation was successful")
	Attribute("error_message", String, "Error Message")
	Required("id", "file_id", "activation_attempt_at")
})
