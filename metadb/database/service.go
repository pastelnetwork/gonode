package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

const (
	logPrefix       = "database"
	queryLevelNone  = "nono"
	schemaDelimiter = "---"

	userInfoWriteTemplate            = "user_info_write"
	artInstanceAskingPriceTemplate   = "art_instance_asking_price"
	userInfoQueryTemplate            = "user_info_query"
	artInfoWriteTemplate             = "art_info_write"
	artInstanceInfoWriteTemplate     = "art_instance_info_write"
	artLikeWriteTemplate             = "art_like_write"
	artPlaceBidTemplate              = "art_place_bid"
	cumulatedSalePriceByUserTemplate = "cumulated_sale_price_by_user"
	endArtAuctionTemplate            = "end_art_auction"
	getAuctionInfoTemplate           = "get_auction_info"
	getFolloweesTemplate             = "get_followees"
	getFollowersTemplate             = "get_followers"
	getFriendTemplate                = "get_friend"
	getInstanceInfoTemplate          = "get_instance_info"
	getTopSNActivitiesTemplate       = "get_top_sn_activities"
	getUserBriefInfoTemplate         = "get_user_brief_info"
	highestSalePriceByUserTemplate   = "highest_sale_price_by_user"
	newArtAuctionTemplate            = "new_art_auction"
	nftCopiesExistTemplate           = "nft_copies_exist"
	nftCreatedByArtistTemplate       = "nft_created_by_artist"
	nftForSaleByArtistTemplate       = "nft_for_sale_by_artist"
	nftOwnedByUserTemplate           = "nft_owned_by_user"
	nftSoldByArtIDTemplate           = "nft_sold_by_art_id"
	snActivityWriteTemplate          = "sn_activity_write"
	transactionWriteTemplate         = "transaction_write"
	uniqueNftByUserTemplate          = "unique_nft_by_user"
	userFollowDeleteTemplate         = "user_follow_delete"
	userFollowWriteTemplate          = "user_follow_write"
	usersLikeNftTemplate             = "users_like_nft"
)

// Config is rqlite database config
type Config struct {
	SchemaPath   string `mapstructure:"schema-path" json:"schema-path,omitempty"`
	TemplatePath string `mapstructure:"template-path" json:"template-path,omitempty"`
}

// NewConfig return the new Config
func NewConfig() *Config {
	return &Config{}
}

// Ops is database operator
type Ops struct {
	metaDB    metadb.MetaDB
	templates *TemplateKeeper
	config    *Config
}

// IsLeader check if current supernode is having rqlite cluster leader
func (db *Ops) IsLeader() bool {
	return db.metaDB.IsLeader()
}

// LeaderAddress return the ipaddress of the supernode contain the leader
func (db *Ops) LeaderAddress() string {
	return db.metaDB.LeaderAddress()
}

