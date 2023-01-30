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
	// OpenAPISubsetID is field
	OpenAPISubsetID string `json:"open_api_subset_id"`
	// FileName to save in ticket
	FileName string `json:"file_name"`
	//MakePubliclyAccessible to save in ticket
	MakePubliclyAccessible bool `json:"make_publicly_accessible"`
}
