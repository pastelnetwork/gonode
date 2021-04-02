package models

type IdTicket struct {
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

type PastelID struct {
	PastelID string `json:"PastelID"`
}
