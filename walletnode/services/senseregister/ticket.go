package senseregister

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
)

// SenseRegistrationRequest represents a request for the registration sense
type SenseRegistrationRequest struct {
	// Image is field
	Image *files.File `json:"image"`
	// BurnTxID is field
	BurnTxID string `json:"burn_txid"`
	// AppPastelID is field
	AppPastelID string `json:"app_pastel_id"`
	// AppPastelIDPassphrase is field
	AppPastelIDPassphrase string `json:"app_pastel_id_passphrase"`
}

// FromSenseRegisterPayload convert StartProcessingPayload to SenseRegistrationRequest
func FromSenseRegisterPayload(payload *sense.StartProcessingPayload) *SenseRegistrationRequest {
	return &SenseRegistrationRequest{
		BurnTxID:              payload.BurnTxid,
		AppPastelID:           payload.AppPastelID,
		AppPastelIDPassphrase: payload.AppPastelidPassphrase,
	}
}
