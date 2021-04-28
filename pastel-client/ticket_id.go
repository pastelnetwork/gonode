package pastel

// IDTickets is multiple IDTicket.
type IDTickets []IDTicket

// IDTicket represensts pastel id ticket.
type IDTicket struct {
	Height int `json:"height"`
	Ticket struct {
		Address   string `json:"address"`
		IDType    string `json:"id_type"`
		PastelID  string `json:"pastelID"`
		Signature string `json:"signature"`
		TimeStamp string `json:"timeStamp"`
		Type      string `json:"type"`
	} `json:"ticket"`
	Txid string `json:"txid"`
}
