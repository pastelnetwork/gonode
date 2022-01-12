package senseregister

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
)

// Request represents a request for the registration sense
type Request struct {
	Image                 *files.File `json:"image"`
	BurnTxID              string      `json:"burn_txid"`
	AppPastelID           string      `json:"app_pastel_id"`
	AppPastelIDPassphrase string      `json:"app_pastel_id_passphrase"`
}
