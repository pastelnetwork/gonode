package helper

import "fmt"

const ()

type UploadImageReq struct {
	File     []byte
	Filename string
}

type UploadImageResp struct {
	ImageID string `json:"image_id"`
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

type SenseActionDetailReq struct {
	AppPastelID         string `json:"app_pastelid"`
	ActionDataHash      string `json:"action_data_hash"`
	ActionDataSignature string `json:"action_data_signature"`
}

type SenseActionDetailResp struct {
	EstimatedFee float64 `json:"estimated_fee"`
}

type SenseStartTaskReq struct {
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
