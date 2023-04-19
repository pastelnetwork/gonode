package collectionregister

import (
	"github.com/pastelnetwork/gonode/walletnode/api/gen/collection"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// FromCollectionRegistrationPayload convert CollectionRegistrationPayload to CollectionRegistrationRequest
func FromCollectionRegistrationPayload(payload *collection.RegisterCollectionPayload) *common.CollectionRegistrationRequest {
	c := &common.CollectionRegistrationRequest{
		CollectionName:                          payload.CollectionName,
		ListOfPastelIDsOfAuthorizedContributors: payload.ListOfPastelidsOfAuthorizedContributors,
		MaxCollectionEntries:                    payload.MaxCollectionEntries,
		CollectionFinalAllowedBlockHeight:       payload.CollectionFinalAllowedBlockHeight,
		AppPastelID:                             payload.AppPastelID,
	}

	if payload.Key != nil {
		c.AppPastelIDPassphrase = *payload.Key
	}

	if payload.Royalty != nil {
		c.Royalty = *payload.Royalty
	}

	if payload.Green != nil {
		c.Green = *payload.Green
	}

	if payload.CollectionItemCopyCount != nil {
		c.CollectionItemCopyCount = *payload.CollectionItemCopyCount
	}

	return c
}
