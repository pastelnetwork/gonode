package cascaderegister

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// FromTaskPayload convert TaskPayload to ActionRegistrationRequest
func FromTaskPayload(payload *common.AddTaskPayload) *common.ActionRegistrationRequest {
	req := &common.ActionRegistrationRequest{
		BurnTxID:               *payload.BurnTxid,
		AppPastelID:            payload.AppPastelID,
		AppPastelIDPassphrase:  payload.Key,
		MakePubliclyAccessible: payload.MakePubliclyAccessible,
	}

	if payload.SpendableAddress != nil {
		req.SpendableAddress = *payload.SpendableAddress
	}

	return req
}
