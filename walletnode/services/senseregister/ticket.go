package senseregister

import "github.com/pastelnetwork/gonode/common/service/artwork"

// Request represents a request for the registration sense
type Request struct {
	Image                 *artwork.File `json:"image"`
	BurnTxID              string        `json:"burn_txid"`
	AppPastelID           string        `json:"app_pastel_id"`
	AppPastelIDPassphrase string        `json:"app_pastel_id_passphrase"`

	// FIXME: change it later
	SpendableAddress string  `json:"spendable_address"`
	MaximumFee       float64 `json:"maximum_fee"`
}
