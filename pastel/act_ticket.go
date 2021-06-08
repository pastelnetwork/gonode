package pastel

// ActivationTicketType represents a type of the activation ticket.
type ActivationTicketType string

// List of types of activation ticket.
const (
	ActTicketAll       ActivationTicketType = "all"
	ActTicketAvailable ActivationTicketType = "available"
	ActTicketSold      ActivationTicketType = "sold"
)

// ActivationTickets is a collection of ActivationTicket
type ActivationTickets []ActivationTicket

// ActivationTicket represents pastel activation ticket.
type ActivationTicket struct {
	PastelID     string `json:"pastelID"`
	Signature    string `json:"signature"`
	Type         string `json:"type"`
	ArtistHeight int    `json:"artist_height"`
	RegTXID      string `json:"reg_txid"`
	RegFee       int    `json:"reg_fee"`
}
