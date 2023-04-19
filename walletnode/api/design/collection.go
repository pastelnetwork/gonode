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
		Path("/openapi/collection")
	})

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
			GET("/register")

			Response("BadRequest", StatusBadRequest)
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

	Attribute("collection_final_allowed_block_height", Int, "final allowed block height in days", func() {
		TypeName("collectionFinalAllowedBlockHeight")
		Minimum(1)
		Maximum(7)
		Default(7)
		Example(5)
	})

	Attribute("collection_item_copy_count", Int, "item copy count in the collection", func() {
		TypeName("collectionItemCopyCount")
		Example(10)
	})

	Attribute("royalty", Float32, "royalty fee", func() {
		TypeName("royalty")
		Example(2.32)
	})

	Attribute("green", Boolean, "green", func() {
		TypeName("green")
		Example(false)
	})

	Attribute("app_pastelid", String, func() {
		Meta("struct:field:name", "AppPastelID")
		Description("App PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Required("app_pastelid")

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
