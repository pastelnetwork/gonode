package database

import (
	"fmt"
	"strings"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

const (
	SNActivityThumbnailRequest = "thumbnail_request"
	SNActivityNftSearch        = "nft_search"
	SNActivityCreatorSearch    = "creator_search"
	SNActivityUserSearch       = "user_search"
	SNActivityKeywordSearch    = "keyword_search"
)

// UserdataWriteCommand represents userdata record in DB
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

	return fmt.Sprintf("%x", img.GetContent()), processInputString(img.GetFilename())
}

func processEscapeString(s string) string {
	s = strings.Replace(s, `'`, `''`, -1)
	return s
}

func processInputString(s string) string {
	if s == "" {
		return "NULL"
	}
	return processEscapeString(s)
}

func pbToWriteCommand(d *pb.UserdataRequest) UserdataWriteCommand {
	avatarImageHex, avatarImageFilename := extractImageInfo(d.GetAvatarImage())
	coverPhotoHex, coverPhotoFilename := extractImageInfo(d.GetCoverPhoto())

	return UserdataWriteCommand{
		Realname:           processInputString(d.GetRealname()),
		FacebookLink:       processInputString(d.GetFacebookLink()),
		TwitterLink:        processInputString(d.GetTwitterLink()),
		NativeCurrency:     processInputString(d.GetNativeCurrency()),
		Location:           processInputString(d.GetLocation()),
		PrimaryLanguage:    processInputString(d.GetPrimaryLanguage()),
		Categories:         processInputString(d.GetCategories()),
		Biography:          processInputString(d.GetBiography()),
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

// ToUserData return the ProcessRequest
func (d *UserdataReadResult) ToUserData() userdata.ProcessRequest {
	return userdata.ProcessRequest{
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

// type ArtInfo struct {
// 	ArtID                 string  `mapstructure:"art_id",json:"art_id"`
// 	ArtistPastelID        string  `mapstructure:"artist_pastel_id",json:"artist_pastel_id"`
// 	Copies                int64   `mapstructure:"copies",json:"copies"`
// 	CreatedTimestamp      int64   `mapstructure:"created_timestamp",json:"created_timestamp"`
// 	GreenNft              bool    `mapstructure:"green_nft",json:"green_nft"`
// 	RarenessScore         float64 `mapstructure:"rareness_score",json:"rareness_score"`
// 	RoyaltyRatePercentage float64 `mapstructure:"royalty_rate_percentage",json:"royalty_rate_percentage"`
// }

// type ArtInstanceInfo struct {
// 	InstanceID    string   `mapstructure:"instance_id",json:"instance_id"`
// 	ArtID         string   `mapstructure:"art_id",json:"art_id"`
// 	OwnerPastelID string   `mapstructure:"owner_pastel_id",json:"owner_pastel_id"`
// 	Price         float64  `mapstructure:"price",json:"price"`
// 	AskingPrice   *float64 `mapstructure:"asking_price,omitempty",json:"asking_price,omitempty"`
// }

// type ArtLike struct {
// 	ArtID    string `mapstructure:"art_id",json:"art_id"`
// 	PastelID string `mapstructure:"pastel_id",json:"pastel_id"`
// }

// type ArtTransaction struct {
// 	TransactionID  string  `mapstructure:"transaction_id",json:"transaction_id"`
// 	InstanceID     string  `mapstructure:"instance_id",json:"instance_id"`
// 	Timestamp      int64   `mapstructure:"timestamp",json:"timestamp"`
// 	SellerPastelID string  `mapstructure:"seller_pastel_id",json:"seller_pastel_id"`
// 	BuyerPastelID  string  `mapstructure:"buyer_pastel_id",json:"buyer_pastel_id"`
// 	Price          float64 `mapstructure:"price",json:"price"`
// }

// type UserFollow struct {
// 	FollowerPastelID string `mapstructure:"follower_pastel_id",json:"follower_pastel_id"`
// 	FolloweePastelID string `mapstructure:"followee_pastel_id",json:"followee_pastel_id"`
// }

// type NftCreatedByArtistQueryResult struct {
// 	InstanceID            string  `mapstructure:"instance_id",json:"instance_id"`
// 	ArtID                 string  `mapstructure:"art_id",json:"art_id"`
// 	Copies                int64   `mapstructure:"copies",json:"copies"`
// 	CreatedTimestamp      int64   `mapstructure:"created_timestamp",json:"created_timestamp"`
// 	GreenNft              bool    `mapstructure:"green_nft",json:"green_nft"`
// 	RarenessScore         float64 `mapstructure:"rareness_score",json:"rareness_score"`
// 	RoyaltyRatePercentage float64 `mapstructure:"royalty_rate_percentage",json:"royalty_rate_percentage"`
// }

// type NftForSaleByArtistQueryResult struct {
// 	InstanceID string  `mapstructure:"instance_id",json:"instance_id"`
// 	ArtID      string  `mapstructure:"art_id",json:"art_id"`
// 	Price      float64 `mapstructure:"price",json:"price"`
// }

// type NftOwnedByUserQueryResult struct {
// 	ArtID string `mapstructure:"art_id",json:"art_id"`
// 	Count int    `mapstructure:"cnt",json:"cnt"`
// }

// type NftSoldByArtIDQueryResult struct {
// 	TotalCopies int `mapstructure:"total_copies",json:"total_copies"`
// 	SoldCopies  int `mapstructure:"sold_copies",json:"sold_copies"`
// }

// type UniqueNftByUserQuery struct {
// 	ArtistPastelID string `mapstructure:"artist_pastel_id",json:"artist_pastel_id"`
// 	LimitTimestamp int64  `mapstructure:"limit_timestamp",json:"limit_timestamp"`
// }

// type AskingPriceUpdateRequest struct {
// 	InstanceID  string  `mapstructure:"instance_id",json:"instance_id"`
// 	AskingPrice float64 `mapstructure:"asking_price",json:"asking_price"`
// }

// type ArtPlaceBidRequest struct {
// 	AuctionID int64   `mapstructure:"auction_id",json:"auction_id"`
// 	PastelID  string  `mapstructure:"pastel_id",json:"pastel_id"`
// 	BidPrice  float64 `mapstructure:"bid_price",json:"bid_price"`
// }

// type NewArtAuctionRequest struct {
// 	InstanceID  string  `mapstructure:"instance_id",json:"instance_id"`
// 	LowestPrice float64 `mapstructure:"lowest_price",json:"lowest_price"`
// }

// type ArtAuctionInfo struct {
// 	AuctionID   int64      `mapstructure:"auction_id",json:"auction_id"`
// 	InstanceID  string     `mapstructure:"instance_id",json:"instance_id"`
// 	LowestPrice float64    `mapstructure:"lowest_price",json:"lowest_price"`
// 	IsOpen      bool       `mapstructure:"is_open",json:"is_open"`
// 	StartTime   *time.Time `mapstructure:"start_time",json:"start_time"`
// 	EndTime     *time.Time `mapstructure:"end_time",json:"end_time"`
// 	FirstPrice  *float64   `mapstructure:"first_price",json:"first_price"`
// 	SecondPrice *float64   `mapstructure:"second_price",json:"second_price"`
// }

// type SNActivityInfo struct {
// 	Query string `mapstructure:"query",json:"query"`

// 	// SNActivityThumbnailRequest = "thumbnail_request"
// 	// SNActivityNftSearch        = "nft_search"
// 	// SNActivityCreatorSearch    = "creator_search"
// 	// SNActivityUserSearch       = "user_search"
// 	// SNActivityKeywordSearch    = "keyword_search"
// 	ActivityType string `mapstructure:"activity_type",json:"activity_type"`

// 	SNPastelID string `mapstructure:"sn_pastel_id",json:"sn_pastel_id"`

// 	// Cnt is the number of query activities, ignored on write request
// 	Cnt int `mapstructure:"cnt",json:"cnt"`
// }

// type SNTopActivityRequest struct {
// 	ActivityType string `mapstructure:"activity_type",json:"activity_type"`
// 	SNPastelID   string `mapstructure:"sn_pastel_id",json:"sn_pastel_id"`
// 	NRecords     int    `mapstructure:"n_records",json:"n_records"`
// }
