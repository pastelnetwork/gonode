package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("letters", func() {
	Description("The letters serves letters data.")

	Error("invalid-scopes", String, "Credentials scopes are invalid")

	HTTP(func() {
		Path("/letters")
	})

	Method("add", func() {
		Description("Add a new letter and return its ID.")
		Payload(func() {
			Extend(Letter)
			Required("token", "first_name", "last_name", "letter")
		})
		Result(String)

		HTTP(func() {
			POST("/")
			Response(StatusCreated)
			Response("invalid-scopes", StatusForbidden)
		})
	})

	Method("list", func() {
		Description("List all stored letters.")

		Security(OAuth2Auth)

		Payload(func() {
			AccessTokenField(1, "token", String)
			Required("token")
		})
		Result(CollectionOf(StoredLetter))

		HTTP(func() {
			GET("/")
			Header("token:Authorization", String, "Auth token", func() {
				Pattern("^Bearer [^ ]+$")
			})
			Response(StatusOK)
		})
	})
})

// Letter describes an letter from an customer to be stored.
var Letter = Type("Letter", func() {
	Description("Letter describes an letter from an customer to be stored.")
	Attribute("first_name", String, "First Name", func() {
		MaxLength(32)
		Example("John")
	})
	Attribute("last_name", String, "Last Name", func() {
		MaxLength(32)
		Example("Smith")
	})
	Attribute("letter", String, "Letter", func() {
		Example("Order a new product.")
	})
})

// StoredLetter describes a letter retrieved by the letter service.
var StoredLetter = ResultType("application/vnd.biptec.stored-letter", func() {
	Description("StoredLetter describes a letter retrieved by the letter service.")
	Reference(Letter)
	TypeName("StoredLetter")

	Attributes(func() {
		Attribute("id", String, "ID is the unique id of the letter.", func() {
			Example("13243")
		})
		Field(2, "first_name")
		Field(3, "last_name")
		Field(4, "letter")
	})

	Required("first_name", "last_name", "letter")
})
