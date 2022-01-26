package nftdownload

import "github.com/pastelnetwork/gonode/walletnode/api/gen/nft"

// NftDownloadRequest represents registered nft ticket.
type NftDownloadRequest struct {
	Txid               string `json:"txid"`
	PastelID           string `json:"pastel_id"`
	PastelIDPassphrase string `json:"pastel_id_passphrase"`
}

// NFT Download
func FromDownloadPayload(payload *nft.NftDownloadPayload) *NftDownloadRequest {
	return &NftDownloadRequest{
		Txid:               payload.Txid,
		PastelID:           payload.Pid,
		PastelIDPassphrase: payload.Key,
	}
}
