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
			Example(20)
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
	Attribute("collection_act_txid", String, func() {
		Description("Act Collection TxID to add given ticket in collection ")
		Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
	})
	Attribute("open_api_group_id", String, func() {
		Description("OpenAPI GroupID string")
		Default("PASTEL")
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

// DDFPResultFile is NFT DDFP file result.
var DDFPResultFile = Type("DDFPResultFile", func() {
	Description("Asset download response")
	Attribute("file", String, func() {
		Description("File downloaded")
	})
	Required("file")
})

// DDServiceOutputFileResult is DD service output file result details.
var DDServiceOutputFileResult = Type("DDServiceOutputFileResult", func() {
	Description("Complete details of dd service output file result")

	Attribute("pastel_block_hash_when_request_submitted", String, func() {
		Description("block hash when request submitted")
	})
	Attribute("pastel_block_height_when_request_submitted", String, func() {
		Description("block Height when request submitted")
	})
	Attribute("utc_timestamp_when_request_submitted", String, func() {
		Description("timestamp of request when submitted")
	})
	Attribute("pastel_id_of_submitter", String, func() {
		Description("pastel id of the submitter")
	})
	Attribute("pastel_id_of_registering_supernode_1", String, func() {
		Description("pastel id of registering SN1")
	})
	Attribute("pastel_id_of_registering_supernode_2", String, func() {
		Description("pastel id of registering SN2")
	})
	Attribute("pastel_id_of_registering_supernode_3", String, func() {
		Description("pastel id of registering SN3")
	})
	Attribute("is_pastel_openapi_request", Boolean, func() {
		Description("is pastel open API request")
	})
	Attribute("dupe_detection_system_version", String, func() {
		Description("system version of dupe detection")
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
	Attribute("open_nsfw_score", Float32, func() {
		Description("open nsfw score")
	})
	Attribute("image_fingerprint_of_candidate_image_file", ArrayOf(Float64), func() {
		Description("Image fingerprint of candidate image file")
	})
	Attribute("hash_of_candidate_image_file", String, func() {
		Description("hash of candidate image file")
	})
	Attribute("rareness_scores_table_json_compressed_b64", String, func() {
		Description("rareness scores table json compressed b64")
	})
	Attribute("collection_name_string", String, func() {
		Description("name of the collection")
	})
	Attribute("open_api_group_id_string", String, func() {
		Description("open api group id string")
	})
	Attribute("group_rareness_score", Float32, func() {
		Description("rareness score of the group")
	})
	Attribute("candidate_image_thumbnail_webp_as_base64_string", String, func() {
		Description("candidate image thumbnail as base64 string")
	})
	Attribute("does_not_impact_the_following_collection_strings", String, func() {
		Description("does not impact collection strings")
	})
	Attribute("similarity_score_to_first_entry_in_collection", Float32, func() {
		Description("similarity score to first entry in collection")
	})
	Attribute("cp_probability", Float32, func() {
		Description("probability of CP")
	})
	Attribute("child_probability", Float32, func() {
		Description("child probability")
	})
	Attribute("image_file_path", String, func() {
		Description("file path of the image")
	})
	Attribute("internet_rareness", InternetRareness, func() {
		Description("internet rareness")
	})
	Attribute("alternative_nsfw_scores", AlternativeNSFWScores, func() {
		Description("alternative NSFW scores")
	})
	Attribute("creator_name", String, func() {
		Description("name of the creator")
	})
	Attribute("creator_website", String, func() {
		Description("website of creator")
	})
	Attribute("creator_written_statement", String, func() {
		Description("written statement of creator")
	})
	Attribute("nft_title", String, func() {
		Description("title of NFT")
	})
	Attribute("nft_series_name", String, func() {
		Description("series name of NFT")
	})
	Attribute("nft_creation_video_youtube_url", String, func() {
		Description("nft creation video youtube url")
	})
	Attribute("nft_keyword_set", String, func() {
		Description("keywords for NFT")
	})
	Attribute("total_copies", Int, func() {
		Description("total copies of NFT")
	})
	Attribute("preview_hash", Bytes, func() {
		Description("preview hash of NFT")
	})
	Attribute("thumbnail1_hash", Bytes, func() {
		Description("thumbnail1 hash of NFT")
	})
	Attribute("thumbnail2_hash", Bytes, func() {
		Description("thumbnail2 hash of NFT")
	})
	Attribute("original_file_size_in_bytes", Int, func() {
		Description("original file size in bytes")
	})
	Attribute("file_type", String, func() {
		Description("type of the file")
	})
	Attribute("max_permitted_open_nsfw_score", Float64, func() {
		Description("max permitted open NSFW score")
	})
	Required("creator_name", "creator_website", "creator_written_statement", "nft_title", "nft_series_name",
		"nft_creation_video_youtube_url", "nft_keyword_set", "total_copies", "preview_hash", "thumbnail1_hash", "thumbnail2_hash",
		"original_file_size_in_bytes", "file_type", "max_permitted_open_nsfw_score")
})

// InternetRareness is the InternetRareness Response Type
var InternetRareness = Type("InternetRareness", func() {
	Attribute("rare_on_internet_summary_table_as_json_compressed_b64", String, func() {
		Description("Base64 Compressed JSON Table of Rare On Internet Summary")
	})
	Attribute("rare_on_internet_graph_json_compressed_b64", String, func() {
		Description("Base64 Compressed JSON of Rare On Internet Graph")
	})
	Attribute("alternative_rare_on_internet_dict_as_json_compressed_b64", String, func() {
		Description("Base64 Compressed Json of Alternative Rare On Internet Dict")
	})
	Attribute("min_number_of_exact_matches_in_page", UInt32, func() {
		Description("Minimum Number of Exact Matches on Page")
	})
	Attribute("earliest_available_date_of_internet_results", String, func() {
		Description("Earliest Available Date of Internet Results")
	})
})

// AlternativeNSFWScores is the AlternativeNSFWScores Response Type
var AlternativeNSFWScores = Type("AlternativeNSFWScores", func() {
	Attribute("drawings", Float32, func() {
		Description("drawings nsfw score")
	})
	Attribute("hentai", Float32, func() {
		Description("hentai nsfw score")
	})
	Attribute("sexy", Float32, func() {
		Description("sexy nsfw score")
	})
	Attribute("porn", Float32, func() {
		Description("porn nsfw score")
	})
	Attribute("neutral", Float32, func() {
		Description("neutral nsfw score")
	})
})
