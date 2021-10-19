package externalstorage

import "github.com/pastelnetwork/gonode/common/service/artwork"

// Request represents external dupe detection request.
type Request struct {
	Image              *artwork.File `json:"image"`
	PastelID           string        `json:"pastel_id"`
	PastelIDPassphrase string        `json:"pastel_id_passphrase"`
	ImageHash          string        `json:"image_hash"`
	SendingAddress     string        `json:"spendable_address"`
	MaximumFee         float64       `json:"maximum_fee"`
	RequestIdentifier  string        `json:"request_identifier"`
}
