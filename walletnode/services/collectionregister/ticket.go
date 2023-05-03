package collectionregister

import (
	"github.com/pastelnetwork/gonode/walletnode/api/gen/collection"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// FromCollectionRegistrationPayload convert CollectionRegistrationPayload to CollectionRegistrationRequest
func FromCollectionRegistrationPayload(payload *collection.RegisterCollectionPayload) *common.CollectionRegistrationRequest {
	c := &common.CollectionRegistrationRequest{
		CollectionName:                          payload.CollectionName,
		ItemType:                                payload.ItemType,
		ListOfPastelIDsOfAuthorizedContributors: payload.ListOfPastelidsOfAuthorizedContributors,
		MaxCollectionEntries:                    payload.MaxCollectionEntries,
		CollectionFinalAllowedBlockHeight:       payload.CollectionFinalAllowedBlockHeight,
		AppPastelID:                             payload.AppPastelID,
		Green:                                   payload.Green,
		MaxPermittedOpenNSFWScore:               payload.MaxPermittedOpenNsfwScore,
		Royalty:                                 payload.Royalty,
		CollectionItemCopyCount:                 payload.CollectionItemCopyCount,
		MinimumSimilarityScoreToFirstEntryInCollection: payload.MinimumSimilarityScoreToFirstEntryInCollection,
	}

	if payload.Key != nil {
		c.AppPastelIDPassphrase = *payload.Key
	}

	return c
}
