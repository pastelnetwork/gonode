package models

// PastelIDVerify represents the result of the data verification by signature, command `pastelid verify "text" "signature" "PastelID"`.
type PastelIDVerify struct {
	Verification string `json:"verification"`
}
