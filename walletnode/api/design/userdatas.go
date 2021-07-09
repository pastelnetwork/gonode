package design

import (
	"time"

	// TODO: Generate this services
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"

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

	Method("processUserdata", func() {
		Description("Runs a user data update process")
		Meta("swagger:summary", "Process user data")

		Payload(func() {
			Extend(UserSpecifiedData)
		})
		Result(UserdataProcessResult)

		HTTP(func() {
			POST("/userdata")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})

	Method("processUserdataTaskState", func() {
		Description("Streams the state of the userdata process.")
		Meta("swagger:summary", "Streams state by task ID")

		Payload(func() {
			Extend(ProcessTaskPayload)
		})
		StreamingResult(UserdataProcessTaskState)

		HTTP(func() {
			GET("/task/{taskId}/state")
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("processUserdataTask", func() {
		Description("Returns a single task.")
		Meta("swagger:summary", "Find task by ID")

		Payload(func() {
			Extend(ProcessTaskPayload)
		})
		Result(UserdataProcessTaskResult, func() {
			View("default")
		})

		HTTP(func() {
			GET("/task/{taskId}")
			Response("NotFound", StatusNotFound)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("processUserdataTasks", func() {
		Description("List of all tasks.")
		Meta("swagger:summary", "Returns list of tasks")

		Result(CollectionOf(UserdataProcessTaskResult), func() {
			View("tiny")
		})

		HTTP(func() {
			GET("/task")
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

})

// UserdataProcessResult is user data process result.
var UserdataProcessResult = ResultType("application/userdata.process", func() {
	TypeName("processResult")
	Attributes(func() {
		Attribute("task_id", String, func() {
			Description("Task ID of the user data process")
			MinLength(8)
			MaxLength(8)
			Example("n6Qn6TFM")
		})
	})
	Required("task_id")
})

// UserdataProcessTaskResult is task streaming of the userdate process.
var UserdataProcessTaskResult = ResultType("application/userdata.process", func() {
	TypeName("Task")
	Attributes(func() {
		Attribute("id", String, func() {
			Description("JOb ID of the userdata process")
			MinLength(8)
			MaxLength(8)
			Example("n6Qn6TFM")
		})
		Attribute("status", String, func() {
			Description("Status of the userdata process")
			Example(userdataprocess.StatusNames()[0])
			Enum(InterfaceSlice(userdataprocess.StatusNames())...)
		})
		Attribute("states", ArrayOf(UserdataProcessTaskState), func() {
			Description("List of states from the very beginning of the process")
		})
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("status")
	})

	Required("id", "status")
})

// UserdataProcessTaskState is task streaming of the userdata process.
var UserdataProcessTaskState = Type("TaskState", func() {
	Attribute("date", String, func() {
		Description("Date of the status creation")
		Example(time.RFC3339)
	})
	Attribute("status", String, func() {
		Description("Status of the registration process")
		Example(userdataprocess.StatusNames()[0])
		Enum(InterfaceSlice(userdataprocess.StatusNames())...)
	})
	Required("date", "status")
})

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
		Description("Real name of the user")
		MaxLength(256)
		Example("Williams Scottish")
	})

	Attribute("facebook_link", String, func() {
		Description("Facebook link of the user")
		MaxLength(128)
		Example("https://www.facebook.com/Williams_Scottish")
	})

	Attribute("twitter_link", String, func() {
		Description("Twitter link of the user")
		MaxLength(128)
		Example("https://www.twitter.com/@Williams_Scottish")
	})

	Attribute("native_currency", String, func() {
		Description("Native currency of user in ISO 4217 Alphabetic Code")
		MinLength(3)
		MaxLength(3)
		Example("USD")
	})

	Attribute("location", String, func() {
		Description("Location of the user")
		MaxLength(256)
		Example("New York, US")
	})

	Attribute("primary_language", String, func() {
		Description("Primary language of the user")
		MaxLength(30)
		Example("English")
	})

	Attribute("categories", ArrayOf(String), func() {
		Description("The categories of user's work")
	})

	Attribute("biography", String, func() {
		Description("Biography of the user")
		MaxLength(1024)
		Example("I'm a digital artist based in Paris, France. ...")
	})

	Attribute("avatar_image", ImageUploadPayload, func() {
		Description("Avatar image of the user")
	})

	Attribute("cover_photo", ImageUploadPayload, func() {
		Description("Cover photo of the user")
	})

	Attribute("biography", String, func() {
		Description("Biography of the user")
		MaxLength(1024)
		Example("I'm a digital artist based in Paris, France. ...")
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

	Attribute("biography", String, func() {
		Description("Error detail on biography")
		MaxLength(256)
	})

	Attribute("cover_photo", String, func() {
		Description("Error detail on biography")
		MaxLength(256)
	})
	Required("response_code", "detail")
})

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