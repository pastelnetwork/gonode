package pastel

// TradeTicket represents pastel Trade ticket
type TradeTicket struct {
	Height int             `json:"height"`
	TXID   string          `json:"txid"`
	Ticket TradeTicketData `json:"ticket"`
}

type TradeTicketData struct {
	Type             string `json:"type"` // "art-trade",
	Version          int    `json:"version"`
	PastelID         string `json:"pastelID"`          // PastelID of the buyer
	SellTXID         string `json:"sell_txid"`         // txid with sale ticket
	BuyTXID          string `json:"buy_txid"`          // txid with buy ticket
	ArtTXID          string `json:"art_txid"`          // txid with either 1) art activation ticket or 2) trade ticket in it
	RegistrationTXID string `json:"registration_txid"` // txid with registration ticket
	Price            string `json:"price"`
	CopySerialNR     string `json:"copy_serial_nr"`
	Signature        string `json:"signature"`
}
