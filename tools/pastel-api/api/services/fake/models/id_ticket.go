package models

// IDTicketRecord represents a record of the ticket in the blockchain.
type IDTicketRecord struct {
	Height int      `json:"height"`
	Ticket IDTicket `json:"ticket"`
	TXID   string   `json:"txid"`
}

// IDTicket represents properties of the ticket.
type IDTicket struct {
	Address   string `json:"address"`
	IDType    string `json:"id_type"`
	Outpoint  string `json:"outpoint"`
	PastelID  string `json:"pastelID"`
	Signature string `json:"signature"`
	TimeStamp string `json:"timeStamp"`
	Type      string `json:"type"`
}

// IDTicketRecords represents the API response that can be retrieved using the command `tickets list id mine`.
type IDTicketRecords []IDTicketRecord
