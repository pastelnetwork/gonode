package common

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
)

// ActionRegistrationRequest represents a request for the registration sense
type ActionRegistrationRequest struct {
	// Image is field
	Image *files.File `json:"image"`
	// BurnTxID is field
	BurnTxID string `json:"burn_txid"`
	// AppPastelID is field
	AppPastelID string `json:"app_pastel_id"`
	// AppPastelIDPassphrase is field
	AppPastelIDPassphrase string `json:"app_pastel_id_passphrase"`
}
