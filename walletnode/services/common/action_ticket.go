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
	// FileName to save in ticket
	FileName string `json:"file_name"`
	//MakePubliclyAccessible to save in ticket
	MakePubliclyAccessible bool `json:"make_publicly_accessible"`
	//CollectionTxID to add it into the collection
	CollectionTxID string `json:"collection_txid"`
	// GroupID to add it into the collection
	GroupID string `json:"group_id"`
	// SpendableAddress to use for registration fee
	SpendableAddress string `json:"spendable_address"`
}
