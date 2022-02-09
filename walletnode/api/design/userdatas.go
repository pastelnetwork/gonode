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
			Response(StatusOK)
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
			Response(StatusOK)
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
		Description("Primary language of the user, follow ISO 639-2 standard")
		MaxLength(30)
		Example("en")
	})

	Attribute("categories", String, func() {
		Meta("struct:field:name", "Categories")
		Description("The categories of user's work, separate by ,")
		Example("Manga&Anime,3D,Comics")
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

	Attribute("user_pastelid", String, func() {
		Meta("struct:field:name", "UserPastelID")
		Description("User's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("user_pastelid_passphrase", String, func() {
		Meta("struct:field:name", "UserPastelIDPassphrase")
		Description("Passphrase of the user's PastelID")
		Example("qwerasdf1234")
	})

	Required("user_pastelid", "user_pastelid_passphrase")
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
		Example("")
		MaxLength(256)
	})

	Attribute("facebook_link", String, func() {
		Description("Error detail on facebook_link")
		Example("")
		MaxLength(256)
	})

	Attribute("twitter_link", String, func() {
		Description("Error detail on twitter_link")
		Example("")
		MaxLength(256)
	})

	Attribute("native_currency", String, func() {
		Description("Error detail on native_currency")
		Example("")
		MaxLength(256)
	})

	Attribute("location", String, func() {
		Description("Error detail on location")
		Example("")
		MaxLength(256)
	})

	Attribute("primary_language", String, func() {
		Description("Error detail on primary_language")
		Example("")
		MaxLength(256)
	})

	Attribute("categories", String, func() {
		Description("Error detail on categories")
		Example("")
		MaxLength(256)
	})

	Attribute("biography", String, func() {
		Description("Error detail on biography")
		Example("")
		MaxLength(256)
	})

	Attribute("avatar_image", String, func() {
		Description("Error detail on avatar")
		Example("")
		MaxLength(256)
	})

	Attribute("cover_photo", String, func() {
		Description("Error detail on cover photo")
		Example("")
		MaxLength(256)
	})
	Required("response_code", "detail")
})

// UserImageUploadPayload represents a payload for uploading image.
var UserImageUploadPayload = Type("UserImageUploadPayload", func() {
	Description("User image upload payload")
	Attribute("content", Bytes, func() {
		Meta("struct:field:name", "Content")
		Example("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")
		Description("File to upload (byte array of the file content)")
	})
	Attribute("filename", String, func() {
		Meta("struct:field:name", "Filename")
		Description("File name of the user image")
		Example("image_name.png")
		Pattern(`^.*\.(png|PNG|jpeg|JPEG|jpg|JPG)$`)
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
