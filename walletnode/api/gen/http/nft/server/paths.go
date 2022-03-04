// Code generated by goa v3.4.3, DO NOT EDIT.
//
// HTTP request path constructors for the nft service.
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package server

import (
	"fmt"
)

// RegisterNftPath returns the URL path to the nft service register HTTP endpoint.
func RegisterNftPath() string {
	return "/nfts/register"
}

// RegisterTaskStateNftPath returns the URL path to the nft service registerTaskState HTTP endpoint.
func RegisterTaskStateNftPath(taskID string) string {
	return fmt.Sprintf("/nfts/register/%v/state", taskID)
}

// RegisterTaskNftPath returns the URL path to the nft service registerTask HTTP endpoint.
func RegisterTaskNftPath(taskID string) string {
	return fmt.Sprintf("/nfts/register/%v", taskID)
}

// RegisterTasksNftPath returns the URL path to the nft service registerTasks HTTP endpoint.
func RegisterTasksNftPath() string {
	return "/nfts/register"
}

// UploadImageNftPath returns the URL path to the nft service uploadImage HTTP endpoint.
func UploadImageNftPath() string {
	return "/nfts/register/upload"
}

// NftSearchNftPath returns the URL path to the nft service nftSearch HTTP endpoint.
func NftSearchNftPath() string {
	return "/nfts/search"
}

// NftGetNftPath returns the URL path to the nft service nftGet HTTP endpoint.
func NftGetNftPath(txid string) string {
	return fmt.Sprintf("/nfts/%v", txid)
}

// DownloadNftPath returns the URL path to the nft service download HTTP endpoint.
func DownloadNftPath() string {
	return "/nfts/download"
}
