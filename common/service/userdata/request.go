package userdata

import "time"

const (
	SNActivityThumbnailRequest = "thumbnail_request"
	SNActivityNftSearch        = "nft_search"
	SNActivityCreatorSearch    = "creator_search"
	SNActivityUserSearch       = "user_search"
	SNActivityKeywordSearch    = "keyword_search"

	CommandUserInfoWrite            = "user_info_write"
	CommandArtInstanceAskingPrice   = "art_instance_asking_price"
	CommandUserInfoQuery            = "user_info_query"
	CommandArtInfoWrite             = "art_info_write"
	CommandArtInstanceInfoWrite     = "art_instance_info_write"
	CommandArtLikeWrite             = "art_like_write"
	CommandArtPlaceBid              = "art_place_bid"
	CommandCumulatedSalePriceByUser = "cumulated_sale_price_by_user"
	CommandEndArtAuction            = "end_art_auction"
	CommandGetAuctionInfo           = "get_auction_info"
	CommandGetFollowees             = "get_followees"
	CommandGetFollowers             = "get_followers"
	CommandGetFriend                = "get_friend"
	CommandGetInstanceInfo          = "get_instance_info"
	CommandGetTopSNActivities       = "get_top_sn_activities"
	CommandHighestSalePriceByUser   = "highest_sale_price_by_user"
	CommandNewArtAuction            = "new_art_auction"
	CommandNftCopiesExist           = "nft_copies_exist"
	CommandNftCreatedByArtist       = "nft_created_by_artist"
	CommandNftForSaleByArtist       = "nft_for_sale_by_artist"
	CommandNftOwnedByUser           = "nft_owned_by_user"
	CommandNftSoldByArtID           = "nft_sold_by_art_id"
	CommandSnActivityWrite          = "sn_activity_write"
	CommandTransactionWrite         = "transaction_write"
	CommandUniqueNftByUser          = "unique_nft_by_user"
	CommandUserFollowWrite          = "user_follow_write"
	CommandUsersLikeNft             = "users_like_nft"
)

// ProcessRequest represents userdata created by wallet node
type ProcessRequest struct {
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
	// Passphrase of the artist's PastelID
	ArtistPastelIDPassphrase string
	// Epoch Timestamp of the request (generated, not sending by UI)
	Timestamp int64
	// Previous block hash in the chain (generated, not sending by UI)
	PreviousBlockHash string
}

// ProcessRequestSigned represents userdata and signature created by wallet node
type ProcessRequestSigned struct {
	Userdata *ProcessRequest
	// UserdataHash represents ProcessRequest's hash value, to make sure ProcessRequest's integrity
	UserdataHash string
	// Signature of the UserdataHash (generated by walletnode)
	Signature string
}

// UserImageUpload store user's upload image
type UserImageUpload struct {
	// File to upload
	Content []byte
	// File name of the user image
	Filename string
}

// ProcessResult is the response walletnode receive from supernode, about the storing's result
type ProcessResult struct {
	// Result of the request is success or not
	ResponseCode int32
	// The detail of why result is success/fail, depend on response_code
	Detail string
	// Error detail on realname
	Realname string
	// Error detail on facebook_link
	FacebookLink string
	// Error detail on twitter_link
	TwitterLink string
	// Error detail on native_currency
	NativeCurrency string
	// Error detail on location
	Location string
	// Error detail on primary_language
	PrimaryLanguage string
	// Error detail on categories
	Categories string
	// Error detail on biography
	Biography string
	// Error detail on avatar
	AvatarImage string
	// Error detail on cover photo
	CoverPhoto string
}

// SuperNodeRequest represents the ProcessResult request
// together with the SuperNode info which signed the request
type SuperNodeRequest struct {
	// UserdataHash represents ProcessRequest's hash value, to make sure ProcessRequest's integrity
	UserdataHash string
	// UserdataResultHash represents ProcessResult's hash value, to make sure ProcessResult's integrity
	UserdataResultHash string
	// SuperNodeSignature is the digital signature created by supernode for the [UserdataResultHash+UserdataHash]
	HashSignature string
	// NodeID is supernode's pastelID of this supernode generate this SuperNodeRequest
	NodeID string
}

// SuperNodeReply represents the response for SuperNodeRequest
type SuperNodeReply struct {
	// Result of the SuperNodeRequest is success or not
	ResponseCode int32
	// The detail of why SuperNodeRequest is success/fail, depend on response_code
	Detail string
}

// Metric represents the metric need to be stored
type Metric struct {
	// Command the the predefine database operation name
	Command string
	// Data is the general structure accept many kind of metric, but it need to be match with Command
	Data interface{}
}

type ArtInfo struct {
	ArtID                 string  `mapstructure:"art_id",json:"art_id"`
	ArtistPastelID        string  `mapstructure:"artist_pastel_id",json:"artist_pastel_id"`
	Copies                int64   `mapstructure:"copies",json:"copies"`
	CreatedTimestamp      int64   `mapstructure:"created_timestamp",json:"created_timestamp"`
	GreenNft              bool    `mapstructure:"green_nft",json:"green_nft"`
	RarenessScore         float64 `mapstructure:"rareness_score",json:"rareness_score"`
	RoyaltyRatePercentage float64 `mapstructure:"royalty_rate_percentage",json:"royalty_rate_percentage"`
}

