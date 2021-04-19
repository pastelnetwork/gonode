package design

import (
	//revive:disable:dot-imports
	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
	. "goa.design/goa/v3/dsl"

	//revive:enable:dot-imports
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("artworks", func() {
	Description("Pastel Artwork")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/artworks")
	})

	Error("BadRequest", ErrorResult)
	Error("NotFound", ErrorResult)
	Error("InternalServerError", ErrorResult)

	Method("register", func() {
		Description("Runs a new registration process for the new artwork.")
		Meta("swagger:summary", "Registers a new artwork")

		Payload(func() {
			Extend(ArtworkTicket)
		})
		Result(ArtworkRegisterResult)

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
		StreamingResult(ArtworkRegisterTaskState)

		HTTP(func() {
			GET("/register/{taskId}/state")
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
		Result(ArtworkRegisterTaskResult, func() {
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

		Result(CollectionOf(ArtworkRegisterTaskResult), func() {
			View("tiny")
		})

		HTTP(func() {
			GET("/register")
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusOK)
		})
	})

	Method("uploadImage", func() {
		Description("Upload the image that is used when registering a new artwork.")
		Meta("swagger:summary", "Uploads an image")

		Payload(func() {
			Extend(ImageUploadPayload)
		})
		Result(ImageUploadResult)

		HTTP(func() {
			POST("/register/upload")
			MultipartRequest()

			// Define error HTTP statuses.
			Response("BadRequest", StatusBadRequest)
			Response("InternalServerError", StatusInternalServerError)
			Response(StatusCreated)
		})
	})
})

// ArtworkTicket is artwork register payload.
var ArtworkTicket = Type("ArtworkTicket", func() {
	Description("Ticket of the registration artwork")

	Attribute("name", String, func() {
		Description("Name of the artwork")
		MaxLength(256)
		Example("Mona Lisa")
	})
	Attribute("description", String, func() {
		Description("Description of the artwork")
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
	Attribute("image_id", Int, func() {
		Description("Uploaded image ID")
		Minimum(1)
		Example(1)
	})
	Attribute("youtube_url", String, func() {
		Description("Artwork creation video youtube URL")
		MaxLength(128)
		Example("https://www.youtube.com/watch?v=0xl6Ufo4ZX0")
	})

	Attribute("artist_pastelid", String, func() {
		Meta("struct:field:name", "ArtistPastelID")
		Description("Artist's PastelID")
		MinLength(86)
		MaxLength(86)
		Pattern(`^[a-zA-Z0-9]+$`)
		Example("jXYJud3rmrR1Sk2scvR47N4E4J5Vv48uCC6se2nzHrBRdjaKj3ybPoi1Y2VVoRqi1GnQrYKjSxQAC7NBtvtEdS")
	})
	Attribute("artist_name", String, func() {
		Description("Name of the artist")
		MaxLength(256)
		Example("Leonardo da Vinci")
	})
	Attribute("artist_website_url", String, func() {
		Description("Artist website URL")
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
	Attribute("network_fee", Float32, func() {
		Minimum(0.00001)
		Default(1)
		Example(100)
	})

	Required("artist_name", "name", "issued_copies", "image_id", "artist_pastelid", "spendable_address", "network_fee")
})

// ArtworkRegisterResult is artwork registeration result.
var ArtworkRegisterResult = ResultType("application/vnd.artwork.register", func() {
	TypeName("RegisterResult")
	Attributes(func() {
		Attribute("task_id", Int, func() {
			Description("Task ID of the registration process")
			Example(5)
		})
	})
	Required("task_id")
})

// ArtworkRegisterTaskResult is task streaming of the artwork registration.
var ArtworkRegisterTaskResult = ResultType("application/vnd.artwork.register.task", func() {
	TypeName("Task")
	Attributes(func() {
		Attribute("id", Int, func() {
			Description("JOb ID of the registration process")
			Example(5)
		})
		Attribute("status", String, func() {
			Description("Status of the registration process")
			Example(state.StatusNames()[0])
			Enum(InterfaceSlice(state.StatusNames())...)
		})
		Attribute("states", ArrayOf(ArtworkRegisterTaskState), func() {
			Description("List of states from the very beginning of the process")
		})
		Attribute("txid", String, func() {
			Description("txid")
			MinLength(64)
			MaxLength(64)
			Example("576e7b824634a488a2f0baacf5a53b237d883029f205df25b300b87c8877ab58")
		})
		Attribute("ticket", ArtworkTicket)
	})

	View("tiny", func() {
		Attribute("id")
		Attribute("status")
		Attribute("txid")
	})

	Required("id", "status", "ticket")
})

// ArtworkRegisterTaskState is task streaming of the artwork registration.
var ArtworkRegisterTaskState = Type("TaskState", func() {
	Attribute("date", String, func() {
		Description("Date of the status creation")
		Example("2019-10-12T07:20:50.52Z")
	})
	Attribute("status", String, func() {
		Description("Status of the registration process")
		Example(state.StatusNames()[0])
		Enum(InterfaceSlice(state.StatusNames())...)
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
	Required("file")
})

// ImageUploadResult is image upload result.
var ImageUploadResult = ResultType("application/vnd.artwork.upload-image", func() {
	TypeName("Image")
	Attributes(func() {
		Attribute("image_id", Int, func() {
			Description("Uploaded image ID")
			Example(1)
		})
		Attribute("expires_in", String, func() {
			Description("Image expiration")
			Format(FormatDateTime)
			Example("2019-10-12T07:20:50.52Z")
		})
	})
	Required("image_id", "expires_in")
})

// RegisterTaskPayload represents a payload for returning task.
var RegisterTaskPayload = Type("RegisterTaskPayload", func() {
	Attribute("taskId", Int, "Task ID of the registration process", func() {
		TypeName("taskID")
		Minimum(1)
		Example(5)
	})
	Required("taskId")
})
