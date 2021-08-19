package design

import (
	// "time"

	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports

	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("userdatas", func() {
	Description("Pastel Process User Specified Data")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/userdatas")
	})

	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("createUserdata", func() {
		Description("Create new user data")
		Meta("swagger:summary", "Create new user data")

		Payload(func() {
			Extend(UserSpecifiedData)
		})
		Result(UserdataProcessResult)

		HTTP(func() {
			POST("/create")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})

	Method("updateUserdata", func() {
		Description("Update user data for an existing user")
		Meta("swagger:summary", "Update user data for an existing user")

		Payload(func() {
			Extend(UserSpecifiedData)
		})
		Result(UserdataProcessResult)

		HTTP(func() {
			POST("/update")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})

	Method("userdataGet", func() {
		Description("Gets the Userdata detail")
		Meta("swagger:summary", "Returns the detail of Userdata")

		Payload(UserdataGetParams)
		Result(UserSpecifiedData)

		HTTP(func() {
			GET("/{pastelid}")
			Params(func() {
				Param("pastelid")
			})
			Response("BadRequest", StatusBadRequest)
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("setUserFollowRelation", func() {
		Description("Set a follower, followee relationship to metadb")
		Meta("swagger:summary", "Set a follower, followee relationship")

		Payload(UserdataSetFollowRelationPayload)
		Result(UserdataSetRelationResult)

		HTTP(func() {
			POST("/follow")

			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("getFollowers", func() {
		Description("Get followers of a user")
		Meta("swagger:summary", "Get followers of a user")

		Payload(GetUserRelationshipPayload)
		Result(GetUserRelationshipResult)

		HTTP(func() {
			GET("/follow/followers")

			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("getFollowees", func() {
		Description("Get followees of a user")
		Meta("swagger:summary", "Get followers of a user")

		Payload(GetUserRelationshipPayload)
		Result(GetUserRelationshipResult)

		HTTP(func() {
			GET("/follow/followees")

			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("getFriends", func() {
		Description("Get friends of a user")
		Meta("swagger:summary", "Get followers of a user")

		Payload(GetUserRelationshipPayload)
		Result(GetUserRelationshipResult)

		HTTP(func() {
			GET("/follow/friends")

			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("setUserLikeArt", func() {
		Description("Notify a new like event of an user to an art")
		Meta("swagger:summary", "Notify a new like event of an user to an art")

		Payload(UserdataSetLikeArtPayload)
		Result(UserdataSetRelationResult)

		HTTP(func() {
			POST("/like/art")

			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("getUsersLikeArt", func() {
		Description("Get users that liked an art")
		Meta("swagger:summary", "Get users that liked an art")

		Payload(GetUserLikeArtPayload)
		Result(GetUserRelationshipResult)

		HTTP(func() {
			GET("/like/art")

			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})
})

// ProcessTaskPayload represents a payload for returning task.
var ProcessTaskPayload = Type("ProcessTaskPayload", func() {
	Attribute("taskId", String, "Task ID of the user data process", func() {
		TypeName("taskID")
		MinLength(8)
		MaxLength(8)
		Example("n6Qn6TFM")
	})
	Required("taskId")
})

// UserSpecifiedData is user data payload.
var UserSpecifiedData = Type("UserSpecifiedData", func() {
	Description("User Specified Data storing")

	Attribute("realname", String, func() {
		Meta("struct:field:name", "Realname")
		Description("Real name of the user")
		MaxLength(256)
		Example("Williams Scottish")
	})

	Attribute("facebook_link", String, func() {
		Meta("struct:field:name", "FacebookLink")
		Description("Facebook link of the user")
		MaxLength(128)
		Example("https://www.facebook.com/Williams_Scottish")
	})

	Attribute("twitter_link", String, func() {
		Meta("struct:field:name", "TwitterLink")
		Description("Twitter link of the user")
		MaxLength(128)
		Example("https://www.twitter.com/@Williams_Scottish")
	})

	Attribute("native_currency", String, func() {
		Meta("struct:field:name", "NativeCurrency")
		Description("Native currency of user in ISO 4217 Alphabetic Code")
		MinLength(3)
		MaxLength(3)
		Example("USD")
	})

	Attribute("location", String, func() {
		Meta("struct:field:name", "Location")
		Description("Location of the user")
		MaxLength(256)
		Example("New York, US")
	})

	Attribute("primary_language", String, func() {
		Meta("struct:field:name", "PrimaryLanguage")
		Description("Primary language of the user")
		MaxLength(30)
		Example("English")
	})

	Attribute("categories", String, func() {
		Meta("struct:field:name", "Categories")
		Description("The categories of user's work")
	})

	Attribute("biography", String, func() {
		Meta("struct:field:name", "Biography")
		Description("Biography of the user")
		MaxLength(1024)
		Example("I'm a digital artist based in Paris, France. ...")
	})

	Attribute("avatar_image", UserImageUploadPayload, func() {
		Description("Avatar image of the user")
	})

	Attribute("cover_photo", UserImageUploadPayload, func() {
		Description("Cover photo of the user")
	})

	Attribute("artist_pastelid", String, func() {
		Meta("struct:field:name", "ArtistPastelID")
		Description("Artist's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("artist_pastelid_passphrase", String, func() {
		Meta("struct:field:name", "ArtistPastelIDPassphrase")
		Description("Passphrase of the artist's PastelID")
		Example("qwerasdf1234")
	})

	Required("artist_pastelid", "artist_pastelid_passphrase")
})

// UserdataProcessResult is result of process userdata
var UserdataProcessResult = Type("UserdataProcessResult", func() {
	Description("Userdata process result")

	Attribute("response_code", Int, func() {
		Description("Result of the request is success or not")
		Example(0) // Success
	})

	Attribute("detail", String, func() {
		Description("The detail of why result is success/fail, depend on response_code")
		MaxLength(256)
		Example("All userdata is processed") // In case of Success
	})

	Attribute("realname", String, func() {
		Description("Error detail on realname")
		MaxLength(256)
	})

	Attribute("facebook_link", String, func() {
		Description("Error detail on facebook_link")
		MaxLength(256)
	})

	Attribute("twitter_link", String, func() {
		Description("Error detail on twitter_link")
		MaxLength(256)
	})

	Attribute("native_currency", String, func() {
		Description("Error detail on native_currency")
		MaxLength(256)
	})

	Attribute("location", String, func() {
		Description("Error detail on location")
		MaxLength(256)
	})

	Attribute("primary_language", String, func() {
		Description("Error detail on primary_language")
		MaxLength(256)
	})

	Attribute("categories", String, func() {
		Description("Error detail on categories")
		MaxLength(256)
	})

	Attribute("biography", String, func() {
		Description("Error detail on biography")
		MaxLength(256)
	})

	Attribute("avatar_image", String, func() {
		Description("Error detail on avatar")
		MaxLength(256)
	})

	Attribute("cover_photo", String, func() {
		Description("Error detail on cover photo")
		MaxLength(256)
	})

	Attribute("data", Bytes, func() {
		Description("Metadata Layer process metric response")
	})

	Required("response_code", "detail")
})

// UserImageUploadPayload represents a payload for uploading image.
var UserImageUploadPayload = Type("UserImageUploadPayload", func() {
	Description("User image upload payload")
	Attribute("content", Bytes, func() {
		Meta("struct:field:name", "Content")
		Description("File to upload")
	})
	Attribute("filename", String, func() {
		Meta("swagger:example", "false")
		Description("File name of the user image")
	})
	Required("content")
})

// UserdataGetParams are request params to UserdataGet Params
var UserdataGetParams = func() {
	Attribute("pastelid", String, func() {
		Description("Artist's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Required("pastelid")
}

// UserdataSetFollowRelationPayload are request params set follow relationship
var UserdataSetFollowRelationPayload = func() {
	Attribute("follower_pastel_id", String, func() {
		Description("Follower's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Attribute("followee_pastel_id", String, func() {
		Description("Followee's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VaoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Required("follower_pastel_id", "followee_pastel_id")
}

// UserdataSetLikeArtPayload are request params to notify like event to an art
var UserdataSetLikeArtPayload = func() {
	Attribute("user_pastel_id", String, func() {
		Description("User's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Attribute("art_pastel_id", String, func() {
		Description("Art's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VaoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Required("user_pastel_id", "art_pastel_id")
}

var UserdataSetRelationResult = func() {
	Attribute("response_code", Int, func() {
		Description("Result of the request is success or not")
		Example(0) // Success
	})

	Attribute("detail", String, func() {
		Description("The detail of why result is success/fail, depend on response_code")
		MaxLength(256)
		Example("All userdata is processed") // In case of Success
	})

	Required("response_code", "detail")
}

// GetUserRelationshipPayload is request param to get user relationship
var GetUserRelationshipPayload = func() {
	Attribute("pastelid", String, func() {
		Description("Artist's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Attribute("limit", Int, func() {
		Description("limit for paginated list")
		Example(10)
	})

	Attribute("offset", Int, func() {
		Description("offset for paginated list")
		Example(0)
	})

	Required("pastelid")
}

// GetUserRelationshipPayload is request param to get user relationship
var GetUserLikeArtPayload = func() {
	Attribute("art_id", String, func() {
		Description("art id that we want to get like data")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Attribute("limit", Int, func() {
		Description("limit for paginated list")
		Example(10)
	})

	Attribute("offset", Int, func() {
		Description("offset for paginated list")
		Example(0)
	})

	Required("art_id")
}

// GetUserRelationshipResult is result of getting user relationship
var GetUserRelationshipResult = func() {
	Attribute("total_count", Int, func() {
		Description("total number of users in relationship with this user")
		Example(10)
	})

	Attribute("result", ArrayOf(UserRelationshipInfo), func() {
		Description("Artist's PastelID")
	})

	Required("total_count")
}

// UserRelationshipInfo is detail of the person in relationship with a user
var UserRelationshipInfo = Type("UserRelationshipInfo", func() {
	Attribute("username", String, func() {
		Meta("struct:field:name", "Username")
		Description("Username of the user")
		MaxLength(256)
		Example("w_scottish")
	})

	Attribute("realname", String, func() {
		Meta("struct:field:name", "Realname")
		Description("Real name of the user")
		MaxLength(256)
		Example("Williams Scottish")
	})

	Attribute("followers_count", Int, func() {
		Description("number of users follow this user")
		Example(10)
	})

	Attribute("avatar_thumbnail", Bytes, func() {
		Meta("struct:field:name", "AvatarThumbnail")
		Description("40x40 avatar thumbnail")
	})

	Required("username", "followers_count")
})
