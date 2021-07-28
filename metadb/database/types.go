package database

import (
	"fmt"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/supernode"
)

// UserdataDBRecord represents userdata record in DB
type UserdataWriteCommand struct {
	// Real name of the user
	Realname string
	// Facebook link of the user
	FacebookLink string
	// Twitter link of the user
	TwitterLink string
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency string
	// Location of the user
	Location string
	// Primary language of the user
	PrimaryLanguage string
	// The categories of user's work
	Categories string
	// Biography of the user
	Biography string
	// Avatar image of the user
	AvatarImage string
	// Avatar filename of the user
	AvatarFilename string
	// Cover photo image of the user
	CoverPhotoImage string
	// Cover photo filename of the user
	CoverPhotoFilename string
	// Artist's PastelID
	ArtistPastelID string
	// Epoch Timestamp of the request (generated, not sending by UI)
	Timestamp int
	// Signature of the message (generated, not sending by UI)
	Signature string
	// Previous block hash in the chain (generated, not sending by UI)
	PreviousBlockHash string
	// UserdataHash represents UserdataProcessRequest's hash value, to make sure UserdataProcessRequest's integrity
	UserdataHash string
}

func extractImageInfo(img *pb.UserdataRequest_UserImageUpload) (string, string) {
	if img == nil {
		return "", ""
	}

	if img.GetContent() == nil {
		return "", img.GetFilename()
	}

	return fmt.Sprintf("%x", img.GetContent()), img.GetFilename()
}

func pbToWriteCommand(d pb.UserdataRequest) UserdataWriteCommand {
	avatarImageHex, avatarImageFilename := extractImageInfo(d.GetAvatarImage())
	coverPhotoHex, coverPhotoFilename := extractImageInfo(d.GetCoverPhoto())

	return UserdataWriteCommand{
		Realname:           d.GetRealname(),
		FacebookLink:       d.GetFacebookLink(),
		TwitterLink:        d.GetTwitterLink(),
		NativeCurrency:     d.GetNativeCurrency(),
		Location:           d.GetLocation(),
		PrimaryLanguage:    d.GetPrimaryLanguage(),
		Categories:         d.GetCategories(),
		Biography:          d.GetBiography(),
		AvatarImage:        avatarImageHex,
		AvatarFilename:     avatarImageFilename,
		CoverPhotoImage:    coverPhotoHex,
		CoverPhotoFilename: coverPhotoFilename,
		ArtistPastelID:     d.GetArtistPastelID(),
		Timestamp:          int(d.GetTimestamp()),
		Signature:          d.GetSignature(),
		PreviousBlockHash:  d.GetPreviousBlockHash(),
		UserdataHash:       d.GetUserdataHash(),
	}
}

// UserdataReadResult represents userdata record reading from DB
type UserdataReadResult struct {
	// Real name of the user
	Realname string `mapstructure:"real_name"`
	// Facebook link of the user
	FacebookLink string `mapstructure:"facebook_link"`
	// Twitter link of the user
	TwitterLink string `mapstructure:"twitter_link"`
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency string `mapstructure:"native_currency"`
	// Location of the user
	Location string `mapstructure:"location"`
	// Primary language of the user
	PrimaryLanguage string `mapstructure:"primary_language"`
	// The categories of user's work
	Categories string `mapstructure:"categories"`
	// Biography of the user
	Biography string `mapstructure:"biography"`
	// Avatar image of the user
	AvatarImage []byte `mapstructure:"avatar_image"`
	// Avatar filename of the user
	AvatarFilename string `mapstructure:"avatar_filename"`
	// Cover photo image of the user
	CoverPhotoImage []byte `mapstructure:"cover_photo_image"`
	// Cover photo filename of the user
	CoverPhotoFilename string `mapstructure:"cover_photo_filename"`
	// Artist's PastelID
	ArtistPastelID string `mapstructure:"artist_pastel_id"`
	// Epoch Timestamp of the request (generated, not sending by UI)
	Timestamp int `mapstructure:"timestamp"`
	// Signature of the message (generated, not sending by UI)
	Signature string `mapstructure:"signature"`
	// Previous block hash in the chain (generated, not sending by UI)
	PreviousBlockHash string `mapstructure:"previous_block_hash"`
	// UserdataHash represents UserdataProcessRequest's hash value, to make sure UserdataProcessRequest's integrity
	UserdataHash string `mapstructure:"user_data_hash"`
}

func (d *UserdataReadResult) ToUserData() userdata.UserdataProcessRequest {
	return userdata.UserdataProcessRequest{
		Realname:        d.Realname,
		FacebookLink:    d.FacebookLink,
		TwitterLink:     d.TwitterLink,
		NativeCurrency:  d.NativeCurrency,
		Location:        d.Location,
		PrimaryLanguage: d.PrimaryLanguage,
		Categories:      d.Categories,
		Biography:       d.Biography,
		AvatarImage: userdata.UserImageUpload{
			Content:  d.AvatarImage,
			Filename: d.AvatarFilename,
		},
		CoverPhoto: userdata.UserImageUpload{
			Content:  d.CoverPhotoImage,
			Filename: d.CoverPhotoFilename,
		},
		ArtistPastelID:    d.ArtistPastelID,
		Timestamp:         int64(d.Timestamp),
		PreviousBlockHash: d.PreviousBlockHash,
	}
}