// ProcessCommand handle different type of metric to be store
func (db *Ops) ProcessCommand(ctx context.Context, req *pb.Metric) (interface{}, error) {
	if req == nil {
		return nil, errors.Errorf("receive nil request")
	}
	rawData := req.GetData()
	if rawData == nil {
		return nil, errors.Errorf("receive nil data in request")
	}

	switch req.Command {
	case userdata.CommandUserInfoWrite:
		var data pb.UserdataRequest
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.WriteUserData(ctx, safeStringQueryStruct(&data).(*pb.UserdataRequest))

	case userdata.CommandArtInstanceAskingPrice:
		var data userdata.AskingPriceUpdateRequest
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.UpdateAskingPrice(ctx, *safeStringQueryStruct(&data).(*userdata.AskingPriceUpdateRequest))

	case userdata.CommandUserInfoQuery:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.ReadUserData(ctx, processEscapeString(data.ID))

	case userdata.CommandArtInfoWrite:
		var data userdata.ArtInfo
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.WriteArtInfo(ctx, *safeStringQueryStruct(&data).(*userdata.ArtInfo))

	case userdata.CommandArtInstanceInfoWrite:
		var data userdata.ArtInstanceInfo
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.WriteArtInstanceInfo(ctx, *safeStringQueryStruct(&data).(*userdata.ArtInstanceInfo))

	case userdata.CommandArtLikeWrite:
		var data userdata.ArtLike
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.WriteArtLike(ctx, *safeStringQueryStruct(&data).(*userdata.ArtLike))

	case userdata.CommandArtPlaceBid:
		var data userdata.ArtPlaceBidRequest
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.ArtPlaceBid(ctx, *safeStringQueryStruct(&data).(*userdata.ArtPlaceBidRequest))

	case userdata.CommandCumulatedSalePriceByUser:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetCumulatedSalePriceByUser(ctx, processEscapeString(data.ID))

	case userdata.CommandEndArtAuction:
		var data userdata.IDIntQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.EndArtAuction(ctx, data.ID)

	case userdata.CommandGetAuctionInfo:
		var data userdata.IDIntQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetAuctionInfo(ctx, data.ID)

	case userdata.CommandGetFollowees:
		var data userdata.PaginationIDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetFollowees(ctx, *safeStringQueryStruct(&data).(*userdata.PaginationIDStringQuery))

	case userdata.CommandGetFollowers:
		var data userdata.PaginationIDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetFollowers(ctx, *safeStringQueryStruct(&data).(*userdata.PaginationIDStringQuery))

	case userdata.CommandGetFriend:
		var data userdata.PaginationIDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetFriends(ctx, *safeStringQueryStruct(&data).(*userdata.PaginationIDStringQuery))

	case userdata.CommandGetInstanceInfo:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetArtInstanceInfo(ctx, processEscapeString(data.ID))

	case userdata.CommandGetTopSNActivities:
		var data userdata.SNTopActivityRequest
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetTopSNActivities(ctx, *safeStringQueryStruct(&data).(*userdata.SNTopActivityRequest))

	case userdata.CommandHighestSalePriceByUser:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetHighestSalePriceByUser(ctx, processEscapeString(data.ID))

	case userdata.CommandNewArtAuction:
		var data userdata.NewArtAuctionRequest
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.NewArtAuction(ctx, *safeStringQueryStruct(&data).(*userdata.NewArtAuctionRequest))

	case userdata.CommandNftCopiesExist:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetExistingNftCopies(ctx, processEscapeString(data.ID))

	case userdata.CommandNftCreatedByArtist:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetNftCreatedByArtist(ctx, processEscapeString(data.ID))

	case userdata.CommandNftForSaleByArtist:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetNftForSaleByArtist(ctx, processEscapeString(data.ID))

	case userdata.CommandNftOwnedByUser:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetNftOwnedByUser(ctx, processEscapeString(data.ID))

	case userdata.CommandNftSoldByArtID:
		var data userdata.IDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetNftSoldByArtID(ctx, processEscapeString(data.ID))

	case userdata.CommandSnActivityWrite:
		var data userdata.SNActivityInfo
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.WriteSNActivity(ctx, *safeStringQueryStruct(&data).(*userdata.SNActivityInfo))

	case userdata.CommandTransactionWrite:
		var data userdata.ArtTransaction
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.WriteTransaction(ctx, *safeStringQueryStruct(&data).(*userdata.ArtTransaction))

	case userdata.CommandUniqueNftByUser:
		var data userdata.UniqueNftByUserQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetUniqueNftByUser(ctx, *safeStringQueryStruct(&data).(*userdata.UniqueNftByUserQuery))

	case userdata.CommandUserFollowDelete:
		var data userdata.UserFollow
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.DeleteUserFollow(ctx, *safeStringQueryStruct(&data).(*userdata.UserFollow))

	case userdata.CommandUserFollowWrite:
		var data userdata.UserFollow
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return nil, db.WriteUserFollow(ctx, *safeStringQueryStruct(&data).(*userdata.UserFollow))

	case userdata.CommandUsersLikeNft:
		var data userdata.PaginationIDStringQuery
		if err := json.Unmarshal(rawData, &data); err != nil {
			return nil, errors.Errorf("error while unmarshaling json: %w", err)
		}
		return db.GetUsersLikeNft(ctx, *safeStringQueryStruct(&data).(*userdata.PaginationIDStringQuery))

	default:
		return nil, errors.Errorf("Unsupported command: %s", req.Command)
	}
}

// writeData writes to metadb
func (db *Ops) writeData(ctx context.Context, command string) error {
	if len(command) == 0 {
		return errors.Errorf("invalid command")
	}
	result, err := db.metaDB.Write(ctx, command)
	if err != nil {
		return errors.Errorf("error while writting to db: %w", err)
	}
	if result.Error != "" {
		return errors.Errorf("error while writting to db: %s", result.Error)
	}

	return nil
}

