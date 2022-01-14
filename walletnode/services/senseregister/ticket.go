package senseregister

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
)

// SenseRegisterRequest represents a request for the registration sense
type SenseRegisterRequest struct {
	Image                 *files.File `json:"image"`
	BurnTxID              string      `json:"burn_txid"`
	AppPastelID           string      `json:"app_pastel_id"`
	AppPastelIDPassphrase string      `json:"app_pastel_id_passphrase"`
}
