package nft

import "fmt"

const ()

type UploadImageReq struct {
	File     []byte
	Filename string
}

type UploadImageResp struct {
	ImageID string
}

type ThumbnailCoordinate struct {
	BottomRightX int
	BottomRightY int
	TopLeftX     int
	TopLeftY     int
}

type RegistrationReq struct {
	ArtistName               string
	ArtistPastelid           string
	ArtistPastelidPassphrase string
	ArtistWebsiteURL         string
	Description              string
	Green                    bool
	ImageID                  string
	IssuesCopies             int
	Keywords                 string
	MaximumFee               int
	Name                     string
	Royalty                  int
	SeriesName               string
	SpendableAddress         string
	YoutubeURL               string
	ThumbnailCoordinate      ThumbnailCoordinate
}

type RegistrationResp struct {
	TaskID string
}

func GetUploadImageURI(baseURI string) string {
	return fmt.Sprintf("%s/%s/", baseURI, "artworks/register/upload")
}

func GetRegistrationURI(baseURI string) string {
	return fmt.Sprintf("%s/%s/", baseURI, "artworks/register")
}
