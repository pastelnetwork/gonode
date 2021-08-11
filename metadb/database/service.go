package database

import (
	"context"
	"io/ioutil"
	"regexp"
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
	queryLevelNone  = "none"
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
	userFollowWriteTemplate          = "user_follow_write"
	usersLikeNftTemplate             = "users_like_nft"
)

// Config is rqlite database config
type Config struct {
	SchemaPath   string `mapstructure:"schema-path" json:"schema-path,omitempty"`
	TemplatePath string `mapstructure:"template-path" json:"write-template-path,omitempty"`
}

// NewConfig return the new Config
func NewConfig() *Config {
	return &Config{}
}

type Ops struct {
	metaDB    metadb.MetaDB
	templates *templateKeeper
	config    *Config
}

// IsLeader check if current supernode is having rqlite cluster leader
func (db *Ops) IsLeader() bool {
	return db.metaDB.IsLeader()
}

// LeaderAddress return the ipaddress of the supernode contain the leader
func (db *Ops) LeaderAddress() string {
	re := regexp.MustCompile(":[0-9]+$")
	address := re.Split(db.metaDB.LeaderAddress(), -1)[0]
	return address
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

func (db *Ops) WriteArtInfo(ctx context.Context, data ArtInfo) error {
	command, err := db.templates.GetCommand(artInfoWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

func (db *Ops) WriteArtInstanceInfo(ctx context.Context, data ArtInstanceInfo) error {
	command, err := db.templates.GetCommand(artInstanceInfoWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

func (db *Ops) WriteArtLike(ctx context.Context, data ArtLike) error {
	command, err := db.templates.GetCommand(artLikeWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

func (db *Ops) WriteTransaction(ctx context.Context, data ArtTransaction) error {
	command, err := db.templates.GetCommand(transactionWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

func (db *Ops) WriteUserFollow(ctx context.Context, data UserFollow) error {
	command, err := db.templates.GetCommand(userFollowWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return db.writeData(ctx, command)
}

func (db *Ops) WriteSNActivity(ctx context.Context, data SNActivityInfo) error {
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

func (db *Ops) UpdateAskingPrice(ctx context.Context, data AskingPriceUpdateRequest) error {
	command, err := db.templates.GetCommand(artInstanceAskingPriceTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}

	return db.writeData(ctx, command)
}

func (db *Ops) NewArtAuction(ctx context.Context, data NewArtAuctionRequest) (int64, error) {
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

func (db *Ops) EndArtAuction(ctx context.Context, auctionID int64) error {
	command, err := db.templates.GetCommand(endArtAuctionTemplate, auctionID)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}

	return db.writeData(ctx, command)
}

func (db *Ops) ArtPlaceBid(ctx context.Context, data ArtPlaceBidRequest) error {
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

func (db *Ops) queryPastelID(ctx context.Context, command string) ([]string, error) {
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

	result := make([]string, 0)
	for queryResult.Next() {
		var id string
		if err := queryResult.Scan(&id); err != nil {
			return nil, errors.Errorf("error while scaning db result: %w", err)
		}
		result = append(result, id)
	}

	return result, nil
}

func (db *Ops) GetFollowees(ctx context.Context, pastelID string) ([]string, error) {
	if pastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(getFolloweesTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

func (db *Ops) GetFollowers(ctx context.Context, pastelID string) ([]string, error) {
	if pastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(getFollowersTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

func (db *Ops) GetFriends(ctx context.Context, pastelID string) ([]string, error) {
	if pastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(getFriendTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

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

func (db *Ops) GetNftCreatedByArtist(ctx context.Context, artistPastelID string) ([]NftCreatedByArtistQueryResult, error) {
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

	result := make([]NftCreatedByArtistQueryResult, 0)
	for _, mp := range allResults {
		var record NftCreatedByArtistQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *Ops) GetNftForSaleByArtist(ctx context.Context, artistPastelID string) ([]NftForSaleByArtistQueryResult, error) {
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

	result := make([]NftForSaleByArtistQueryResult, 0)
	for _, mp := range allResults {
		var record NftForSaleByArtistQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *Ops) GetNftOwnedByUser(ctx context.Context, pastelID string) ([]NftOwnedByUserQueryResult, error) {
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

	result := make([]NftOwnedByUserQueryResult, 0)
	for _, mp := range allResults {
		var record NftOwnedByUserQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *Ops) GetNftSoldByArtID(ctx context.Context, artID string) (NftSoldByArtIDQueryResult, error) {
	var result NftSoldByArtIDQueryResult
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

func (db *Ops) GetUniqueNftByUser(ctx context.Context, query UniqueNftByUserQuery) ([]ArtInfo, error) {
	command, err := db.templates.GetCommand(uniqueNftByUserTemplate, query)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]ArtInfo, 0)
	for _, mp := range allResults {
		var record ArtInfo
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *Ops) GetUsersLikeNft(ctx context.Context, artID string) ([]string, error) {
	if artID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(usersLikeNftTemplate, artID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

func (db *Ops) GetArtInstanceInfo(ctx context.Context, instanceID string) (ArtInstanceInfo, error) {
	var result ArtInstanceInfo
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

func (db *Ops) GetAuctionInfo(ctx context.Context, auctionID int64) (ArtAuctionInfo, error) {
	var result ArtAuctionInfo
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

func (db *Ops) GetTopSNActivities(ctx context.Context, query SNTopActivityRequest) ([]SNActivityInfo, error) {
	command, err := db.templates.GetCommand(getTopSNActivitiesTemplate, query)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]SNActivityInfo, 0)
	for _, mp := range allResults {
		var record SNActivityInfo
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
