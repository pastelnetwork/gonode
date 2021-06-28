package artworkdownload

// Ticket represents registered artwork ticket.
type Ticket struct {
	TXID               string `json:"txid"`
	PastelID           string `json:"pastel_id"`
	PastelIDPassphrase string `json:"pastel_id_passphrase"`
}
