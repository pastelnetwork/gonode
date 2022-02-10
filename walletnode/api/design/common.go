package design

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import

	. "goa.design/goa/v3/dsl"
	"time"
	//revive:enable:dot-imports
)

// ImageUploadPayload represents a payload for uploading image.
var ImageUploadPayload = Type("ImageUploadPayload", func() {
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

// ImageUploadResult is image upload result.
var ImageUploadResult = ResultType("application/vnd.nft.upload-image", func() {
	TypeName("Image")
	Attributes(func() {
		Attribute("image_id", String, func() {
			Description("Uploaded image ID")
			MinLength(8)
			MaxLength(8)
			Example("VK7mpAqZ")
		})
		Attribute("expires_in", String, func() {
			Description("Image expiration")
			Format(FormatDateTime)
			Example(time.RFC3339)
		})
	})
	Required("image_id", "expires_in")
})

// RegisterTaskPayload represents a payload for returning task.
var RegisterTaskPayload = Type("RegisterTaskPayload", func() {
	Attribute("taskId", String, "Task ID of the registration process", func() {
		TypeName("taskID")
		MinLength(8)
		MaxLength(8)
		Example("n6Qn6TFM")
	})
	Required("taskId")
})

// RegisterTaskState is task streaming of the NFT registration.
var RegisterTaskState = Type("TaskState", func() {
	Attribute("date", String, func() {
		Description("Date of the status creation")
		Example(time.RFC3339)
	})
	Attribute("status", String, func() {
		Description("Status of the registration process")
		Example(common.StatusNames()[0])
		Enum(InterfaceSlice(common.StatusNames())...)
	})
	Required("date", "status")
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
