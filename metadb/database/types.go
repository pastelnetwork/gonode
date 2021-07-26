package database

import "fmt"

// Userdata represents userdata process payload
type Userdata struct {
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
	AvatarImage UserImageUpload
	// Cover photo of the user
	CoverPhoto UserImageUpload
	// Artist's PastelID
	ArtistPastelID string
	// Epoch Timestamp of the request (generated, not sending by UI)
	Timestamp int
	// Signature of the message (generated, not sending by UI)
	Signature string
	// Previous block hash in the chain (generated, not sending by UI)
	PreviousBlockHash string
}

// Store user's upload image
type UserImageUpload struct {
	// File to upload
	Content []byte
	// File name of the user image
	Filename string
}

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
}

func (d *Userdata) ToWriteCommand() UserdataWriteCommand {
	var avatarImageHex string
	var coverPhotoHex string
	if d.AvatarImage.Content != nil || len(d.AvatarImage.Content) == 0 {
		avatarImageHex = fmt.Sprintf("%x", d.AvatarImage.Content)
	}
	if d.CoverPhoto.Content != nil || len(d.CoverPhoto.Content) == 0 {
		coverPhotoHex = fmt.Sprintf("%x", d.CoverPhoto.Content)
	}
	return UserdataWriteCommand{
		Realname:           d.Realname,
		FacebookLink:       d.FacebookLink,
		TwitterLink:        d.TwitterLink,
		NativeCurrency:     d.NativeCurrency,
		Location:           d.Location,
		PrimaryLanguage:    d.PrimaryLanguage,
		Categories:         d.Categories,
		Biography:          d.Biography,
		AvatarImage:        avatarImageHex,
		AvatarFilename:     d.AvatarImage.Filename,
		CoverPhotoImage:    coverPhotoHex,
		CoverPhotoFilename: d.CoverPhoto.Filename,
		ArtistPastelID:     d.ArtistPastelID,
		Timestamp:          d.Timestamp,
		Signature:          d.Signature,
		PreviousBlockHash:  d.PreviousBlockHash,
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
}

func (d *UserdataReadResult) ToUserData() Userdata {
	return Userdata{
		Realname:        d.Realname,
		FacebookLink:    d.FacebookLink,
		TwitterLink:     d.TwitterLink,
		NativeCurrency:  d.NativeCurrency,
		Location:        d.Location,
		PrimaryLanguage: d.PrimaryLanguage,
		Categories:      d.Categories,
		Biography:       d.Biography,
		AvatarImage: UserImageUpload{
			Content:  d.AvatarImage,
			Filename: d.AvatarFilename,
		},
		CoverPhoto: UserImageUpload{
			Content:  d.CoverPhotoImage,
			Filename: d.CoverPhotoFilename,
		},
		ArtistPastelID:    d.ArtistPastelID,
		Timestamp:         d.Timestamp,
		Signature:         d.Signature,
		PreviousBlockHash: d.PreviousBlockHash,
	}
}
