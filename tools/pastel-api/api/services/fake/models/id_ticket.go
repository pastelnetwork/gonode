package models

// IDTicket represents a record of the ticket in the blockchain.
type IDTicket struct {
	Height int          `json:"height"`
	Ticket IDTicketProp `json:"ticket"`
	TXID   string       `json:"txid"`
}

// IDTicketProp represents properties of the ticket.
type IDTicketProp struct {
	Address   string `json:"address"`
	IDType    string `json:"id_type"`
	Outpoint  string `json:"outpoint"`
	PastelID  string `json:"pastelID"`
	Signature string `json:"signature"`
	TimeStamp string `json:"timeStamp"`
	Type      string `json:"type"`
}

// IDTickets represents pastel id ticket, that can be retrieved using the command `tickets list id mine`.
type IDTickets []IDTicket
