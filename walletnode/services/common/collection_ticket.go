package common

import (
	"context"
)

// CollectionRegistrationRequest represents a request for the registration collection
type CollectionRegistrationRequest struct {
	// CollectionName is name of the collection
	CollectionName string `json:"collection_name"`
	// ItemType is type of item collection will store
	ItemType string `json:"item_type"`
	//BurnTXID is id of burn transaction
	BurnTXID string `json:"burn_txid"`
	// ListOfPastelIDsOfAuthorizedContributors is the list of authorized contributors' pastelIDs
	ListOfPastelIDsOfAuthorizedContributors []string `json:"list_of_pastelids_of_authorized_contributors"`
	// MaxCollectionEntries in the collection's max entries
	MaxCollectionEntries int `json:"max_collection_entries"`
	// CollectionFinalAllowedBlockHeight is the height of final allowed block in days
	NoOfDaysToFinalizeCollection int `json:"NoOfDaysToFinalizeCollection"`
	// CollectionItemCopyCount is the collection item copy count
	CollectionItemCopyCount int `json:"collection_item_copy_count"`
	// Royalty represents the royalty fee
	Royalty float64 `json:"royalty"`
	// Green represents the green
	Green bool `json:"green"`
	//AppPastelID is the PastelID of the owner
	AppPastelID string `json:"app_pastel_id"`
	//AppPastelIDPassphrase is the passphrase of the owner
	AppPastelIDPassphrase string `json:"app_pastel_id_passphrase"`
	//MaxPermittedOpenNSFWScore is the MaxPermittedOpenNFSWScore allowed for collection items
	MaxPermittedOpenNSFWScore float64 `json:"max_permitted_open_nsfw_score"`
	//MaxPermittedOpenNSFWScore is the MinimumSimilarityScoreToFirstEntryInCollection allowed for collection items
	MinimumSimilarityScoreToFirstEntryInCollection float64 `json:"minimum_similarity_score_to_first_entry_in_collection"`
	//SpendableAddress is the SpendableAddress for 10% burning of reg fee
	SpendableAddress string `json:"spendable_address"`
}

// CollectionVerifier represents an interface with the func that checks if nft/sense if valid to enter in collection
type CollectionVerifier interface {
	IsValidForCollection(ctx context.Context) error
}
