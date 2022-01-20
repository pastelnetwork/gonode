package artworkdownload

import "github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"

// NftDownloadRequest represents registered artwork ticket.
type NftDownloadRequest struct {
	Txid               string `json:"txid"`
	PastelID           string `json:"pastel_id"`
	PastelIDPassphrase string `json:"pastel_id_passphrase"`
}

// NFT Download
func FromDownloadPayload(payload *artworks.ArtworkDownloadPayload) *NftDownloadRequest {
	return &NftDownloadRequest{
		Txid:               payload.Txid,
		PastelID:           payload.Pid,
		PastelIDPassphrase: payload.Key,
	}
}
