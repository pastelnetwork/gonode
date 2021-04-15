package design

import (
	"net/http"

	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("arts", func() {
	Description("Pastel Artwork")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/arts")
	})

	Error("BadRequest", BadRequest)
	Error("InternalServerError", InternalServerError)

	Method("register", func() {
		Description("Registers a new art.")
		Meta("swagger:summary", "Registers a new artwork")

		Payload(func() {
			Extend(Ticket)
		})
		Result(String)

		HTTP(func() {
			POST("/register")
			Response(StatusCreated)
		})
	})

	Method("uploadImage", func() {
		Description("Upload an image that is used when registering the artwork.")
		Meta("swagger:summary", "Uploads an image")

		Payload(ImageUploadPayload)

		Result(String)

		HTTP(func() {
			POST("/upload-image")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
		})
	})
})

// Ticket describes an ticket from an customer to be stored.
var Ticket = Type("RegArtTicket", func() {
	Description("Registration art ticket")
	Attribute("artist_name", String, func() {
		Description("Name of the artist")
		MaxLength(32)
		Example("Leonardo da Vinci")
	})
	Attribute("artwork_name", String, func() {
		Description("Name of the artwork")
		MaxLength(32)
		Example("Mona Lisa")
	})
	Attribute("issued_copies", Int, func() {
		Description("Number of copies issued")
		Minimum(1)
		Maximum(1000)
		Default(1)
		Example(5)
	})
	Attribute("fee", Float32, func() {
		Minimum(0.00001)
		Example(100)
	})

	Attribute("pastel_id", String, "PastelID", func() {
		MaxLength(32)
		Example("123456789")
	})

	Attribute("address", String, func() {
		Description("Spendable address")
		MaxLength(32)
		Example("12349231421309dsfdf")
	})
})

// ImageUpload is an image upload element
var ImageUpload = Type("ImageUpload", func() {
	Description("Image upload type")
	Attribute("file", Bytes, func() {
		Description("File to upload")
	})
})

// ImageUploadPayload is a list of files
var ImageUploadPayload = Type("ImageUploadPayload", func() {
	Extend(ImageUpload)
	Description("Image upload payload")
	Required("file")
})

var InnerError = func(code int) {
	Attribute("code", Int, func() {
		Description("Code refers to a code number in the response header that indicates the general classification of the response.")
		Example(code)
		Default(code)
	})
	Attribute("message", String, func() {
		Description("Message is a human-readable explanation specific to this occurrence of the problem.")
		Example(http.StatusText(code))
		Default(http.StatusText(code))
	})
	Required("code")
}

var BadRequest = Type("BadRequest", func() {
	Attribute("error", func() {
		InnerError(http.StatusBadRequest)
		Meta("struct:field:name", "InnerError")
	})
	Required("error")

})

var InternalServerError = Type("InternalServerError", func() {
	Attribute("error", func() {
		InnerError(http.StatusInternalServerError)
		Meta("struct:field:name", "InnerError")
	})
	Required("error")

})

// var ServiceError = Type("Error", func() {
// 	Attribute("error", InnerError(http.StatusBadRequest), func() {
// 		Meta("struct:field:name", "InnerError")
// 	})
// 	Required("error")
// })
