package design

import (
	//revive:disable:dot-imports
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports
)

var _ = Service("tickets", func() {
	Description("The tickets serves tickets data.")

	HTTP(func() {
		Path("/tickets")
	})

	Method("add", func() {
		Description("Add a new ticket and return its ID.")
		Payload(func() {
			Extend(Ticket)
			Required("first_name", "last_name", "ticket")
		})
		Result(String)

		HTTP(func() {
			POST("/")
			Response(StatusCreated)
		})
	})

	Method("list", func() {
		Description("List all stored tickets.")
		Result(CollectionOf(StoredTicket))

		HTTP(func() {
			GET("/")
			Response(StatusOK)
		})
	})
})

// Ticket describes an ticket from an customer to be stored.
var Ticket = Type("Ticket", func() {
	Description("Ticket describes an ticket from an customer to be stored.")
	Attribute("address", String, "Address", func() {
		MaxLength(32)
		// Example("")
	})
	Attribute("pastel_id", String, "PastelID", func() {
		MaxLength(32)
		Example("123456789")
	})
	Attribute("signature", String, "Signature", func() {
		Example("wdk231r23kkl23rn23jk423nkl.")
	})
})

// StoredTicket describes a ticket retrieved by the ticket service.
var StoredTicket = ResultType("application/vnd.stored-ticket", func() {
	Description("StoredTicket describes a ticket retrieved by the ticket service.")
	Reference(Ticket)
	TypeName("StoredTicket")

	Attributes(func() {
		// Attribute("id", String, "ID is the unique id of the ticket.", func() {
		// 	Example("13243")
		// })
		Field(2, "address")
		Field(3, "pastel_id")
		Field(4, "signature")
	})

	Required("address", "pastel_id", "signature")
})
