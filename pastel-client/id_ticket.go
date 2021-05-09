package pastel

// IDTickets represents multiple IDTicket.
type IDTickets []IDTicket

// IDTicket represensts pastel id ticket.
type IDTicket struct {
	Height int          `json:"height"`
	Ticket IDTicketProp `json:"ticket"`
	TXID   string       `json:"txid"`
}

// IDTicketProp represents properties of the id ticket.
type IDTicketProp struct {
	Address   string `json:"address"`
	IDType    string `json:"id_type"`
	Outpoint  string `json:"outpoint"`
	PastelID  string `json:"pastelID"`
	Signature string `json:"signature"`
	TimeStamp string `json:"timeStamp"`
	Type      string `json:"type"`
}
