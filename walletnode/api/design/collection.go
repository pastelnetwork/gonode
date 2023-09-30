package design

import (

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("collection", func() {
	Description("OpenAPI Collection service")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/collection")
	})

	Error("UnAuthorized", ErrorResult)
	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("registerCollection", func() {
		Description("Streams the state of the registration process.")
		Meta("swagger:summary", "Streams state by task ID")

		Security(APIKeyAuth)

		Payload(RegisterCollectionPayload)

		Result(RegisterCollectionResponse)

		HTTP(func() {
			POST("/register")

			Response("UnAuthorized", StatusUnauthorized)
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
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
			GET("/{taskId}/state")
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
})

// RegisterCollectionPayload represents a payload for registering collection.
var RegisterCollectionPayload = Type("RegisterCollectionPayload", func() {
	Attribute("collection_name", String, "name of the collection", func() {
		TypeName("collectionName")
		Example("galaxies")
	})
	Required("collection_name")

	Attribute("item_type", String, "type of items, store by collection", func() {
		TypeName("itemType")
		Example("nft", "sense")
		Enum("sense", "nft")
	})
	Required("item_type")

	Attribute("list_of_pastelids_of_authorized_contributors", ArrayOf(String), "list of authorized contributors", func() {
		TypeName("listOfPastelIDsOfAuthorizedContributors")
		Example([]string{"apple", "banana", "orange"})

	})
	Required("list_of_pastelids_of_authorized_contributors")

	Attribute("max_collection_entries", Int, "max no of entries in the collection", func() {
		TypeName("maxCollectionEntries")
		Minimum(1)
		Maximum(10000)
		Example(5000)
	})
	Required("max_collection_entries")

	Attribute("no_of_days_to_finalize_collection", Int, "no of days to finalize collection", func() {
		TypeName("no_of_days_to_finalize_collection")
		Minimum(1)
		Maximum(7)
		Default(7)
		Example(5)
	})

	Attribute("collection_item_copy_count", Int, "item copy count in the collection", func() {
		TypeName("collectionItemCopyCount")
		Minimum(1)
		Maximum(1000)
		Default(1)
		Example(10)
	})

	Attribute("royalty", Float64, "royalty fee", func() {
		TypeName("royalty")
		Minimum(0.0)
		Maximum(20.0)
		Default(0.0)
		Example(2.32)
	})

	Attribute("green", Boolean, "green", func() {
		TypeName("green")
		Default(false)
		Example(false)
	})

	Attribute("max_permitted_open_nsfw_score", Float64, "max open nfsw score sense and nft items can have", func() {
		TypeName("maxPermittedOpenNSFWScore")
		Minimum(0)
		Maximum(1)
		Example(0.5)
	})
	Required("max_permitted_open_nsfw_score")

	Attribute("minimum_similarity_score_to_first_entry_in_collection", Float64, "min similarity for 1st entry to have", func() {
		TypeName("minimumSimilarityScoreToFirstEntryInCollection")
		Minimum(0)
		Maximum(1)
		Example(0.5)
	})
	Required("minimum_similarity_score_to_first_entry_in_collection")

	Attribute("app_pastelid", String, func() {
		Meta("struct:field:name", "AppPastelID")
		Description("App PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Required("app_pastelid")

	Attribute("spendable_address", String, func() {
		Description("Spendable address")
		MinLength(35)
		MaxLength(35)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("PtiqRXn2VQwBjp1K8QXR2uW2w2oZ3Ns7N6j")
	})
	Required("spendable_address")

	APIKey("api_key", "key", String, func() {
		Description("Passphrase of the owner's PastelID")
		Example("Basic abcdef12345")
	})
})

// RegisterCollectionResponse is the response of register collection
var RegisterCollectionResponse = ResultType("application/collection-registration", func() {
	TypeName("RegisterCollectionResponse")
	Attributes(func() {
		Attribute("task_id", String, func() {
			Description("Uploaded file ID")
			MinLength(8)
			MaxLength(8)
			Example("VK7mpAqZ")
		})
	})
	Required("task_id")
})
