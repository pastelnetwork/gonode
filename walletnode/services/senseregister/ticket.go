package senseregister

import (
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// FromStartProcessingPayload convert StartProcessingPayload to ActionRegistrationRequest
func FromStartProcessingPayload(payload *sense.StartProcessingPayload) *common.ActionRegistrationRequest {
	return &common.ActionRegistrationRequest{
		BurnTxID:              payload.BurnTxid,
		AppPastelID:           payload.AppPastelID,
		AppPastelIDPassphrase: payload.AppPastelidPassphrase,
	}
}
