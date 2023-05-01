package design

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports
	"time"
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
			Example(100)
		})
	})
	Required("image_id", "expires_in", "total_estimated_fee")
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

// Details represents status log details with additional fields
var Details = Type("Details", func() {
	Attribute("message", String, func() {
		Description("details regarding the status")
		Example("Image has been downloaded...")
	})
	Attribute("fields", MapOf(String, Any), func() {
		Description("important fields regarding status history")
	})
})

// TaskHistory is a list of strings containing all past task histories.
var TaskHistory = Type("TaskHistory", func() {
	Attribute("timestamp", String, func() {
		Description("Timestamp of the status creation")
		Example(time.RFC3339)
	})

	Attribute("status", String, func() {
		Description("past status string")
		Example("Started, Image Probed, Downloaded...")
	})

	Attribute("message", String, func() {
		Description("message string (if any)")
		Example("Balance less than maximum fee provied in the request, could not gather enough confirmations...")
	})

	Attribute("details", Details, func() {
		Description("details of the status")
	})

	Required("status")
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
	APIKey("api_key", "key", String, func() {
		Description("Passphrase of the owner's PastelID")
		Example("Basic abcdef12345")
	})

	Required("image_id", "burn_txid", "app_pastelid", "key")
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

// DownloadPayload is asset download payload.
var DownloadPayload = Type("DownloadPayload", func() {
	Attribute("txid", String, func() {
		Description("Nft Registration Request transaction ID")
		MinLength(64)
		MaxLength(64)
		Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
	})
	Attribute("pid", String, func() {
		Meta("struct:field:name", "Pid")
		Description("Owner's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	APIKey("api_key", "key", String, func() {
		Description("Passphrase of the owner's PastelID")
		Example("Basic abcdef12345")
	})
	Required("txid", "pid", "key")
})

// DownloadResult is Asset download result.
var DownloadResult = Type("DownloadResult", func() {
	Description("Asset download response")
	Attribute("file", Bytes, func() {
		Description("File downloaded")
	})
	Required("file")
})

// DDServiceOutputFileResult is DD service output file result details.
var DDServiceOutputFileResult = Type("DDServiceOutputFileResult", func() {
	Description("Complete details of dd service output file result")

	Extend(NftSummary)

	Attribute("version", Int, func() {
		Description("version")
		Example(1)
	})
	Attribute("storage_fee", Int, func() {
		Description("Storage fee %")
		Example(100)
	})
	Attribute("block_height", String, func() {
		Description("Block Height When request submitted")
	})
	Attribute("timestamp_of_request", String, func() {
		Description("Timestamp of request when submitted")
	})
	Attribute("submitter_pastel_id", String, func() {
		Description("Pastel id of the submitter")
	})
	Attribute("sn1_pastel_id", String, func() {
		Description("Pastel id of SN1")
	})
	Attribute("sn2_pastel_id", String, func() {
		Description("Pastel id of SN2")
	})
	Attribute("sn3_pastel_id", String, func() {
		Description("Pastel id of SN3")
	})
	Attribute("is_open_api_request", Boolean, func() {
		Description("Is Open API request")
	})
	Attribute("open_api_subset_id", String, func() {
		Description("Subset id of the open API")
	})
	Attribute("dupe_detection_system_version", String, func() {
		Description("System version of dupe detection")
	})
	Attribute("is_likely_dupe", Boolean, func() {
		Description("is this nft likely a duplicate")
		Example(false)
	})
	Attribute("is_rare_on_internet", Boolean, func() {
		Description("is this nft rare on the internet")
		Example(false)
	})
	Attribute("overall_rareness_score", Float32, func() {
		Description("pastel rareness score")
	})
	Attribute("pct_of_top_10_most_similar_with_dupe_prob_above_25pct", Float32, func() {
		Description("PCT of top 10 most similar with dupe probe above 25 PCT")
	})
	Attribute("pct_of_top_10_most_similar_with_dupe_prob_above_33pct", Float32, func() {
		Description("PCT of top 10 most similar with dupe probe above 33 PCT")
	})
	Attribute("pct_of_top_10_most_similar_with_dupe_prob_above_50pct", Float32, func() {
		Description("PCT of top 10 most similar with dupe probe above 50 PCT")
	})
	Attribute("rareness_scores_table_json_compressed_b64", String, func() {
		Description("Rareness scores table json compressed b64 encoded")
	})
	Attribute("internet_rareness_score", Float32, func() {
		Description("internet rareness score")
	})
	Attribute("open_nsfw_score", Float32, func() {
		Description("internet rareness score")
	})
	Attribute("image_fingerprint_of_candidate_image_file", ArrayOf(Float32), func() {
		Description("Image fingerprint of candidate image file")
	})
	Attribute("drawing_nsfw_score", Float32, func() {
		Description("nsfw score")
	})
	Attribute("neutral_nsfw_score", Float32, func() {
		Description("nsfw score")
	})
	Attribute("sexy_nsfw_score", Float32, func() {
		Description("nsfw score")
	})
	Attribute("porn_nsfw_score", Float32, func() {
		Description("nsfw score")
	})
	Attribute("hentai_nsfw_score", Float32, func() {
		Description("nsfw score")
	})
	Attribute("preview_thumbnail", Bytes, func() {
		Description("Preview Image")
	})
	Attribute("rare_on_internet_summary_table_json_b64", String, func() {
		Description("Base64 Compressed JSON Table of Rare On Internet Summary")
	})
	Attribute("rare_on_internet_graph_json_b64", String, func() {
		Description("Base64 Compressed JSON of Rare On Internet Graph")
	})
	Attribute("alt_rare_on_internet_dict_json_b64", String, func() {
		Description("Base64 Compressed Json of Alternative Rare On Internet Dict")
	})
	Attribute("min_num_exact_matches_on_page", UInt32, func() {
		Description("Minimum Number of Exact Matches on Page")
	})
	Attribute("earliest_date_of_results", String, func() {
		Description("Earliest Available Date of Internet Results")
	})
})
