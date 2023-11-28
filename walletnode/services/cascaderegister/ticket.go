package cascaderegister

import (
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// FromStartProcessingPayload convert StartProcessingPayload to ActionRegistrationRequest
func FromStartProcessingPayload(payload *cascade.StartProcessingPayload) *common.ActionRegistrationRequest {
	return &common.ActionRegistrationRequest{
		BurnTxID:               payload.BurnTxid,
		AppPastelID:            payload.AppPastelID,
		AppPastelIDPassphrase:  payload.Key,
		MakePubliclyAccessible: payload.MakePubliclyAccessible,
		SpendableAddress:       *payload.SpendableAddress,
	}
}
