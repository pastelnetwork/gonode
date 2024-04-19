package download

import "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"

// NftDownloadingRequest represents registered nft ticket.
type NftDownloadingRequest struct {
	// Txid is field
	Txid string `json:"txid"`
	// PastelID is field
	PastelID string `json:"pastel_id"`
	// PastelIDPassphrase is field
	PastelIDPassphrase string `json:"pastel_id_passphrase"`
	// Type
	Type string `json:"type"`

	HashOnly bool `json:"hash_only"`
}

// FromDownloadPayload is NFT Download request
func FromDownloadPayload(payload *nft.DownloadPayload, ticketType string, hashOnly bool) *NftDownloadingRequest {
	return &NftDownloadingRequest{
		Txid:               payload.Txid,
		PastelID:           payload.Pid,
		PastelIDPassphrase: payload.Key,
		Type:               ticketType,
		HashOnly:           hashOnly,
	}
}
