package pastel

// TradeTicket represents pastel Trade ticket
type TradeTicket struct {
	Type      string `json:"type"`      // "trade",
	PastelID  string `json:"pastelID"`  // PastelID of the buyer
	SellTXID  string `json:"sell_txid"` // txid with sale ticket
	BueTXID   string `json:"buy_txid"`  // txid with buy ticket
	ArtTXID   string `json:"art_txid"`  // txid with either 1) art activation ticket or 2) trade ticket in it
	Price     string `json:"price"`
	Reserved  string `json:"reserved"`
	Signature string `json:"signature"`
}
