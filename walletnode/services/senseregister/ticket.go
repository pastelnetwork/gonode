package senseregister

import (
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
)

// SenseRegisterRequest represents a request for the registration sense
type SenseRegisterRequest struct {
	Image                 *files.File `json:"image"`
	BurnTxID              string      `json:"burn_txid"`
	AppPastelID           string      `json:"app_pastel_id"`
	AppPastelIDPassphrase string      `json:"app_pastel_id_passphrase"`
}

func FromSenseRegisterPayload(payload *sense.StartProcessingPayload) *SenseRegisterRequest {
	return &SenseRegisterRequest{
		BurnTxID:              payload.BurnTxid,
		AppPastelID:           payload.AppPastelID,
		AppPastelIDPassphrase: payload.AppPastelidPassphrase,
	}
}

//func toNftRegisterTicket(ticket *SenseRegisterRequest) *sense.SenseTicket {
//}
