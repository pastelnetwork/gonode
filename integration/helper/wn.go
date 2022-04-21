package helper

import (
	"fmt"
	"strings"
)

type UploadImageReq struct {
	File     []byte
	Filename string
}

type UploadImageResp struct {
	ImageID      string  `json:"image_id"`
	EstimatedFee float64 `json:"estimated_fee"`
}

type UploadAssetResp struct {
	FileID       string  `json:"file_id"`
	EstimatedFee float64 `json:"estimated_fee"`
}

type DownloadResp struct {
	File string `json:"file"`
}

type ThumbnailCoordinate struct {
	BottomRightX int `json:"bottom_right_x"`
	BottomRightY int `json:"bottom_right_y"`
	TopLeftX     int `json:"top_left_x"`
	TopLeftY     int `json:"top_left_y"`
}

type RegistrationReq struct {
	CreatorName               string              `json:"creator_name"`
	CreatorPastelid           string              `json:"creator_pastelid"`
	CreatorPastelidPassphrase string              `json:"creator_pastelid_passphrase"`
	CreatorWebsiteURL         string              `json:"creator_website_url"`
	Description               string              `json:"description"`
	Green                     bool                `json:"green"`
	ImageID                   string              `json:"image_id"`
	IssuedCopies              int                 `json:"issued_copies"`
	Keywords                  string              `json:"keywords"`
	MaximumFee                int                 `json:"maximum_fee"`
	Name                      string              `json:"name"`
	Royalty                   int                 `json:"royalty"`
	SeriesName                string              `json:"series_name"`
	SpendableAddress          string              `json:"spendable_address"`
	YoutubeURL                string              `json:"youtube_url"`
	ThumbnailCoordinate       ThumbnailCoordinate `json:"thumbnail_coordinate"`
}

type SenseCascadeStartTaskReq struct {
	AppPastelID           string `json:"app_pastelid"`
	BurnTXID              string `json:"burn_txid"`
	AppPastelIDPassphrase string `json:"app_pastelid_passphrase"`
}

type RegistrationResp struct {
	TaskID string `json:"task_id"`
}

func GetUploadImageURI(baseURI string) string {
	return fmt.Sprintf("%s/%s", baseURI, "nfts/register/upload")
}

func GetRegistrationURI(baseURI string) string {
	return fmt.Sprintf("%s/%s", baseURI, "nfts/register")
}

func GetNFTDetailURI(baseURI, txid string) string {
	return fmt.Sprintf("%s/%s/%s", baseURI, "nfts", txid)
}

func GetDownloadURI(baseURI, pid, txid string) string {
	return fmt.Sprintf("%s/%s?pid=%s&txid=%s", baseURI, "nfts/download", pid, txid)
}

// Sense

func GetSenseUploadImageURI(baseURI string) string {
	return fmt.Sprintf("%s/%s", baseURI, "openapi/sense/upload")
}

func GetSenseActionURI(baseURI, imageID string) string {
	return fmt.Sprintf("%s/%s/%s", baseURI, "openapi/sense/details", imageID)
}

func GetSenseStartTaskURI(baseURI, imageID string) string {
	return fmt.Sprintf("%s/%s/%s", baseURI, "openapi/sense/start", imageID)
}

func GetSenseTaskStateURI(taskID string) string {
	return fmt.Sprintf("openapi/sense/start/%s/state", taskID)
}

func GetNFTSearchURI(queryParams map[string]string) string {
	uri := "nfts/search?"
	for key, val := range queryParams {
		uri = uri + key + "=" + val + "&"
	}
	uri = strings.TrimSuffix(uri, "&")

	return uri
}

func GetSenseDownloadURI(baseURI, pid, txid string) string {
	return fmt.Sprintf("%s/%s?pid=%s&txid=%s", baseURI, "openapi/sense/download", pid, txid)
}

// Cascade
func GetCascadeUploadImageURI(baseURI string) string {
	return fmt.Sprintf("%s/%s", baseURI, "openapi/cascade/upload")
}

func GetCascadeActionURI(baseURI, imageID string) string {
	return fmt.Sprintf("%s/%s/%s", baseURI, "openapi/cascade/details", imageID)
}

func GetCascadeStartTaskURI(baseURI, imageID string) string {
	return fmt.Sprintf("%s/%s/%s", baseURI, "openapi/cascade/start", imageID)
}

func GetCascadeTaskStateURI(taskID string) string {
	return fmt.Sprintf("openapi/cascade/start/%s/state", taskID)
}

func GetCascadeDownloadURI(baseURI, pid, txid string) string {
	return fmt.Sprintf("%s/%s?pid=%s&txid=%s", baseURI, "openapi/cascade/download", pid, txid)
}
