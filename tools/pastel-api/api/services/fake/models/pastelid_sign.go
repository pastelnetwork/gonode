package models

// PastelIDSign represents the signature of the data by pastelID, command `pastelid sign "text" "PastelID" "passphrase"`.
type PastelIDSign struct {
	Signature string `json:"signature"`
}