type ArtInstanceInfo struct {
	InstanceID    string   `mapstructure:"instance_id",json:"instance_id"`
	ArtID         string   `mapstructure:"art_id",json:"art_id"`
	OwnerPastelID string   `mapstructure:"owner_pastel_id",json:"owner_pastel_id"`
	Price         float64  `mapstructure:"price",json:"price"`
	AskingPrice   *float64 `mapstructure:"asking_price,omitempty",json:"asking_price,omitempty"`
}

type ArtLike struct {
	ArtID    string `mapstructure:"art_id",json:"art_id"`
	PastelID string `mapstructure:"pastel_id",json:"pastel_id"`
}

type ArtTransaction struct {
	TransactionID  string  `mapstructure:"transaction_id",json:"transaction_id"`
	InstanceID     string  `mapstructure:"instance_id",json:"instance_id"`
	Timestamp      int64   `mapstructure:"timestamp",json:"timestamp"`
	SellerPastelID string  `mapstructure:"seller_pastel_id",json:"seller_pastel_id"`
	BuyerPastelID  string  `mapstructure:"buyer_pastel_id",json:"buyer_pastel_id"`
	Price          float64 `mapstructure:"price",json:"price"`
}

type UserFollow struct {
	FollowerPastelID string `mapstructure:"follower_pastel_id",json:"follower_pastel_id"`
	FolloweePastelID string `mapstructure:"followee_pastel_id",json:"followee_pastel_id"`
}

type NftCreatedByArtistQueryResult struct {
	InstanceID            string  `mapstructure:"instance_id",json:"instance_id"`
	ArtID                 string  `mapstructure:"art_id",json:"art_id"`
	Copies                int64   `mapstructure:"copies",json:"copies"`
	CreatedTimestamp      int64   `mapstructure:"created_timestamp",json:"created_timestamp"`
	GreenNft              bool    `mapstructure:"green_nft",json:"green_nft"`
	RarenessScore         float64 `mapstructure:"rareness_score",json:"rareness_score"`
	RoyaltyRatePercentage float64 `mapstructure:"royalty_rate_percentage",json:"royalty_rate_percentage"`
}

type NftForSaleByArtistQueryResult struct {
	InstanceID string  `mapstructure:"instance_id",json:"instance_id"`
	ArtID      string  `mapstructure:"art_id",json:"art_id"`
	Price      float64 `mapstructure:"price",json:"price"`
}

type NftOwnedByUserQueryResult struct {
	ArtID string `mapstructure:"art_id",json:"art_id"`
	Count int    `mapstructure:"cnt",json:"cnt"`
}

type NftSoldByArtIDQueryResult struct {
	TotalCopies int `mapstructure:"total_copies",json:"total_copies"`
	SoldCopies  int `mapstructure:"sold_copies",json:"sold_copies"`
}

type UniqueNftByUserQuery struct {
	ArtistPastelID string `mapstructure:"artist_pastel_id",json:"artist_pastel_id"`
	LimitTimestamp int64  `mapstructure:"limit_timestamp",json:"limit_timestamp"`
}

type AskingPriceUpdateRequest struct {
	InstanceID  string  `mapstructure:"instance_id",json:"instance_id"`
	AskingPrice float64 `mapstructure:"asking_price",json:"asking_price"`
}

type ArtPlaceBidRequest struct {
	AuctionID int64   `mapstructure:"auction_id",json:"auction_id"`
	PastelID  string  `mapstructure:"pastel_id",json:"pastel_id"`
	BidPrice  float64 `mapstructure:"bid_price",json:"bid_price"`
}

type NewArtAuctionRequest struct {
	InstanceID  string  `mapstructure:"instance_id",json:"instance_id"`
	LowestPrice float64 `mapstructure:"lowest_price",json:"lowest_price"`
}

type ArtAuctionInfo struct {
	AuctionID   int64      `mapstructure:"auction_id",json:"auction_id"`
	InstanceID  string     `mapstructure:"instance_id",json:"instance_id"`
	LowestPrice float64    `mapstructure:"lowest_price",json:"lowest_price"`
	IsOpen      bool       `mapstructure:"is_open",json:"is_open"`
	StartTime   *time.Time `mapstructure:"start_time",json:"start_time"`
	EndTime     *time.Time `mapstructure:"end_time",json:"end_time"`
	FirstPrice  *float64   `mapstructure:"first_price",json:"first_price"`
	SecondPrice *float64   `mapstructure:"second_price",json:"second_price"`
}

type SNActivityInfo struct {
	Query string `mapstructure:"query",json:"query"`

	// SNActivityThumbnailRequest = "thumbnail_request"
	// SNActivityNftSearch        = "nft_search"
	// SNActivityCreatorSearch    = "creator_search"
	// SNActivityUserSearch       = "user_search"
	// SNActivityKeywordSearch    = "keyword_search"
	ActivityType string `mapstructure:"activity_type",json:"activity_type"`

	SNPastelID string `mapstructure:"sn_pastel_id",json:"sn_pastel_id"`

	// Cnt is the number of query activities, ignored on write request
	Cnt int `mapstructure:"cnt",json:"cnt"`
}

type SNTopActivityRequest struct {
	ActivityType string `mapstructure:"activity_type",json:"activity_type"`
	SNPastelID   string `mapstructure:"sn_pastel_id",json:"sn_pastel_id"`
	NRecords     int    `mapstructure:"n_records",json:"n_records"`
}

type IDStringQuery struct {
	ID string `mapstructure:"id",json:"id"`
}

type IDIntQuery struct {
	ID int64 `mapstructure:"id",json:"id"`
}
