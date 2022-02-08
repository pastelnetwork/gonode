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
	ArtistName               string              `json:"nft_name"`
	ArtistPastelid           string              `json:"nft_pastelid"`
	ArtistPastelidPassphrase string              `json:"nft_pastelid_passphrase"`
	ArtistWebsiteURL         string              `json:"nft_website_url"`
	Description              string              `json:"description"`
	Green                    bool                `json:"green"`
	ImageID                  string              `json:"image_id"`
	IssuedCopies             int                 `json:"issued_copies"`
	Keywords                 string              `json:"keywords"`
	MaximumFee               int                 `json:"maximum_fee"`
	Name                     string              `json:"name"`
	Royalty                  int                 `json:"royalty"`
	SeriesName               string              `json:"series_name"`
	SpendableAddress         string              `json:"spendable_address"`
	YoutubeURL               string              `json:"youtube_url"`
	ThumbnailCoordinate      ThumbnailCoordinate `json:"thumbnail_coordinate"`
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
