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

	Method("getUserdata", func() {
		Description("Gets the Userdata detail")
		Meta("swagger:summary", "Returns the detail of Userdata")

		Payload(GetUserdataParams)
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
		Meta("struct:field:name", "RealName")
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
		Meta("struct:field:name", "RealName")
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

// GetUserdataParams are request params to GetUserdata Params
var GetUserdataParams = func() {
	Attribute("pastelid", String, func() {
		Description("Artist's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})

	Required("pastelid")
}
