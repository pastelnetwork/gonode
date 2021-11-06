package pastel

// IDTicketType represents a type of the id ticket.
type IDTicketType string

// List of types of id ticket.
const (
	IDTicketAll      IDTicketType = "all"
	IDTicketMine     IDTicketType = "mine"
	IDTicketMN       IDTicketType = "mn"
	IDTicketPersonal IDTicketType = "personal"
)

// IDTickets represents multiple IDTicket.
type IDTickets []IDTicket

// IDTicket represents pastel id ticket.
type IDTicket struct {
	IDTicketProp `json:"ticket"`
	Height       int    `json:"height"`
	TXID         string `json:"txid"`
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
	PqKey     string `json:"pq_key"`
}
