package design

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("nft", func() {
	Description("Pastel NFT")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/nfts")
	})

	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("register", func() {
		Description("Runs a new registration process for the new NFT.")
		Meta("swagger:summary", "Registers a new NFT")

		Payload(func() {
			Extend(NftRegisterPayload)
			Attribute("image_id", String, func() {
				Description("Uploaded image ID")
				MinLength(8)
				MaxLength(8)
				Example("VK7mpAqZ")
			})
			Required("image_id")
		})
		Result(NftRegisterResult)

		HTTP(func() {
			POST("/register")
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
			GET("/register/{taskId}/state")
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

	Method("registerTask", func() {
		Description("Returns a single task.")
		Meta("swagger:summary", "Find task by ID")

		Payload(func() {
			Extend(RegisterTaskPayload)
		})
		Result(NftRegisterTaskResult, func() {
			View("default")
		})

		HTTP(func() {
			GET("/register/{taskId}")
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("registerTasks", func() {
		Description("List of all tasks.")
		Meta("swagger:summary", "Returns list of tasks")

		Result(CollectionOf(NftRegisterTaskResult), func() {
			View("tiny")
		})

		HTTP(func() {
			GET("/register")
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("uploadImage", func() {
		Description("Upload the image that is used when registering a new NFT.")
		Meta("swagger:summary", "Uploads an image")

		Payload(func() {
			Extend(ImageUploadPayload)
		})
		Result(NFTImageUploadResult)

		HTTP(func() {
			POST("/register/upload")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})

	Method("nftSearch", func() {
		Description("Streams the search result for NFT")
		Meta("swagger:summary", "Streams the search result for NFT")

		Payload(SearchNftParams)
		StreamingResult(NftSearchResult)

		HTTP(func() {
			GET("/search")
			Params(func() {
				Param("artist")
				Param("limit")
				Param("query")
				Param("creator_name")
				Param("art_title")
				Param("series")
				Param("descr")
				Param("keyword")
				Param("min_copies")
				Param("max_copies")
				Param("min_block")
				Param("max_block")
				Param("is_likely_dupe")
				Param("min_rareness_score")
				Param("max_rareness_score")
				Param("min_nsfw_score")
				Param("max_nsfw_score")
				Param("is_likely_dupe")
			})
			Header("user_pastelid")
			Header("user_passphrase")
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("nftGet", func() {
		Description("Gets the NFT detail")
		Meta("swagger:summary", "Returns the detail of NFT")

		Security(APIKeyAuth)

		Payload(NftGetParams)
		Result(NftDetail)

		HTTP(func() {
			GET("/")
			Params(func() {
				Param("txid")
				Param("pid")
			})
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})
	Method("download", func() {
		Description("Download registered NFT.")
		Meta("swagger:summary", "Downloads NFT")

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

	Method("ddServiceOutputFileDetail", func() {
		Description("Duplication detection output file details")
		Meta("swagger:summary", "Duplication detection output file details")

		Security(APIKeyAuth)
		Payload(DownloadPayload)
		Result(DDServiceOutputFileResult)

		HTTP(func() {
			GET("/get_dd_results")
			Param("txid")
			Param("pid")

			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})
})

// NftSearchResult is NFT search result.
var NftSearchResult = Type("NftSearchResult", func() {
	Description("Result of NFT search call")

	Attribute("nft", NftSummary, func() {
		Description("NFT data")
	})

	Attribute("match_index", Int, func() {
		Description("Sort index of the match based on score.This must be used to sort results on UI.")
	})
	Attribute("matches", ArrayOf(FuzzyMatch), func() {
		Description("Match result details")
	})

	Required("nft", "matches", "match_index")
})

// NftRegisterPayload is NFT register payload.
var NftRegisterPayload = Type("NftRegisterPayload", func() {
	Description("Request of the registration NFT")

	Attribute("name", String, func() {
		Description("Name of the NFT")
		MaxLength(256)
		Example("Mona Lisa")
	})
	Attribute("description", String, func() {
		Description("Description of the NFT")
		MaxLength(1024)
		Example("The Mona Lisa is an oil painting by Italian artist, inventor, and writer Leonardo da Vinci. Likely completed in 1506, the piece features a portrait of a seated woman set against an imaginary landscape.")
	})
	Attribute("keywords", String, func() {
		Description("Keywords")
		MaxLength(256)
		Example("Renaissance, sfumato, portrait")
	})
	Attribute("series_name", String, func() {
		Description("Series name")
		MaxLength(256)
		Example("Famous artist")
	})
	Attribute("issued_copies", Int, func() {
		Description("Number of copies issued")
		Minimum(1)
		Maximum(1000)
		Default(1)
		Example(1)
	})
	Attribute("youtube_url", String, func() {
		Description("NFT creation video youtube URL")
		MaxLength(128)
		Example("https://www.youtube.com/watch?v=0xl6Ufo4ZX0")
	})

	Attribute("creator_pastelid", String, func() {
		Meta("struct:field:name", "CreatorPastelID")
		Description("Creator's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("creator_pastelid_passphrase", String, func() {
		Meta("struct:field:name", "CreatorPastelIDPassphrase")
		Description("Passphrase of the artist's PastelID")
		Example("qwerasdf1234")
	})
	Attribute("creator_name", String, func() {
		Description("Name of the NFT creator")
		MaxLength(256)
		Example("Leonardo da Vinci")
	})
	Attribute("creator_website_url", String, func() {
		Description("NFT creator website URL")
		MaxLength(256)
		Example("https://www.leonardodavinci.net")
	})

	Attribute("spendable_address", String, func() {
		Description("Spendable address")
		MinLength(35)
		MaxLength(35)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j")
	})
	Attribute("maximum_fee", Float64, func() {
		Description("Used to find a suitable masternode with a fee equal or less")
		Minimum(0.00001)
		Default(1)
		Example(100)
	})

	Attribute("royalty", Float64, func() {
		Description("Percentage the artist received in future sales. If set to 0% he only get paids for the first sale on each copy of the NFT")
		Default(0.0)
		Example(12.0)
		Minimum(0.0)
		Maximum(20.0)
	})

	Attribute("green", Boolean, func() {
		Description("To donate 2% of the sale proceeds on every sale to TeamTrees which plants trees")
		Example(false)
		Default(false)
	})

	Attribute("thumbnail_coordinate", ThumbnailCoordinate)

	Attribute("make_publicly_accessible", Boolean, func() {
		Meta("struct:field:name", "MakePubliclyAccessible")
		Description("To make it publicly accessible")
		Example(false)
		Default(false)

	})

	Attribute("collection_act_txid", String, func() {
		Description("Act Collection TxID to add given ticket in collection ")
		Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
	})

	Required("creator_name", "name", "issued_copies", "creator_pastelid", "creator_pastelid_passphrase", "spendable_address", "maximum_fee")
})

// NftRegisterResult is NFT registeration result.
var NftRegisterResult = ResultType("application/vnd.nft.register", func() {
	TypeName("RegisterResult")
	Attributes(func() {
		Attribute("task_id", String, func() {
			Description("Task ID of the registration process")
			MinLength(8)
			MaxLength(8)
			Example("n6Qn6TFM")
		})
	})
	Required("task_id")
})

// NftRegisterTaskResult is task streaming of the NFT registration.
var NftRegisterTaskResult = ResultType("application/vnd.nft.register.task", func() {
	TypeName("Task")
	Attributes(func() {
		Attribute("id", String, func() {
			Description("JOb ID of the registration process")
			MinLength(8)
			MaxLength(8)
			Example("n6Qn6TFM")
		})
		Attribute("status", String, func() {
			Description("Status of the registration process")
			Example(common.StatusNames()[0])
			Enum(InterfaceSlice(common.StatusNames())...)
		})
		Attribute("states", ArrayOf(RegisterTaskState), func() {
			Description("List of states from the very beginning of the process")
		})
		Attribute("txid", String, func() {
			Description("txid")
			MinLength(64)
			MaxLength(64)
			Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
		})
		Attribute("ticket", NftRegisterPayload)
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("status")
		Attribute("txid")
		Attribute("ticket")
	})

	Required("id", "status", "ticket")
})

// FuzzyMatch is search results detail
var FuzzyMatch = Type("FuzzyMatch", func() {
	Attribute("str", String, func() {
		Description("String that is matched")
	})
	Attribute("field_type", String, func() {
		Description("Field that is matched")
		Enum(InterfaceSlice(nftsearch.NftSearchQueryFields)...) //TODO: Rename after renaming package
	})
	Attribute("matched_indexes", ArrayOf(Int), func() {
		Description("The indexes of matched characters. Useful for highlighting matches")
	})
	Attribute("score", Int, func() {
		Description("Score used to rank matches")
	})
})

// NftGetParams are request params to nftGet Params
var NftGetParams = func() {
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
}

// SearchNftParams are query params to searchNft request
var SearchNftParams = func() {
	Attribute("artist", String, func() {
		Description("Artist PastelID or special value; mine")
		MaxLength(256)
	})
	Attribute("limit", Int, func() {
		Description("Number of results to be return")
		Minimum(10)
		Maximum(200)
		Default(10)
		Example(10)
	})
	Attribute("query", String, func() {
		Description("Query is search query entered by user")
	})
	Attribute("creator_name", Boolean, func() {
		Description("Name of the nft creator")
		Default(true)
	})
	Attribute("art_title", Boolean, func() {
		Description("Title of NFT")
		Default(true)
	})
	Attribute("series", Boolean, func() {
		Description("NFT series name")
		Default(true)
	})
	Attribute("descr", Boolean, func() {
		Description("Artist written statement")
		Default(true)
	})
	Attribute("keyword", Boolean, func() {
		Description("Keyword that Artist assigns to NFT")
		Default(true)
	})
	Attribute("min_block", Int, func() {
		Description("Minimum blocknum")
		Minimum(1)
		Default(1)
	})
	Attribute("max_block", Int, func() {
		Description("Maximum blocknum")
		Minimum(1)
	})
	Attribute("min_copies", Int, func() {
		Description("Minimum number of created copies")
		Minimum(1)
		Maximum(1000)
		Example(1)
	})
	Attribute("max_copies", Int, func() {
		Description("Maximum number of created copies")
		Minimum(1)
		Maximum(1000)
		Example(1000)
	})
	Attribute("min_nsfw_score", Float64, func() {
		Description("Minimum nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("max_nsfw_score", Float64, func() {
		Description("Maximum nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("min_rareness_score", Float64, func() {
		Description("Minimum pastel rareness score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("max_rareness_score", Float64, func() {
		Description("Maximum pastel rareness score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("is_likely_dupe", Boolean, func() {
		Description("Is this image likely a duplicate of another known image")
		Example(false)
	})
	Attribute("user_pastelid", String, func() {
		Description("User's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("user_passphrase", String, func() {
		Description("Passphrase of the User PastelID")
		Example("qwerasdf1234")
	})
	Required("query")
}

// NftSummary is part of NFT search response.
var NftSummary = Type("NftSummary", func() {
	Description("NFT response")

	Attribute("thumbnail_1", Bytes, func() {
		Description("Thumbnail_1 image")
	})

	Attribute("thumbnail_2", Bytes, func() {
		Description("Thumbnail_2 image")
	})

	Attribute("txid", String, func() {
		Description("txid")
		MinLength(64)
		MaxLength(64)
		Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
	})

	Attribute("title", String, func() {
		Description("Name of the NFT")
		MaxLength(256)
		Example("Mona Lisa")
	})
	Attribute("description", String, func() {
		Description("Description of the NFT")
		MaxLength(1024)
		Example("The Mona Lisa is an oil painting by Italian artist, inventor, and writer Leonardo da Vinci. Likely completed in 1506, the piece features a portrait of a seated woman set against an imaginary landscape.")
	})
	Attribute("keywords", String, func() {
		Description("Keywords")
		MaxLength(256)
		Example("Renaissance, sfumato, portrait")
	})
	Attribute("series_name", String, func() {
		Description("Series name")
		MaxLength(256)
		Example("Famous artist")
	})
	Attribute("copies", Int, func() {
		Description("Number of copies")
		Minimum(1)
		Maximum(1000)
		Default(1)
		Example(1)
	})
	Attribute("youtube_url", String, func() {
		Description("NFT creation video youtube URL")
		MaxLength(128)
		Example("https://www.youtube.com/watch?v=0xl6Ufo4ZX0")
	})

	Attribute("creator_pastelid", String, func() {
		Meta("struct:field:name", "CreatorPastelID")
		Description("Artist's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Attribute("creator_name", String, func() {
		Description("Name of the artist")
		MaxLength(256)
		Example("Leonardo da Vinci")
	})
	Attribute("creator_website_url", String, func() {
		Description("Artist website URL")
		MaxLength(256)
		Example("https://www.leonardodavinci.net")
	})
	Attribute("nsfw_score", Float32, func() {
		Description("NSFW Average score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("rareness_score", Float32, func() {
		Description("Average pastel rareness score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("is_likely_dupe", Boolean, func() {
		Description("Is this image likely a duplicate of another known image")
		Example(false)
	})
	Required("title", "description", "creator_name", "copies", "creator_pastelid", "txid")
})

// NFTImageUploadResult is image upload result.
var NFTImageUploadResult = ResultType("application/vnd.nft.upload-image-result", func() {
	TypeName("ImageRes")
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
		Attribute("estimated_fee", Float64, func() {
			Description("Estimated fee")
			Minimum(0.00001)
			Default(1)
			Example(100)
		})
	})
	Required("image_id", "expires_in", "estimated_fee")
})

// NftDetail is NFT get response.
var NftDetail = Type("NftDetail", func() {
	Description("NFT detail response")

	Extend(NftSummary)

	Attribute("version", Int, func() {
		Description("version")
		Example(1)
	})
	Attribute("green_address", Boolean, func() {
		Description("Green address")
	})
	Attribute("royalty", Float64, func() {
		Description("how much artist should get on all future resales")
	})
	Attribute("storage_fee", Int, func() {
		Description("Storage fee %")
		Example(100)
	})
	Attribute("nsfw_score", Float32, func() {
		Description("nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("rareness_score", Float32, func() {
		Description("pastel rareness score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	// Attribute("internet_rareness_score", Float32, func() {
	// 	Description("internet rareness score")
	// 	Minimum(0)
	// 	Maximum(1)
	// 	Example(1)
	// })
	Attribute("is_likely_dupe", Boolean, func() {
		Description("is this nft likely a duplicate")
		Example(false)
	})
	Attribute("is_rare_on_internet", Boolean, func() {
		Description("is this nft rare on the internet")
		Example(false)
	})
	// Attribute("matches_found_on_first_page", UInt32, func() {
	// 	Description("How many matches the scraper found on the first page of the google results")
	// 	Example(1)
	// })
	// Attribute("number_of_pages_of_results", UInt32, func() {
	// 	Description("How many pages of search results the scraper found when searching for this image")
	// 	Example(1)
	// })
	// Attribute("URL_of_first_match_in_page", String, func() {
	// 	Description("URL of the first match on the first page of search results")
	// 	Example("https://pastel.network")
	// })
	Attribute("drawing_nsfw_score", Float32, func() {
		Description("nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("neutral_nsfw_score", Float32, func() {
		Description("nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("sexy_nsfw_score", Float32, func() {
		Description("nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("porn_nsfw_score", Float32, func() {
		Description("nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
	})
	Attribute("hentai_nsfw_score", Float32, func() {
		Description("nsfw score")
		Minimum(0)
		Maximum(1)
		Example(1)
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

	Required("rareness_score", "nsfw_score", "is_likely_dupe", "is_rare_on_internet")
})

// ThumbnailCoordinate is the cordinate of the cropped region selectd by user
var ThumbnailCoordinate = ResultType("ThumbnailCoordinate", func() {
	Description("Coordinate of the thumbnail")
	Attribute("top_left_x", Int64, func() {
		Description("X coordinate of the thumbnail's top left conner")
		Example(0)
		Default(0)
	})
	Attribute("top_left_y", Int64, func() {
		Description("Y coordinate of the thumbnail's top left conner")
		Example(0)
		Default(0)
	})
	Attribute("bottom_right_x", Int64, func() {
		Description("X coordinate of the thumbnail's bottom right conner")
		Example(640)
		Default(0)
	})
	Attribute("bottom_right_y", Int64, func() {
		Description("Y coordinate of the thumbnail's bottom right conner")
		Example(480)
		Default(0)
	})
	View("default", func() {
		Attribute("top_left_x")
		Attribute("top_left_y")
		Attribute("bottom_right_x")
		Attribute("bottom_right_y")
	})
	Required("top_left_x", "top_left_y", "bottom_right_x", "bottom_right_y")
})

// APIKeyAuth is donwload security schemes.
var APIKeyAuth = APIKeySecurity("api_key", func() {
	Description("Nft Owner's passphrase to authenticate")
})
