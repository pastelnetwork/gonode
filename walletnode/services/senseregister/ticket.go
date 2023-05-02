package senseregister

import (
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// FromStartProcessingPayload convert StartProcessingPayload to ActionRegistrationRequest
func FromStartProcessingPayload(payload *sense.StartProcessingPayload) *common.ActionRegistrationRequest {
	req := &common.ActionRegistrationRequest{
		BurnTxID:              payload.BurnTxid,
		AppPastelID:           payload.AppPastelID,
		AppPastelIDPassphrase: payload.Key,
		OpenAPISubsetID:       payload.OpenAPISubsetID,
	}

	if payload.OpenAPIGroupID != nil {
		req.GroupID = *payload.OpenAPIGroupID
	}

	if payload.CollectionActTxid != nil {
		req.CollectionTxID = *payload.CollectionActTxid
	}

	return req
}
