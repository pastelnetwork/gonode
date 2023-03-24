package pastel

// ActTicketType represents a type of the activation ticket.
type ActTicketType string

// List of types of activation ticket.
const (
	ActTicketAll       ActTicketType = "all"
	ActTicketAvailable ActTicketType = "available"
	ActTicketSold      ActTicketType = "sold"
)

// ActTickets is a collection of ActTicket
type ActTickets []ActTicket

// ActTicket represents pastel activation ticket.
type ActTicket struct {
	Height        int           `json:"height"`
	TXID          string        `json:"txid"`
	ActTicketData ActTicketData `json:"ticket"`
}

// ActTicketData represents activation ticket properties
type ActTicketData struct {
	PastelID      string `json:"pastelID"`
	Signature     string `json:"signature"`
	Type          string `json:"type"`
	CreatorHeight int    `json:"creator_height"`
	RegTXID       string `json:"reg_txid"`
	StorageFee    int    `json:"storage_fee"`
	Version       int    `json:"version"`
	CalledAt      int    `json:"called_at"`
}