func (db *Ops) writeDataReturning(ctx context.Context, command string) (*metadb.WriteResult, error) {
	if len(command) == 0 {
		return nil, errors.Errorf("invalid command")
	}
	result, err := db.metaDB.Write(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while writting to db: %w", err)
	}
	if result.Error != "" {
		return nil, errors.Errorf("error while writting to db: %s", result.Error)
	}

	return result, nil
}

// WriteArtInfo writes art info to db
func (db *Ops) WriteArtInfo(ctx context.Context, data userdata.ArtInfo) error {
	command, err := db.templates.GetCommand(artInfoWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

// WriteArtInstanceInfo writes art instace info to db
func (db *Ops) WriteArtInstanceInfo(ctx context.Context, data userdata.ArtInstanceInfo) error {
	command, err := db.templates.GetCommand(artInstanceInfoWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

// WriteArtLike writes art like info to db
func (db *Ops) WriteArtLike(ctx context.Context, data userdata.ArtLike) error {
	command, err := db.templates.GetCommand(artLikeWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

// WriteTransaction writes transaction info to db
func (db *Ops) WriteTransaction(ctx context.Context, data userdata.ArtTransaction) error {
	command, err := db.templates.GetCommand(transactionWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

// WriteUserFollow writes user follow to db
func (db *Ops) WriteUserFollow(ctx context.Context, data userdata.UserFollow) error {
	command, err := db.templates.GetCommand(userFollowWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

// DeleteUserFollow delete user follow
func (db *Ops) DeleteUserFollow(ctx context.Context, data userdata.UserFollow) error {
	command, err := db.templates.GetCommand(userFollowDeleteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

// WriteSNActivity write SN activity
func (db *Ops) WriteSNActivity(ctx context.Context, data userdata.SNActivityInfo) error {
	if data.Query == "" {
		return errors.Errorf("invalid query")
	}
	command, err := db.templates.GetCommand(snActivityWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

// WriteUserData writes metadata in the struct UserdataProcessRequest to metadb
func (db *Ops) WriteUserData(ctx context.Context, data *pb.UserdataRequest) error {
	if data == nil {
		return errors.Errorf("input nil data")
	}
	command, err := db.templates.GetCommand(userInfoWriteTemplate, pbToWriteCommand(data))
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}

	return db.writeData(ctx, command)
}

// UpdateAskingPrice updates art asking price to db
func (db *Ops) UpdateAskingPrice(ctx context.Context, data userdata.AskingPriceUpdateRequest) error {
	command, err := db.templates.GetCommand(artInstanceAskingPriceTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}

	return db.writeData(ctx, command)
}

// NewArtAuction creates new art auction in db
func (db *Ops) NewArtAuction(ctx context.Context, data userdata.NewArtAuctionRequest) (int64, error) {
	command, err := db.templates.GetCommand(newArtAuctionTemplate, data)
	if err != nil {
		return 0, errors.Errorf("error while subtitute template: %w", err)
	}

	result, err := db.writeDataReturning(ctx, command)
	if err != nil {
		return 0, errors.Errorf("error while writing to db: %w", err)
	}
	return result.LastInsertID, nil
}

// EndArtAuction ends art auction in db
func (db *Ops) EndArtAuction(ctx context.Context, auctionID int64) error {
	command, err := db.templates.GetCommand(endArtAuctionTemplate, auctionID)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}

	return db.writeData(ctx, command)
}

// ArtPlaceBid place bid for an art
func (db *Ops) ArtPlaceBid(ctx context.Context, data userdata.ArtPlaceBidRequest) error {
	command, err := db.templates.GetCommand(artPlaceBidTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}

	return db.writeData(ctx, command)
}

// ReadUserData read metadata in the struct UserdataProcessRequest to metadb
func (db *Ops) ReadUserData(ctx context.Context, artistPastelID string) (userdata.ProcessRequest, error) {
	command, err := db.templates.GetCommand(userInfoQueryTemplate, artistPastelID)
	if err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while subtitute template: %w", err)
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while querying db: %w", err)
	}
	if queryResult.Err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while querying db: %w", queryResult.Err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return userdata.ProcessRequest{}, errors.Errorf("no artist with pastel id = %s", artistPastelID)
	} else if nrows > 1 {
		return userdata.ProcessRequest{}, errors.Errorf("upto %d records are returned", nrows)
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()
	resultMap, err := queryResult.Map()
	if err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while extracting result: %w", err)
	}

	var dbResult UserdataReadResult
	if err := mapstructure.Decode(resultMap, &dbResult); err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while decoding result: %w", err)
	}

	return dbResult.ToUserData(), nil
}

func (db *Ops) salePriceByUserQuery(ctx context.Context, command string) (float64, error) {
	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return 0.0, errors.Errorf("error while querying db: %w", err)
	}
	if queryResult.Err != nil {
		return 0.0, errors.Errorf("error while querying db: %w", queryResult.Err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return 0.0, nil
	} else if nrows > 1 {
		return 0.0, errors.Errorf("upto %d records are returned", nrows)
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()

	var price float64
	if err := queryResult.Scan(&price); err != nil {
		return 0.0, errors.Errorf("error occurs while scanning query result: %w", err)
	}
	return price, nil
}

// GetCumulatedSalePriceByUser get cumulated sale price by a user
func (db *Ops) GetCumulatedSalePriceByUser(ctx context.Context, pastelID string) (float64, error) {
	if pastelID == "" {
		return 0.0, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(cumulatedSalePriceByUserTemplate, pastelID)
	if err != nil {
		return 0.0, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.salePriceByUserQuery(ctx, command)
}

func (db *Ops) queryPastelIDPagination(ctx context.Context, command string) ([]string, int, error) {
	if len(command) == 0 {
		return nil, 0, errors.Errorf("invalid command")
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return nil, 0, errors.Errorf("error while querying db: %w", err)
	}
	if queryResult.Err != nil {
		return nil, 0, errors.Errorf("error while querying db: %w", queryResult.Err)
	}

	result := make([]string, 0)
	var totalCount int
	for queryResult.Next() {
		var id string
		if err := queryResult.Scan(&id, &totalCount); err != nil {
			return nil, 0, errors.Errorf("error while scaning db result: %w", err)
		}
		result = append(result, id)
	}

	return result, totalCount, nil
}

func (db *Ops) queryUserRelationshipItems(ctx context.Context, ids []string) ([]userdata.UserRelationshipItem, error) {
	if ids == nil {
		return nil, errors.Errorf("recieve nil ids")
	}
	if len(ids) == 0 {
		return make([]userdata.UserRelationshipItem, 0), nil
	}

	processedIds := make([]string, 0)
	for _, id := range ids {
		id = fmt.Sprintf("'%s'", id)
		processedIds = append(processedIds, id)
	}

	command, err := db.templates.GetCommand(getUserBriefInfoTemplate, strings.Join(processedIds, ","))
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]userdata.UserRelationshipItem, 0)
	for _, mp := range allResults {
		var record userdata.UserRelationshipItem
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *Ops) getPaginationUsersWithTemplate(ctx context.Context, data userdata.PaginationIDStringQuery, template string) (userdata.UserRelationshipQueryResult, error) {
	var result userdata.UserRelationshipQueryResult
	if data.ID == "" {
		return result, errors.Errorf("invalid pastel ID")
	}

	if data.Limit == 0 && data.Offset == 0 {
		data.Limit = 1000000000
	}

	command, err := db.templates.GetCommand(template, data)
	if err != nil {
		return result, errors.Errorf("error while subtitute template: %w", err)
	}

	ids, totalCount, err := db.queryPastelIDPagination(ctx, command)
	if err != nil {
		return result, errors.Errorf("error while subtitute template: %w", err)
	}

	items, err := db.queryUserRelationshipItems(ctx, ids)
	if err != nil {
		return result, errors.Errorf("error while querying user relationship items: %w", err)
	}

	result.TotalCount = totalCount
	result.Items = items
	return result, nil
}

// GetFollowees get followee relationship
func (db *Ops) GetFollowees(ctx context.Context, data userdata.PaginationIDStringQuery) (userdata.UserRelationshipQueryResult, error) {
	return db.getPaginationUsersWithTemplate(ctx, data, getFolloweesTemplate)
}

// GetFollowers get follower relationship
func (db *Ops) GetFollowers(ctx context.Context, data userdata.PaginationIDStringQuery) (userdata.UserRelationshipQueryResult, error) {
	return db.getPaginationUsersWithTemplate(ctx, data, getFollowersTemplate)
}

// GetFriends get friend relationship
func (db *Ops) GetFriends(ctx context.Context, data userdata.PaginationIDStringQuery) (userdata.UserRelationshipQueryResult, error) {
	return db.getPaginationUsersWithTemplate(ctx, data, getFriendTemplate)
}

// GetHighestSalePriceByUser get highest sale price by user
func (db *Ops) GetHighestSalePriceByUser(ctx context.Context, pastelID string) (float64, error) {
	if pastelID == "" {
		return 0.0, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(highestSalePriceByUserTemplate, pastelID)
	if err != nil {
		return 0.0, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.salePriceByUserQuery(ctx, command)
}

// GetExistingNftCopies get existing copies of an NFT
func (db *Ops) GetExistingNftCopies(ctx context.Context, artID string) (int, error) {
	if artID == "" {
		return 0, errors.Errorf("invalid art ID")
	}

	command, err := db.templates.GetCommand(nftCopiesExistTemplate, artID)
	if err != nil {
		return 0, errors.Errorf("error while subtitute template: %w", err)
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return 0, errors.Errorf("error while querying db: %w", err)
	}
	if queryResult.Err != nil {
		return 0, errors.Errorf("error while querying db: %w", queryResult.Err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return 0, errors.Errorf("art id does not exist or no transaction have been made: %s", artID)
	} else if nrows > 1 {
		return 0, errors.Errorf("upto %d records are returned", nrows)
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()

	var copies int
	if err := queryResult.Scan(&copies); err != nil {
		return 0, errors.Errorf("error occurs while scanning query result: %w", err)
	}
	return copies, nil
}

func (db *Ops) queryToInterface(ctx context.Context, command string) ([]map[string]interface{}, error) {
	if len(command) == 0 {
		return nil, errors.Errorf("invalid command")
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return nil, errors.Errorf("error while querying db: %w", err)
	}
	if queryResult.Err != nil {
		return nil, errors.Errorf("error while querying db: %w", queryResult.Err)
	}

	result := make([]map[string]interface{}, 0)
	for queryResult.Next() {
		resultMap, err := queryResult.Map()
		if err != nil {
			return nil, errors.Errorf("error while extracting result: %w", err)
		}
		result = append(result, resultMap)
	}

	return result, nil
}

// GetNftCreatedByArtist acquire all NFTs created by an artist
func (db *Ops) GetNftCreatedByArtist(ctx context.Context, artistPastelID string) ([]userdata.NftCreatedByArtistQueryResult, error) {
	if artistPastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(nftCreatedByArtistTemplate, artistPastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]userdata.NftCreatedByArtistQueryResult, 0)
	for _, mp := range allResults {
		var record userdata.NftCreatedByArtistQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

// GetNftForSaleByArtist get NFT for sale byan artist
func (db *Ops) GetNftForSaleByArtist(ctx context.Context, artistPastelID string) ([]userdata.NftForSaleByArtistQueryResult, error) {
	if artistPastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(nftForSaleByArtistTemplate, artistPastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]userdata.NftForSaleByArtistQueryResult, 0)
	for _, mp := range allResults {
		var record userdata.NftForSaleByArtistQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

// GetNftOwnedByUser get NFTs owned by an user
func (db *Ops) GetNftOwnedByUser(ctx context.Context, pastelID string) ([]userdata.NftOwnedByUserQueryResult, error) {
	if pastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(nftOwnedByUserTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]userdata.NftOwnedByUserQueryResult, 0)
	for _, mp := range allResults {
		var record userdata.NftOwnedByUserQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

// GetNftSoldByArtID get sold NFT
func (db *Ops) GetNftSoldByArtID(ctx context.Context, artID string) (userdata.NftSoldByArtIDQueryResult, error) {
	var result userdata.NftSoldByArtIDQueryResult
	if artID == "" {
		return result, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(nftSoldByArtIDTemplate, artID)
	if err != nil {
		return result, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return result, errors.Errorf("error while query db: %w", err)
	}

	nrows := len(allResults)
	if nrows == 0 {
		return result, errors.Errorf("no instance found for id: %s", artID)
	} else if nrows > 1 {
		return result, errors.Errorf("upto %d records returned for id: %s", nrows, artID)
	}

	if err := mapstructure.Decode(allResults[0], &result); err != nil {
		return result, errors.Errorf("error while decoding result: %w", err)
	}

	return result, nil
}

// GetUniqueNftByUser get unique NFT by user
func (db *Ops) GetUniqueNftByUser(ctx context.Context, query userdata.UniqueNftByUserQuery) ([]userdata.ArtInfo, error) {
	command, err := db.templates.GetCommand(uniqueNftByUserTemplate, query)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]userdata.ArtInfo, 0)
	for _, mp := range allResults {
		var record userdata.ArtInfo
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

// GetUsersLikeNft get users like the NFT
func (db *Ops) GetUsersLikeNft(ctx context.Context, data userdata.PaginationIDStringQuery) (userdata.UserRelationshipQueryResult, error) {
	return db.getPaginationUsersWithTemplate(ctx, data, usersLikeNftTemplate)
}

// GetArtInstanceInfo acquires art instance info
func (db *Ops) GetArtInstanceInfo(ctx context.Context, instanceID string) (userdata.ArtInstanceInfo, error) {
	var result userdata.ArtInstanceInfo
	if instanceID == "" {
		return result, errors.Errorf("invalid instance ID")
	}

	command, err := db.templates.GetCommand(getInstanceInfoTemplate, instanceID)
	if err != nil {
		return result, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return result, errors.Errorf("error while query db: %w", err)
	}

	nrows := len(allResults)
	if nrows == 0 {
		return result, errors.Errorf("no instance found for id: %s", instanceID)
	} else if nrows > 1 {
		return result, errors.Errorf("upto %d records returned for id: %s", nrows, instanceID)
	}

	if err := mapstructure.Decode(allResults[0], &result); err != nil {
		return result, errors.Errorf("error while decoding result: %w", err)
	}

	return result, nil
}

// GetAuctionInfo acuire auction info
func (db *Ops) GetAuctionInfo(ctx context.Context, auctionID int64) (userdata.ArtAuctionInfo, error) {
	var result userdata.ArtAuctionInfo
	command, err := db.templates.GetCommand(getAuctionInfoTemplate, auctionID)
	if err != nil {
		return result, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return result, errors.Errorf("error while query db: %w", err)
	}

	nrows := len(allResults)
	if nrows == 0 {
		return result, errors.Errorf("no instance found for id: %d", auctionID)
	} else if nrows > 1 {
		return result, errors.Errorf("upto %d records returned for id: %d", nrows, auctionID)
	}

	if err := mapstructure.Decode(allResults[0], &result); err != nil {
		return result, errors.Errorf("error while decoding result: %w", err)
	}

	return result, nil
}

// GetTopSNActivities get top SN activities
func (db *Ops) GetTopSNActivities(ctx context.Context, query userdata.SNTopActivityRequest) ([]userdata.SNActivityInfo, error) {
	command, err := db.templates.GetCommand(getTopSNActivitiesTemplate, query)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]userdata.SNActivityInfo, 0)
	for _, mp := range allResults {
		var record userdata.SNActivityInfo
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

// Run run the rqlite database service
func (db *Ops) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)
	log.WithContext(ctx).Info("start initialization")

	content, err := ioutil.ReadFile(db.config.SchemaPath)
	if err != nil {
		return errors.Errorf("error while reading schema file: %w", err)
	}

	db.metaDB.WaitForStarting()
	if err := db.metaDB.EnableFKConstraints(true); err != nil {
		return errors.Errorf("error while enabling db fk constraint: %w", err)
	}
	if db.metaDB.IsLeader() {
		listOfCommands := strings.Split(string(content), schemaDelimiter)
		for _, cmd := range listOfCommands {
			result, err := db.metaDB.Write(ctx, cmd)
			if err != nil {
				return errors.Errorf("error while creating db schema: %w", err)
			}
			if result.Error != "" {
				return errors.Errorf("error while creating db schema: %s", result.Error)
			}
		}
	}

	db.templates, err = NewTemplateKeeper(db.config.TemplatePath)
	if err != nil {
		fmt.Println(db.config.TemplatePath)
		return errors.Errorf("error while creating new template keeper: %w", err)
	}

	log.WithContext(ctx).Info("done initialization")
	// block until context is done
	<-ctx.Done()
	log.WithContext(ctx).Info("userdata service is stopped")
	return nil
}

// NewDatabaseOps return the Ops
func NewDatabaseOps(metaDB metadb.MetaDB, config *Config) *Ops {
	return &Ops{
		metaDB: metaDB,
		config: config,
	}
}
